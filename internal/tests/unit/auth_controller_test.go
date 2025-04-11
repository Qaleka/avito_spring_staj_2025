package unit

import (
	"avito_spring_staj_2025/domain/responses"
	"avito_spring_staj_2025/internal/auth/controller"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/tests/mocks/usecase_mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDummyLogin(t *testing.T) {
	logger.AccessLogger = zap.NewNop()
	tests := []struct {
		name           string
		body           string
		mockSetup      func(m *usecase_mocks.AuthUsecaseMock)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"role": "employee"}`,
			mockSetup: func(m *usecase_mocks.AuthUsecaseMock) {
				m.On("DummyLogin", mock.Anything, "employee").
					Return("token", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid body",
			body: `{invalid-json`,
			mockSetup: func(m *usecase_mocks.AuthUsecaseMock) {
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "service error",
			body: `{"role": "aboba"}`,
			mockSetup: func(m *usecase_mocks.AuthUsecaseMock) {
				m.On("DummyLogin", mock.Anything, "aboba").
					Return("", errors.New("unauthorized"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(usecase_mocks.AuthUsecaseMock)
			if tt.mockSetup != nil {
				tt.mockSetup(mockUsecase)
			}
			handler := controller.NewAuthHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodPost, "/api/dummyLogin", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler.DummyLogin(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestRegister(t *testing.T) {
	logger.AccessLogger = zap.NewNop()
	tests := []struct {
		name           string
		body           string
		headers        map[string]string
		mockSetup      func(m *usecase_mocks.AuthUsecaseMock)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"email":"test@example.com","password":"123456","role":"employee"}`,
			mockSetup: func(m *usecase_mocks.AuthUsecaseMock) {
				m.On("Register", mock.Anything, mock.Anything).
					Return(responses.RegisterResponse{Id: "1", Email: "test@example.com", Role: "employee"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid body",
			body:           `{bad json`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "authorization token exists",
			body: `{"email":"x","password":"x","role":"x"}`,
			headers: map[string]string{
				"Authorization": "Bearer something",
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "usecase error",
			body: `{"email":"e","password":"p","role":"r"}`,
			mockSetup: func(m *usecase_mocks.AuthUsecaseMock) {
				m.On("Register", mock.Anything, mock.Anything).Return(responses.RegisterResponse{}, errors.New("err"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(usecase_mocks.AuthUsecaseMock)
			if tt.mockSetup != nil {
				tt.mockSetup(mockUsecase)
			}
			handler := controller.NewAuthHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			handler.Register(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	logger.AccessLogger = zap.NewNop()
	tests := []struct {
		name           string
		body           string
		headers        map[string]string
		mockSetup      func(m *usecase_mocks.AuthUsecaseMock)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"email":"test@example.com","password":"123456"}`,
			mockSetup: func(m *usecase_mocks.AuthUsecaseMock) {
				m.On("Login", mock.Anything, mock.Anything).Return("token", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid body",
			body:           `{invalid`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "token already exists",
			body: `{"email":"e","password":"p"}`,
			headers: map[string]string{
				"Authorization": "Bearer token",
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "usecase error",
			body: `{"email":"e","password":"p"}`,
			mockSetup: func(m *usecase_mocks.AuthUsecaseMock) {
				m.On("Login", mock.Anything, mock.Anything).Return("", errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(usecase_mocks.AuthUsecaseMock)
			if tt.mockSetup != nil {
				tt.mockSetup(mockUsecase)
			}
			handler := controller.NewAuthHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}
