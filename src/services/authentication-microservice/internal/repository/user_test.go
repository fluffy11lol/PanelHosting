package repository

import (
	"context"
	"sync"
	"testing"
	"time"

	"authentication-microservice/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("create new user", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
		}

		_, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)

		var found models.User
		err = db.Where("username = ?", user.Username).First(&found).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Username, found.Username)
		assert.Equal(t, user.Email, found.Email)
	})

	t.Run("create duplicate user", func(t *testing.T) {
		user := models.User{
			Username: "duplicate",
			Password: "password123",
			Email:    "duplicate@example.com",
		}

		_, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)

		_, err = repo.CreateUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, models.ErrUserExist, err)
	})

	t.Run("create user with special characters", func(t *testing.T) {
		user := models.User{
			Username: "test@user#$%",
			Password: "password123!@#",
			Email:    "special.test@example.com",
		}

		_, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)

		var found models.User
		err = db.Where("username = ?", user.Username).First(&found).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Username, found.Username)
	})

	t.Run("create user with very long values", func(t *testing.T) {
		user := models.User{
			Username: string(make([]byte, 50)),
			Password: string(make([]byte, 72)),
			Email:    "very.long.email@really.long.domain.example.com",
		}

		_, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)

		var found models.User
		err = db.Where("username = ?", user.Username).First(&found).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Username, found.Username)
	})

	t.Run("create user with unicode characters", func(t *testing.T) {
		user := models.User{
			Username: "测试用户",
			Password: "密码123",
			Email:    "test@例子.com",
		}

		_, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)

		var found models.User
		err = db.Where("username = ?", user.Username).First(&found).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Username, found.Username)
	})

	t.Run("create user with empty fields", func(t *testing.T) {
		testCases := []struct {
			name    string
			user    models.User
			wantErr error
		}{
			{
				name: "empty username",
				user: models.User{
					Username: "",
					Password: "password123",
					Email:    "test@example.com",
				},
				wantErr: models.ErrEmptyField,
			},
			{
				name: "empty password",
				user: models.User{
					Username: "testuser",
					Password: "",
					Email:    "test@example.com",
				},
				wantErr: models.ErrEmptyField,
			},
			{
				name: "empty email",
				user: models.User{
					Username: "testuser",
					Password: "password123",
					Email:    "",
				},
				wantErr: models.ErrEmptyField,
			},
			{
				name: "all empty",
				user: models.User{
					Username: "",
					Password: "",
					Email:    "",
				},
				wantErr: models.ErrEmptyField,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := repo.CreateUser(ctx, tc.user)
				assert.ErrorIs(t, err, tc.wantErr)
			})
		}
	})
}

func TestCreateUserExtended(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("very long password", func(t *testing.T) {
		user := models.User{
			Username: "longpass",
			Password: string(make([]byte, 73)),
			Email:    "long@example.com",
		}
		_, err := repo.CreateUser(ctx, user)
		assert.Error(t, err)
	})

	t.Run("minimal valid values", func(t *testing.T) {
		user := models.User{
			Username: "a",
			Password: "1",
			Email:    "a@b.c",
		}
		_, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)

		var found models.User
		err = db.Where("username = ?", user.Username).First(&found).Error
		assert.NoError(t, err)
		assert.NotEmpty(t, found.ID)
	})
}

func TestLoginUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	testUser := models.User{
		Username: "testlogin",
		Password: "password123",
		Email:    "testlogin@example.com",
	}
	_, err := repo.CreateUser(ctx, testUser)
	assert.NoError(t, err)

	t.Run("successful login", func(t *testing.T) {
		loginUser := models.User{
			Username: "testlogin",
			Password: "password123",
		}

		user, err := repo.LoginUser(ctx, loginUser)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, user.Token)
	})

	t.Run("wrong password", func(t *testing.T) {
		loginUser := models.User{
			Username: "testlogin",
			Password: "wrongpassword",
		}

		user, err := repo.LoginUser(ctx, loginUser)
		assert.Error(t, err)
		assert.Equal(t, models.ErrInvalidPassword, err)
		assert.Nil(t, user)
	})

	t.Run("user not found", func(t *testing.T) {
		loginUser := models.User{
			Username: "nonexistent",
			Password: "password123",
		}

		user, err := repo.LoginUser(ctx, loginUser)
		assert.Error(t, err)
		assert.Equal(t, models.ErrUserNotExist, err)
		assert.Nil(t, user)
	})

	t.Run("login with case sensitivity check", func(t *testing.T) {
		originalUser := models.User{
			Username: "casesensitive",
			Password: "password123",
			Email:    "case@example.com",
		}
		_, err := repo.CreateUser(ctx, originalUser)
		assert.NoError(t, err)

		loginUser := models.User{
			Username: "CASESENSITIVE",
			Password: "password123",
		}

		_, err = repo.LoginUser(ctx, loginUser)
		assert.Error(t, err)
	})

	t.Run("login with trimmed whitespace", func(t *testing.T) {
		originalUser := models.User{
			Username: "whitespace",
			Password: "password123",
			Email:    "whitespace@example.com",
		}
		_, err := repo.CreateUser(ctx, originalUser)
		assert.NoError(t, err)

		loginUser := models.User{
			Username: "  whitespace  ",
			Password: "password123",
		}

		_, err = repo.LoginUser(ctx, loginUser)
		assert.Error(t, err)
	})

	t.Run("concurrent logins", func(t *testing.T) {
		user := models.User{
			Username: "concurrent",
			Password: "password123",
			Email:    "concurrent@example.com",
		}
		_, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)

		var mu sync.Mutex
		concurrentUsers := 5
		var wg sync.WaitGroup
		wg.Add(concurrentUsers)

		tokens := make([]string, 0, concurrentUsers)
		var tokensMu sync.Mutex

		for i := 0; i < concurrentUsers; i++ {
			go func(i int) {
				defer wg.Done()

				mu.Lock()
				loginUser := models.User{
					Username: "concurrent",
					Password: "password123",
				}
				user, err := repo.LoginUser(ctx, loginUser)
				mu.Unlock()

				if assert.NoError(t, err) && assert.NotNil(t, user) {
					tokensMu.Lock()
					tokens = append(tokens, user.Token)
					tokensMu.Unlock()
					time.Sleep(time.Millisecond * 10)
				}
			}(i)
		}

		wg.Wait()

		assert.Equal(t, concurrentUsers, len(tokens), "Should have generated %d tokens", concurrentUsers)
		uniqueTokens := make(map[string]bool)
		for _, token := range tokens {
			assert.NotEmpty(t, token)
			uniqueTokens[token] = true
		}
		assert.Equal(t, concurrentUsers, len(uniqueTokens), "All tokens should be unique")
	})

	t.Run("verify token expiration", func(t *testing.T) {
		loginUser := models.User{
			Username: "testlogin",
			Password: "password123",
		}

		user1, err := repo.LoginUser(ctx, loginUser)
		assert.NoError(t, err)
		time.Sleep(time.Millisecond * 10)
		user2, err := repo.LoginUser(ctx, loginUser)
		assert.NoError(t, err)

		assert.NotEqual(t, user1.Token, user2.Token, "Tokens should be different")
	})

	t.Run("login with empty fields", func(t *testing.T) {
		testCases := []struct {
			name     string
			username string
			password string
		}{
			{"empty username", "", "password123"},
			{"empty password", "testuser", ""},
			{"both empty", "", ""},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				loginUser := models.User{
					Username: tc.username,
					Password: tc.password,
				}

				user, err := repo.LoginUser(ctx, loginUser)
				assert.Error(t, err)
				assert.Nil(t, user)
			})
		}
	})
}

func TestLoginUserExtended(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("invalid password hash in db", func(t *testing.T) {
		user := models.User{
			ID:       uuid.New().String(),
			Username: "invalidhash",
			Password: "invalid_hex_string",
			Email:    "invalid@hash.com",
		}
		result := db.Create(&user)
		assert.NoError(t, result.Error)

		loginUser := models.User{
			Username: "invalidhash",
			Password: "password123",
		}
		_, err := repo.LoginUser(ctx, loginUser)
		assert.Error(t, err)
	})

}
