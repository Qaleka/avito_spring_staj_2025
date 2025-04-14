package usecase

import (
	"avito_spring_staj_2025/domain/models"
	jwt_package "avito_spring_staj_2025/internal/service/jwt"
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type JwtTokenService interface {
	Create(role string, tokenExpTime int64) (string, error)
	Validate(tokenString string) (*jwt_package.JwtCsrfClaims, error)
	ParseSecretGetter(token *jwt.Token) (interface{}, error)
}
