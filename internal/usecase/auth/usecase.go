package auth

import (
	"context"
	"errors"
	"golang-pinjaman-api/internal/delivery/request"
	"golang-pinjaman-api/internal/domain"
	"golang-pinjaman-api/internal/repository/auth"
	"golang-pinjaman-api/internal/repository/redis"
	"golang-pinjaman-api/pkg/util"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user tidak ditemukan")
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrUserExist          = errors.New("username sudah terdaftar")
)

type AuthUsecase interface {
	Register(ctx context.Context, req request.RegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req request.LoginRequest) (string, time.Time, *domain.User, error)
	Logout(ctx context.Context, tokenString string) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type authUsecase struct {
	userRepo  auth.UserRepository
	tokenRepo redis.TokenRepository
	jwtSecret string
}

func NewAuthUsecase(userRepo auth.UserRepository, tokenRepo redis.TokenRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtSecret: jwtSecret,
	}
}

func (uc *authUsecase) Register(ctx context.Context, req request.RegisterRequest) (*domain.User, error) {

	_, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, ErrUserExist
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal menghash password")
	}

	newUser := &domain.User{

		Username:    req.Username,
		Password:    hashedPassword,
		NamaLengkap: req.NamaLengkap,
		Role:        "nasabah",
	}

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *authUsecase) Login(ctx context.Context, req request.LoginRequest) (string, time.Time, *domain.User, error) {

	user, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", time.Time{}, nil, ErrInvalidCredentials
		}
		return "", time.Time{}, nil, err
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return "", time.Time{}, nil, ErrInvalidCredentials
	}

	token, expiryTime, err := util.GenerateToken(user, uc.jwtSecret, 24*time.Hour)
	if err != nil {
		return "", time.Time{}, nil, errors.New("gagal menghasilkan token")
	}

	return token, expiryTime, user, nil
}

func (uc *authUsecase) Logout(ctx context.Context, tokenString string) error {

	claims := &util.Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.jwtSecret), nil
	})

	var expiry time.Duration
	if claims.ExpiresAt != nil {
		expiry = time.Until(claims.ExpiresAt.Time)
	}

	if expiry > 0 {
		return uc.tokenRepo.BlacklistToken(ctx, tokenString, expiry)
	}

	if err != nil {
		return errors.New("gagal logout")
	}

	return nil
}

func (uc *authUsecase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
