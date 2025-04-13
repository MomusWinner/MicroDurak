package jwt

import (
	"crypto/rsa"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
)

func GetPublicKey(encoded string) (*rsa.PublicKey, error) {
	pemBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(pemBytes)
}
