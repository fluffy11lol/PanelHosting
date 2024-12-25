package modelsControlPanel

import "github.com/golang-jwt/jwt"

type TokenClaim struct {
	jwt.StandardClaims
	UserID string `json:"userid"`
}
