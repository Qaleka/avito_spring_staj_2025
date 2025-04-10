package usecase

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/domain/responses"
	"avito_spring_staj_2025/internal/pvz/repository"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type pvzUsecase struct {
	pvzRepository repository.PvzRepository
}

func NewPvzUsecase(pvzRepository repository.PvzRepository) PvzUsecase {
	return &pvzUsecase{
		pvzRepository: pvzRepository,
	}
}

func (pu *pvzUsecase) CreatePvz(ctx context.Context, data *requests.CreatePvzRequest) error {
	if data.City != "Москва" && data.City != "Санкт-Петербург" && data.City != "Казань" {
		return errors.New("this city is not allowed")
	}
	if ctx.Value(middleware.ContextKeyRole).(string) != "moderator" {
		return errors.New("this role is not allowed")
	}
	err := pu.pvzRepository.CreatePvz(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (pu *pvzUsecase) CreateReception(ctx context.Context, data requests.CreateReceptionRequest) (responses.CreateReceptionResponse, error) {
	if ctx.Value(middleware.ContextKeyRole).(string) != "employee" {
		return responses.CreateReceptionResponse{}, errors.New("this role is not allowed")
	}

	_, err := pu.pvzRepository.GetPvzById(ctx, data.Id)
	if err != nil {
		return responses.CreateReceptionResponse{}, err
	}

	reception := models.Reception{
		Id:       uuid.New().String(),
		DateTime: time.Now(),
		PvzId:    data.Id,
		Status:   "in_progress",
	}

	err = pu.pvzRepository.CreateReception(ctx, reception)
	if err != nil {
		return responses.CreateReceptionResponse{}, err
	}
	return responses.CreateReceptionResponse{
		Id:       reception.Id,
		DateTime: reception.DateTime,
		PvzId:    reception.PvzId,
		Status:   reception.Status,
	}, nil
}

func (pu *pvzUsecase) AddProductToReception(ctx context.Context, data requests.AddProductRequest) (responses.AddProductResponse, error) {
	if ctx.Value(middleware.ContextKeyRole).(string) != "employee" {
		return responses.AddProductResponse{}, errors.New("this role is not allowed")
	}

	if data.Type != models.CLOTHES_TYPE && data.Type != models.BOOTS_TYPE && data.Type != models.ELECTRONIC_TYPE {
		return responses.AddProductResponse{}, errors.New("this type is not allowed")
	}

	_, err := pu.pvzRepository.GetPvzById(ctx, data.PvzId)
	if err != nil {
		return responses.AddProductResponse{}, err
	}

	reception, err := pu.pvzRepository.GetCurrentReception(ctx, data.PvzId)
	if err != nil {
		return responses.AddProductResponse{}, err
	}
	product := models.Product{
		Id:          uuid.New().String(),
		Type:        data.Type,
		ReceptionId: reception.Id,
		DateTime:    time.Now(),
	}

	err = pu.pvzRepository.AddProductToReception(ctx, product)
	if err != nil {
		return responses.AddProductResponse{}, err
	}

	return responses.AddProductResponse{
		Id:          product.Id,
		Type:        product.Type,
		ReceptionId: product.ReceptionId,
		DateTime:    product.DateTime,
	}, nil
}

func (pu *pvzUsecase) DeleteLastProduct(ctx context.Context, pvzId string) error {
	if ctx.Value(middleware.ContextKeyRole).(string) != "employee" {
		return errors.New("this role is not allowed")
	}
	_, err := pu.pvzRepository.GetPvzById(ctx, pvzId)
	if err != nil {
		return err
	}

	reception, err := pu.pvzRepository.GetCurrentReception(ctx, pvzId)
	if err != nil {
		return err
	}

	product, err := pu.pvzRepository.GetLastProductInReception(ctx, reception.Id)
	if err != nil {
		return err
	}

	err = pu.pvzRepository.DeleteProductById(ctx, product.Id)
	if err != nil {
		return err
	}

	return nil
}

func (pu *pvzUsecase) CloseReception(ctx context.Context, pvzId string) (responses.CloseReceptionResponse, error) {
	if ctx.Value(middleware.ContextKeyRole).(string) != "employee" {
		return responses.CloseReceptionResponse{}, errors.New("this role is not allowed")
	}
	_, err := pu.pvzRepository.GetPvzById(ctx, pvzId)
	if err != nil {
		return responses.CloseReceptionResponse{}, err
	}

	reception, err := pu.pvzRepository.GetCurrentReception(ctx, pvzId)
	if err != nil {
		return responses.CloseReceptionResponse{}, err
	}

	err = pu.pvzRepository.CloseReception(ctx, reception)
	if err != nil {
		return responses.CloseReceptionResponse{}, err
	}

	return responses.CloseReceptionResponse{
		Id:       reception.Id,
		DateTime: reception.DateTime,
		PvzId:    reception.PvzId,
		Status:   reception.Status,
	}, nil
}
