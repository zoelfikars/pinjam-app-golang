package pinjaman

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"golang-pinjaman-api/internal/domain"
	"gorm.io/gorm"
)

type PinjamanRepository interface {
	Create(ctx context.Context, pinjaman *domain.PengajuanPinjaman) error
	GetByNasabahID(ctx context.Context, nasabahID uuid.UUID) ([]domain.PengajuanPinjaman, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PengajuanPinjaman, error)
	GetAll(ctx context.Context) ([]domain.PengajuanPinjaman, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, catatan string, inspectedBy uuid.UUID) error
}
type pinjamanRepository struct {
	DB *gorm.DB
}

func NewPinjamanRepository(db *gorm.DB) PinjamanRepository {
	return &pinjamanRepository{DB: db}
}
func (r *pinjamanRepository) Create(ctx context.Context, pinjaman *domain.PengajuanPinjaman) error {
	return r.DB.WithContext(ctx).Create(pinjaman).Error
}
func (r *pinjamanRepository) GetByNasabahID(ctx context.Context, nasabahID uuid.UUID) ([]domain.PengajuanPinjaman, error) {
	var pengajuans []domain.PengajuanPinjaman
	err := r.DB.WithContext(ctx).Preload("Inspector").Where("nasabah_id = ?", nasabahID).Order("created_at desc").Find(&pengajuans).Error
	return pengajuans, err
}
func (r *pinjamanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PengajuanPinjaman, error) {
	var pinjaman domain.PengajuanPinjaman
	err := r.DB.WithContext(ctx).Preload("Nasabah").First(&pinjaman, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pinjaman, nil
}
func (r *pinjamanRepository) GetAll(ctx context.Context) ([]domain.PengajuanPinjaman, error) {
	var pengajuans []domain.PengajuanPinjaman
	err := r.DB.WithContext(ctx).Preload("Nasabah").Preload("Inspector").Order("created_at desc").Find(&pengajuans).Error
	return pengajuans, err
}
func (r *pinjamanRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, catatan string, inspectedBy uuid.UUID) error {
	updateFields := map[string]interface{}{
		"status":        status,
		"catatan_admin": catatan,
		"inspected_by":   inspectedBy,
		"inspected_at":   gorm.Expr("NOW()"),
	}
	result := r.DB.WithContext(ctx).Model(&domain.PengajuanPinjaman{}).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
