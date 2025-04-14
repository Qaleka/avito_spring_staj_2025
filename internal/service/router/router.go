package router

import (
	auth "avito_spring_staj_2025/internal/auth/handler"
	pvz "avito_spring_staj_2025/internal/pvz/handler"
	reception "avito_spring_staj_2025/internal/reception/handler"
	"avito_spring_staj_2025/internal/service/jwt"
	"avito_spring_staj_2025/internal/service/metrics"
	"avito_spring_staj_2025/internal/service/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func SetUpRoutes(authHandler *auth.AuthHandler, pvzHandler *pvz.PvzHandler, receptionHandler *reception.ReceptionHandler, jwtService jwt.JwtToken) *mux.Router {
	router := mux.NewRouter()
	api := "/api"

	withLogging := middleware.WithLoggingAndMetrics
	withAuth := middleware.RoleMiddleware(jwtService)

	withCreatedPvzMetric := middleware.WithCustomMetric(metrics.AmountOfCreatedPvz)
	withCreatedReceptionMetric := middleware.WithCustomMetric(metrics.AmountOfCreatedReceptions)
	withAddedProductMetric := middleware.WithCustomMetric(metrics.AmountOfAddedProducts)

	router.Handle(api+"/dummyLogin", middleware.ChainMiddlewares(http.HandlerFunc(authHandler.DummyLogin), withLogging)).Methods("POST")
	router.Handle(api+"/register", middleware.ChainMiddlewares(http.HandlerFunc(authHandler.Register), withLogging)).Methods("POST")
	router.Handle(api+"/login", middleware.ChainMiddlewares(http.HandlerFunc(authHandler.Login), withLogging)).Methods("POST")

	router.Handle(api+"/pvz", middleware.ChainMiddlewares(http.HandlerFunc(pvzHandler.CreatePvz), withLogging, withAuth, withCreatedPvzMetric)).Methods("POST")
	router.Handle(api+"/receptions", middleware.ChainMiddlewares(http.HandlerFunc(receptionHandler.CreateReception), withLogging, withAuth, withCreatedReceptionMetric)).Methods("POST")
	router.Handle(api+"/products", middleware.ChainMiddlewares(http.HandlerFunc(receptionHandler.AddProductToReception), withLogging, withAuth, withAddedProductMetric)).Methods("POST")
	router.Handle(api+"/pvz/{pvzId}/delete_last_product", middleware.ChainMiddlewares(http.HandlerFunc(receptionHandler.DeleteLastProduct), withLogging, withAuth)).Methods("POST")
	router.Handle(api+"/pvz/{pvzId}/close_last_reception", middleware.ChainMiddlewares(http.HandlerFunc(receptionHandler.CloseLastReception), withLogging, withAuth)).Methods("POST")
	router.Handle(api+"/pvz", middleware.ChainMiddlewares(http.HandlerFunc(pvzHandler.GetPvzsInformation), withLogging, withAuth)).Methods("GET")
	router.Handle(api+"/metrics", promhttp.Handler())

	router.HandleFunc(api+"/pvz/grpc", pvzHandler.GetPvzListFromGrpc).Methods("GET")
	return router
}
