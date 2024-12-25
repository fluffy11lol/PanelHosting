package repository

import (
	"context"
	"encoding/hex"
	"log"
	"time"

	"authentication-microservice/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (s *UserRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return nil, models.ErrEmptyField
	}
	var existingUser models.User
	if err := s.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		log.Println(existingUser.ID, existingUser.Username, existingUser.Password)
		return nil, models.ErrUserExist
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	s.db.Create(&models.User{ID: uuid.New().String(), Password: hex.EncodeToString(hash), Username: user.Username, Email: user.Email})
	return nil, nil
}

func (s *UserRepository) LoginUser(ctx context.Context, user models.User) (*models.User, error) {
	var existingUser models.User
	err := s.db.Where("username = ?", user.Username).First(&existingUser).Error
	if err != nil {
		return nil, models.ErrUserNotExist
	}

	hashed, err := hex.DecodeString(existingUser.Password)
	if err != nil {
		log.Println(err)
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte(user.Password))
	if err != nil {
		return nil, models.ErrInvalidPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			Id:        uuid.New().String(),
		},
		UserID: existingUser.ID,
	})

	TokenStr, err := token.SignedString([]byte("My Key"))
	if err != nil {
		return nil, err
	}
	return &models.User{Token: TokenStr}, nil
}
