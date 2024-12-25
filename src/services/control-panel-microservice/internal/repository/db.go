package databaseContronPanel

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	models "control-panel/internal/models"
	"control-panel/pkg/api/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PanelStorage struct {
	DB      *gorm.DB // Для PostgreSQL
	DBMySQL *gorm.DB
}

func (p *PanelStorage) StartDB(ctx context.Context, dsnPSQL, dsnMySQL string) (*PanelStorage, error) {
	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsnPSQL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// DBMySQL, err := gorm.Open(mysql.Open(dsnMySQL), &gorm.Config{})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to MySQL database: %v", err)
	// }
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("MySQL Info:", dsnMySQL))

	return &PanelStorage{DB: db}, nil
}

func (p *PanelStorage) CreateUser(username, password string) (int, error) {
	user := &models.User{
		Username: username,
		Password: password,
	}

	err := p.DB.Create(user).Error
	if err != nil {
		return -1, err
	}

	return int(user.ID), nil
}

func (p *PanelStorage) GetUserByName(username string) (int, string, error) {
	var user models.User
	result := p.DB.Preload("Servers").Where("username = ?", username).First(&user)
	if result.Error != nil {
		fmt.Printf("Error fetching user: %v\n", result.Error)
		return 0, "", result.Error
	}

	return int(user.ID), user.Password, nil
}

func (p *PanelStorage) CreateServer(name, address string, userID int) (*models.Server, error) {
	id_hash, err := generateRandomID(12)
	if err != nil {
		return nil, err
	}
	server := &models.Server{
		ID:           id_hash,
		Name:         name,
		Address:      address,
		UserID:       uint(userID),
		Status:       "Offline",
		TariffStatus: "Active",
	}

	result := p.DB.Create(server)
	if result.Error != nil {
		return nil, result.Error
	}
	return server, nil
}

func (p *PanelStorage) GetServers(userID uint) ([]models.Server, error) {
	var servers []models.Server
	result := p.DB.Where("user_id = ?", userID).Find(&servers)
	if result.Error != nil {
		return nil, result.Error
	}
	return servers, nil
}

func (p *PanelStorage) GetServerDetails(serverID string) (*models.Server, error) {
	var server models.Server
	result := p.DB.Where("id =?", serverID).First(&server)
	if result.Error != nil {
		return nil, result.Error
	}

	return &server, nil
}

func generateServerID(serverName, userID string) string {
	hash := sha256.New()
	hash.Write([]byte(serverName + userID))
	return hex.EncodeToString(hash.Sum(nil))
}

func generateRandomID(byteLength int) (string, error) {
	// Создаем буфер для случайных байтов
	bytes := make([]byte, byteLength)

	// Заполняем буфер случайными байтами
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	// Преобразуем байты в строку в формате hex
	return hex.EncodeToString(bytes), nil
}
