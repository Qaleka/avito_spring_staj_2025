package controller

import (
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/pvz/usecase"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type PvzHandler struct {
	usecase usecase.PvzUsecase
}

func NewPvzHandler(usecase usecase.PvzUsecase) *PvzHandler {
	return &PvzHandler{
		usecase: usecase,
	}
}

func (h *PvzHandler) CreatePvz(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()

	logger.AccessLogger.Info("Received SendCoins request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	var data requests.CreatePvzRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.handleError(w, err, requestID)
		return
	}
	data = requests.CreatePvzRequest{
		Id:               sanitizer.Sanitize(data.Id),
		RegistrationDate: data.RegistrationDate,
		City:             sanitizer.Sanitize(data.City),
	}

	err := h.usecase.CreatePvz(ctx, &data)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.handleError(w, err, requestID)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed CreatePvz request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusCreated))
}

func (h *PvzHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()
	logger.AccessLogger.Info("Received SendCoins request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()))

	var data requests.CreateReceptionRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.handleError(w, err, requestID)
		return
	}
	data = requests.CreateReceptionRequest{
		Id: sanitizer.Sanitize(data.Id),
	}
	response, err := h.usecase.CreateReception(ctx, data)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.handleError(w, err, requestID)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed CreateReception request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusCreated))
}

func (h *PvzHandler) AddProductToReception(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()

	logger.AccessLogger.Info("Received AddProductToReception request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()))

	var data requests.AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.handleError(w, err, requestID)
		return
	}

	data = requests.AddProductRequest{
		Type:  sanitizer.Sanitize(data.Type),
		PvzId: sanitizer.Sanitize(data.PvzId),
	}

	response, err := h.usecase.AddProductToReception(ctx, data)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.handleError(w, err, requestID)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed AddProductToReception request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusCreated))

}

func (h *PvzHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()

	logger.AccessLogger.Info("Received DeleteLastProduct request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()))

	rpzId := sanitizer.Sanitize(mux.Vars(r)["rpzId"])

	err := h.usecase.DeleteLastProduct(ctx, rpzId)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed DeleteLastProduct request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK))
}

func (h *PvzHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()

	logger.AccessLogger.Info("Received CloseReception request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()))

	rpzId := sanitizer.Sanitize(mux.Vars(r)["rpzId"])

	response, err := h.usecase.CloseReception(ctx, rpzId)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.handleError(w, err, requestID)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed CloseReception request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK))
}

func (h *PvzHandler) handleError(w http.ResponseWriter, err error, requestID string) {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"errors": err.Error()}

	switch err.Error() {
	case "not correct username", "not correct password",
		"jwt_token already exists", "Input contains invalid characters",
		"Input exceeds character limit":
		w.WriteHeader(http.StatusBadRequest)
	case "invalid credentials":
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
