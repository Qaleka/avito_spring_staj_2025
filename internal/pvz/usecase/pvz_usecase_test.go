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

func TestPvzUsecase_CreatePvz(t *testing.T) {
	tests := []struct {
		name        string
		ctx         func() context.Context
		data        requests.CreatePvzRequest
		mockSetup   func(*repositoryMocks.MockPvzRepository)
		expectedErr error
	}{
		{
			name: "success",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "moderator")
			},
			data: requests.CreatePvzRequest{City: "Москва"},
			mockSetup: func(m *repositoryMocks.MockPvzRepository) {
				m.On("CreatePvz", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "invalid city",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "moderator")
			},
			data:        requests.CreatePvzRequest{City: "Новосибирск"},
			mockSetup:   func(_ *repositoryMocks.MockPvzRepository) {},
			expectedErr: errors.New("this city is not allowed"),
		},
		{
			name: "invalid role - wrong role",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")
			},
			data:        requests.CreatePvzRequest{City: "Москва"},
			mockSetup:   func(_ *repositoryMocks.MockPvzRepository) {},
			expectedErr: errors.New("this role is not allowed"),
		},
		{
			name: "repository error",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), middleware.ContextKeyRole, "moderator")
			},
			data: requests.CreatePvzRequest{City: "Москва"},
			mockSetup: func(m *repositoryMocks.MockPvzRepository) {
				m.On("CreatePvz", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositoryMocks.MockPvzRepository)
			uc := NewPvzUsecase(mockRepo)
			tt.mockSetup(mockRepo)
			err := uc.CreatePvz(tt.ctx(), &tt.data)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPvzUsecase_GetPvzsInformation(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		fromDate    time.Time
		toDate      time.Time
		limit       int
		page        int
		mockSetup   func(*repositoryMocks.MockPvzRepository)
		expectedRes []models.Pvz
		expectedErr error
	}{
		{
			name:     "success",
			fromDate: now.Add(-24 * time.Hour),
			toDate:   now,
			limit:    10,
			page:     1,
			mockSetup: func(m *repositoryMocks.MockPvzRepository) {
				pvz := models.Pvz{Id: "pvz123", City: "Москва"}
				reception := models.Reception{Id: "reception123", PvzId: "pvz123"}
				product := models.Product{Id: "product123", ReceptionId: "reception123"}

				m.On("GetPvzsFilteredByReceptionDate", mock.Anything, now.Add(-24*time.Hour), now, 10, 0).
					Return([]models.Pvz{pvz}, nil)
				m.On("GetPvzReceptions", mock.Anything, "pvz123").
					Return([]models.Reception{reception}, nil)
				m.On("GetReceptionProducts", mock.Anything, "reception123").
					Return([]models.Product{product}, nil)
			},
			expectedRes: []models.Pvz{
				{
					Id:   "pvz123",
					City: "Москва",
					Receptions: []models.Reception{
						{
							Id:       "reception123",
							PvzId:    "pvz123",
							Products: []models.Product{{Id: "product123", ReceptionId: "reception123"}},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:     "no pvzs found",
			fromDate: now.Add(-24 * time.Hour),
			toDate:   now,
			limit:    10,
			page:     1,
			mockSetup: func(m *repositoryMocks.MockPvzRepository) {
				m.On("GetPvzsFilteredByReceptionDate", mock.Anything, now.Add(-24*time.Hour), now, 10, 0).
					Return([]models.Pvz{}, nil)
			},
			expectedRes: []models.Pvz{},
			expectedErr: nil,
		},
		{
			name:     "get pvzs error",
			fromDate: now.Add(-24 * time.Hour),
			toDate:   now,
			limit:    10,
			page:     1,
			mockSetup: func(m *repositoryMocks.MockPvzRepository) {
				m.On("GetPvzsFilteredByReceptionDate", mock.Anything, now.Add(-24*time.Hour), now, 10, 0).
					Return([]models.Pvz{}, errors.New("db error"))
			},
			expectedRes: nil,
			expectedErr: errors.New("db error"),
		},
		{
			name:     "get receptions error",
			fromDate: now.Add(-24 * time.Hour),
			toDate:   now,
			limit:    10,
			page:     1,
			mockSetup: func(m *repositoryMocks.MockPvzRepository) {
				pvz := models.Pvz{Id: "pvz123", City: "Москва"}
				m.On("GetPvzsFilteredByReceptionDate", mock.Anything, now.Add(-24*time.Hour), now, 10, 0).
					Return([]models.Pvz{pvz}, nil)
				m.On("GetPvzReceptions", mock.Anything, "pvz123").
					Return([]models.Reception{}, errors.New("receptions error"))
			},
			expectedRes: nil,
			expectedErr: errors.New("receptions error"),
		},
		{
			name:     "get products error",
			fromDate: now.Add(-24 * time.Hour),
			toDate:   now,
			limit:    10,
			page:     1,
			mockSetup: func(m *repositoryMocks.MockPvzRepository) {
				pvz := models.Pvz{Id: "pvz123", City: "Москва"}
				reception := models.Reception{Id: "reception123", PvzId: "pvz123"}
				m.On("GetPvzsFilteredByReceptionDate", mock.Anything, now.Add(-24*time.Hour), now, 10, 0).
					Return([]models.Pvz{pvz}, nil)
				m.On("GetPvzReceptions", mock.Anything, "pvz123").
					Return([]models.Reception{reception}, nil)
				m.On("GetReceptionProducts", mock.Anything, "reception123").
					Return([]models.Product{}, errors.New("products error"))
			},
			expectedRes: nil,
			expectedErr: errors.New("products error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repositoryMocks.MockPvzRepository)
			uc := NewPvzUsecase(mockRepo)
			tt.mockSetup(mockRepo)

			res, err := uc.GetPvzsInformation(context.Background(), tt.fromDate, tt.toDate, tt.limit, tt.page)
			assert.Equal(t, tt.expectedRes, res)
			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
