package modelsControlPanel

import (
	"crypto/rand"
	"log"
)

type ServerDatabases struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `gorm:"not null"`                 // ID пользователя из таблицы users
	ServerID string `gorm:"not null"`                 // ID сервера
	DBName   string `gorm:"size:255;not null;unique"` // Имя базы данных MySQL
	Username string `gorm:"size:255;not null;unique"` // Имя пользователя MySQL
	Password string `gorm:"size:255;not null"`
}

// Генерация случайной строки для использования в dbname, username, password
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Error generating random string: %v", err)
	}
	for i := 0; i < n; i++ {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}

// Функция генерации данных для новой базы данных и пользователя
func GenerateDatabaseCredentials() (dbname string, username string, password string) {
	dbname = "db_" + generateRandomString(8)
	username = "user_" + generateRandomString(8)
	password = generateRandomString(12) // Делаем пароль достаточно сложным
	return dbname, username, password
}
