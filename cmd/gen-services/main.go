package main

import (
	"gitea/pcp-inariam/inariam/core/config"
	"gitea/pcp-inariam/inariam/pkgs/log"
	"gitea/pcp-inariam/inariam/pkgs/storage/postgres/db"
	"gitea/pcp-inariam/inariam/pkgs/storage/postgres/generator"
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
	// TODO: change this to use the config package to get the db config
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
	var outpath string

	if len(os.Args) > 1 {
		outpath = os.Args[1]
	} else {
		outpath = "core/models/repositories"

	}

	generator.GenerateServices(db_connection, outpath)
	log.Logger.Infoln("Services are Generated")
}
