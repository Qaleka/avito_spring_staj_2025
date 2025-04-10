package repository

import (
	"avito_spring_staj_2025/domain/models"
	"context"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
