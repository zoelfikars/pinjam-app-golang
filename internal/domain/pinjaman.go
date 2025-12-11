package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type PengajuanPinjaman struct {
	BaseModel
	NasabahID      uuid.UUID     `gorm:"type:uuid;not null" json:"nasabah_id"`
	Nik            string        `gorm:"type:varchar(20);not null" json:"nik"`
	NamaLengkap    string        `gorm:"type:varchar(255);not null" json:"nama_lengkap"`
	Alamat         string        `gorm:"type:text;not null" json:"alamat"`
	NoTelepon      string        `gorm:"type:varchar(20);not null" json:"no_telepon"`
	JumlahPinjaman int64         `gorm:"type:bigint;not null" json:"jumlah_pinjaman"`
	Status         string        `gorm:"type:status_type;not null" json:"status"`
	CatatanAdmin   string        `gorm:"type:text" json:"catatan_admin"`
	InspectedBy    uuid.NullUUID `gorm:"type:uuid" json:"inspected_by"`
	InspectedAt    *time.Time    `gorm:"type:timestamp with time zone" json:"inspected_at"`
	Nasabah        *User         `gorm:"foreignKey:NasabahID" json:"nasabah,omitempty"`
	Inspector      *User         `gorm:"foreignKey:InspectedBy" json:"inspector,omitempty"`
}

func (PengajuanPinjaman) TableName() string {
	return "pengajuan_pinjaman"
}
