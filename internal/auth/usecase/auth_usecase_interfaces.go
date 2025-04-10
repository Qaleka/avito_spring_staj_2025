package usecase

import (
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/domain/responses"
	"context"
)

type AuthUsecase interface {
	DummyLogin(ctx context.Context, role string) (string, error)
	Register(ctx context.Context, credentials requests.RegisterRequest) (responses.RegisterResponse, error)
	Login(ctx context.Context, credentials requests.LoginRequest) (string, error)
}
