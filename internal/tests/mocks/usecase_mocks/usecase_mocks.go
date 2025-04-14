package usecaseMocks

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
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

func (m *MockPvzUsecase) GetPvzsInformation(ctx context.Context, fromDate, endDate time.Time, limit, page int) ([]models.Pvz, error) {
	args := m.Called(ctx, fromDate, endDate, limit, page)
	return args.Get(0).([]models.Pvz), args.Error(1)
}

func (m *MockPvzUsecase) GetAllPvzs(ctx context.Context) ([]models.Pvz, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Pvz), args.Error(1)
}

type AuthUsecaseMock struct {
	mock.Mock
}

func (m *AuthUsecaseMock) DummyLogin(ctx context.Context, role string) (string, error) {
	args := m.Called(ctx, role)
	return args.String(0), args.Error(1)
}

func (m *AuthUsecaseMock) Register(ctx context.Context, credentials requests.RegisterRequest) (models.User, error) {
	args := m.Called(ctx, credentials)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *AuthUsecaseMock) Login(ctx context.Context, credential requests.LoginRequest) (string, error) {
	args := m.Called(ctx, credential)
	return args.String(0), args.Error(1)
}

type ReceptionUsecaseMock struct {
	mock.Mock
}

func (m *ReceptionUsecaseMock) CreateReception(ctx context.Context, req requests.CreateReceptionRequest) (models.Reception, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(models.Reception), args.Error(1)
}

func (m *ReceptionUsecaseMock) AddProductToReception(ctx context.Context, req requests.AddProductRequest) (models.Product, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *ReceptionUsecaseMock) DeleteLastProduct(ctx context.Context, pvzID string) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

func (m *ReceptionUsecaseMock) CloseReception(ctx context.Context, pvzID string) (models.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}
