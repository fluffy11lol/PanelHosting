package middleware

import (
	"billing-microservice/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthorized_Success(t *testing.T) {
	token := createTestToken("12345")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	recorder := httptest.NewRecorder()

	// Оборачиваем тестовый хендлер
	handler := Authorized(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("userID")
		assert.Equal(t, "12345", userID) // Проверяем, что userID успешно извлечен
		w.WriteHeader(http.StatusOK)
	})
	handler(recorder, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAuthorized_NoToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil) // Без cookie
	recorder := httptest.NewRecorder()

	handler := Authorized(func(w http.ResponseWriter, r *http.Request) {})
	handler(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthorized_InvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalidToken"})

	recorder := httptest.NewRecorder()

	handler := Authorized(func(w http.ResponseWriter, r *http.Request) {})
	handler(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
func createTestToken(userID string) string {
	claims := models.TokenClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("My Key"))
	return signedToken
}

func TestParseToken_ValidToken(t *testing.T) {
	token := createTestToken("12345")

	parsedToken, userID, err := ParseToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken)
	assert.Equal(t, "12345", userID)
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, _, err := ParseToken("invalidToken")
	assert.Error(t, err)
}

func TestReadCookie_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "testValue"})

	value, err := ReadCookie("token", req)
	assert.NoError(t, err)
	assert.Equal(t, "testValue", value)
}

func TestReadCookie_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := ReadCookie("token", req)
	assert.Error(t, err)
}

func TestReadCookie_EmptyName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := ReadCookie("", req)
	assert.Error(t, err)
	assert.Equal(t, "you are trying to read an empty cookie", err.Error())
}

func TestRespondWithError(t *testing.T) {
	recorder := httptest.NewRecorder()
	respondWithError(recorder, http.StatusBadRequest, "Test error")

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.JSONEq(t, `{"error":"Test error"}`, recorder.Body.String())
}

func TestRespondWithJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	payload := map[string]string{"message": "success"}
	respondWithJSON(recorder, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"message":"success"}`, recorder.Body.String())
}
