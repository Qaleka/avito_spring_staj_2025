package usecase

import (
	"avito_spring_staj_2025/domain/models"
	"context"
)

type ReceptionRepository interface {
	CreateReception(ctx context.Context, data models.Reception) error
	GetPvzById(ctx context.Context, pvzId string) (*models.Pvz, error)
	GetCurrentReception(ctx context.Context, pvzId string) (*models.Reception, error)
	AddProductToReception(ctx context.Context, product models.Product) error
	GetLastProductInReception(ctx context.Context, receptionId string) (*models.Product, error)
	DeleteProductById(ctx context.Context, productId string) error
	CloseReception(ctx context.Context, reception *models.Reception) error
}
