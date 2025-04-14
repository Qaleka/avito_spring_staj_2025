package repository

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
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

			repo := NewPvzRepository(db)
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

			repo := NewPvzRepository(db)
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

			repo := NewPvzRepository(db)
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

			repo := NewPvzRepository(db)
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

func TestPvzRepository_GetAllPvzs(t *testing.T) {
	logger.DBLogger = zap.NewNop()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	now := time.Now()
	testCases := []struct {
		name          string
		mock          func(sqlmock.Sqlmock)
		expectedPvzs  []models.Pvz
		expectedError string
	}{
		{
			name: "Success - multiple rows",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow("1", now, "Moscow").
					AddRow("2", now.Add(time.Hour), "London").
					AddRow("3", now.Add(2*time.Hour), "Paris")
				mock.ExpectQuery("SELECT id, registration_date, city FROM pvzs").
					WillReturnRows(rows)
			},
			expectedPvzs: []models.Pvz{
				{Id: "1", RegistrationDate: now, City: "Moscow"},
				{Id: "2", RegistrationDate: now.Add(time.Hour), City: "London"},
				{Id: "3", RegistrationDate: now.Add(2 * time.Hour), City: "Paris"},
			},
		},
		{
			name: "Success - empty result",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"})
				mock.ExpectQuery("SELECT id, registration_date, city FROM pvzs").
					WillReturnRows(rows)
			},
			expectedPvzs: []models.Pvz(nil),
		},
		{
			name: "Query error",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, registration_date, city FROM pvzs").
					WillReturnError(errors.New("database error"))
			},
			expectedError: "failed to query pvzs: database error",
		},
		{
			name: "Scan error",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow("1", now, "Moscow").
					AddRow(nil, now.Add(time.Hour), "London") // invalid id
				mock.ExpectQuery("SELECT id, registration_date, city FROM pvzs").
					WillReturnRows(rows)
			},
			expectedError: "failed to scan row: sql: Scan error on column index 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewPvzRepository(db)
			tc.mock(mock)

			pvzs, err := repo.GetAllPvzs(ctx)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedPvzs, pvzs)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
