package unit

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/pvz/repository"
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
	"time"
)

func TestPvzRepository_CreatePvz(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	tests := []struct {
		name        string
		data        *requests.CreatePvzRequest
		mock        func(sqlmock.Sqlmock)
		expectedErr string
	}{
		{
			name: "Success",
			data: &requests.CreatePvzRequest{
				Id:               "pvz123",
				RegistrationDate: time.Now(),
				City:             "Moscow",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO pvzs \(id,registration_date,city\) VALUES \(\$1,\$2,\$3\)`).
					WithArgs("pvz123", sqlmock.AnyArg(), "Moscow").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: "",
		},
		{
			name: "Insert Error",
			data: &requests.CreatePvzRequest{
				Id:               "pvz456",
				RegistrationDate: time.Now(),
				City:             "Spb",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO pvzs \(id,registration_date,city\) VALUES \(\$1,\$2,\$3\)`).
					WithArgs("pvz456", sqlmock.AnyArg(), "Spb").
					WillReturnError(errors.New("insert failed"))
			},
			expectedErr: "insert failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			err = repo.CreatePvz(ctx, tt.data)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_CreateReception(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	tests := []struct {
		name        string
		data        models.Reception
		mock        func(sqlmock.Sqlmock)
		expectedErr string
	}{
		{
			name: "Success",
			data: models.Reception{
				Id:       "r1",
				DateTime: time.Now(),
				PvzId:    "pvz1",
				Status:   "ACTIVE",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO receptions \(id,date_time,pvz_id,status\) VALUES \(\$1,\$2,\$3,\$4\)`).
					WithArgs("r1", sqlmock.AnyArg(), "pvz1", "ACTIVE").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: "",
		},
		{
			name: "Insert Error",
			data: models.Reception{
				Id:       "r2",
				DateTime: time.Now(),
				PvzId:    "pvz2",
				Status:   "ACTIVE",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO receptions \(id,date_time,pvz_id,status\) VALUES \(\$1,\$2,\$3,\$4\)`).
					WithArgs("r2", sqlmock.AnyArg(), "pvz2", "ACTIVE").
					WillReturnError(errors.New("insert failed"))
			},
			expectedErr: "insert failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			err = repo.CreateReception(ctx, tt.data)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_GetPvzById(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	tests := []struct {
		name        string
		pvzId       string
		mock        func(sqlmock.Sqlmock)
		expected    *models.Pvz
		expectedErr string
	}{
		{
			name:  "Success",
			pvzId: "pvz1",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, registration_date, city FROM pvzs WHERE id = \$1`).
					WithArgs("pvz1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
						AddRow("pvz1", time.Now(), "Moscow"))
			},
			expected: &models.Pvz{
				Id:   "pvz1",
				City: "Moscow",
			},
			expectedErr: "",
		},
		{
			name:  "Not Found",
			pvzId: "pvz2",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, registration_date, city FROM pvzs WHERE id = \$1`).
					WithArgs("pvz2").
					WillReturnError(sql.ErrNoRows)
			},
			expected:    nil,
			expectedErr: "pvz not found",
		},
		{
			name:  "Query Error",
			pvzId: "pvz3",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, registration_date, city FROM pvzs WHERE id = \$1`).
					WithArgs("pvz3").
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

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			result, err := repo.GetPvzById(ctx, tt.pvzId)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Id, result.Id)
				assert.Equal(t, tt.expected.City, result.City)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_GetCurrentReception(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	tests := []struct {
		name        string
		pvzId       string
		mock        func(sqlmock.Sqlmock)
		expected    *models.Reception
		expectedErr string
	}{
		{
			name:  "Success",
			pvzId: "pvz1",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = \$1 AND status = \$2`).
					WithArgs("pvz1", models.STATUS_ACTIVE).
					WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
						AddRow("r1", time.Now(), "pvz1", models.STATUS_ACTIVE))
			},
			expected: &models.Reception{
				Id:     "r1",
				PvzId:  "pvz1",
				Status: models.STATUS_ACTIVE,
			},
			expectedErr: "",
		},
		{
			name:  "No Active Reception",
			pvzId: "pvz2",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = \$1 AND status = \$2`).
					WithArgs("pvz2", models.STATUS_ACTIVE).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    nil,
			expectedErr: "no active reception",
		},
		{
			name:  "Query Error",
			pvzId: "pvz3",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = \$1 AND status = \$2`).
					WithArgs("pvz3", models.STATUS_ACTIVE).
					WillReturnError(errors.New("no active reception"))
			},
			expected:    nil,
			expectedErr: "no active reception",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			result, err := repo.GetCurrentReception(ctx, tt.pvzId)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Id, result.Id)
				assert.Equal(t, tt.expected.PvzId, result.PvzId)
				assert.Equal(t, tt.expected.Status, result.Status)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_AddProductToReception(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name        string
		product     models.Product
		mock        func(sqlmock.Sqlmock)
		expectedErr string
	}{
		{
			name: "Success",
			product: models.Product{
				Id:          "prod1",
				DateTime:    time.Now(),
				Type:        "type1",
				ReceptionId: "rec1",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO products \(id,date_time,type,reception_id\) VALUES \(\$1,\$2,\$3,\$4\)`).
					WithArgs("prod1", sqlmock.AnyArg(), "type1", "rec1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: "",
		},
		{
			name: "Database Error",
			product: models.Product{
				Id:          "prod2",
				DateTime:    time.Now(),
				Type:        "type2",
				ReceptionId: "rec2",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO products \(id,date_time,type,reception_id\) VALUES \(\$1,\$2,\$3,\$4\)`).
					WithArgs("prod2", sqlmock.AnyArg(), "type2", "rec2").
					WillReturnError(errors.New("database error"))
			},
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			err = repo.AddProductToReception(ctx, tt.product)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_GetLastProductInReception(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name        string
		receptionId string
		mock        func(sqlmock.Sqlmock)
		expected    *models.Product
		expectedErr string
	}{
		{
			name:        "Success",
			receptionId: "rec1",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"}).
					AddRow("prod1", time.Now(), "type1", "rec1")
				mock.ExpectQuery(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = \$1 ORDER BY date_time DESC LIMIT 1`).
					WithArgs("rec1").
					WillReturnRows(rows)
			},
			expected: &models.Product{
				Id:          "prod1",
				Type:        "type1",
				ReceptionId: "rec1",
			},
			expectedErr: "",
		},
		{
			name:        "No Products",
			receptionId: "rec2",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = \$1 ORDER BY date_time DESC LIMIT 1`).
					WithArgs("rec2").
					WillReturnError(sql.ErrNoRows)
			},
			expected:    nil,
			expectedErr: "no products in reception",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			product, err := repo.GetLastProductInReception(ctx, tt.receptionId)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, product)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Id, product.Id)
				assert.Equal(t, tt.expected.Type, product.Type)
				assert.Equal(t, tt.expected.ReceptionId, product.ReceptionId)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_DeleteProductById(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name      string
		productId string
		mock      func(sqlmock.Sqlmock)
		expectErr bool
		errMsg    string
	}{
		{
			name:      "Success",
			productId: "prod1",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
					WithArgs("prod1").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectErr: false,
		},
		{
			name:      "Not Found",
			productId: "prod2",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
					WithArgs("prod2").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectErr: true,
			errMsg:    "product not found or already deleted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			err = repo.DeleteProductById(ctx, tt.productId)

			if tt.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_CloseReception(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name       string
		reception  *models.Reception
		mock       func(sqlmock.Sqlmock)
		expectErr  bool
		errMsg     string
		expectStat string
	}{
		{
			name: "Success",
			reception: &models.Reception{
				Id:     "rec1",
				Status: models.STATUS_ACTIVE,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE receptions SET status = \$1 WHERE id = \$2`).
					WithArgs(models.STATUS_CLOSED, "rec1").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectErr:  false,
			expectStat: models.STATUS_CLOSED,
		},
		{
			name: "Update Error",
			reception: &models.Reception{
				Id:     "rec2",
				Status: models.STATUS_ACTIVE,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE receptions SET status = \$1 WHERE id = \$2`).
					WithArgs(models.STATUS_CLOSED, "rec2").
					WillReturnError(errors.New("update failed"))
			},
			expectErr: true,
			errMsg:    "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			err = repo.CloseReception(ctx, tt.reception)

			if tt.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectStat, tt.reception.Status)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_GetPvzsFilteredByReceptionDate(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.Background()

	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	tests := []struct {
		name     string
		from     time.Time
		to       time.Time
		limit    int
		offset   int
		mock     func(sqlmock.Sqlmock)
		expected []models.Pvz
	}{
		{
			name:   "With Date Range",
			from:   from,
			to:     to,
			limit:  10,
			offset: 0,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow("pvz1", now, "Moscow").
					AddRow("pvz2", now, "SPb")
				mock.ExpectQuery(`SELECT id, registration_date, city FROM pvzs WHERE \(registration_date >= \$1 AND registration_date <= \$2\) LIMIT 10 OFFSET 0`).
					WithArgs(from, to).
					WillReturnRows(rows)
			},
			expected: []models.Pvz{
				{Id: "pvz1", City: "Moscow"},
				{Id: "pvz2", City: "SPb"},
			},
		},
		{
			name:   "Without Date Range",
			from:   time.Time{},
			to:     time.Time{},
			limit:  5,
			offset: 0,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"})
				mock.ExpectQuery(`SELECT id, registration_date, city FROM pvzs LIMIT 5 OFFSET 0`).
					WillReturnRows(rows)
			},
			expected: []models.Pvz{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			result, err := repo.GetPvzsFilteredByReceptionDate(ctx, tt.from, tt.to, tt.limit, tt.offset)
			require.NoError(t, err)

			assert.Equal(t, len(tt.expected), len(result))
			for i := range tt.expected {
				assert.Equal(t, tt.expected[i].Id, result[i].Id)
				assert.Equal(t, tt.expected[i].City, result[i].City)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_GetPvzReceptions(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name        string
		pvzId       string
		mock        func(sqlmock.Sqlmock)
		expected    []models.Reception
		expectedErr string
	}{
		{
			name:  "Success",
			pvzId: "pvz1",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
					AddRow("rec1", time.Now(), "pvz1", "ACTIVE").
					AddRow("rec2", time.Now(), "pvz1", "CLOSED")
				mock.ExpectQuery(`SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = \$1`).
					WithArgs("pvz1").
					WillReturnRows(rows)
			},
			expected: []models.Reception{
				{Id: "rec1", PvzId: "pvz1", Status: "ACTIVE"},
				{Id: "rec2", PvzId: "pvz1", Status: "CLOSED"},
			},
			expectedErr: "",
		},
		{
			name:  "No Receptions",
			pvzId: "pvz2",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = \$1`).
					WithArgs("pvz2").
					WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}))
			},
			expected:    []models.Reception{},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			result, err := repo.GetPvzReceptions(ctx, tt.pvzId)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(result))
				for i := range tt.expected {
					assert.Equal(t, tt.expected[i].Id, result[i].Id)
					assert.Equal(t, tt.expected[i].PvzId, result[i].PvzId)
					assert.Equal(t, tt.expected[i].Status, result[i].Status)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepository_GetReceptionProducts(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name        string
		receptionId string
		mock        func(sqlmock.Sqlmock)
		expected    []models.Product
		expectedErr string
	}{
		{
			name:        "Success",
			receptionId: "rec1",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"}).
					AddRow("prod1", time.Now(), "type1", "rec1").
					AddRow("prod2", time.Now(), "type2", "rec1")
				mock.ExpectQuery(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = \$1`).
					WithArgs("rec1").
					WillReturnRows(rows)
			},
			expected: []models.Product{
				{Id: "prod1", Type: "type1", ReceptionId: "rec1"},
				{Id: "prod2", Type: "type2", ReceptionId: "rec1"},
			},
			expectedErr: "",
		},
		{
			name:        "No Products",
			receptionId: "rec2",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, date_time, type, reception_id FROM products WHERE reception_id = \$1`).
					WithArgs("rec2").
					WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"}))
			},
			expected:    []models.Product{},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewPvzRepository(db)
			tt.mock(mock)

			result, err := repo.GetReceptionProducts(ctx, tt.receptionId)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(result))
				for i := range tt.expected {
					assert.Equal(t, tt.expected[i].Id, result[i].Id)
					assert.Equal(t, tt.expected[i].Type, result[i].Type)
					assert.Equal(t, tt.expected[i].ReceptionId, result[i].ReceptionId)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
