package main

import (
	"challange/internal/grpcApi"
	"challange/internal/repository"
	"log"
)

func main() {
	repo := repository.NewRepository()
	server := grpcApi.NewServer(repo)
	log.Fatal(server.StartServing(9091))
}
