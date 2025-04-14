package handler

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"context"
)

type AuthUsecase interface {
	DummyLogin(ctx context.Context, role string) (string, error)
	Register(ctx context.Context, credentials requests.RegisterRequest) (models.User, error)
	Login(ctx context.Context, credentials requests.LoginRequest) (string, error)
}
