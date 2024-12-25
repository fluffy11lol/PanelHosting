package middleware

import (
	"net/http"
	"net/http/httptest"
	"storage-microservice/internal/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(userID string) string {
	claims := &models.TokenClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("My Key"))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func TestParseToken(t *testing.T) {
	tokenString := generateTestToken("12345")

	token, userID, err := ParseToken(tokenString)
	assert.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, "12345", userID)

	_, _, err = ParseToken("invalid_token")
	assert.Error(t, err)
}

func TestReadCookie(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookie := &http.Cookie{
		Name:  "token",
		Value: "test_cookie_value",
	}
	req.AddCookie(cookie)

	value, err := ReadCookie("token", req)
	assert.NoError(t, err)
	assert.Equal(t, "test_cookie_value", value)

	req, err = http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ReadCookie("token", req)
	assert.Error(t, err)
}

func TestAuthorizedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/protected", Authorized(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	validToken := generateTestToken("12345")
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: validToken})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: ""})

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid_token"})

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
