package utils

import (
	"crypto/rsa"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetPrivateKey(encoded string) (*rsa.PrivateKey, error) {
	pemBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
}

func GenerateToken(encoded string, userID string) (string, error) {
	privateKey, err := GetPrivateKey(encoded)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(1000, 0, 0)),
		Issuer:    "micro-durak",
	})
	return token.SignedString(privateKey)
}
