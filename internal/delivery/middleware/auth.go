package middleware

import (
	"golang-pinjaman-api/internal/repository/redis"
	"golang-pinjaman-api/pkg/util"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

func JWTMiddleware(jwtSecret string, tokenRepo redis.TokenRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return util.FailResponse(c, http.StatusUnauthorized, "Header Authorization tidak ditemukan", nil)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return util.FailResponse(c, http.StatusUnauthorized, "Format Authorization tidak valid", nil)
			}

			tokenString := parts[1]

			isBlacklisted, err := tokenRepo.IsTokenBlacklisted(c.Request().Context(), tokenString)
			if err != nil {
				return util.ErrorResponse(c, http.StatusInternalServerError, "Server error ketika memeriksa token")
			}
			if isBlacklisted {
				return util.FailResponse(c, http.StatusUnauthorized, "Token tidak valid", nil)
			}

			claims, err := util.ExtractClaims(tokenString, jwtSecret)
			if err != nil {
				return util.FailResponse(c, http.StatusUnauthorized, "Token tidak valid", nil)
			}

			c.Set(UserIDKey, claims.UserID)
			c.Set(RoleKey, claims.Role)

			return next(c)
		}
	}
}

func RoleMiddleware(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get(RoleKey).(string)
			if !ok || userRole != requiredRole {
				return util.FailResponse(c, http.StatusForbidden, "Akses dilarang: izin peran tidak mencukupi", nil)
			}
			return next(c)
		}
	}
}
