package pinjaman

import (
	"errors"
	"fmt"
	"golang-pinjaman-api/internal/delivery/middleware"
	"golang-pinjaman-api/internal/delivery/request"
	"golang-pinjaman-api/internal/delivery/response"
	"golang-pinjaman-api/internal/usecase/pinjaman"
	"golang-pinjaman-api/pkg/util"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type PinjamanHandler struct {
	pinjamanUsecase pinjaman.PinjamanUsecase
	validator       *validator.Validate
}

func NewPinjamanHandler(uc pinjaman.PinjamanUsecase) *PinjamanHandler {
	v := validator.New()
	util.RegisterCustomValidators(v)
	return &PinjamanHandler{
		pinjamanUsecase: uc,
		validator:       v,
	}
}
func (h *PinjamanHandler) AjukanPinjaman(c echo.Context) error {
	var req request.AjukanPinjamanRequest
	if err := c.Bind(&req); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Permintaan tidak valid")
	}
	if err := h.validator.Struct(req); err != nil {
		formattedErrors := util.TranslateValidationErrors(err)
		return util.FailResponse(c, http.StatusBadRequest, "Data yang diberikan tidak valid", formattedErrors)
	}
	nasabahID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan di context")
	}
	pinjaman, err := h.pinjamanUsecase.AjukanPinjaman(c.Request().Context(), nasabahID, req)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengajukan pinjaman")
	}
	return util.SuccessResponse(c, http.StatusCreated, "Pinjaman berhasil diajukan", map[string]interface{}{
		"id":              pinjaman.ID,
		"status":          pinjaman.Status,
		"jumlah_pinjaman": pinjaman.JumlahPinjaman,
	})
}
func (h *PinjamanHandler) GetMyPinjamanList(c echo.Context) error {
	nasabahID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan di context")
	}
	list, err := h.pinjamanUsecase.GetMyPinjamanList(c.Request().Context(), nasabahID)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil list pinjaman")
	}
	var responses []response.PinjamanNasabahResponse
	for _, p := range list {
		responses = append(responses, *response.MapToNasabahResponse(&p))
	}
	return util.SuccessResponse(c, http.StatusOK, "List pinjaman berhasil diambil", responses)
}
func (h *PinjamanHandler) GetAllPinjamanList(c echo.Context) error {
	list, err := h.pinjamanUsecase.GetAllPinjamanList(c.Request().Context())
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Gagal mengambil list pinjaman e : %s", err.Error()))
	}
	var responses []response.PinjamanAdminResponse
	for _, p := range list {
		responses = append(responses, *response.MapToAdminResponse(&p))
	}
	return util.SuccessResponse(c, http.StatusOK, "List pinjaman berhasil diambil", responses)
}
func (h *PinjamanHandler) ApproveOrRejectPinjaman(c echo.Context) error {
	pinjamanIDStr := c.Param("id")
	pinjamanID, err := uuid.FromString(pinjamanIDStr)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Format ID pinjaman tidak valid")
	}
	var req request.StatusUpdateRequest
	if err := c.Bind(&req); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Permintaan tidak valid")
	}
	if err := h.validator.Struct(req); err != nil {
		formattedErrors := util.TranslateValidationErrors(err)
		return util.FailResponse(c, http.StatusBadRequest, "Data yang diberikan tidak valid", formattedErrors)
	}
	adminID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return util.ErrorResponse(c, http.StatusUnauthorized, "Admin ID tidak ditemukan di context")
	}
	err = h.pinjamanUsecase.ApproveOrRejectPinjaman(c.Request().Context(), pinjamanID, adminID, req)
	if err != nil {
		if errors.Is(err, pinjaman.ErrPinjamanNotFound) {
			return util.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, pinjaman.ErrAlreadyProcessed) {
			return util.ErrorResponse(c, http.StatusConflict, err.Error())
		}
		return util.ErrorResponse(c, http.StatusInternalServerError, "Gagal memproses status pinjaman")
	}
	return util.SuccessResponse(c, http.StatusOK, "Status pinjaman berhasil diperbarui", nil)
}
