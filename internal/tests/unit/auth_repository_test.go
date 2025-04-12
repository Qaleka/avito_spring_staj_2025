package unit

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/internal/auth/repository"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestAuthRepository_CreateUser(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	tests := []struct {
		name        string
		user        *models.User
		mock        func(sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name: "Success",
			user: &models.User{
				Id:       "123",
				Email:    "test@example.com",
				Password: "hashed_password",
				Role:     "user",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(id,email,hash_password,role\) VALUES \(\$1,\$2,\$3,\$4\)`).
					WithArgs("123", "test@example.com", "hashed_password", "user").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Database Error",
			user: &models.User{
				Id:       "123",
				Email:    "test@example.com",
				Password: "hashed_password",
				Role:     "user",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(id,email,hash_password,role\) VALUES \(\$1,\$2,\$3,\$4\)`).
					WithArgs("123", "test@example.com", "hashed_password", "user").
					WillReturnError(errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewAuthRepository(db)
			tt.mock(mock)

			err = repo.CreateUser(ctx, tt.user)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthRepository_GetUserByEmail(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	tests := []struct {
		name        string
		email       string
		mock        func(sqlmock.Sqlmock)
		expected    *models.User
		expectedErr string
	}{
		{
			name:  "Success",
			email: "test@example.com",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, email, hash_password, role FROM users WHERE email = \$1`).
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "hash_password", "role"}).
						AddRow("123", "test@example.com", "hashed_password", "user"))
			},
			expected: &models.User{
				Id:       "123",
				Email:    "test@example.com",
				Password: "hashed_password",
				Role:     "user",
			},
			expectedErr: "",
		},
		{
			name:  "User Not Found",
			email: "missing@example.com",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, email, hash_password, role FROM users WHERE email = \$1`).
					WithArgs("missing@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			expected:    nil,
			expectedErr: "user not found",
		},
		{
			name:  "Query Error",
			email: "error@example.com",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, email, hash_password, role FROM users WHERE email = \$1`).
					WithArgs("error@example.com").
					WillReturnError(errors.New("query failed"))
			},
			expected:    nil,
			expectedErr: "query failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewAuthRepository(db)
			tt.mock(mock)

			user, err := repo.GetUserByEmail(ctx, tt.email)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
