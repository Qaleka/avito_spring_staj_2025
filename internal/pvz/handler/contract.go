package handler

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"context"
	"time"
)

type PvzUsecase interface {
	CreatePvz(ctx context.Context, data *requests.CreatePvzRequest) error
	GetPvzsInformation(ctx context.Context, fromDate, toDate time.Time, limit, page int) ([]models.Pvz, error)
	GetAllPvzs(ctx context.Context) ([]models.Pvz, error)
}
