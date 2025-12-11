package auth

import (
	"errors"
	"golang-pinjaman-api/internal/delivery/middleware"
	"golang-pinjaman-api/internal/delivery/request"
	"golang-pinjaman-api/internal/usecase/auth"
	"golang-pinjaman-api/pkg/util"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase auth.AuthUsecase
	validator   *validator.Validate
}

func NewAuthHandler(e *echo.Echo, uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: uc,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req request.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Permintaan tidak valid")
	}

	if err := h.validator.Struct(req); err != nil {
		formattedErrors := util.TranslateValidationErrors(err)
		return util.FailResponse(c, http.StatusBadRequest, "Data yang diberikan tidak valid", formattedErrors)
	}

	user, err := h.authUsecase.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrUserExist) {
			return util.ErrorResponse(c, http.StatusConflict, "User sudah terdaftar")
		}
		return util.ErrorResponse(c, http.StatusInternalServerError, "Gagal mendaftarkan nasabah")
	}
	user.Password = ""
	return util.SuccessResponse(c, http.StatusCreated, "Pendaftaran nasabah berhasil. Silahkan login.", user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Permintaan tidak valid")
	}

	if err := h.validator.Struct(req); err != nil {
		formattedErrors := util.TranslateValidationErrors(err)
		return util.FailResponse(c, http.StatusBadRequest, "Data yang diberikan tidak valid", formattedErrors)
	}

	token, expiryTime, user, err := h.authUsecase.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return util.ErrorResponse(c, http.StatusUnauthorized, "Username atau password salah")
		}
		return util.ErrorResponse(c, http.StatusInternalServerError, "Gagal login")
	}

	return util.SuccessResponse(c, http.StatusOK, "Login berhasil", map[string]interface{}{
		"token":      token,
		"expired_at": expiryTime,
		"user":       user,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return util.FailResponse(c, http.StatusUnauthorized, "Header Authorization tidak ditemukan", nil)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return util.FailResponse(c, http.StatusUnauthorized, "Format Authorization tidak valid", nil)
	}

	tokenString := parts[1]

	err := h.authUsecase.Logout(c.Request().Context(), tokenString)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, "Gagal logout")
	}

	return util.SuccessResponse(c, http.StatusOK, "Logout berhasil dan token dihapus.", nil)
}

func (h *AuthHandler) GetUser(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return util.FailResponse(c, http.StatusUnauthorized, "User ID tidak valid", nil)
	}

	user, err := h.authUsecase.GetUser(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return util.ErrorResponse(c, http.StatusNotFound, "User tidak ditemukan")
		}
		return util.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data user")
	}

	user.Password = ""
	return util.SuccessResponse(c, http.StatusOK, "Berhasil mengambil data user", user)
}
