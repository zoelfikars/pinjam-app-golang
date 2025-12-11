package request

type AjukanPinjamanRequest struct {
	Nik            string `json:"nik" validate:"required,min=16"`
	Alamat         string `json:"alamat" validate:"required"`
	NoTelepon      string `json:"no_telepon" validate:"required,min=12,phone_id"`
	JumlahPinjaman int64  `json:"jumlah_pinjaman" validate:"required,gte=1000000"`
}
