package router

import (
	authHandler "golang-pinjaman-api/internal/delivery/auth"
	"golang-pinjaman-api/internal/delivery/middleware"
	"golang-pinjaman-api/internal/repository/redis"
	authUseCase "golang-pinjaman-api/internal/usecase/auth"

	"github.com/labstack/echo/v4"
)

func AuthRouterGroup(e *echo.Echo, uc authUseCase.AuthUsecase, tokenRepo redis.TokenRepository, jwtSecret string) {
	authHandler := authHandler.NewAuthHandler(e, uc)

	g := e.Group("/api/auth")

	g.POST("/register", authHandler.Register)
	g.POST("/login", authHandler.Login)
	g.POST("/logout", authHandler.Logout, middleware.JWTMiddleware(jwtSecret, tokenRepo))
	g.GET("/user", authHandler.GetUser, middleware.JWTMiddleware(jwtSecret, tokenRepo))
}
