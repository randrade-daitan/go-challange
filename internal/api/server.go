package api

import (
	"challange/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const authorizationKey = "Authorization"

var validPath = regexp.MustCompile("^/(tasks)/([0-9]+)$")

type Server struct {
	repo repository.Repository
	http.Handler
}

func NewServer(repo repository.Repository) *Server {
	server := new(Server)

	server.repo = repo

	router := http.NewServeMux()
	router.Handle("/tasks", authenticatedHandler(server.handleTasks))
	router.Handle("/tasks/", authenticatedHandler(server.handleTask))

	server.Handler = router

	return server
}

func authenticatedHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get(authorizationKey)
		token := "Bearer " + os.Getenv("BEARER_TOKEN")

		if auth != token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fn(w, r)
	}
}

func (server *Server) StartServing(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%v", port), server)
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
	matches := validPath.FindStringSubmatch(r.URL.Path)
	if matches == nil || len(matches) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err := strconv.Atoi(matches[2])
	if handleError(err, w, http.StatusNotFound) {
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
	t, e := server.repo.GetAllTasks()
	handleResponse(t, e, w)
}

func (server *Server) getTasksByCompletion(isCompleted string, w http.ResponseWriter, r *http.Request) {
	completed, err := strconv.ParseBool(isCompleted)
	if handleError(err, w, http.StatusBadRequest) {
		return
	}

	t, e := server.repo.GetTasksByCompletion(completed)
	handleResponse(t, e, w)
}

func (server *Server) getTaskById(id int64, w http.ResponseWriter, r *http.Request) {
	t, e := server.repo.GetTaskByID(id)
	handleResponse(t, e, w)
}

func (server *Server) createNewTask(w http.ResponseWriter, r *http.Request) {
	task := repository.Task{
		ID:        0,
		Name:      "",
		Completed: false,
	}
	id, e := server.repo.AddTask(task)
	if handleError(e, w, http.StatusInternalServerError) {
		return
	}
	fmt.Fprint(w, id)
}

func (server *Server) updateTaskById(id int64, w http.ResponseWriter, r *http.Request) {
	var task repository.Task
	decodeErr := json.NewDecoder(r.Body).Decode(&task)
	if handleError(decodeErr, w, http.StatusBadRequest) {
		return
	}

	task.ID = id
	editErr := server.repo.EditTask(task)
	if handleError(editErr, w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleResponse(r any, err error, w http.ResponseWriter) {
	if handleError(err, w, http.StatusInternalServerError) {
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
