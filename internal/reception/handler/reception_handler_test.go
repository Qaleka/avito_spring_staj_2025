package handler

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/service/logger"
	usecaseMocks "avito_spring_staj_2025/internal/tests/mocks/usecase_mocks"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
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

func TestPvzHandler_CreateReception(t *testing.T) {
	logger.AccessLogger = zap.NewNop()
	type mockBehavior func(_ *usecaseMocks.ReceptionUsecaseMock, req requests.CreateReceptionRequest)

	testTime := time.Date(2025, 4, 11, 23, 59, 0, 0, time.Local)

	tests := []struct {
		name           string
		inputBody      string
		inputRequest   requests.CreateReceptionRequest
		mockBehavior   mockBehavior
		expectedStatus int
		expectedBody   string
	}{
		{
			name:         "success",
			inputBody:    `{"pvzId":"123"}`,
			inputRequest: requests.CreateReceptionRequest{PvzId: "123"},
			mockBehavior: func(usecase *usecaseMocks.ReceptionUsecaseMock, req requests.CreateReceptionRequest) {
				usecase.On("CreateReception", mock.Anything, req).Return(models.Reception{
					Id:       "rec1",
					DateTime: testTime,
					PvzId:    "123",
					Status:   "in_progress",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   fmt.Sprintf(`{"id":"rec1","dateTime":"%s","pvzId":"123","status":"in_progress"}`, testTime.Format(time.RFC3339)),
		},
		{
			name:           "bad json",
			inputBody:      `{"pvzId":123}`,
			mockBehavior:   func(_ *usecaseMocks.ReceptionUsecaseMock, req requests.CreateReceptionRequest) {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"json: cannot unmarshal number into Go struct field CreateReceptionRequest.pvzId of type string"}`,
		},
		{
			name:         "usecase error",
			inputBody:    `{"pvzId":"fail"}`,
			inputRequest: requests.CreateReceptionRequest{PvzId: "fail"},
			mockBehavior: func(usecase *usecaseMocks.ReceptionUsecaseMock, req requests.CreateReceptionRequest) {
				usecase.On("CreateReception", mock.Anything, req).Return(models.Reception{}, errors.New("reception creation failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"reception creation failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(usecaseMocks.ReceptionUsecaseMock)
			handler := NewReceptionHandler(mockUsecase)

			tt.mockBehavior(mockUsecase, tt.inputRequest)

			req := httptest.NewRequest(http.MethodPost, "/api/reception", strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()

			handler.CreateReception(w, req)

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

func TestPvzHandler_AddProductToReception(t *testing.T) {
	logger.AccessLogger = zap.NewNop()
	type mockBehavior func(usecase *usecaseMocks.ReceptionUsecaseMock, req requests.AddProductRequest)

	testTime := time.Date(2025, 4, 11, 23, 59, 30, 0, time.Local)

	tests := []struct {
		name           string
		inputBody      string
		inputRequest   requests.AddProductRequest
		mockBehavior   mockBehavior
		expectedStatus int
		expectedBody   string
	}{
		{
			name:         "success",
			inputBody:    `{"type":"book","pvzId":"123"}`,
			inputRequest: requests.AddProductRequest{Type: "book", PvzId: "123"},
			mockBehavior: func(usecase *usecaseMocks.ReceptionUsecaseMock, req requests.AddProductRequest) {
				usecase.On("AddProductToReception", mock.Anything, req).Return(models.Product{
					Id:          "prod1",
					DateTime:    testTime,
					Type:        "обувь",
					ReceptionId: "rec1",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   fmt.Sprintf(`{"id":"prod1","dateTime":"%s","type":"обувь","receptionId":"rec1"}`, testTime.Format(time.RFC3339)),
		},
		{
			name:           "bad json",
			inputBody:      `{"type":123}`,
			mockBehavior:   func(_ *usecaseMocks.ReceptionUsecaseMock, req requests.AddProductRequest) {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"json: cannot unmarshal number into Go struct field AddProductRequest.type of type string"}`,
		},
		{
			name:         "usecase error",
			inputBody:    `{"type":"fail","pvzId":"123"}`,
			inputRequest: requests.AddProductRequest{Type: "fail", PvzId: "123"},
			mockBehavior: func(usecase *usecaseMocks.ReceptionUsecaseMock, req requests.AddProductRequest) {
				usecase.On("AddProductToReception", mock.Anything, req).Return(models.Product{}, errors.New("product add failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"product add failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(usecaseMocks.ReceptionUsecaseMock)
			handler := NewReceptionHandler(mockUsecase)

			tt.mockBehavior(mockUsecase, tt.inputRequest)

			req := httptest.NewRequest(http.MethodPost, "/api/products", strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()

			handler.AddProductToReception(w, req)

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

func TestPvzHandler_CloseLastReception(t *testing.T) {
	logger.AccessLogger = zap.NewNop()
	type mockBehavior func(usecase *usecaseMocks.ReceptionUsecaseMock, pvzId string)

	testTime := time.Date(2025, 4, 11, 23, 59, 50, 0, time.Local)

	tests := []struct {
		name           string
		pathParam      string
		mockPvzId      string
		mockBehavior   mockBehavior
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "success",
			pathParam: "123",
			mockPvzId: "123",
			mockBehavior: func(usecase *usecaseMocks.ReceptionUsecaseMock, pvzId string) {
				usecase.On("CloseReception", mock.Anything, pvzId).Return(models.Reception{
					Id:       "rec1",
					DateTime: testTime,
					PvzId:    "123",
					Status:   "close",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf(`{"id":"rec1","dateTime":"%s","pvzId":"123","status":"close"}`, testTime.Format(time.RFC3339)),
		},
		{
			name:      "usecase error",
			pathParam: "123",
			mockPvzId: "123",
			mockBehavior: func(usecase *usecaseMocks.ReceptionUsecaseMock, pvzId string) {
				usecase.On("CloseReception", mock.Anything, pvzId).Return(models.Reception{}, errors.New("close error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"close error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(usecaseMocks.ReceptionUsecaseMock)
			handler := NewReceptionHandler(mockUsecase)

			tt.mockBehavior(mockUsecase, tt.mockPvzId)

			path := "/pvz/close_last_reception"
			if tt.pathParam != "" {
				path = "/pvz/" + tt.pathParam + "/close_last_reception"
			}

			req := httptest.NewRequest(http.MethodPatch, path, nil)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"pvzId": tt.pathParam})
			handler.CloseLastReception(w, req)

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

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestPvzHandler_DeleteLastProduct(t *testing.T) {
	logger.AccessLogger = zap.NewNop()
	type mockBehavior func(usecase *usecaseMocks.ReceptionUsecaseMock, pvzId string)

	tests := []struct {
		name           string
		pvzId          string
		mockBehavior   mockBehavior
		expectedStatus int
		expectedBody   string
	}{
		{
			name:  "success",
			pvzId: "123",
			mockBehavior: func(usecase *usecaseMocks.ReceptionUsecaseMock, pvzId string) {
				usecase.On("DeleteLastProduct", mock.Anything, pvzId).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   ``,
		},
		{
			name:  "usecase error",
			pvzId: "123",
			mockBehavior: func(usecase *usecaseMocks.ReceptionUsecaseMock, pvzId string) {
				usecase.On("DeleteLastProduct", mock.Anything, pvzId).Return(errors.New("delete error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"delete error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(usecaseMocks.ReceptionUsecaseMock)
			handler := NewReceptionHandler(mockUsecase)

			tt.mockBehavior(mockUsecase, tt.pvzId)

			req := httptest.NewRequest(http.MethodDelete, "/api/pvz/"+tt.pvzId+"/delete_last_product", nil)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"pvzId": tt.pvzId})
			handler.DeleteLastProduct(w, req)

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
		})
	}
}
