package unit

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/domain/responses"
	"avito_spring_staj_2025/internal/auth/usecase"
	"avito_spring_staj_2025/internal/service/auth"
	"avito_spring_staj_2025/internal/tests/mocks/jwt_mocks"
	"avito_spring_staj_2025/internal/tests/mocks/repository_mocks"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
	"time"
)

func TestAuthUsecase_DummyLogin(t *testing.T) {
	expTime := time.Now().Add(24 * time.Hour).Unix()

	tests := []struct {
		name          string
		role          string
		mockSetup     func(*repository_mocks.MockAuthRepository, *jwt_mocks.MockJwtService)
		expectedToken string
		expectedErr   string
	}{
		{
			name: "success employee",
			role: "employee",
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
				mj.On("Create", "employee", expTime).Return("employee_token", nil)
			},
			expectedToken: "employee_token",
			expectedErr:   "",
		},
		{
			name: "invalid role",
			role: "admin",
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
			},
			expectedToken: "",
			expectedErr:   "invalid role",
		},
		{
			name: "jwt error",
			role: "employee",
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
				mj.On("Create", "employee", expTime).Return("", errors.New("jwt error"))
			},
			expectedToken: "",
			expectedErr:   "jwt error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository_mocks.MockAuthRepository)
			mockJwt := new(jwt_mocks.MockJwtService)
			uc := usecase.NewAuthUsecase(mockRepo, mockJwt)

			tt.mockSetup(mockRepo, mockJwt)

			token, err := uc.DummyLogin(context.Background(), tt.role)

			assert.Equal(t, tt.expectedToken, token)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockJwt.AssertExpectations(t)
		})
	}
}

func TestAuthUsecase_Register(t *testing.T) {
	testEmail := "test@example.com"
	testPassword := "password123"
	testUserID := uuid.NewString()

	tests := []struct {
		name        string
		credentials requests.RegisterRequest
		mockSetup   func(*repository_mocks.MockAuthRepository, *jwt_mocks.MockJwtService)
		expectedRes responses.RegisterResponse
		expectedErr string
	}{
		{
			name: "success registration",
			credentials: requests.RegisterRequest{
				Email:    testEmail,
				Password: testPassword,
				Role:     "employee",
			},
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
				mr.On("GetUserByEmail", mock.Anything, testEmail).
					Return(&models.User{}, errors.New("not found"))
				mr.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
					return u.Email == testEmail && u.Role == "employee"
				})).Run(func(args mock.Arguments) {
					u := args.Get(1).(*models.User)
					u.Id = testUserID
				}).
					Return(nil)
			},
			expectedRes: responses.RegisterResponse{
				Id:    testUserID,
				Email: testEmail,
				Role:  "employee",
			},
			expectedErr: "",
		},
		{
			name: "invalid role",
			credentials: requests.RegisterRequest{
				Email:    testEmail,
				Password: testPassword,
				Role:     "invalid",
			},
			mockSetup:   func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {},
			expectedRes: responses.RegisterResponse{},
			expectedErr: "invalid role",
		},
		{
			name: "user already exists",
			credentials: requests.RegisterRequest{
				Email:    testEmail,
				Password: testPassword,
				Role:     "employee",
			},
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
				mr.On("GetUserByEmail", mock.Anything, testEmail).
					Return(&models.User{Email: testEmail}, nil)
			},
			expectedRes: responses.RegisterResponse{},
			expectedErr: "user with this email already exists",
		},
		{
			name: "password hash error",
			credentials: requests.RegisterRequest{
				Email:    testEmail,
				Password: strings.Repeat("a", 73),
				Role:     "employee",
			},
			mockSetup:   func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {},
			expectedRes: responses.RegisterResponse{},
			expectedErr: "bcrypt: password length exceeds 72 bytes",
		},
		{
			name: "create user error",
			credentials: requests.RegisterRequest{
				Email:    testEmail,
				Password: testPassword,
				Role:     "employee",
			},
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
				mr.On("GetUserByEmail", mock.Anything, testEmail).
					Return(&models.User{}, errors.New("not found"))
				mr.On("CreateUser", mock.Anything, mock.Anything).
					Return(errors.New("database error"))
			},
			expectedRes: responses.RegisterResponse{},
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository_mocks.MockAuthRepository)
			mockJwt := new(jwt_mocks.MockJwtService)
			uc := usecase.NewAuthUsecase(mockRepo, mockJwt)

			tt.mockSetup(mockRepo, mockJwt)

			res, err := uc.Register(context.Background(), tt.credentials)

			assert.Equal(t, tt.expectedRes, res)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockJwt.AssertExpectations(t)
		})
	}
}

func TestAuthUsecase_Login(t *testing.T) {
	expTime := time.Now().Add(24 * time.Hour).Unix()
	testEmail := "test@example.com"
	testPassword := "password123"
	hashedPassword, _ := auth.HashPassword(testPassword)
	testUser := &models.User{
		Email:    testEmail,
		Password: hashedPassword,
		Role:     "employee",
	}

	tests := []struct {
		name          string
		credentials   requests.LoginRequest
		mockSetup     func(*repository_mocks.MockAuthRepository, *jwt_mocks.MockJwtService)
		expectedToken string
		expectedErr   string
	}{
		{
			name: "success login",
			credentials: requests.LoginRequest{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
				mr.On("GetUserByEmail", mock.Anything, testEmail).
					Return(testUser, nil)
				mj.On("Create", "employee", expTime).
					Return("valid_token", nil)
			},
			expectedToken: "valid_token",
			expectedErr:   "",
		},
		{
			name: "user not found",
			credentials: requests.LoginRequest{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
				mr.On("GetUserByEmail", mock.Anything, testEmail).
					Return(&models.User{}, errors.New("not found"))
			},
			expectedToken: "",
			expectedErr:   "not found",
		},
		{
			name: "invalid password",
			credentials: requests.LoginRequest{
				Email:    testEmail,
				Password: "wrong_password",
			},
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
				mr.On("GetUserByEmail", mock.Anything, testEmail).
					Return(testUser, nil)
			},
			expectedToken: "",
			expectedErr:   "invalid password",
		},
		{
			name: "jwt error",
			credentials: requests.LoginRequest{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func(mr *repository_mocks.MockAuthRepository, mj *jwt_mocks.MockJwtService) {
				mr.On("GetUserByEmail", mock.Anything, testEmail).
					Return(testUser, nil)
				mj.On("Create", "employee", expTime).
					Return("", errors.New("jwt error"))
			},
			expectedToken: "",
			expectedErr:   "jwt error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository_mocks.MockAuthRepository)
			mockJwt := new(jwt_mocks.MockJwtService)
			uc := usecase.NewAuthUsecase(mockRepo, mockJwt)

			tt.mockSetup(mockRepo, mockJwt)

			token, err := uc.Login(context.Background(), tt.credentials)

			assert.Equal(t, tt.expectedToken, token)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockJwt.AssertExpectations(t)
		})
	}
}
