package request

type StatusUpdateRequest struct {
	Status       string `json:"status" validate:"required,oneof=approved rejected"`
	CatatanAdmin string `json:"catatan_admin" validate:"required_if=Status rejected"`
}
