package handler

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/domain/responses"
	"avito_spring_staj_2025/internal/pvz/handler/gen"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"strconv"
	"time"
)

type PvzHandler struct {
	usecase PvzUsecase
}

func NewPvzHandler(usecase PvzUsecase) *PvzHandler {
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

	pvzs, err := h.usecase.GetPvzsInformation(ctx, startDate, endDate, limit, page)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	response := make([]responses.GetPvzsInformationResponse, 0, len(pvzs))
	for _, pvz := range pvzs {
		receptionsWithProducts := make([]responses.GetReceptionWithProducts, 0, len(pvz.Receptions))
		for _, reception := range pvz.Receptions {
			receptionsWithProducts = append(receptionsWithProducts, responses.GetReceptionWithProducts{
				Reception: reception,
				Products:  reception.Products,
			})
		}

		response = append(response, responses.GetPvzsInformationResponse{
			Pvz: models.Pvz{
				Id:               pvz.Id,
				RegistrationDate: pvz.RegistrationDate,
				City:             pvz.City,
			},
			Receptions: receptionsWithProducts,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.handleError(w, err, requestID)
		return
	}
}

func (h *PvzHandler) GetPvzListFromGrpc(w http.ResponseWriter, _ *http.Request) {
	_ = godotenv.Load()
	conn, err := grpc.NewClient(os.Getenv("GRPC_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := gen.NewPVZServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetPVZList(ctx, &gen.GetPVZListRequest{})
	if err != nil {
		http.Error(w, fmt.Sprintf("gRPC error: %v", err), http.StatusInternalServerError)
		return
	}

	type HttpPvz struct {
		ID               string    `json:"id"`
		RegistrationDate time.Time `json:"registration_date"`
		City             string    `json:"city"`
	}

	var result []HttpPvz
	for _, p := range resp.Pvzs {
		t := time.Unix(p.RegistrationDate.Seconds, int64(p.RegistrationDate.Nanos)).UTC()
		result = append(result, HttpPvz{
			ID:               p.Id,
			RegistrationDate: t,
			City:             p.City,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("json encode error: %v", err), http.StatusInternalServerError)
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
