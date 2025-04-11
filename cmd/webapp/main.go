package main

import (
	authController "avito_spring_staj_2025/internal/auth/controller"
	authRepository "avito_spring_staj_2025/internal/auth/repository"
	authUsecase "avito_spring_staj_2025/internal/auth/usecase"
	"avito_spring_staj_2025/internal/service/db"
	"avito_spring_staj_2025/internal/service/metrics"

	pvzController "avito_spring_staj_2025/internal/pvz/controller"
	pvzRepository "avito_spring_staj_2025/internal/pvz/repository"
	pvzUsecase "avito_spring_staj_2025/internal/pvz/usecase"
	"avito_spring_staj_2025/internal/service/jwt"
	"avito_spring_staj_2025/internal/service/logger"
	"avito_spring_staj_2025/internal/service/middleware"
	"avito_spring_staj_2025/internal/service/router"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load()
	db := db.DbConnect()
	jwtToken, err := jwt.NewJwtToken("secret-key")
	if err != nil {
		log.Fatalf("Failed to create JWT token: %v", err)
	}

	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			log.Fatalf("Failed to sync loggers: %v", err)
		}
	}()

	metrics.InitMetrics()

	authRepository := authRepository.NewAuthRepository(db)
	authUseCase := authUsecase.NewAuthUsecase(authRepository, jwtToken)
	authHandler := authController.NewAuthHandler(authUseCase)

	pvzRepository := pvzRepository.NewPvzRepository(db)
	pvzUseCase := pvzUsecase.NewPvzUsecase(pvzRepository)
	pvzHandler := pvzController.NewPvzHandler(pvzUseCase)

	mainRouter := router.SetUpRoutes(authHandler, pvzHandler, jwtToken)
	mainRouter.Use(middleware.RequestIDMiddleware)
	mainRouter.Use(middleware.RateLimitMiddleware)
	http.Handle("/", middleware.EnableCORS(mainRouter))
	fmt.Printf("Starting HTTP server on address %s\n", os.Getenv("BACKEND_URL"))
	if err := http.ListenAndServe(os.Getenv("BACKEND_URL"), nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}
