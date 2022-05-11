package main

import (
	"challange/internal/api/grpcapi"
	"challange/internal/repository"
	"log"
)

func main() {
	repo := repository.NewRepository()
	server := grpcapi.NewServer(repo)
	log.Fatal(server.StartServing(9091))
}
