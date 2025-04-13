package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(publicKey string, tokenString string) (jwt.Claims, error) {
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
