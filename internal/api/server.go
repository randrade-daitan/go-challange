package api

import (
	"challange/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

var validPath = regexp.MustCompile("^/(tasks)/([0-9]+)$")

type Server struct {
	database repository.Repository
	http.Handler
}

func NewServer(db repository.Repository) Api {
	server := new(Server)

	server.database = db

	router := http.NewServeMux()
	router.Handle("/tasks", http.HandlerFunc(server.handleTasks))
	router.Handle("/tasks/", http.HandlerFunc(server.handleTask))

	server.Handler = router

	return server
}

func (server *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		server.handleGetTasks(w, r)
	case http.MethodPost:
		server.createNewTask(w, r)
	}
}

func (server *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil || len(m) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, e := strconv.Atoi(m[2])
	if handleError(e, w, http.StatusNotFound) {
		return
	}

	id := int64(i)
	switch r.Method {
	case http.MethodGet:
		server.getTaskById(id, w, r)
	case http.MethodPut:
		server.updateTaskById(id, w, r)
	}
}

func (server *Server) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	if isCompleted := r.URL.Query().Get("completed"); isCompleted != "" {
		server.getTasksByCompletion(isCompleted, w, r)
		return
	}

	server.getAllTasks(w, r)
}

func (server *Server) getAllTasks(w http.ResponseWriter, r *http.Request) {
	t, e := server.database.GetAllTasks()
	handleResponse(t, e, w)
}

func (server *Server) getTasksByCompletion(isCompleted string, w http.ResponseWriter, r *http.Request) {
	c, err := strconv.ParseBool(isCompleted)
	if handleError(err, w, http.StatusBadRequest) {
		return
	}

	t, e := server.database.GetTasksByCompletion(c)
	handleResponse(t, e, w)
}

func (server *Server) getTaskById(id int64, w http.ResponseWriter, r *http.Request) {
	t, e := server.database.GetTaskByID(id)
	handleResponse(t, e, w)
}

func (server *Server) createNewTask(w http.ResponseWriter, r *http.Request) {
	t := repository.Task{
		ID:        0,
		Name:      "",
		Completed: false,
	}
	id, e := server.database.AddTask(t)
	if handleError(e, w, http.StatusInternalServerError) {
		return
	}
	fmt.Fprint(w, id)
}

func (server *Server) updateTaskById(id int64, w http.ResponseWriter, r *http.Request) {
	var t repository.Task
	decoreErr := json.NewDecoder(r.Body).Decode(&t)
	if handleError(decoreErr, w, http.StatusBadRequest) {
		return
	}

	t.ID = id
	editErr := server.database.EditTask(t)
	if handleError(editErr, w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleResponse(r any, e error, w http.ResponseWriter) {
	if handleError(e, w, http.StatusInternalServerError) {
		return
	}
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(r)
}

func handleError(err error, w http.ResponseWriter, code int) (hadError bool) {
	if err != nil {
		hadError = true
		http.Error(w, err.Error(), code)
	}
	return
}
