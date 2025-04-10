package usecase

import (
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/domain/responses"
	"context"
)

type PvzUsecase interface {
	CreatePvz(ctx context.Context, data *requests.CreatePvzRequest) error
	CreateReception(ctx context.Context, data requests.CreateReceptionRequest) (responses.CreateReceptionResponse, error)
	AddProductToReception(ctx context.Context, data requests.AddProductRequest) (responses.AddProductResponse, error)
	DeleteLastProduct(ctx context.Context, pvdId string) error
	CloseReception(ctx context.Context, pvdId string) (responses.CloseReceptionResponse, error)
}
