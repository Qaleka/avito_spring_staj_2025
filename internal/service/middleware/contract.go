package middleware

import (
	jwt_package "avito_spring_staj_2025/internal/service/jwt"
	"github.com/golang-jwt/jwt/v4"
)

type JwtTokenService interface {
	Create(role string, tokenExpTime int64) (string, error)
	Validate(tokenString string) (*jwt_package.JwtCsrfClaims, error)
	ParseSecretGetter(token *jwt.Token) (interface{}, error)
}
