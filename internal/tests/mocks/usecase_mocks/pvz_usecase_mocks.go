package usecase_mocks

import (
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/domain/responses"
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockPvzUsecase struct {
	mock.Mock
}

func (m *MockPvzUsecase) CreatePvz(ctx context.Context, req *requests.CreatePvzRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockPvzUsecase) CreateReception(ctx context.Context, req requests.CreateReceptionRequest) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func (m *MockPvzUsecase) AddProductToReception(ctx context.Context, req requests.AddProductRequest) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func (m *MockPvzUsecase) DeleteLastProduct(ctx context.Context, pvzID string) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

func (m *MockPvzUsecase) CloseReception(ctx context.Context, pvzID string) (interface{}, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0), args.Error(1)
}

func (m *MockPvzUsecase) GetPvzsInformation(ctx context.Context, start, end time.Time, limit, page int) (interface{}, error) {
	args := m.Called(ctx, start, end, limit, page)
	return args.Get(0), args.Error(1)
}

type AuthUsecaseMock struct {
	mock.Mock
}

func (m *AuthUsecaseMock) DummyLogin(ctx context.Context, role string) (string, error) {
	args := m.Called(ctx, role)
	return args.String(0), args.Error(1)
}

func (m *AuthUsecaseMock) Register(ctx context.Context, credentials requests.RegisterRequest) (responses.RegisterResponse, error) {
	args := m.Called(ctx, credentials)
	return args.Get(0).(responses.RegisterResponse), args.Error(1)
}

func (m *AuthUsecaseMock) Login(ctx context.Context, credential requests.LoginRequest) (string, error) {
	args := m.Called(ctx, credential)
	return args.String(0), args.Error(1)
}
