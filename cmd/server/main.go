package main

import (
	"challange/internal/api"
	"challange/internal/repository"
	"log"
	"os"
)

func main() {
	repo := newRepository()
	server := api.NewRestServer(repo)
	log.Fatal(server.StartServing(9090))
}

func newRepository() (repo repository.Repository) {
	switch os.Getenv("DB_IMPL") {
	case "vanilla":
		repo = repository.NewMySqlRepository()
	case "orm":
		repo = repository.NewOrmRepository()
	default:
		log.Fatal("could not init the database")
	}
	return
}
