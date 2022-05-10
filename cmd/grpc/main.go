package main

import (
	"challange/internal/grpcApi"
	"challange/internal/repository"
	"log"
)

func main() {
	db := repository.NewDatabase()
	server := grpcApi.NewServer(db)
	log.Fatal(server.StartServing(9091))
}
