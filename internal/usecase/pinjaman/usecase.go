package pinjaman

import (
	"context"
	"errors"
	"golang-pinjaman-api/internal/delivery/request"
	"golang-pinjaman-api/internal/domain"
	"golang-pinjaman-api/internal/repository/auth"
	"golang-pinjaman-api/internal/repository/pinjaman"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

var (
	ErrPinjamanNotFound = errors.New("pengajuan pinjaman tidak ditemukan")
	ErrAlreadyProcessed = errors.New("pinjaman sudah diproses (disetujui atau ditolak)")
)

type PinjamanUsecase interface {
	AjukanPinjaman(ctx context.Context, nasabahID uuid.UUID, req request.AjukanPinjamanRequest) (*domain.PengajuanPinjaman, error)
	GetMyPinjamanList(ctx context.Context, nasabahID uuid.UUID) ([]domain.PengajuanPinjaman, error)
	GetAllPinjamanList(ctx context.Context) ([]domain.PengajuanPinjaman, error)
	ApproveOrRejectPinjaman(ctx context.Context, pinjamanID uuid.UUID, adminID uuid.UUID, updateReq request.StatusUpdateRequest) error
}
type pinjamanUsecase struct {
	repo     pinjaman.PinjamanRepository
	userRepo auth.UserRepository
}

func NewPinjamanUsecase(repo pinjaman.PinjamanRepository, userRepo auth.UserRepository) PinjamanUsecase {
	return &pinjamanUsecase{repo: repo, userRepo: userRepo}
}
func (uc *pinjamanUsecase) AjukanPinjaman(ctx context.Context, nasabahID uuid.UUID, req request.AjukanPinjamanRequest) (*domain.PengajuanPinjaman, error) {
	user, err := uc.userRepo.GetByID(ctx, nasabahID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data pengguna")
	}
	newPinjaman := &domain.PengajuanPinjaman{
		NasabahID:      nasabahID,
		Nik:            req.Nik,
		NamaLengkap:    user.NamaLengkap,
		Alamat:         req.Alamat,
		NoTelepon:      req.NoTelepon,
		JumlahPinjaman: req.JumlahPinjaman,
		Status:         "pending",
	}
	if err := uc.repo.Create(ctx, newPinjaman); err != nil {
		return nil, err
	}
	return newPinjaman, nil
}
func (uc *pinjamanUsecase) GetMyPinjamanList(ctx context.Context, nasabahID uuid.UUID) ([]domain.PengajuanPinjaman, error) {
	list, err := uc.repo.GetByNasabahID(ctx, nasabahID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []domain.PengajuanPinjaman{}, nil
		}
		return nil, err
	}
	return list, nil
}
func (uc *pinjamanUsecase) GetAllPinjamanList(ctx context.Context) ([]domain.PengajuanPinjaman, error) {
	return uc.repo.GetAll(ctx)
}
func (uc *pinjamanUsecase) ApproveOrRejectPinjaman(ctx context.Context, pinjamanID uuid.UUID, adminID uuid.UUID, updateReq request.StatusUpdateRequest) error {
	pinjaman, err := uc.repo.GetByID(ctx, pinjamanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPinjamanNotFound
		}
		return err
	}
	if pinjaman.Status != "pending" {
		return ErrAlreadyProcessed
	}
	err = uc.repo.UpdateStatus(ctx, pinjamanID, updateReq.Status, updateReq.CatatanAdmin, adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPinjamanNotFound
		}
		return err
	}
	return nil
}
