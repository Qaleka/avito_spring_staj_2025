package integration

import (
	"avito_spring_staj_2025/domain/models"
	"avito_spring_staj_2025/domain/requests"
	pvzController "avito_spring_staj_2025/internal/pvz/handler"
	pvzRepository "avito_spring_staj_2025/internal/pvz/repository"
	pvzUsecase "avito_spring_staj_2025/internal/pvz/usecase"
	"avito_spring_staj_2025/internal/service/logger"
	"go.uber.org/zap"

	receptionController "avito_spring_staj_2025/internal/reception/handler"
	receptionRepository "avito_spring_staj_2025/internal/reception/repository"
	receptionUsecase "avito_spring_staj_2025/internal/reception/usecase"
	"avito_spring_staj_2025/internal/service/middleware"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFullPVZWorkflow(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()
	logger.AccessLogger = zap.NewNop()
	logger.DBLogger = zap.NewNop()
	pvzRepo := pvzRepository.NewPvzRepository(db)
	receptionRepo := receptionRepository.NewReceptionRepository(db)

	pvzUsecase := pvzUsecase.NewPvzUsecase(pvzRepo)
	receptionUsecase := receptionUsecase.NewReceptionUsecase(receptionRepo)

	pvzHandler := pvzController.NewPvzHandler(pvzUsecase)
	receptionHandler := receptionController.NewReceptionHandler(receptionUsecase)

	ctx := context.WithValue(context.Background(), middleware.ContextKeyRole, "employee")

	t.Run("Create PVZ - moderator", func(t *testing.T) {

		mock.ExpectExec(`^INSERT INTO pvzs \(id,registration_date,city\) VALUES \(\$1,\$2,\$3\)$`).
			WithArgs("pvz-1", sqlmock.AnyArg(), "Москва").
			WillReturnResult(sqlmock.NewResult(1, 1))

		req := httptest.NewRequest("POST", "/pvz", mockJSONBody(t, requests.CreatePvzRequest{
			Id:               "pvz-1",
			RegistrationDate: time.Now(),
			City:             "Москва",
		}))

		ctx := context.WithValue(req.Context(), middleware.ContextKeyRole, "moderator")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		pvzHandler.CreatePvz(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create Reception", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM receptions").
			WithArgs("pvz-1", models.STATUS_ACTIVE).
			WillReturnError(sql.ErrNoRows) // Нет активной приёмки

		mock.ExpectQuery("SELECT.*FROM pvzs").
			WithArgs("pvz-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
				AddRow("pvz-1", time.Now(), "Москва"))

		mock.ExpectExec("INSERT INTO receptions").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "pvz-1", "in_progress").
			WillReturnResult(sqlmock.NewResult(1, 1))

		req := httptest.NewRequest("POST", "/receptions", nil).WithContext(ctx)
		req.Body = mockJSONBody(t, requests.CreateReceptionRequest{
			PvzId: "pvz-1",
		})

		rr := httptest.NewRecorder()
		receptionHandler.CreateReception(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("Add 50 Products", func(t *testing.T) {
		for i := 0; i < 50; i++ {
			mock.ExpectQuery(`^SELECT id, registration_date, city FROM pvzs WHERE id = \$1$`).
				WithArgs("pvz-1").
				WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow("pvz-1", time.Now(), "Москва"))

			mock.ExpectQuery(`^SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = \$1 AND status = \$2$`).
				WithArgs("pvz-1", models.STATUS_ACTIVE).
				WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
					AddRow("rec-1", time.Now(), "pvz-1", "in_progress"))

			mock.ExpectExec(`^INSERT INTO products \(id,date_time,type,reception_id\) VALUES \(\$1,\$2,\$3,\$4\)$`).
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "одежда", "rec-1").
				WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
		}

		for i := 0; i < 50; i++ {
			req := httptest.NewRequest("POST", "/products", mockJSONBody(t, requests.AddProductRequest{
				Type:  "одежда",
				PvzId: "pvz-1",
			})).WithContext(ctx)

			rr := httptest.NewRecorder()
			receptionHandler.AddProductToReception(rr, req)

			assert.Equal(t, http.StatusCreated, rr.Code)
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Close Reception", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM pvzs").
			WithArgs("pvz-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
				AddRow("pvz-1", time.Now(), "Москва"))

		mock.ExpectQuery("SELECT.*FROM receptions").
			WithArgs("pvz-1", models.STATUS_ACTIVE).
			WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
				AddRow("rec-1", time.Now(), "pvz-1", "in_progress"))

		mock.ExpectExec("UPDATE receptions").
			WithArgs(models.STATUS_CLOSED, "rec-1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		req := httptest.NewRequest("PUT", "/pvz/pvz-1/close_last_reception", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"pvzId": "pvz-1"})

		rr := httptest.NewRecorder()
		receptionHandler.CloseLastReception(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func mockJSONBody(t *testing.T, data interface{}) *mockBody {
	body, err := json.Marshal(data)
	require.NoError(t, err)
	return &mockBody{data: body}
}

type mockBody struct {
	data []byte
	read int
}

func (m *mockBody) Read(p []byte) (n int, err error) {
	if m.read >= len(m.data) {
		return 0, io.EOF
	}
	n = copy(p, m.data[m.read:])
	m.read += n
	return n, nil
}

func (m *mockBody) Close() error {
	return nil
}
