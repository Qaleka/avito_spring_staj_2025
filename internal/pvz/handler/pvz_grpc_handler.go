package handler

import (
	"avito_spring_staj_2025/internal/pvz/handler/gen"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PvzGrpcHandler struct {
	gen.UnimplementedPVZServiceServer
	usecase PvzUsecase
}

func NewPvzGrpcHandler(usecase PvzUsecase) *PvzGrpcHandler {
	return &PvzGrpcHandler{usecase: usecase}
}

func (h *PvzGrpcHandler) GetPVZList(ctx context.Context, _ *gen.GetPVZListRequest) (*gen.GetPVZListResponse, error) {
	pvzs, err := h.usecase.GetAllPvzs(ctx)
	if err != nil {
		return nil, err
	}

	var responsePvzs []*gen.PVZ
	for _, pvz := range pvzs {
		responsePvzs = append(responsePvzs, &gen.PVZ{
			Id:               pvz.Id,
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             pvz.City,
		})
	}

	return &gen.GetPVZListResponse{
		Pvzs: responsePvzs,
	}, nil
}
