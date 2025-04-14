package handler

import (
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/domain/responses"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"net/http"
)

type ReceptionHandler struct {
	usecase ReceptionUsecase
}

func NewReceptionHandler(usecase ReceptionUsecase) *ReceptionHandler {
	return &ReceptionHandler{
		usecase: usecase,
	}
}

func (h *ReceptionHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
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

	reception, err := h.usecase.CreateReception(ctx, data)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	response := responses.CreateReceptionResponse{
		Id:       reception.Id,
		DateTime: reception.DateTime,
		PvzId:    reception.PvzId,
		Status:   reception.Status,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *ReceptionHandler) AddProductToReception(w http.ResponseWriter, r *http.Request) {
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

	product, err := h.usecase.AddProductToReception(ctx, data)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}

	response := responses.AddProductResponse{
		Id:          product.Id,
		DateTime:    product.DateTime,
		Type:        product.Type,
		ReceptionId: product.ReceptionId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *ReceptionHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
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

func (h *ReceptionHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()

	pvzId := sanitizer.Sanitize(mux.Vars(r)["pvzId"])

	reception, err := h.usecase.CloseReception(ctx, pvzId)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}

	response := responses.CloseReceptionResponse{
		Id:       reception.Id,
		DateTime: reception.DateTime,
		PvzId:    reception.PvzId,
		Status:   reception.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *ReceptionHandler) handleError(w http.ResponseWriter, err error, requestID string) {
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
