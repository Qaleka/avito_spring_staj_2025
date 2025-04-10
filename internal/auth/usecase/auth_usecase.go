package usecase

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/domain/responses"
	"avito_spring_staj_2025/internal/auth/repository"
	"avito_spring_staj_2025/internal/service/jwt"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type authUsecase struct {
	authRepository repository.AuthRepository
	jwtService     jwt.JwtTokenService
}

func NewAuthUsecase(authRepository repository.AuthRepository, jwtService jwt.JwtTokenService) AuthUsecase {
	return &authUsecase{
		authRepository: authRepository,
		jwtService:     jwtService,
	}
}

func (au *authUsecase) DummyLogin(ctx context.Context, role string) (string, error) {
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

func (au *authUsecase) Register(ctx context.Context, credentials requests.RegisterRequest) (responses.RegisterResponse, error) {
	hashedPassword, err := middleware.HashPassword(credentials.Password)
	if err != nil {
		return responses.RegisterResponse{}, err
	}

	if credentials.Role != "employee" && credentials.Role != "moderator" {
		return responses.RegisterResponse{}, errors.New("invalid role")
	}

	_, err = au.authRepository.GetUserByEmail(ctx, credentials.Email)
	if err == nil {
		return responses.RegisterResponse{}, errors.New("user with this email already exists")
	}

	user := models.User{Id: uuid.New().String(), Email: credentials.Email, Password: hashedPassword, Role: credentials.Role}
	err = au.authRepository.CreateUser(ctx, &user)
	if err != nil {
		return responses.RegisterResponse{}, err
	}
	return responses.RegisterResponse{Id: user.Id, Email: user.Email, Role: user.Role}, nil
}

func (au *authUsecase) Login(ctx context.Context, credentials requests.LoginRequest) (string, error) {
	user, err := au.authRepository.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return "", err
	}
	if ok := middleware.CheckPassword(user.Password, credentials.Password); !ok {
		return "", errors.New("invalid password")
	}
	tokenExpTime := time.Now().Add(24 * time.Hour).Unix()
	jwtToken, err := au.jwtService.Create(user.Role, tokenExpTime)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
