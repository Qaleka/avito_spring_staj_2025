package jwt_mocks

import (
	jwtService "avito_spring_staj_2025/internal/service/jwt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
)

type MockJwtService struct {
	mock.Mock
}

func (m *MockJwtService) Create(role string, exp int64) (string, error) {
	args := m.Called(role, exp)
	return args.String(0), args.Error(1)
}

func (m *MockJwtService) Validate(tokenString string) (*jwtService.JwtCsrfClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*jwtService.JwtCsrfClaims), args.Error(1)
}

func (m *MockJwtService) ParseSecretGetter(token *jwt.Token) (interface{}, error) {
	args := m.Called(token)
	return args.Get(0), args.Error(1)
}
