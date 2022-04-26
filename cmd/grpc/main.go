package main

import (
	"challange/internal/grpcApi"
	"challange/internal/repository"
	"log"
)

func main() {
	repo := repository.NewDatabase()
	server := grpcApi.NewServer(repo)
	log.Fatal(server.StartServing(9091))
}
