package main

import (
	"gitea/pcp-inariam/inariam/core/config"
	"gitea/pcp-inariam/inariam/pkgs/log"
	"gitea/pcp-inariam/inariam/pkgs/storage/postgres/db"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {

	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Logger.Panicln("Error parsing database port")
	}

	db_connection, err := db.ConnectDB(&config.DbConfig{
		Username:     os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DatabaseName: os.Getenv("POSTGRES_DB"),
		Host:         os.Getenv("POSTGRES_HOST"),
		Port:         uint(port),
	})

	if err != nil {
		panic(err)
	}

	err = db.MigrateDB(db_connection)

	if err != nil {
		panic(err)
	}

	log.Logger.Infoln("Migrated")
}
