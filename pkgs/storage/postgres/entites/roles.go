package entites

import (
	uuid "github.com/gofrs/uuid"

	"gorm.io/gorm"
)

// Roles has and belongs to many Permissions, `role_permissions` is the join table
type Roles struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;"`
	RoleName    string        `gorm:"type:varchar(255);not null"`
	Permissions []Permissions `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Users       []Users       `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (user *Roles) BeforeCreate(*gorm.DB) error {
	user.ID = uuid.Must(uuid.NewV4())
	return nil
}
