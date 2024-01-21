package db

import (
	"gitea/pcp-inariam/inariam/pkgs/storage/postgres/entites"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {

	err := db.AutoMigrate(
		&entites.Users{},
		&entites.Groups{},
		&entites.Roles{},
		&entites.Permissions{},
		&entites.Teams{},
		&entites.Accounts{},
	)

	if err != nil {
		return err
	}

	return nil
}
