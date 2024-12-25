package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"project-microservice/internal/models"

	"github.com/golang-jwt/jwt"
	"net/http"
	"net/url"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Authorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, err := ReadCookie("token", r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if auth == "" {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		token, id, err := ParseToken(auth)
		if err != nil || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		r.Header.Set("userID", fmt.Sprintf("%v", id))
		ctx := context.WithValue(r.Context(), UserIDKey, id)
		next(w, r.WithContext(ctx))
	}
}

func ParseToken(accessToken string) (*jwt.Token, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to encode JSON response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
