package api

import (
	"challange/internal/repository"
	"net/http"
)

type Server struct {
	database *repository.Database
	http.Handler
}
