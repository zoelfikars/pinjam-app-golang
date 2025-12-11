package response

import (
	"golang-pinjaman-api/internal/domain"
	"time"

	"github.com/gofrs/uuid/v5"
)

type NasabahResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	NamaLengkap string    `json:"nama_lengkap"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InspectorResponse struct {
	ID          uuid.UUID `json:"id"`
	NamaLengkap string    `json:"nama_lengkap"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PinjamanNasabahResponse struct {
	ID             uuid.UUID          `json:"id"`
	Nik            string             `json:"nik"`
	NamaLengkap    string             `json:"nama_lengkap"`
	Alamat         string             `json:"alamat"`
	NoTelepon      string             `json:"no_telepon"`
	JumlahPinjaman int64              `json:"jumlah_pinjaman"`
	Status         string             `json:"status"`
	CatatanAdmin   string             `json:"catatan_admin,omitempty"`
	InspectedAt    *time.Time         `json:"inspected_at,omitempty"`
	Inspector      *InspectorResponse `json:"inspector,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
type PinjamanAdminResponse struct {
	ID             uuid.UUID          `json:"id"`
	Nik            string             `json:"nik"`
	NamaLengkap    string             `json:"nama_lengkap"`
	Alamat         string             `json:"alamat"`
	NoTelepon      string             `json:"no_telepon"`
	JumlahPinjaman int64              `json:"jumlah_pinjaman"`
	Status         string             `json:"status"`
	CatatanAdmin   string             `json:"catatan_admin,omitempty"`
	InspectedAt    *time.Time         `json:"inspected_at,omitempty"`
	Inspector      *InspectorResponse `json:"inspector,omitempty"`
	Nasabah        *NasabahResponse   `json:"nasabah"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

func MapToAdminResponse(p *domain.PengajuanPinjaman) *PinjamanAdminResponse {
	resp := &PinjamanAdminResponse{
		ID:             p.ID,
		Nik:            p.Nik,
		NamaLengkap:    p.NamaLengkap,
		Alamat:         p.Alamat,
		NoTelepon:      p.NoTelepon,
		JumlahPinjaman: p.JumlahPinjaman,
		Status:         p.Status,
		CatatanAdmin:   p.CatatanAdmin,
		InspectedAt:    p.InspectedAt,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	if p.Nasabah != nil {
		resp.Nasabah = &NasabahResponse{
			ID:          p.Nasabah.ID,
			Username:    p.Nasabah.Username,
			NamaLengkap: p.Nasabah.NamaLengkap,
			Role:        p.Nasabah.Role,
			CreatedAt:   p.Nasabah.CreatedAt,
			UpdatedAt:   p.Nasabah.UpdatedAt,
		}
	}
	if p.Inspector != nil {
		resp.Inspector = &InspectorResponse{
			ID:          p.Inspector.ID,
			NamaLengkap: p.Inspector.NamaLengkap,
			Role:        p.Inspector.Role,
			CreatedAt:   p.Inspector.CreatedAt,
			UpdatedAt:   p.Inspector.UpdatedAt,
		}
	}
	return resp
}
func MapToNasabahResponse(p *domain.PengajuanPinjaman) *PinjamanNasabahResponse {
	resp := &PinjamanNasabahResponse{
		ID:             p.ID,
		Nik:            p.Nik,
		NamaLengkap:    p.NamaLengkap,
		Alamat:         p.Alamat,
		NoTelepon:      p.NoTelepon,
		JumlahPinjaman: p.JumlahPinjaman,
		Status:         p.Status,
		CatatanAdmin:   p.CatatanAdmin,
		InspectedAt:    p.InspectedAt,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	if p.Inspector != nil {
		resp.Inspector = &InspectorResponse{
			ID:          p.Inspector.ID,
			NamaLengkap: p.Inspector.NamaLengkap,
			Role:        p.Inspector.Role,
			CreatedAt:   p.Inspector.CreatedAt,
			UpdatedAt:   p.Inspector.UpdatedAt,
		}
	}
	return resp
}
