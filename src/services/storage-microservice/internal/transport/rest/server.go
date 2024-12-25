package rest

import (
	"archive/zip"
	"context"
	"fmt"
	_ "github.com/distribution/reference"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"storage-microservice/internal/middleware"
	"storage-microservice/pkg/logger"
	"strings"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Устанавливаем CORS-заголовки
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082") // Укажите конкретный домен, если нужно
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Если это preflight-запрос (OPTIONS), отправляем пустой ответ
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Передаем управление следующему обработчику
		c.Next()
	}
}

func EnsureBucketExists(client *minio.Client, bucketName string) error {
	ctx := context.Background()

	// Проверяем, существует ли бакет
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	// Создаем бакет, если он не существует
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("could not create bucket: %w", err)
		}
		fmt.Printf("Bucket '%s' created successfully\n", bucketName)
	} else {
		fmt.Printf("Bucket '%s' already exists\n", bucketName)
	}

	return nil
}

type Server struct {
	httpServer *http.Server
}

func NewS3Client(endpoint, accessKeyID, secretAccessKey string) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func New(ctx context.Context, port int, s3Client *minio.Client, bucketName string) (*Server, error) {
	// Убедиться, что бакет существует
	err := EnsureBucketExists(s3Client, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	// Инициализация роутера
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery(), gin.Logger())
	router.Use(CORSMiddleware())
	// Middleware авторизации
	router.Use(middleware.Authorized())

	// Регистрация маршрутов
	RegisterRoutes(router, s3Client, bucketName)

	// HTTP сервер
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &Server{httpServer}, nil
}

// Start запускает HTTP сервер
func (s *Server) Start(ctx context.Context) error {
	logger.GetLoggerFromCtx(ctx).Info(ctx, "starting HTTP server", zap.String("port", s.httpServer.Addr))

	// Запуск HTTP сервера
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to start server", zap.Error(err))
		}
	}()

	// Ожидание завершения сервера
	<-ctx.Done()
	return nil
}

// Stop останавливает HTTP сервер
func (s *Server) Stop(ctx context.Context) error {
	logger.GetLoggerFromCtx(ctx).Info(ctx, "stopping HTTP server")

	// Устанавливаем тайм-аут для завершения операций
	ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Завершаем работу сервера
	return s.httpServer.Shutdown(ctxShutdown)
}

