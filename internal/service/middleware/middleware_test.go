package middleware

import (
	jwt_service "avito_spring_staj_2025/internal/service/jwt"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/tests/mocks/jwt_mocks"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("should add request ID to header and context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := GetRequestID(r.Context())
			assert.NotEmpty(t, requestID)
			assert.Equal(t, requestID, w.Header().Get("X-Request-ID"))
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	t.Run("should allow requests under limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		rr := httptest.NewRecorder()

		handler := RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		for i := 0; i < 10; i++ {
			handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
		}
	})

	t.Run("should block requests over limit", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "127.0.0.1:54321"
		rr := httptest.NewRecorder()

		handler := RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		limiter := GetLimiter(req.RemoteAddr)
		for i := 0; i < BurstLimit+1; i++ {
			limiter.Allow()
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	})
}

func TestEnableCORS(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "should set CORS headers for regular request",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "should handle OPTIONS request",
			method:         "OPTIONS",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			rr := httptest.NewRecorder()

			handler := EnableCORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		})
	}
}

func TestRoleMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *jwtMocks.MockJwtService)
		authHeader     string
		expectedStatus int
		expectedRole   string
	}{
		{
			name: "success with valid token",
			setupMock: func(m *jwtMocks.MockJwtService) {
				m.On("Validate", "valid.token").Return(
					&jwt_service.JwtCsrfClaims{Role: "admin"}, nil)
			},
			authHeader:     "Bearer valid.token",
			expectedStatus: http.StatusOK,
			expectedRole:   "admin",
		},
		{
			name:           "missing authorization header",
			setupMock:      func(_ *jwtMocks.MockJwtService) {},
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid token format",
			setupMock: func(m *jwtMocks.MockJwtService) {
				m.On("Validate", "invalid.token").Return(
					(*jwt_service.JwtCsrfClaims)(nil), jwt.ErrSignatureInvalid)
			},
			authHeader:     "Bearer invalid.token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "expired token",
			setupMock: func(m *jwtMocks.MockJwtService) {
				m.On("Validate", "expired.token").Return(
					(*jwt_service.JwtCsrfClaims)(nil), jwt.NewValidationError("token expired", jwt.ValidationErrorExpired))
			},
			authHeader:     "Bearer expired.token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJwt := new(jwtMocks.MockJwtService)
			tt.setupMock(mockJwt)

			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()

			handler := RoleMiddleware(mockJwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				role := r.Context().Value(ContextKeyRole)
				assert.Equal(t, tt.expectedRole, role)
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockJwt.AssertExpectations(t)
		})
	}
}

func TestWithLoggingAndMetrics(t *testing.T) {
	oldLogger := logger.AccessLogger
	defer func() { logger.AccessLogger = oldLogger }()

	var loggedMessages []string
	logger.AccessLogger = zap.NewNop()

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler := WithLoggingAndMetrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Deadline()
		assert.True(t, ok)
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 0, len(loggedMessages))
}

func TestWithCustomMetric(t *testing.T) {
	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_metric",
			Help: "Test metric",
		},
		[]string{"path"},
	)

	registry := prometheus.NewRegistry()
	err := registry.Register(metric)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test-path", nil)
	rr := httptest.NewRecorder()

	handler := WithCustomMetric(metric)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	metrics, err := registry.Gather()
	require.NoError(t, err)
	require.Len(t, metrics, 1)

	assert.Equal(t, "test_metric", metrics[0].GetName())
	require.Len(t, metrics[0].GetMetric(), 1)

	counter := metrics[0].GetMetric()[0].GetCounter()
	assert.NotNil(t, counter)
	assert.Equal(t, 1.0, counter.GetValue())
}

func TestChainMiddlewares(t *testing.T) {
	var calls []string

	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "m1")
			next.ServeHTTP(w, r)
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "m2")
			next.ServeHTTP(w, r)
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls = append(calls, "handler")
		w.WriteHeader(http.StatusOK)
	})

	chained := ChainMiddlewares(handler, m1, m2)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	chained.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, []string{"m2", "m1", "handler"}, calls)
}

func TestWithTimeout(t *testing.T) {
	t.Run("should return context with timeout", func(t *testing.T) {
		ctx := context.Background()
		newCtx, cancel := WithTimeout(ctx)
		defer cancel()

		deadline, ok := newCtx.Deadline()
		assert.True(t, ok)
		assert.WithinDuration(t, time.Now().Add(RequestTimeout), deadline, time.Second)
	})
}

func TestGetRequestID(t *testing.T) {
	t.Run("should return request ID from context", func(t *testing.T) {
		expectedID := "test-request-id"
		ctx := context.WithValue(context.Background(), RequestIDKey, expectedID)

		assert.Equal(t, expectedID, GetRequestID(ctx))
	})

	t.Run("should return empty string when no request ID in context", func(t *testing.T) {
		assert.Empty(t, GetRequestID(context.Background()))
	})
}
