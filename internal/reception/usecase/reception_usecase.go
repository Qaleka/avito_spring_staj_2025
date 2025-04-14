package usecase

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type ReceptionUsecase struct {
	pvzRepository ReceptionRepository
}

func NewReceptionUsecase(receptionRepository ReceptionRepository) ReceptionUsecase {
	return ReceptionUsecase{
		pvzRepository: receptionRepository,
	}
}

func (pu ReceptionUsecase) CreateReception(ctx context.Context, data requests.CreateReceptionRequest) (models.Reception, error) {
	if ctx.Value(middleware.ContextKeyRole).(string) != "employee" {
		return models.Reception{}, errors.New("this role is not allowed")
	}

	_, err := pu.pvzRepository.GetCurrentReception(ctx, data.PvzId)
	if err != nil && err.Error() != "no active reception" {
		return models.Reception{}, err
	}
	if err == nil {
		return models.Reception{}, errors.New("active reception already exists")
	}

	_, err = pu.pvzRepository.GetPvzById(ctx, data.PvzId)
	if err != nil {
		return models.Reception{}, err
	}

	reception := models.Reception{
		Id:       uuid.New().String(),
		DateTime: time.Now(),
		PvzId:    data.PvzId,
		Status:   "in_progress",
	}

	err = pu.pvzRepository.CreateReception(ctx, reception)
	if err != nil {
		return models.Reception{}, err
	}
	return reception, nil
}

func (pu ReceptionUsecase) AddProductToReception(ctx context.Context, data requests.AddProductRequest) (models.Product, error) {
	if ctx.Value(middleware.ContextKeyRole).(string) != "employee" {
		return models.Product{}, errors.New("this role is not allowed")
	}

	if data.Type != models.CLOTHES_TYPE && data.Type != models.BOOTS_TYPE && data.Type != models.ELECTRONIC_TYPE {
		return models.Product{}, errors.New("this type is not allowed")
	}

	_, err := pu.pvzRepository.GetPvzById(ctx, data.PvzId)
	if err != nil {
		return models.Product{}, err
	}

	reception, err := pu.pvzRepository.GetCurrentReception(ctx, data.PvzId)
	if err != nil {
		return models.Product{}, err
	}
	product := models.Product{
		Id:          uuid.New().String(),
		Type:        data.Type,
		ReceptionId: reception.Id,
		DateTime:    time.Now(),
	}

	err = pu.pvzRepository.AddProductToReception(ctx, product)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (pu ReceptionUsecase) DeleteLastProduct(ctx context.Context, pvzId string) error {
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

func (pu ReceptionUsecase) CloseReception(ctx context.Context, pvzId string) (models.Reception, error) {
	if ctx.Value(middleware.ContextKeyRole).(string) != "employee" {
		return models.Reception{}, errors.New("this role is not allowed")
	}
	_, err := pu.pvzRepository.GetPvzById(ctx, pvzId)
	if err != nil {
		return models.Reception{}, err
	}

	reception, err := pu.pvzRepository.GetCurrentReception(ctx, pvzId)
	if err != nil {
		return models.Reception{}, err
	}

	err = pu.pvzRepository.CloseReception(ctx, reception)
	if err != nil {
		return models.Reception{}, err
	}

	return *reception, nil
}
