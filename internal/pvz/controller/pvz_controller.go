package controller

import (
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/pvz/usecase"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()

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
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *PvzHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()

	var data requests.CreateReceptionRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.handleError(w, err, requestID)
		return
	}
	data = requests.CreateReceptionRequest{
		PvzId: sanitizer.Sanitize(data.PvzId),
	}

	response, err := h.usecase.CreateReception(ctx, data)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *PvzHandler) AddProductToReception(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()

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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *PvzHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()

	pvzId := sanitizer.Sanitize(mux.Vars(r)["pvzId"])

	err := h.usecase.DeleteLastProduct(ctx, pvzId)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PvzHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()

	pvzId := sanitizer.Sanitize(mux.Vars(r)["pvzId"])

	response, err := h.usecase.CloseReception(ctx, pvzId)
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
}

func (h *PvzHandler) GetPvzsInformation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()

	queryParams := r.URL.Query()
	startDateStr := sanitizer.Sanitize(queryParams.Get("startDate"))
	endDateStr := sanitizer.Sanitize(queryParams.Get("endDate"))

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			h.handleError(w, fmt.Errorf("invalid startDate: %v", err), requestID)
			return
		}
	}
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			h.handleError(w, fmt.Errorf("invalid endDate: %v", err), requestID)
			return
		}
	}

	pageStr := sanitizer.Sanitize(queryParams.Get("page"))
	limitStr := sanitizer.Sanitize(queryParams.Get("limit"))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	response, err := h.usecase.GetPvzsInformation(ctx, startDate, endDate, limit, page)
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
}

func (h *PvzHandler) handleError(w http.ResponseWriter, err error, requestID string) {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"errors": err.Error()}

	switch err.Error() {
	case "this city is not allowed", "active reception already exists",
		"this type is not allowed", "pvz not found",
		"no active reception", "invalid startDate", "invalid endDate", "no products in reception":
		w.WriteHeader(http.StatusBadRequest)
	case "this role is not allowed":
		w.WriteHeader(http.StatusForbidden)
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
