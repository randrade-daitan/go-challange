package main

import (
	"challange/internal/api"
	"challange/internal/orm"
	"challange/internal/repository"
	"log"
	"net/http"
	"os"
)

func main() {
	var db repository.Repository

	switch os.Getenv("DB_IMPL") {
	case "vanilla":
		db = repository.NewDatabase()
	case "orm":
		db = orm.NewOrm()
	default:
		log.Fatal("could not init the database")
	}

	server := api.NewServer(db)
	log.Fatal(http.ListenAndServe(":9090", server))
}
