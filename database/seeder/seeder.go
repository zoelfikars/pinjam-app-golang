package seeder

import (
	"errors"
	"golang-pinjaman-api/internal/domain"
	"golang-pinjaman-api/pkg/util"
	"log"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	hashedPassword, _ := util.HashPassword("admin123")

	adminUser := domain.User{
		Username:    "admin@pinjaman.com",
		Password:    hashedPassword,
		NamaLengkap: "Administrator Sistem",
		Role:        "admin",
	}

	var existingUser domain.User
	result := db.Where("username = ?", adminUser.Username).First(&existingUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatalf("Failed to seed admin user: %v", err)
		}
		log.Println("Admin user created successfully (Username: admin@pinjaman.com, Pass: admin123)")
	} else if result.Error == nil {
		log.Println("Admin user already exists. Skipping seeding.")
	} else {
		log.Fatalf("Error checking for existing admin: %v", result.Error)
	}
}
