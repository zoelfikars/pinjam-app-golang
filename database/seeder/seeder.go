package seeder

import (
	"errors"
	"github.com/gofrs/uuid/v5"
	"golang-pinjaman-api/internal/domain"
	"golang-pinjaman-api/pkg/util"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

var dummyUsers = []struct {
	Username    string
	Password    string
	NamaLengkap string
	Role        string
}{
	{"admin", "admin123", "Administrator Sistem", "admin"},
	{"nasabah1", "nasabah1", "Budi Santoso", "nasabah"},
	{"nasabah2", "nasabah2", "Citra Dewi", "nasabah"},
	{"nasabah3", "nasabah3", "Dedi Kurniawan", "nasabah"},
	{"inspector", "inspector123", "Inspektur Lapangan", "admin"},
}
var statuses = []string{"pending", "approved", "rejected"}

func newUUID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}
func randomLoanAmount() int64 {
	min := int64(1000000)
	max := int64(100000000)
	amount := min + rand.Int63n(max-min+1)
	return (amount / 500000) * 500000
}
func randomPastTime() *time.Time {
	t := time.Now().AddDate(0, 0, -rand.Intn(30))
	return &t
}
func Seed(db *gorm.DB) {
	log.Println("--- Starting Database Seeding ---")
	var users []domain.User
	for _, du := range dummyUsers {
		var existingUser domain.User
		result := db.Where("username = ?", du.Username).First(&existingUser)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			hashedPassword, _ := util.HashPassword(du.Password)
			user := domain.User{
				BaseModel:   domain.BaseModel{ID: newUUID()},
				Username:    du.Username,
				Password:    hashedPassword,
				NamaLengkap: du.NamaLengkap,
				Role:        du.Role,
			}
			if err := db.Create(&user).Error; err != nil {
				log.Fatalf("Failed to seed user %s: %v", du.Username, err)
			}
			users = append(users, user)
			log.Printf("User created: %s (Role: %s)", du.Username, du.Role)
		} else if result.Error == nil {
			users = append(users, existingUser)
			log.Printf("User already exists. Skipping seeding for %s.", du.Username)
		} else {
			log.Fatalf("Error checking for existing user %s: %v", du.Username, result.Error)
		}
	}
	var inspectorUser domain.User
	for _, user := range users {
		if user.Role == "admin" && user.Username == "inspector" {
			inspectorUser = user
			break
		}
	}
	for _, user := range users {
		if user.Role == "nasabah" {
			numPinjaman := 1 + rand.Intn(10)
			for i := 0; i < numPinjaman; i++ {
				status := statuses[rand.Intn(len(statuses))]
				var catatanAdmin string
				var inspectedBy uuid.NullUUID
				var inspectedAt *time.Time
				if status != "pending" {
					if status == "rejected" {
						catatanAdmin = "Tidak memenuhi skor kredit minimum."
					} else {
						catatanAdmin = ""
					}
					inspectedBy = uuid.NullUUID{UUID: inspectorUser.ID, Valid: true}
					inspectedAt = randomPastTime()
				}
				pinjaman := domain.PengajuanPinjaman{
					BaseModel: domain.BaseModel{
						ID:        newUUID(),
						CreatedAt: *randomPastTime(),
					},
					NasabahID:      user.ID,
					Nik:            "3201" + util.RandomString(12, "0123456789"),
					NamaLengkap:    user.NamaLengkap,
					Alamat:         "Jl. Merdeka No. " + util.RandomString(2, "0123456789"),
					NoTelepon:      "08" + util.RandomString(10, "0123456789"),
					JumlahPinjaman: randomLoanAmount(),
					Status:         status,
					CatatanAdmin:   catatanAdmin,
					InspectedBy:    inspectedBy,
					InspectedAt:    inspectedAt,
				}
				if err := db.Create(&pinjaman).Error; err != nil {
					log.Fatalf("Failed to seed pinjaman for %s: %v", user.Username, err)
				}
			}
			log.Printf("Created %d pinjaman for nasabah %s.", numPinjaman, user.Username)
		}
	}
	log.Println("--- Database Seeding Complete ---")
}
