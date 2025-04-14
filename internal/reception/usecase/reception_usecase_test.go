package usecase

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/service/middleware"
	"avito_spring_staj_2025/internal/tests/mocks/repository_mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestPvzUsecase_CreateReception(t *testing.T) {
	tests := []struct {
		name        string
		ctx         func() context.Context
		data        requests.CreateReceptionRequest
		mockSetup   func(repository *repositoryMocks.MockReceptionRepository)
		expectedRes models.Reception
		expectedErr error
	}{
		{
			name: "success",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			data: requests.CreateReceptionRequest{PvzId: "pvz123"},
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(nil, errors.New("no active reception"))
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, nil)
				m.On("CreateReception", mock.Anything, mock.MatchedBy(func(r models.Reception) bool {
					return r.PvzId == "pvz123" && r.Status == "in_progress"
				})).Return(nil)
			},
			expectedRes: models.Reception{
				PvzId:  "pvz123",
				Status: "in_progress",
			},
			expectedErr: nil,
		},
		{
			name: "invalid role",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "wrong_role")
			},
			data:        requests.CreateReceptionRequest{PvzId: "pvz123"},
			mockSetup:   func(_ *repositoryMocks.MockReceptionRepository) {},
			expectedErr: errors.New("this role is not allowed"),
		},
		{
			name: "active reception exists",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			data: requests.CreateReceptionRequest{PvzId: "pvz123"},
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(&models.Reception{}, nil)
			},
			expectedErr: errors.New("active reception already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositoryMocks.MockReceptionRepository)
			uc := NewReceptionUsecase(mockRepo)
			tt.mockSetup(mockRepo)

			res, err := uc.CreateReception(tt.ctx(), tt.data)

			assert.Equal(t, tt.expectedRes.PvzId, res.PvzId)
			assert.Equal(t, tt.expectedRes.Status, res.Status)

			if tt.expectedErr == nil {
				assert.NotEmpty(t, res.Id)
				assert.False(t, res.DateTime.IsZero())
			}

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPvzUsecase_AddProductToReception(t *testing.T) {
	tests := []struct {
		name        string
		ctx         func() context.Context
		data        requests.AddProductRequest
		mockSetup   func(*repositoryMocks.MockReceptionRepository)
		expectedRes models.Product
		expectedErr error
	}{
		{
			name: "success",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			data: requests.AddProductRequest{
				PvzId: "pvz123",
				Type:  models.CLOTHES_TYPE,
			},
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, nil)
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(&models.Reception{Id: "reception123"}, nil)
				m.On("AddProductToReception", mock.Anything, mock.Anything).
					Return(nil)
			},
			expectedRes: models.Product{
				ReceptionId: "reception123",
				Type:        models.CLOTHES_TYPE,
			},
			expectedErr: nil,
		},
		{
			name: "invalid product type",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			data: requests.AddProductRequest{
				PvzId: "pvz123",
				Type:  "invalid_type",
			},
			mockSetup:   func(_ *repositoryMocks.MockReceptionRepository) {},
			expectedErr: errors.New("this type is not allowed"),
		},
		{
			name: "invalid role",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "wrong_role")
			},
			data: requests.AddProductRequest{
				PvzId: "pvz123",
				Type:  models.CLOTHES_TYPE,
			},
			mockSetup:   func(_ *repositoryMocks.MockReceptionRepository) {},
			expectedErr: errors.New("this role is not allowed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositoryMocks.MockReceptionRepository)
			uc := NewReceptionUsecase(mockRepo)
			tt.mockSetup(mockRepo)

			res, err := uc.AddProductToReception(tt.ctx(), tt.data)
			if tt.expectedErr == nil {
				assert.NotEmpty(t, res.Id)
				assert.False(t, res.DateTime.IsZero())
				res.Id = ""
				res.DateTime = time.Time{}
			}

			assert.Equal(t, tt.expectedRes, res)
			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPvzUsecase_DeleteLastProduct(t *testing.T) {
	tests := []struct {
		name        string
		ctx         func() context.Context
		pvzId       string
		mockSetup   func(*repositoryMocks.MockReceptionRepository)
		expectedErr error
	}{
		{
			name: "success",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			pvzId: "pvz123",
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, nil)
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(&models.Reception{Id: "reception123"}, nil)
				m.On("GetLastProductInReception", mock.Anything, "reception123").
					Return(&models.Product{Id: "product123"}, nil)
				m.On("DeleteProductById", mock.Anything, "product123").
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "invalid role",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "wrong_role")
			},
			pvzId:       "pvz123",
			mockSetup:   func(_ *repositoryMocks.MockReceptionRepository) {},
			expectedErr: errors.New("this role is not allowed"),
		},
		{
			name: "pvz not found",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			pvzId: "pvz123",
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, errors.New("pvz not found"))
			},
			expectedErr: errors.New("pvz not found"),
		},
		{
			name: "no active reception",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			pvzId: "pvz123",
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, nil)
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(nil, errors.New("no active reception"))
			},
			expectedErr: errors.New("no active reception"),
		},
		{
			name: "no products in reception",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			pvzId: "pvz123",
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, nil)
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(&models.Reception{Id: "reception123"}, nil)
				m.On("GetLastProductInReception", mock.Anything, "reception123").
					Return(&models.Product{}, errors.New("no products"))
			},
			expectedErr: errors.New("no products"),
		},
		{
			name: "delete product error",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			pvzId: "pvz123",
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, nil)
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(&models.Reception{Id: "reception123"}, nil)
				m.On("GetLastProductInReception", mock.Anything, "reception123").
					Return(&models.Product{Id: "product123"}, nil)
				m.On("DeleteProductById", mock.Anything, "product123").
					Return(errors.New("delete error"))
			},
			expectedErr: errors.New("delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositoryMocks.MockReceptionRepository)
			uc := NewReceptionUsecase(mockRepo)
			tt.mockSetup(mockRepo)

			err := uc.DeleteLastProduct(tt.ctx(), tt.pvzId)
			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPvzUsecase_CloseReception(t *testing.T) {
	tests := []struct {
		name        string
		ctx         func() context.Context
		pvzId       string
		mockSetup   func(*repositoryMocks.MockReceptionRepository)
		expectedRes models.Reception
		expectedErr error
	}{
		{
			name: "success",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			pvzId: "pvz123",
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				reception := &models.Reception{
					Id:     "reception123",
					PvzId:  "pvz123",
					Status: "closed",
				}
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, nil)
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(reception, nil)
				m.On("CloseReception", mock.Anything, reception).
					Return(nil)
			},
			expectedRes: models.Reception{
				Id:     "reception123",
				PvzId:  "pvz123",
				Status: "closed",
			},
			expectedErr: nil,
		},
		{
			name: "invalid role",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "wrong_role")
			},
			pvzId:       "pvz123",
			mockSetup:   func(_ *repositoryMocks.MockReceptionRepository) {},
			expectedRes: models.Reception{},
			expectedErr: errors.New("this role is not allowed"),
		},
		{
			name: "pvz not found",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			pvzId: "pvz123",
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, errors.New("pvz not found"))
			},
			expectedRes: models.Reception{},
			expectedErr: errors.New("pvz not found"),
		},
		{
			name: "no active reception",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			pvzId: "pvz123",
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, nil)
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(nil, errors.New("no active reception"))
			},
			expectedRes: models.Reception{},
			expectedErr: errors.New("no active reception"),
		},
		{
			name: "close reception error",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			pvzId: "pvz123",
			mockSetup: func(m *repositoryMocks.MockReceptionRepository) {
				reception := &models.Reception{
					Id:     "reception123",
					PvzId:  "pvz123",
					Status: "in_progress",
				}
				m.On("GetPvzById", mock.Anything, "pvz123").
					Return(&models.Pvz{}, nil)
				m.On("GetCurrentReception", mock.Anything, "pvz123").
					Return(reception, nil)
				m.On("CloseReception", mock.Anything, mock.Anything).
					Return(errors.New("close error"))
			},
			expectedRes: models.Reception{},
			expectedErr: errors.New("close error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositoryMocks.MockReceptionRepository)
			uc := NewReceptionUsecase(mockRepo)
			tt.mockSetup(mockRepo)

			res, err := uc.CloseReception(tt.ctx(), tt.pvzId)
			assert.Equal(t, tt.expectedRes, res)
			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
