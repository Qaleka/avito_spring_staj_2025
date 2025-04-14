package handler

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"context"
)

type ReceptionUsecase interface {
	CreateReception(ctx context.Context, data requests.CreateReceptionRequest) (models.Reception, error)
	AddProductToReception(ctx context.Context, data requests.AddProductRequest) (models.Product, error)
	DeleteLastProduct(ctx context.Context, pvdId string) error
	CloseReception(ctx context.Context, pvdId string) (models.Reception, error)
}
