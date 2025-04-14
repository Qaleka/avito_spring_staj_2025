package usecase

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"context"
	"time"
)

type PvzRepository interface {
	CreatePvz(ctx context.Context, data *requests.CreatePvzRequest) error
	GetPvzReceptions(ctx context.Context, pvzId string) ([]models.Reception, error)
	GetReceptionProducts(ctx context.Context, receptionId string) ([]models.Product, error)
	GetPvzsFilteredByReceptionDate(ctx context.Context, from, to time.Time, limit, offset int) ([]models.Pvz, error)
	GetAllPvzs(ctx context.Context) ([]models.Pvz, error)
}
