package handler

import (
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/domain/responses"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"encoding/json"
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	usecase AuthUsecase
}

func NewAuthHandler(usecase AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	var role requests.DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		h.handleError(w, err, requestID)
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	role = requests.DummyLoginRequest{
		Role: sanitizer.Sanitize(role.Role),
	}

	token, err := h.usecase.DummyLogin(ctx, role.Role)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		h.handleError(w, errors.New("authorization token already exists"), requestID)
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	var credentials requests.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.handleError(w, err, requestID)
		return
	}

	credentials = requests.RegisterRequest{
		Email:    sanitizer.Sanitize(credentials.Email),
		Password: sanitizer.Sanitize(credentials.Password),
		Role:     sanitizer.Sanitize(credentials.Role),
	}

	user, err := h.usecase.Register(ctx, credentials)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	response := responses.RegisterResponse{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		h.handleError(w, errors.New("authorization token already exists"), requestID)
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	var credentials requests.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.handleError(w, err, requestID)
		return
	}

	credentials = requests.LoginRequest{
		Email:    sanitizer.Sanitize(credentials.Email),
		Password: sanitizer.Sanitize(credentials.Password),
	}

	token, err := h.usecase.Login(ctx, credentials)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error, requestID string) {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"errors": err.Error()}

	switch err.Error() {
	case "invalid role", "invalid credentials",
		"user with this email already exists", "user not found",
		"Input exceeds character limit":
		w.WriteHeader(http.StatusBadRequest)
	case "invalid password":
		w.WriteHeader(http.StatusUnauthorized)
	case "failed to generate error response":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if jsonErr := json.NewEncoder(w).Encode(errorResponse); jsonErr != nil {
		logger.AccessLogger.Error("Failed to encode error response",
			zap.String("request_id", requestID),
			zap.Error(jsonErr),
		)
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
}
