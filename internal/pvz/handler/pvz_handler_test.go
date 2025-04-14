package handler

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/service/logger"
	usecaseMocks "avito_spring_staj_2025/internal/tests/mocks/usecase_mocks"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPvzHandler_CreatePvz(t *testing.T) {
	logger.AccessLogger = zap.NewNop()
	type mockBehavior func(usecase *usecaseMocks.MockPvzUsecase, req requests.CreatePvzRequest)

	testTime := time.Date(2025, 4, 11, 23, 58, 7, 0, time.Local)

	tests := []struct {
		name           string
		inputBody      string
		inputRequest   requests.CreatePvzRequest
		mockBehavior   mockBehavior
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "success",
			inputBody: fmt.Sprintf(`{"id":"123","registrationDate":"%s","city":"Moscow"}`, testTime.Format(time.RFC3339)),
			inputRequest: requests.CreatePvzRequest{
				Id:               "123",
				RegistrationDate: testTime,
				City:             "Moscow",
			},
			mockBehavior: func(usecase *usecaseMocks.MockPvzUsecase, req requests.CreatePvzRequest) {
				usecase.On("CreatePvz", mock.Anything, &req).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   fmt.Sprintf(`{"id":"123","registrationDate":"%s","city":"Moscow"}`, testTime.Format(time.RFC3339)),
		},
		{
			name:           "bad json",
			inputBody:      `{"id":123}`,
			mockBehavior:   func(_ *usecaseMocks.MockPvzUsecase, _ requests.CreatePvzRequest) {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"json: cannot unmarshal number into Go struct field CreatePvzRequest.id of type string"}`,
		},
		{
			name:      "usecase error",
			inputBody: fmt.Sprintf(`{"id":"err","registrationDate":"%s","city":"Moscow"}`, testTime.Format(time.RFC3339)),
			inputRequest: requests.CreatePvzRequest{
				Id:               "err",
				RegistrationDate: testTime,
				City:             "Moscow",
			},
			mockBehavior: func(usecase *usecaseMocks.MockPvzUsecase, req requests.CreatePvzRequest) {
				usecase.On("CreatePvz", mock.Anything, &req).Return(errors.New("create error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"create error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(usecaseMocks.MockPvzUsecase)
			handler := NewPvzHandler(mockUsecase)

			tt.mockBehavior(mockUsecase, tt.inputRequest)

			req := httptest.NewRequest(http.MethodPost, "/api/pvz", strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()

			handler.CreatePvz(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				if err != nil {
					return
				}
			}()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}

func TestPvzHandler_GetPvzsInformation(t *testing.T) {
	logger.AccessLogger = zap.NewNop()
	type mockBehavior func(usecase *usecaseMocks.MockPvzUsecase, ctx context.Context, startDate, endDate time.Time, limit, page int)

	testTime := time.Date(2025, 4, 11, 23, 59, 0, 0, time.Local)
	startDate := time.Date(2025, 4, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 4, 12, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		query          string
		mockBehavior   mockBehavior
		expectedStatus int
		expectedBody   string
	}{
		{
			name:  "success",
			query: "?startDate=2025-04-10T00:00:00Z&endDate=2025-04-12T00:00:00Z&page=0&limit=0",
			mockBehavior: func(usecase *usecaseMocks.MockPvzUsecase, ctx context.Context, startDate, endDate time.Time, limit, page int) {
				usecase.On("GetPvzsInformation", ctx, startDate, endDate, limit, page).
					Return([]models.Pvz{
						{
							Id:               "123",
							RegistrationDate: testTime,
							City:             "Moscow",
							Receptions:       []models.Reception{},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf(`[{"pvz":{"Id":"123","RegistrationDate":"%s","City":"Moscow"},"receptions":[]}]`, testTime.Format(time.RFC3339)),
		},
		{
			name:  "usecase error",
			query: "?startDate=2025-04-10T00:00:00Z&endDate=2025-04-12T00:00:00Z&page=1&limit=10",
			mockBehavior: func(usecase *usecaseMocks.MockPvzUsecase, ctx context.Context, startDate, endDate time.Time, limit, page int) {
				usecase.On("GetPvzsInformation", ctx, startDate, endDate, limit, page).
					Return([]models.Pvz{}, errors.New("get error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"get error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(usecaseMocks.MockPvzUsecase)
			handler := NewPvzHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodGet, "/api/pvz"+tt.query, nil)
			w := httptest.NewRecorder()

			tt.mockBehavior(mockUsecase, req.Context(), startDate, endDate, 10, 1)

			handler.GetPvzsInformation(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				if err != nil {
					return
				}
			}()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, string(body))
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}
