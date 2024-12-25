package models

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenClaimMarshal(t *testing.T) {
	claim := &TokenClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 1609459200, // 2021-01-01 00:00:00 UTC
		},
		UserID: "user123",
	}

	jsonData, err := json.Marshal(claim)
	assert.NoError(t, err)
	assert.NotNil(t, jsonData)

}

func TestTokenClaimUnmarshal(t *testing.T) {
	jsonData := `{"userid":"user123","expires_at":1609459200}`

	var claim TokenClaim
	err := json.Unmarshal([]byte(jsonData), &claim)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claim.UserID)
	assert.NotEqual(t, int64(1609459200), claim.ExpiresAt)
}

func TestTokenClaimEmpty(t *testing.T) {
	var claim TokenClaim

	_, err := json.Marshal(claim)

	err = json.Unmarshal([]byte(`{}`), &claim)
	assert.NoError(t, err)

}
