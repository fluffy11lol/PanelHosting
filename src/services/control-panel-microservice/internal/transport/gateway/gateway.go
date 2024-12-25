package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	models "control-panel/internal/models"

	"github.com/golang-jwt/jwt/v4"
)

type ServerJSON struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Port         string `json:"port"`
	ID           string `json:"id"`
	Status       string `json:"status"`
	TariffStatus string `json:"tariffstatus"`
	CreatedAt    string `json:"created_at"`
}

type DashboardData struct {
	Servers []ServerJSON `json:"servers"`
}

func renderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func Dashboardhandler(w http.ResponseWriter, r *http.Request) {

	// Проверяем, есть ли токен в куке
	auth, err := ReadCookie("token", r)
	if err != nil {
		http.Redirect(w, r, "http://localhost:8081/authentication", http.StatusFound)
		return
	}
	if auth == "" {
		http.Redirect(w, r, "http://localhost:8081/authentication", http.StatusFound)
		return
	}
	_, _, err = ParseToken(auth)
	if err != nil {
		logErr := fmt.Sprint("Token validating failed:", err)
		log.Println(logErr)
		http.Redirect(w, r, "http://localhost:8081/authentication", http.StatusFound)
		return
	}
	renderTemplate(w, "/app/static/dashboard.html", nil)

}

func ServerHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID сервера из URL
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Server ID is missing", http.StatusBadRequest)
		return
	}
	serverID := parts[2] // Берем ID из третьего элемента пути

	// Проверяем, есть ли токен в куке
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value == "" {
		// Если токен отсутствует, редиректим на страницу логина
		http.Redirect(w, r, "http://localhost:8081/authentication", http.StatusFound)
		return
	}

	_, err = jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("invalid token format")
		}

		id, ok := claims["id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid field id")
		}

		if err != nil {
			return nil, fmt.Errorf("invalid field id: %v", err)
		}
		return []byte(id), nil
	})
	if err != nil {
		logErr := fmt.Sprint("Token validating failed:", err)
		log.Println(logErr)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method == http.MethodGet {
		// Выполняем запрос к API Gateway с токеном для получения сервера
		apiGatewayURL := fmt.Sprintf("http://localhost:%s/v1/servers/%s", models.EnvsVars.Rest_port, serverID) // Укажите правильный URL вашего API Gateway

		// Делаем GET запрос к API Gateway
		req, err := http.NewRequest("GET", apiGatewayURL, nil)
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			renderTemplate(w, "./ui/static/dashboard.html", nil)
			return
		}

		// Добавляем заголовок Authorization с токеном Bearer
		bearerToken := "Bearer " + cookie.Value // Используем токен из cookie
		req.Header.Add("Authorization", bearerToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Ошибка при выполнении запроса к API Gateway:", err)
			renderTemplate(w, "./ui/static/server.html", nil)
			return
		}
		// Проверяем, что API вернул успешный ответ
		if resp.StatusCode != http.StatusOK {
			logErr := fmt.Sprint("Ошибка ответа от API Gateway: ", resp.Status)
			log.Println(logErr)
			renderTemplate(w, "./ui/static/server.html", nil)
			return
		}
		// Декодируем ответ
		var response struct {
			Server ServerJSON `json:"server"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			http.Error(w, "Failed to parse server data", http.StatusInternalServerError)
			return
		}
		serverData := response.Server

		// Рендерим шаблон страницы сервера
		renderTemplate(w, "./ui/static/server.html", serverData)
	}

	if r.Method == http.MethodPost {
		uploadFileURL := fmt.Sprintf("http://localhost:%s/v1/servers/upload", models.EnvsVars.Rest_port)

		// Делаем POST запрос к API Gateway
		req, err := http.NewRequest("POST", uploadFileURL, nil)
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			renderTemplate(w, "./ui/static/dashboard.html", nil)
			return
		}

		// Добавляем заголовок Authorization с токеном Bearer
		bearerToken := "Bearer " + cookie.Value // Используем токен из cookie
		req.Header.Add("Authorization", bearerToken)
	}
}
func ParseToken(accessToken string) (*jwt.Token, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("My Key"), nil
	})
	if err != nil {
		return nil, "", err
	}
	claims, ok := token.Claims.(*models.TokenClaim)
	if !ok {
		return nil, "", errors.New("failed to parse claims")
	}
	return token, claims.UserID, nil
}

func ReadCookie(name string, r *http.Request) (string, error) {
	if name == "" {
		return "", errors.New("you are trying to read an empty cookie")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	str := cookie.Value
	value, _ := url.QueryUnescape(str)
	return value, nil
}
