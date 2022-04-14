package api

import (
	"challange/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type Server struct {
	database repository.Repository
	http.Handler
}

func NewServer(db repository.Repository) Api {
	server := new(Server)

	server.database = db

	router := http.NewServeMux()
	router.Handle("/tasks", http.HandlerFunc(server.handleTasks))

	server.Handler = router

	return server
}

func (server *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if isCompleted := query.Get("completed"); isCompleted != "" {
		c, err := strconv.ParseBool(isCompleted)
		handleRequestError(err, w)
		server.getTasksByCompletion(c, w, r)
		return
	}

	server.getAllTasks(w, r)
}

func (server *Server) getAllTasks(w http.ResponseWriter, r *http.Request) {
	t, e := server.database.GetAllTasks()
	handleTasksRequest(t, e, w)
}

func (server *Server) getTasksByCompletion(isCompleted bool, w http.ResponseWriter, r *http.Request) {
	t, e := server.database.GetTasksByCompletion(isCompleted)
	handleTasksRequest(t, e, w)
}

func handleTasksRequest(tasks []repository.Task, err error, w http.ResponseWriter) {
	handleDatabaseError(err, w)
	encodeTasks(tasks, w)
}

func handleRequestError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleDatabaseError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func encodeTasks(tasks []repository.Task, w http.ResponseWriter) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(tasks)
}
