package server

import (
	"challange/internal/api"
	"challange/internal/repository"
	"log"
	"net/http"
)

func Main() {
	db := repository.NewDatabase()
	server := api.NewServer(db)
	log.Fatal(http.ListenAndServe(":9090", server))
}
