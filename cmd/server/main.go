package main

import (
	"challange/internal/api/restapi"
	"challange/internal/repository"
	"log"
)

func main() {
	repo := repository.NewRepository()
	server := restapi.NewServer(repo)
	log.Fatal(server.StartServing(9090))
}
