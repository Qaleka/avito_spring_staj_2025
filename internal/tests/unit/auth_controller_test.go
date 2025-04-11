package unit

import (
	"avito_spring_staj_2025/internal/auth/controller"
	"avito_spring_staj_2025/internal/tests/mocks/usecase_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDummyLogin_Success(t *testing.T) {
	mockUC := new(usecase_mocks.AuthUsecaseMock)
	handler := controller.NewAuthHandler(mockUC)

	reqBody := `{"role":"employee"}`
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	mockUC.On("DummyLogin", mock.Anything, "employee").
		Return("token", nil)

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}
