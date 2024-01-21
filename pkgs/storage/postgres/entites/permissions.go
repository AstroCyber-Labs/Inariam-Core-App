package entites

import (
	uuid "github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Permissions struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;"`
	PermissionName string    `gorm:"type:varchar(255);not null"`
	Description    string    `gorm:"type:varchar(255);not null"`
	Provider       string    `gorm:"type:varchar(255);not null"`
	Roles          []Roles   `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Groups         []Groups  `gorm:"many2many:group_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (user *Permissions) BeforeCreate(*gorm.DB) error {
	user.ID = uuid.Must(uuid.NewV4())
	return nil
}
