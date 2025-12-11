package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"golang-pinjaman-api/config"
	"golang-pinjaman-api/database/seeder"
	"golang-pinjaman-api/internal/app"
	"golang-pinjaman-api/internal/delivery/router"
	authRepository "golang-pinjaman-api/internal/repository/auth"
	pinjamanRepository "golang-pinjaman-api/internal/repository/pinjaman"
	redisRepository "golang-pinjaman-api/internal/repository/redis"
	authUsecase "golang-pinjaman-api/internal/usecase/auth"
	pinjamanUsecase "golang-pinjaman-api/internal/usecase/pinjaman"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := app.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	seeder.Seed(db)

	redisClient := app.NewRedisClient(cfg)
	tokenRepo := redisRepository.NewTokenRepository(redisClient)

	authRepo := authRepository.NewUserRepository(db)
	pinjamanRepo := pinjamanRepository.NewPinjamanRepository(db)

	authUC := authUsecase.NewAuthUsecase(authRepo, tokenRepo, cfg.JWTSecret)
	pinjamanUC := pinjamanUsecase.NewPinjamanUsecase(pinjamanRepo, authRepo)

	e := echo.New()

	router.AuthRouterGroup(e, authUC, tokenRepo, cfg.JWTSecret)
	router.PinjamanRouterGroup(e, pinjamanUC, tokenRepo, cfg.JWTSecret)

	log.Printf("Server running on port %s", cfg.ServerPort)
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}
