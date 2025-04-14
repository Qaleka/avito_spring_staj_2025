package usecase

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"errors"
	"time"
)

type PvzUsecase struct {
	pvzRepository PvzRepository
}

func NewPvzUsecase(pvzRepository PvzRepository) PvzUsecase {
	return PvzUsecase{
		pvzRepository: pvzRepository,
	}
}

func (pu PvzUsecase) CreatePvz(ctx context.Context, data *requests.CreatePvzRequest) error {
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

func (pu PvzUsecase) GetPvzsInformation(ctx context.Context, fromDate, toDate time.Time, limit, page int) ([]models.Pvz, error) {
	offset := (page - 1) * limit
	pvzs, err := pu.pvzRepository.GetPvzsFilteredByReceptionDate(ctx, fromDate, toDate, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range pvzs {
		pvz := &pvzs[i]
		receptions, err := pu.pvzRepository.GetPvzReceptions(ctx, pvz.Id)
		if err != nil {
			return nil, err
		}
		for j := range receptions {
			reception := &receptions[j]
			products, err := pu.pvzRepository.GetReceptionProducts(ctx, reception.Id)
			if err != nil {
				return nil, err
			}
			reception.Products = products
		}
		pvz.Receptions = receptions
	}
	return pvzs, nil
}

func (pu PvzUsecase) GetAllPvzs(ctx context.Context) ([]models.Pvz, error) {
	return pu.pvzRepository.GetAllPvzs(ctx)
}
