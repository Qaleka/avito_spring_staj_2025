package usecase

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/service/auth"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type AuthUsecase struct {
	authRepository AuthRepository
	jwtService     JwtTokenService
}

func NewAuthUsecase(authRepository AuthRepository, jwtService JwtTokenService) AuthUsecase {
	return AuthUsecase{
		authRepository: authRepository,
		jwtService:     jwtService,
	}
}

func (au AuthUsecase) DummyLogin(_ context.Context, role string) (string, error) {
	if role != "employee" && role != "moderator" {
		return "", errors.New("invalid role")
	}
	tokenExpTime := time.Now().Add(24 * time.Hour).Unix()
	jwtToken, err := au.jwtService.Create(role, tokenExpTime)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func (au AuthUsecase) Register(ctx context.Context, credentials requests.RegisterRequest) (models.User, error) {
	hashedPassword, err := auth.HashPassword(credentials.Password)
	if err != nil {
		return models.User{}, err
	}

	if credentials.Role != "employee" && credentials.Role != "moderator" {
		return models.User{}, errors.New("invalid role")
	}

	_, err = au.authRepository.GetUserByEmail(ctx, credentials.Email)
	if err == nil {
		return models.User{}, errors.New("user with this email already exists")
	}

	user := models.User{Id: uuid.New().String(), Email: credentials.Email, Password: hashedPassword, Role: credentials.Role}
	err = au.authRepository.CreateUser(ctx, &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (au AuthUsecase) Login(ctx context.Context, credentials requests.LoginRequest) (string, error) {
	user, err := au.authRepository.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return "", err
	}
	if ok := auth.CheckPassword(user.Password, credentials.Password); !ok {
		return "", errors.New("invalid password")
	}
	tokenExpTime := time.Now().Add(24 * time.Hour).Unix()
	jwtToken, err := au.jwtService.Create(user.Role, tokenExpTime)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
