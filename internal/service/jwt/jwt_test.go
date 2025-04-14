package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJwtToken_CreateAndValidate(t *testing.T) {
	secret := "test-secret"
	service, err := NewJwtToken(secret)
	require.NoError(t, err)

	tests := []struct {
		name            string
		role            string
		expTime         int64
		wantCreateErr   bool
		wantValidateErr bool
		errContains     string
	}{
		{
			name:            "Success with admin role",
			role:            "admin",
			expTime:         time.Now().Add(1 * time.Hour).Unix(),
			wantCreateErr:   false,
			wantValidateErr: false,
		},
		{
			name:            "Success with user role",
			role:            "user",
			expTime:         time.Now().Add(1 * time.Hour).Unix(),
			wantCreateErr:   false,
			wantValidateErr: false,
		},
		{
			name:            "Empty role",
			role:            "",
			expTime:         time.Now().Add(1 * time.Hour).Unix(),
			wantCreateErr:   true,
			wantValidateErr: true,
			errContains:     "role is empty",
		},
		{
			name:            "Expired token",
			role:            "user",
			expTime:         time.Now().Add(-1 * time.Hour).Unix(),
			wantCreateErr:   false,
			wantValidateErr: true,
			errContains:     "invalid token: token is expired by 1h0m0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Create(tt.role, tt.expTime)
			if tt.wantCreateErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, token)

			claims, err := service.Validate(token)
			if tt.wantValidateErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.role, claims.Role)
			assert.Equal(t, tt.expTime, claims.ExpiresAt)
		})
	}
}

func TestJwtToken_ValidateInvalidTokens(t *testing.T) {
	secret := "test-secret"
	service, err := NewJwtToken(secret)
	require.NoError(t, err)

	validToken, err := service.Create("admin", time.Now().Add(1*time.Hour).Unix())
	require.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		wantErr     bool
		errContains string
	}{
		{
			name:        "Empty token",
			token:       "",
			wantErr:     true,
			errContains: "invalid token",
		},
		{
			name:        "Malformed token",
			token:       "malformed.token.here",
			wantErr:     true,
			errContains: "invalid token",
		},
		{
			name:        "Wrong signature",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJleHAiOjE2MTk5MjQwMDB9.wrong-signature-here",
			wantErr:     true,
			errContains: "invalid token",
		},
		{
			name:        "Valid token but wrong secret",
			token:       validToken,
			wantErr:     true,
			errContains: "invalid token",
		},
	}

	wrongSecretService, err := NewJwtToken("wrong-secret")
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var claims *JwtCsrfClaims

			if tt.name == "Valid token but wrong secret" {
				claims, err = wrongSecretService.Validate(tt.token)
			} else {
				claims, err = service.Validate(tt.token)
			}

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, claims)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestJwtToken_ParseSecretGetter(t *testing.T) {
	secret := "test-secret"
	service, err := NewJwtToken(secret)
	require.NoError(t, err)

	tests := []struct {
		name        string
		method      jwt.SigningMethod
		wantErr     bool
		errContains string
	}{
		{
			name:    "Valid HS256 method",
			method:  jwt.SigningMethodHS256,
			wantErr: false,
		},
		{
			name:        "Invalid method",
			method:      jwt.SigningMethodRS256,
			wantErr:     true,
			errContains: "bad sign method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := &jwt.Token{
				Method: tt.method,
			}

			key, err := service.ParseSecretGetter(token)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, key)
			} else {
				require.NoError(t, err)
				assert.Equal(t, service.Secret, key)
			}
		})
	}
}

func TestNewJwtToken(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{
			name:    "Success with secret",
			secret:  "valid-secret",
			wantErr: false,
		},
		{
			name:    "Empty secret",
			secret:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewJwtToken(tt.secret)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
