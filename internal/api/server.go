package api

import (
	"challange/internal/repository"
	"encoding/json"
	"net/http"
)

type Server struct {
	database repository.Repository
	http.Handler
}

func NewServer(db repository.Repository) Api {
	server := new(Server)

	server.database = db

	router := http.NewServeMux()
	router.Handle("/tasks", http.HandlerFunc(server.getAllTasks))

	server.Handler = router

	return server
}

func (server *Server) getAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := server.database.GetAllTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(tasks)
}
