package repository

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("AuthUser called", zap.String("request_id", requestID), zap.String("username", user.Email))
	queryBuilder := sq.Insert("users").
		Columns("id", "email", "hash_password", "role").
		Values(user.Id, user.Email, user.Password, user.Role).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build SQL", zap.Error(err))
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.DBLogger.Error("failed to insert user", zap.Error(err))
		return err
	}

	logger.DBLogger.Info("User successfully registered",
		zap.String("request_id", requestID),
		zap.String("user_id", user.Id),
	)

	return nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("AuthUser called", zap.String("request_id", requestID), zap.String("email", email))
	queryBuilder := sq.Select("id", "email", "hash_password", "role").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.DBLogger.Error("failed to build SQL", zap.Error(err))
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, query, args...)

	var user models.User
	err = row.Scan(&user.Id, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.DBLogger.Info("user not found",
				zap.String("request_id", requestID),
				zap.String("email", email),
			)
			return nil, errors.New("user not found")
		}
		logger.DBLogger.Error("failed to scan user", zap.Error(err))
		return nil, err
	}

	return &user, nil
}
