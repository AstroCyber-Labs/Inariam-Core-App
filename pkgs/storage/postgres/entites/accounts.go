package entites

import (
	uuid "github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Accounts struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;"`
	Provider string    `gorm:"type:varchar(255);not null"`
	Creds    string    `gorm:"type:text;not null"`
	UserID   uuid.UUID `gorm:"type:uuid;"`
}

func (user *Accounts) BeforeCreate(*gorm.DB) error {
	user.ID = uuid.Must(uuid.NewV4())
	return nil
}
