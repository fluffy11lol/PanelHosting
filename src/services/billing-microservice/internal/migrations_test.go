package internal

import (
	"billing-microservice/pkg/logger"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return db, mock
}

func TestInsertData_WithLogger(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := logger.New("TestService")
	ctx := logger.WithLogger(context.Background(), log)

	mock.ExpectQuery(`SELECT \* FROM "tariffs" WHERE id = \$1`).
		WithArgs("1").
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "tariffs"`).
		WithArgs("1", "Host-0", 1, 1, 4, 500).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	InsertData(ctx, db)
}
