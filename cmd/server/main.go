package main

import (
	"challange/internal/api/restapi"
	"challange/internal/repository"
	"challange/internal/repository/mysqlrepo"
	"challange/internal/repository/ormrepo"
	"log"
	"os"
)

func main() {
	var repo repository.Repository

	switch os.Getenv("DB_IMPL") {
	case "vanilla":
		repo = mysqlrepo.NewRepository()
	case "orm":
		repo = ormrepo.NewRepository()
	default:
		log.Fatal("could not init the database")
	}

	server := restapi.NewServer(repo)
	log.Fatal(server.StartServing(9090))
}
