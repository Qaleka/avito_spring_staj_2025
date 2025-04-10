package repository

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type pvzRepository struct {
	db *sql.DB
}

func NewPvzRepository(db *sql.DB) PvzRepository {
	return &pvzRepository{
		db: db,
	}
}

func (r *pvzRepository) CreatePvz(ctx context.Context, data *requests.CreatePvzRequest) error {
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

func (r *pvzRepository) CreateReception(ctx context.Context, data models.Reception) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CreateReception called", zap.String("request_id", requestID))
	queryBuilder := sq.Insert("receptions").
		Columns("id", "date_time", "pvz_id", "status").
		Values(data.Id, data.DateTime, data.PvzId, data.Status).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build SQL", zap.Error(err))
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.DBLogger.Error("failed to insert reception", zap.Error(err))
		return err
	}
	logger.DBLogger.Info("Reception successfully created",
		zap.String("request_id", requestID))

	return nil
}

func (r *pvzRepository) GetPvzById(ctx context.Context, pvzId string) (*models.Pvz, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetPvzById called", zap.String("request_id", requestID), zap.String("pvz_id", pvzId))
	queryBuilder := sq.Select("id", "registration_date", "city").
		From("pvzs").
		Where(sq.Eq{"id": pvzId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build SQL", zap.Error(err))
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, query, args...)

	var pvz models.Pvz
	err = row.Scan(&pvz.Id, &pvz.RegistrationDate, &pvz.City)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.DBLogger.Info("pvz not found",
				zap.String("request_id", requestID),
				zap.String("pvz_id", pvzId),
			)
			return nil, errors.New("pvz not found")
		}
		logger.DBLogger.Error("failed to scan user", zap.Error(err))
		return nil, err
	}

	return &pvz, nil
}

func (r *pvzRepository) GetCurrentReception(ctx context.Context, pvzId string) (*models.Reception, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetCurrentReception called", zap.String("request_id", requestID), zap.String("pvz_id", pvzId))
	queryBuilder := sq.Select("id", "date_time", "pvz_id", "status").
		From("receptions").
		Where(sq.Eq{"id": pvzId}).
		Where(sq.Eq{"status": models.STATUS_ACTIVE}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build SQL", zap.Error(err))
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	var reception models.Reception

	err = row.Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.DBLogger.Info("reception not found",
				zap.String("request_id", requestID),
				zap.String("pvz_id", pvzId))
		}
		return nil, errors.New("no active reception")
	}
	return &reception, nil
}

func (r *pvzRepository) AddProductToReception(ctx context.Context, product models.Product) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("AddProductToReception called", zap.String("request_id", requestID))
	queryBuilder := sq.Insert("products").
		Columns("id", "date_time", "type", "reception_id").
		Values(product.Id, product.DateTime, product.Type, product.ReceptionId).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build SQL", zap.Error(err))
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.DBLogger.Error("failed to insert product", zap.Error(err))
		return err
	}
	logger.DBLogger.Info("Product successfully added",
		zap.String("request_id", requestID),
		zap.String("product_id", product.Id))

	return nil
}

func (r *pvzRepository) GetLastProductInReception(ctx context.Context, receptionId string) (*models.Product, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetLastProductInReception called",
		zap.String("request_id", requestID),
		zap.String("reception_id", receptionId),
	)

	queryBuilder := sq.
		Select("id", "date_time", "type", "reception_id").
		From("products").
		Where(sq.Eq{"reception_id": receptionId}).
		OrderBy("date_time DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build SQL", zap.Error(err))
		return nil, err
	}

	var product models.Product
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&product.Id,
		&product.DateTime,
		&product.Type,
		&product.ReceptionId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.DBLogger.Info("no product found for reception",
				zap.String("request_id", requestID),
				zap.String("reception_id", receptionId),
			)
			return nil, nil
		}
		logger.DBLogger.Error("failed to scan product", zap.Error(err))
		return nil, err
	}

	return &product, nil
}

func (r *pvzRepository) DeleteProductById(ctx context.Context, productId string) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("DeleteProductByID called",
		zap.String("request_id", requestID),
		zap.String("product_id", productId),
	)

	queryBuilder := sq.
		Delete("products").
		Where(sq.Eq{"id": productId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build delete SQL", zap.Error(err))
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.DBLogger.Error("failed to execute delete", zap.Error(err))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found or already deleted")
	}

	logger.DBLogger.Info("Product deleted successfully",
		zap.String("request_id", requestID),
		zap.String("product_id", productId),
	)

	return nil
}

func (r *pvzRepository) CloseReception(ctx context.Context, reception *models.Reception) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CloseReception called",
		zap.String("request_id", requestID),
		zap.String("reception_id", reception.Id),
	)

	queryBuilder := sq.
		Update("receptions").
		Set("status", models.STATUS_CLOSED).
		Where(sq.Eq{"id": reception.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build update SQL", zap.Error(err))
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.DBLogger.Error("failed to update reception status", zap.Error(err))
		return err
	}

	reception.Status = "closed"

	logger.DBLogger.Info("Reception successfully closed",
		zap.String("request_id", requestID),
		zap.String("reception_id", reception.Id),
	)

	return nil
}
