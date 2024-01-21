package entites

import (
	uuid "github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Teams struct {
	TeamID    uuid.UUID `gorm:"type:uuid;primary_key;"`
	TeamName  string    `gorm:"type:varchar(255);not null"`
	Groups    []Groups  `gorm:"many2many:group_teams;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Users     []Users   `gorm:"many2many:user_teams;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TeamAdmin uuid.UUID `gorm:"foreignKey:UserID"`
}

func (user *Teams) BeforeCreate(*gorm.DB) error {
	user.TeamID = uuid.Must(uuid.NewV4())
	return nil
}
