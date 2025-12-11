package domain

type User struct {
	BaseModel
	Username    string    `gorm:"type:varchar(100);unique;not null" json:"username"`
	Password    string    `gorm:"type:text;not null" json:"-"`
	NamaLengkap string    `gorm:"type:varchar(255);not null" json:"nama_lengkap"`
	Role        string    `gorm:"type:role_type;not null" json:"role"`
}


