package repositoryMocks

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.User), args.Error(1)
}

type MockPvzRepository struct {
	mock.Mock
}

func (m *MockPvzRepository) CreatePvz(ctx context.Context, data *requests.CreatePvzRequest) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockPvzRepository) GetPvzReceptions(ctx context.Context, pvzId string) ([]models.Reception, error) {
	args := m.Called(ctx, pvzId)
	return args.Get(0).([]models.Reception), args.Error(1)
}

func (m *MockPvzRepository) GetReceptionProducts(ctx context.Context, receptionId string) ([]models.Product, error) {
	args := m.Called(ctx, receptionId)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockPvzRepository) GetPvzsFilteredByReceptionDate(ctx context.Context, from, to time.Time, limit, offset int) ([]models.Pvz, error) {
	args := m.Called(ctx, from, to, limit, offset)
	return args.Get(0).([]models.Pvz), args.Error(1)
}

func (m *MockPvzRepository) GetAllPvzs(ctx context.Context) ([]models.Pvz, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Pvz), args.Error(1)
}

type MockReceptionRepository struct {
	mock.Mock
}

func (m *MockReceptionRepository) CreateReception(ctx context.Context, data models.Reception) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockReceptionRepository) GetPvzById(ctx context.Context, pvzId string) (*models.Pvz, error) {
	args := m.Called(ctx, pvzId)
	return args.Get(0).(*models.Pvz), args.Error(1)
}

func (m *MockReceptionRepository) GetCurrentReception(ctx context.Context, pvzId string) (*models.Reception, error) {
	args := m.Called(ctx, pvzId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (m *MockReceptionRepository) AddProductToReception(ctx context.Context, product models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockReceptionRepository) GetLastProductInReception(ctx context.Context, receptionId string) (*models.Product, error) {
	args := m.Called(ctx, receptionId)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockReceptionRepository) DeleteProductById(ctx context.Context, productId string) error {
	args := m.Called(ctx, productId)
	return args.Error(0)
}

func (m *MockReceptionRepository) CloseReception(ctx context.Context, reception *models.Reception) error {
	args := m.Called(ctx, reception)
	return args.Error(0)
}
