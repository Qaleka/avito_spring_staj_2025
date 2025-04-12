package middleware

import (
	"avito_spring_staj_2025/internal/service/jwt"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/metrics"
	"context"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
	"sync"
	"time"
)

type contextKey string

const (
	loggerKey contextKey = "logger"
)

const requestTimeout = 1000 * time.Second

func WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, requestTimeout)
}

type key int

const RequestIDKey key = 0

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

const (
	requestsPerSecond = 10000 // Лимит запросов в секунду для каждого IP
	BurstLimit        = 10000 // Максимальный «всплеск» запросов
)

var clientLimiters = sync.Map{}

func GetLimiter(ip string) *rate.Limiter {
	limiter, exists := clientLimiters.Load(ip)
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), BurstLimit)
		clientLimiters.Store(ip, limiter)
	}
	return limiter.(*rate.Limiter)
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		limiter := GetLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Set-Cookie, X-CSRFToken, x-csrftoken, X-CSRF-Token")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

const (
	ContextKeyRole contextKey = "role"
	ContextKeyUser contextKey = "user"
)

func RoleMiddleware(jwtService jwt.JwtTokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtService.Validate(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithLoggingAndMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := GetRequestID(r.Context())
		ctx, cancel := WithTimeout(r.Context())
		defer cancel()

		logger.AccessLogger.Info("Received request",
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
		)

		defer func() {
			duration := time.Since(start).Seconds()
			metrics.TotalRequests.WithLabelValues(r.URL.Path).Inc()
			metrics.AnswerDuration.WithLabelValues(r.URL.Path).Observe(duration)

			logger.AccessLogger.Info("Completed request",
				zap.String("request_id", requestID),
				zap.Duration("duration", time.Since(start)),
				zap.Int("status", http.StatusOK),
			)
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithCustomMetric(metric *prometheus.CounterVec) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			metric.WithLabelValues(r.URL.Path).Inc()
		})
	}
}

func ChainMiddlewares(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
