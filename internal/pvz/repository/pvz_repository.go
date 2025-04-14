package repository

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
	"time"
)

type PvzRepository struct {
	db *sql.DB
}

func NewPvzRepository(db *sql.DB) PvzRepository {
	return PvzRepository{
		db: db,
	}
}

func (r PvzRepository) CreatePvz(ctx context.Context, data *requests.CreatePvzRequest) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CreatePvz called", zap.String("request_id", requestID), zap.String("pvz_id", data.Id))
	queryBuilder := sq.Insert("pvzs").
		Columns("id", "registration_date", "city").
		Values(data.Id, data.RegistrationDate, data.City).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build SQL", zap.Error(err))
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.DBLogger.Error("failed to insert pvz", zap.Error(err))
		return err
	}

	logger.DBLogger.Info("Pvz successfully created",
		zap.String("request_id", requestID),
		zap.String("pvz_id", data.Id),
	)

	return nil
}

func (r PvzRepository) GetPvzsFilteredByReceptionDate(ctx context.Context, from, to time.Time, limit, offset int) ([]models.Pvz, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetPvzsFilteredByReceptionDate called",
		zap.String("request_id", requestID),
	)
	queryBuilder := sq.
		Select("id", "registration_date", "city").
		From("pvzs").
		PlaceholderFormat(sq.Dollar)

	whereClauses := sq.And{}

	if !from.IsZero() {
		whereClauses = append(whereClauses, sq.GtOrEq{"registration_date": from})
	}
	if !to.IsZero() {
		whereClauses = append(whereClauses, sq.LtOrEq{"registration_date": to})
	}

	if len(whereClauses) > 0 {
		queryBuilder = queryBuilder.Where(whereClauses)
	}

	queryBuilder = queryBuilder.
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.DBLogger.Error("failed to close rows", zap.Error(err))
		}
	}()

	var pvzs []models.Pvz
	for rows.Next() {
		var p models.Pvz
		if err := rows.Scan(&p.Id, &p.RegistrationDate, &p.City); err != nil {
			return nil, err
		}
		pvzs = append(pvzs, p)
	}

	return pvzs, nil
}

func (r PvzRepository) GetPvzReceptions(ctx context.Context, pvzId string) ([]models.Reception, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetPvzReceptions called",
		zap.String("request_id", requestID),
		zap.String("pvz_id", pvzId),
	)

	queryBuilder := sq.Select("id", "date_time", "pvz_id", "status").
		From("receptions").
		Where(sq.Eq{"pvz_id": pvzId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build GetPvzReceptions SQL", zap.Error(err))
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.DBLogger.Error("failed to close rows", zap.Error(err))
		}
	}()

	var receptions []models.Reception
	for rows.Next() {
		var r models.Reception
		if err := rows.Scan(&r.Id, &r.DateTime, &r.PvzId, &r.Status); err != nil {
			return nil, err
		}
		receptions = append(receptions, r)
	}

	return receptions, nil
}

func (r PvzRepository) GetReceptionProducts(ctx context.Context, receptionId string) ([]models.Product, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetReceptionProducts called",
		zap.String("request_id", requestID),
		zap.String("reception_id", receptionId),
	)

	queryBuilder := sq.Select("id", "date_time", "type", "reception_id").
		From("products").
		Where(sq.Eq{"reception_id": receptionId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build GetReceptionProducts SQL", zap.Error(err))
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.DBLogger.Error("failed to close rows", zap.Error(err))
		}
	}()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.Id, &p.DateTime, &p.Type, &p.ReceptionId); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r PvzRepository) GetAllPvzs(ctx context.Context) ([]models.Pvz, error) {
	query := "SELECT id, registration_date, city FROM pvzs"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pvzs: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.DBLogger.Error("failed to close rows", zap.Error(err))
		}
	}()

	var pvzs []models.Pvz
	for rows.Next() {
		var pvz models.Pvz
		if err := rows.Scan(&pvz.Id, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		pvzs = append(pvzs, pvz)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return pvzs, nil
}
