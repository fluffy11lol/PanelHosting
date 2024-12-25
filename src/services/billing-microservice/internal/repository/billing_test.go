package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}

	return db, mock
}

func TestBillingRepository_GetTariffs_Success(t *testing.T) {

	db, mock := setupTestDB(t)
	defer func() { _ = mock.ExpectationsWereMet() }()

	repo := NewBillingRepository(db)
	ctx := context.Background()

	expectedRows := sqlmock.NewRows([]string{"id", "name", "ssd", "cpu", "ram", "price"}).
		AddRow("1", "Basic", 256, 2, 4, 100).
		AddRow("2", "Premium", 512, 4, 8, 200)

	mock.ExpectQuery(`SELECT \* FROM "tariffs"`).WillReturnRows(expectedRows)

	tariffs, err := repo.GetTariffs(ctx)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	// Проверка результатов
	if len(*tariffs) != 2 {
		t.Errorf("Expected 2 tariffs, but got %d", len(*tariffs))
	}

	if (*tariffs)[0].Name != "Basic" || (*tariffs)[1].Name != "Premium" {
		t.Errorf("Unexpected tariff names: got %+v", tariffs)
	}
}

func TestBillingRepository_GetTariffs_Error(t *testing.T) {

	db, mock := setupTestDB(t)
	defer func() { _ = mock.ExpectationsWereMet() }()

	repo := NewBillingRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* FROM "tariffs"`).WillReturnError(errors.New("db error"))

	tariffs, err := repo.GetTariffs(ctx)
	if err == nil {
		t.Fatal("Expected error, but got none")
	}

	if tariffs != nil {
		t.Errorf("Expected nil tariffs, but got %+v", tariffs)
	}
}
