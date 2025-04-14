package main

import (
	authController "avito_spring_staj_2025/internal/auth/handler"
	authRepository "avito_spring_staj_2025/internal/auth/repository"
	authUsecase "avito_spring_staj_2025/internal/auth/usecase"
	"avito_spring_staj_2025/internal/db"
	"avito_spring_staj_2025/internal/pvz/handler/gen"
	"avito_spring_staj_2025/internal/service/metrics"
	"google.golang.org/grpc"
	"net"

	pvzController "avito_spring_staj_2025/internal/pvz/handler"
	pvzRepository "avito_spring_staj_2025/internal/pvz/repository"
	pvzUsecase "avito_spring_staj_2025/internal/pvz/usecase"

	receptionController "avito_spring_staj_2025/internal/reception/handler"
	receptionRepository "avito_spring_staj_2025/internal/reception/repository"
	receptionUsecase "avito_spring_staj_2025/internal/reception/usecase"
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
	jwtToken, err := jwt.NewJwtToken(os.Getenv("SECRET_KEY"))
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

	receptionRepository := receptionRepository.NewReceptionRepository(db)
	receptionUseCase := receptionUsecase.NewReceptionUsecase(receptionRepository)
	receptionHandler := receptionController.NewReceptionHandler(receptionUseCase)

	go func() {
		lis, err := net.Listen("tcp", os.Getenv("GRPC_URL"))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		gen.RegisterPVZServiceServer(grpcServer, pvzController.NewPvzGrpcHandler(pvzUseCase))

		fmt.Printf("grpc server listening at %s\n", os.Getenv("GRPC_URL"))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	mainRouter := router.SetUpRoutes(authHandler, pvzHandler, receptionHandler, jwtToken)
	mainRouter.Use(middleware.RequestIDMiddleware)
	mainRouter.Use(middleware.RateLimitMiddleware)
	http.Handle("/", middleware.EnableCORS(mainRouter))
	fmt.Printf("Starting HTTP server on address %s\n", os.Getenv("BACKEND_URL"))
	if err := http.ListenAndServe(os.Getenv("BACKEND_URL"), nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}
