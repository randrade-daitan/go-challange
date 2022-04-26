package main

import (
	"challange/internal/api"
	"challange/internal/orm"
	"challange/internal/repository"
	"log"
	"os"
)

func main() {
	var repo repository.Repository

	switch os.Getenv("DB_IMPL") {
	case "vanilla":
		repo = repository.NewDatabase()
	case "orm":
		repo = orm.NewOrm()
	default:
		log.Fatal("could not init the database")
	}

	server := api.NewServer(repo)
	log.Fatal(server.StartServing(9090))
}
