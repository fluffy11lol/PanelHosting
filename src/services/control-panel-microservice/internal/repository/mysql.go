package databaseContronPanel

import (
	"fmt"

	models "control-panel/internal/models"

	"gorm.io/gorm"
)

func executeMySQLQuery(db *gorm.DB, query string) error {
	if err := db.Exec(query).Error; err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}
	return nil
}

func CreateMySQLDatabase(db *gorm.DB, userID uint, serverID string, dbName, username, password string) error {
	// Команда для создания базы данных в MySQL
	createDBQuery := fmt.Sprintf("CREATE DATABASE %s;", dbName)
	// Команда для создания пользователя
	createUserQuery := fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s';", username, password)
	// Команда для выдачи прав
	grantPrivilegesQuery := fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%';", dbName, username)

	// Выполняем команды в MySQL
	if err := executeMySQLQuery(db, createDBQuery); err != nil {
		return err
	}
	if err := executeMySQLQuery(db, createUserQuery); err != nil {
		return err
	}
	if err := executeMySQLQuery(db, grantPrivilegesQuery); err != nil {
		return err
	}

	// Сохраняем данные о базе в PostgreSQL
	mysqlDB := models.ServerDatabases{
		UserID:   userID,
		ServerID: serverID,
		DBName:   dbName,
		Username: username,
		Password: password,
	}
	if err := db.Create(&mysqlDB).Error; err != nil {
		return err
	}
	return nil
}

func GetMySqlDatabases(db *gorm.DB, userID uint, serverID string) ([]models.ServerDatabases, error) {
	var databases []models.ServerDatabases
	if err := db.Where("user_id =? AND server_id =?", userID, serverID).Find(&databases).Error; err != nil {
		return nil, err
	}
	return databases, nil
}

func DeleteMySqlDatabase(db *gorm.DB, userID uint, serverName string) (string, error) {
	var mysqlDB models.ServerDatabases
	if err := db.Where("user_id =? AND server_id =?", userID, serverName).First(&mysqlDB).Error; err != nil {
		return "", err
	}

	// Удаляем пользователя из MySQL
	dropUserQuery := fmt.Sprintf("DROP USER '%s'@'%%';", mysqlDB.Username)
	if err := executeMySQLQuery(db, dropUserQuery); err != nil {
		return "", err
	}

	// Удаляем базу из MySQL
	dropDBQuery := fmt.Sprintf("DROP DATABASE %s;", mysqlDB.DBName)
	if err := executeMySQLQuery(db, dropDBQuery); err != nil {
		return "", err
	}

	// Удаляем данные о базе из PostgreSQL
	if err := db.Delete(&mysqlDB).Error; err != nil {
		return "", err
	}

	return mysqlDB.DBName, nil
}
