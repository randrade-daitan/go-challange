package main

import (
	"challange/internal/api"
	"challange/internal/repository"
	"log"
)

func main() {
	repo := repository.NewMySqlRepository()
	server := api.NewGrpcServer(repo)
	log.Fatal(server.StartServing(9091))
}
