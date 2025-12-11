package domain

import (
	"time"
	"github.com/gofrs/uuid/v5"
	// "gorm.io/gorm"
)
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}