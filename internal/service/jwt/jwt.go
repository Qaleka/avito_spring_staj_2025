package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtToken struct {
	Secret []byte
}

func NewJwtToken(secret string) (JwtToken, error) {
	if secret == "" {
		return JwtToken{}, errors.New("secret key can't be empty")
	}
	return JwtToken{
		Secret: []byte(secret),
	}, nil
}

type JwtCsrfClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

func (tk JwtToken) Create(role string, tokenExpTime int64) (string, error) {
	if role == "" {
		return "", errors.New("role is empty")
	}
	data := JwtCsrfClaims{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpTime,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(tk.Secret)
}

func (tk JwtToken) Validate(tokenString string) (*JwtCsrfClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCsrfClaims{}, tk.ParseSecretGetter)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JwtCsrfClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

func (tk JwtToken) ParseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return tk.Secret, nil
}
