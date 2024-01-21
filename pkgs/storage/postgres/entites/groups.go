package entites

import (
	uuid "github.com/gofrs/uuid"

	"gorm.io/gorm"
)

// Groups has and belongs to many Teams, `group_teams` is the join table
// Groups has and belongs to many Permissions, `group_permissions` is the join table
type Groups struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;"`
	GroupName   string        `gorm:"type:varchar(255);not null"`
	Teams       []Teams       `gorm:"many2many:group_teams;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Permissions []Permissions `gorm:"many2many:group_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (user *Groups) BeforeCreate(*gorm.DB) error {
	user.ID = uuid.Must(uuid.NewV4())
	return nil
}
