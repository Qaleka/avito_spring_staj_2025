package router

import (
	auth "avito_spring_staj_2025/internal/auth/controller"
	pvz "avito_spring_staj_2025/internal/pvz/controller"
	"avito_spring_staj_2025/internal/service/jwt"
	"avito_spring_staj_2025/internal/service/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func SetUpRoutes(authHandler *auth.AuthHandler, pvzHandler *pvz.PvzHandler, jwtService jwt.JwtTokenService) *mux.Router {
	router := mux.NewRouter()
	api := "/api"

	router.HandleFunc(api+"/dummyLogin", authHandler.DummyLogin).Methods("POST")
	router.HandleFunc(api+"/register", authHandler.Register).Methods("POST")
	router.HandleFunc(api+"/login", authHandler.Login).Methods("POST")
	router.Handle(api+"/pvz", middleware.RoleMiddleware(jwtService)(http.HandlerFunc(pvzHandler.CreatePvz))).Methods("POST")
	router.Handle(api+"/receptions", middleware.RoleMiddleware(jwtService)(http.HandlerFunc(pvzHandler.CreateReception))).Methods("POST")
	router.Handle(api+"/products", middleware.RoleMiddleware(jwtService)(http.HandlerFunc(pvzHandler.AddProductToReception))).Methods("POST")
	router.Handle(api+"/pvz/{pvzId}/delete_last_product", middleware.RoleMiddleware(jwtService)(http.HandlerFunc(pvzHandler.DeleteLastProduct))).Methods("POST")
	router.Handle(api+"/pvz/{pvzId}/close_last_reception", middleware.RoleMiddleware(jwtService)(http.HandlerFunc(pvzHandler.CloseLastReception))).Methods("POST")
	//router.HandleFunc(api+"/pvz/{pvzId}/close_last_reception", pvzHandler.CloseLastReception).Methods("POST")
	//router.HandleFunc(api+"/pvz/{pvzId}/delete_last_product", pvzHandler.DeleteLastProduct).Methods("POST")
	//router.HandleFunc(api+"/receptions", pvzHandler.CreateReceptions).Methods("POST")
	//router.HandleFunc(api+"/products", pvzHandler.AddProduct).Methods("POST")
	return router
}
