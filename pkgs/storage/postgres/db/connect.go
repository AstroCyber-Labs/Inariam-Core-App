package db

import (
	"fmt"
	"gitea/pcp-inariam/inariam/core/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(config *config.DbConfig) (*gorm.DB, error) {

	dsnString := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"

	fmt.Printf("%+v %+v", dsnString, config)
	dsn := fmt.Sprintf(
		dsnString,
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.DatabaseName,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