func RegisterRoutes(router *gin.Engine, client *minio.Client, bucketName string) {
	router.POST("/upload", func(c *gin.Context) {
		auth, err := middleware.ReadCookie("token", c.Request)
		_, userID, _ := middleware.ParseToken(auth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Получаем загруженный архив
		file, err := c.FormFile("archive")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer src.Close()

		// Читаем архив
		zipReader, err := zip.NewReader(src, file.Size)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid zip file"})
			return
		}

		// Проходимся по всем файлам в архиве
		for _, zipFile := range zipReader.File {
			if zipFile.FileInfo().IsDir() {
				continue // Пропускаем папки
			}

			// Открываем файл внутри архива
			zipFileReader, err := zipFile.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read file %s", zipFile.Name)})
				return
			}
			defer zipFileReader.Close()

			// Убираем корневую директорию из пути zipFile.Name
			relativePath := zipFile.Name
			if idx := strings.Index(relativePath, "/"); idx != -1 {
				relativePath = relativePath[idx+1:] // Убираем всё до первого "/"
			}

			// Новый путь с userID на первом уровне
			objectName := fmt.Sprintf("%s/%s", userID, relativePath)

			// Загружаем файл в S3
			_, err = client.PutObject(c.Request.Context(), bucketName, objectName, zipFileReader, int64(zipFile.UncompressedSize64), minio.PutObjectOptions{
				ContentType: "application/octet-stream",
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %s", zipFile.Name)})
				return
			}
		}

		// Ответ при успешной загрузке
		c.JSON(http.StatusOK, gin.H{
			"message": "Archive uploaded and unpacked successfully",
		})
	})

	router.GET("/download", func(c *gin.Context) {
		auth, err := middleware.ReadCookie("token", c.Request)
		_, userID, _ := middleware.ParseToken(auth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Путь к директории в хранилище
		prefix := fmt.Sprintf("%s/", userID)

		// Локальная директория
		localDir := fmt.Sprintf("./projects/%s", userID)

		// Проверяем существование локальной папки
		if _, err := os.Stat(localDir); err == nil {
			os.RemoveAll(localDir) // Удаляем существующую папку перед скачиванием
		}

		err = os.MkdirAll(localDir, os.ModePerm) // Создаем директорию
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create local directory"})
			return
		}

		// Получаем список объектов с указанным префиксом
		objectCh := client.ListObjects(c.Request.Context(), bucketName, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		})

		for object := range objectCh {
			if object.Err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list object: %v", object.Err)})
				return
			}

			// Путь к локальному файлу
			localFilePath := fmt.Sprintf("%s/%s", localDir, object.Key[len(prefix):])

			// Создаем вложенные директории, если необходимо
			if err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subdirectory"})
				return
			}

			// Скачиваем объект
			objectReader, err := client.GetObject(c.Request.Context(), bucketName, object.Key, minio.GetObjectOptions{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to download object: %s", object.Key)})
				return
			}
			defer objectReader.Close()

			// Создаем файл на локальной машине
			localFile, err := os.Create(localFilePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create file: %s", localFilePath)})
				return
			}
			defer localFile.Close()

			// Копируем данные в файл
			if _, err = io.Copy(localFile, objectReader); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %s", localFilePath)})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Folder downloaded successfully",
		})
	})

	router.DELETE("/delete", func(c *gin.Context) {
		auth, err := middleware.ReadCookie("token", c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		_, userID, _ := middleware.ParseToken(auth)
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		// Путь к удалению
		dirPath := fmt.Sprintf("%s/", userID)

		// Получаем список объектов с указанным префиксом
		objectCh := client.ListObjects(c.Request.Context(), bucketName, minio.ListObjectsOptions{
			Prefix:    dirPath,
			Recursive: true,
		})

		for object := range objectCh {
			if object.Err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list object: %v", object.Err)})
				return
			}

			// Удаляем объект
			err := client.RemoveObject(c.Request.Context(), bucketName, object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete object: %s", object.Key)})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Files deleted successfully",
		})
	})

	// Новый маршрут для запуска Docker из docker-compose
	router.POST("/run-docker-compose", func(c *gin.Context) {
		auth, err := middleware.ReadCookie("token", c.Request)
		_, userID, _ := middleware.ParseToken(auth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		err = runDockerFromCompose(c, userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
	})
	router.POST("/stop-docker-compose", func(c *gin.Context) {
		auth, err := middleware.ReadCookie("token", c.Request)
		_, userID, _ := middleware.ParseToken(auth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Останавливаем контейнер
		err = stopDockerCompose(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Docker containers stopped successfully",
		})
	})

}

func runDockerFromCompose(c *gin.Context, userID string) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	projectPath, _ := filepath.Abs(filepath.Join("./projects", userID))
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", projectPath)
	}
	if err := os.Chdir(projectPath); err != nil {
		return err
	}

	// Проверяем наличие файла docker-compose.yaml
	if _, err := os.Stat("docker-compose.yaml"); os.IsNotExist(err) {
		return err
	}

	// Создаем команду docker-compose up --build
	cmd := exec.Command("docker-compose", "up", "--build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Выполняем команду
	if err := cmd.Run(); err != nil {
		return err
	}
	if err := os.Chdir(originalDir); err != nil {
		return fmt.Errorf("failed to return to original directory: %w", err)
	}
	return nil
}

func stopDockerCompose(userID string) error {
	if err := os.Chdir("../.."); err != nil {
		return fmt.Errorf("failed to return to root directory: %w", err)
	}
	projectPath := filepath.Join("./projects", userID)
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	fmt.Println(originalDir)
	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", projectPath, err)
	}
	if _, err := os.Stat("docker-compose.yaml"); os.IsNotExist(err) {
		return fmt.Errorf("docker-compose.yaml not found in %s", projectPath)
	}
	cmd := exec.Command("docker-compose", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Выполняем команду
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop docker-compose project: %w", err)
	}

	return nil
}
