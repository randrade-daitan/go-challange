package main

import (
	"challange/internal/api/grpcapi"
	"challange/internal/repository/mysqlrepo"
	"log"
)

func main() {
	repo := mysqlrepo.NewRepository()
	server := grpcapi.NewServer(repo)
	log.Fatal(server.StartServing(9091))
}
