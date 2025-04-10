package middleware

import (
	"avito_spring_staj_2025/internal/service/dsn"
	"avito_spring_staj_2025/internal/service/jwt"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
	"log"
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
	burstLimit        = 10000 // Максимальный «всплеск» запросов
)

var clientLimiters = sync.Map{}

func getLimiter(ip string) *rate.Limiter {
	limiter, exists := clientLimiters.Load(ip)
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), burstLimit)
		clientLimiters.Store(ip, limiter)
	}
	return limiter.(*rate.Limiter)
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		limiter := getLimiter(ip)

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

func DbConnect() *sql.DB {
	dsnString := dsn.FromEnv()
	if dsnString == "" {
		log.Fatal("DSN not provided. Please check your environment variables.")
	}

	db, err := sql.Open("postgres", dsnString)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	const maxOpenConns = 1000
	const maxIdleConns = 500
	const connMaxLifetime = 30 * time.Minute
	const connMaxIdleTime = 5 * time.Minute

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database")
	return db
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

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
