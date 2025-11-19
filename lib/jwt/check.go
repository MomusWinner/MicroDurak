package jwt

import (
	"encoding/base64"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(pubKeyEncoded string, tokenString string) (jwt.Claims, error) {
	pemBytes, err := base64.StdEncoding.DecodeString(pubKeyEncoded)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}
