package entites

import (
	uuid "github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Users has and belongs to many Teams, `user_teams` is the join table
// Users has and belongs to many Roles, `user_roles` is the join table
// Users has many Accounts, UserID is the foreign key
type Users struct {
	UserID   uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name     string    `gorm:"type:varchar(255);not null"`
	Email    string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"`

	Accounts    []Accounts `gorm:"foreignKey:UserID"`
	Roles       []Roles    `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Teams       []Teams    `gorm:"many2many:user_teams;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OtpEnabled  bool       `gorm:"default:false;"`
	OtpVerified bool       `gorm:"default:false;"`

	OtpSecret  string
	OtpAuthUrl string
}

func (user *Users) BeforeCreate(*gorm.DB) error {
	user.UserID = uuid.Must(uuid.NewV4())
	return nil
}
