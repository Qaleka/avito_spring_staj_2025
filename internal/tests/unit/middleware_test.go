package unit

import (
	"avito_spring_staj_2025/internal/service/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("should add request ID to header and context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler := middleware.RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetRequestID(r.Context())
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

		handler := middleware.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		handler := middleware.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		limiter := middleware.GetLimiter(req.RemoteAddr)
		for i := 0; i < middleware.BurstLimit+1; i++ {
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

			handler := middleware.EnableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		})
	}
}

//func TestRoleMiddleware(t *testing.T) {
//	tests := []struct {
//		name           string
//		authHeader     string
//		mockValidate   func(token string) (*jwt_service.JwtCsrfClaims, error)
//		expectedStatus int
//	}{
//		{
//			name:       "should set role in context for valid token",
//			authHeader: "Bearer valid.token",
//			mockValidate: func(token string) (*jwt_service.JwtCsrfClaims, error) {
//				return &jwt_service.JwtCsrfClaims{Role: "admin"}, nil
//			},
//			expectedStatus: http.StatusOK,
//		},
//		{
//			name:           "should reject request without auth header",
//			authHeader:     "",
//			mockValidate:   nil,
//			expectedStatus: http.StatusUnauthorized,
//		},
//		{
//			name:       "should reject request with invalid token",
//			authHeader: "Bearer invalid.token",
//			mockValidate: func(token string) (*jwt_service.JwtCsrfClaims, error) {
//				return nil, jwt.ErrSignatureInvalid
//			},
//			expectedStatus: http.StatusUnauthorized,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			req := httptest.NewRequest("GET", "/", nil)
//			if tt.authHeader != "" {
//				req.Header.Set("Authorization", tt.authHeader)
//			}
//
//			rr := httptest.NewRecorder()
//
//			mockJwt := &jwt_mocks.MockJwtService{ValidateFunc: tt.mockValidate}
//			handler := middleware.RoleMiddleware(mockJwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				role := r.Context().Value(middleware.ContextKeyRole)
//				assert.Equal(t, "admin", role)
//				w.WriteHeader(http.StatusOK)
//			}))
//
//			handler.ServeHTTP(rr, req)
//			assert.Equal(t, tt.expectedStatus, rr.Code)
//		})
//	}
//}
//
//func TestWithLoggingAndMetrics(t *testing.T) {
//	// Мокируем логгер и метрики для теста
//	oldLogger := logger.AccessLogger
//	defer func() { logger.AccessLogger = oldLogger }()
//
//	var loggedMessages []string
//	logger.AccessLogger = zap.NewNop()
//
//	req := httptest.NewRequest("GET", "/test", nil)
//	rr := httptest.NewRecorder()
//
//	handler := middleware.WithLoggingAndMetrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Проверяем что timeout установлен в контексте
//		_, ok := r.Context().Deadline()
//		assert.True(t, ok)
//		w.WriteHeader(http.StatusOK)
//	}))
//
//	handler.ServeHTTP(rr, req)
//	assert.Equal(t, http.StatusOK, rr.Code)
//	assert.Equal(t, 0, len(loggedMessages)) // Проверяем что вызовы логирования были
//}
//
//func TestWithCustomMetric(t *testing.T) {
//	metric := prometheus.NewCounterVec(
//		prometheus.CounterOpts{
//			Name: "test_metric",
//			Help: "Test metric",
//		},
//		[]string{"path"},
//	)
//
//	req := httptest.NewRequest("GET", "/test-path", nil)
//	rr := httptest.NewRecorder()
//
//	handler := middleware.WithCustomMetric(metric)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//	}))
//
//	handler.ServeHTTP(rr, req)
//	assert.Equal(t, http.StatusOK, rr.Code)
//
//	// Проверяем что метрика была инкрементирована
//	metrics, err := metric.GetMetricWithLabelValues("/test-path")
//	require.NoError(t, err)
//
//	var m prometheus.Metric
//	require.NotPanics(t, func() { m = metrics.(prometheus.Metric) })
//
//	pb := &prometheus.Metric{}
//	m.Write(pb)
//	assert.True(t, pb.Counter != nil && pb.Counter.Value != nil && *pb.Counter.Value > 0)
//}
//
//func TestChainMiddlewares(t *testing.T) {
//	var calls []string
//
//	m1 := func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			calls = append(calls, "m1")
//			next.ServeHTTP(w, r)
//		}
//	}
//
//	m2 := func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			calls = append(calls, "m2")
//			next.ServeHTTP(w, r)
//		}
//	}
//
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		calls = append(calls, "handler")
//		w.WriteHeader(http.StatusOK)
//	})
//
//	chained := middleware.ChainMiddlewares(handler, m1, m2)
//
//	req := httptest.NewRequest("GET", "/", nil)
//	rr := httptest.NewRecorder()
//	chained.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusOK, rr.Code)
//	assert.Equal(t, []string{"m1", "m2", "handler"}, calls)
//}
//
//func TestWithTimeout(t *testing.T) {
//	t.Run("should return context with timeout", func(t *testing.T) {
//		ctx := context.Background()
//		newCtx, cancel := middleware.WithTimeout(ctx)
//		defer cancel()
//
//		deadline, ok := newCtx.Deadline()
//		assert.True(t, ok)
//		assert.WithinDuration(t, time.Now().Add(middleware.requestTimeout), deadline, time.Second)
//	})
//}
//
//func TestGetRequestID(t *testing.T) {
//	t.Run("should return request ID from context", func(t *testing.T) {
//		expectedID := "test-request-id"
//		ctx := context.WithValue(context.Background(), middleware.RequestIDKey, expectedID)
//
//		assert.Equal(t, expectedID, middleware.GetRequestID(ctx))
//	})
//
//	t.Run("should return empty string when no request ID in context", func(t *testing.T) {
//		assert.Empty(t, middleware.GetRequestID(context.Background()))
//	})
//}
