package router

import (
	"github.com/labstack/echo/v4"
	"golang-pinjaman-api/internal/delivery/middleware"
	pinjamanHandler "golang-pinjaman-api/internal/delivery/pinjaman"
	"golang-pinjaman-api/internal/repository/redis"
	pinjamanUseCase "golang-pinjaman-api/internal/usecase/pinjaman"
)

func PinjamanRouterGroup(e *echo.Echo, uc pinjamanUseCase.PinjamanUsecase, tokenRepo redis.TokenRepository, jwtSecret string) {
	handler := pinjamanHandler.NewPinjamanHandler(uc)
	r := e.Group("/api/pinjaman")
	r.Use(middleware.JWTMiddleware(jwtSecret, tokenRepo))
	r.POST("/ajukan", handler.AjukanPinjaman, middleware.RoleMiddleware("nasabah"))
	r.GET("/nasabah", handler.GetMyPinjamanList, middleware.RoleMiddleware("nasabah"))
	r.GET("/list", handler.GetAllPinjamanList, middleware.RoleMiddleware("admin"))
	r.PUT("/:id/approval", handler.ApproveOrRejectPinjaman, middleware.RoleMiddleware("admin"))
}
