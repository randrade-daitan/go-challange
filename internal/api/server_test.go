package api

import (
	"bytes"
	"challange/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOperations(t *testing.T) {
	tasks := []repository.Task{
		{ID: 0, Name: "a", Completed: true},
		{ID: 1, Name: "b", Completed: false},
		{ID: 2, Name: "c", Completed: true},
		{ID: 3, Name: "d", Completed: true},
		{ID: 4, Name: "e", Completed: false},
	}
	server, _ := newTestServer(tasks)

	t.Run("serve all tasks", func(t *testing.T) {
		got := performGetTasksRequest(server, t)
		if len(got) != len(tasks) {
			t.Errorf("did not get correct tasks count")
		}
	})

	t.Run("serve completed tasks", func(t *testing.T) {
		got := performGetTasksByCompletionRequest(true, server, t)
		if len(got) != 3 {
			t.Errorf("did not get correct completed tasks count")
		}
	})

	t.Run("serve incompleted tasks", func(t *testing.T) {
		got := performGetTasksByCompletionRequest(false, server, t)
		if len(got) != 2 {
			t.Errorf("did not get correct completed tasks count")
		}
	})

	t.Run("serve task by id", func(t *testing.T) {
		id := int64(1)
		got := performGetTasksByIdRequest(id, server, t)
		if got.ID != id || got.Name != "b" {
			t.Errorf("did not get correct task by id")
		}
	})

	t.Run("serve task by no id", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/tasks/", nil)
		performRequest(server, r, t, http.StatusBadRequest)
	})

	t.Run("serve task by bad id", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/tasks/bad", nil)
		performRequest(server, r, t, http.StatusBadRequest)
	})
}

func TestPostOperations(t *testing.T) {
	tasks := []repository.Task{
		{ID: 0, Name: "", Completed: false},
	}
	server, _ := newTestServer(tasks)

	t.Run("serve add task", func(t *testing.T) {
		response := performAddTaskRequest(server, t)
		if response.Body.String() != "1" {
			t.Errorf("did not get correct completed tasks count")
		}
	})
}

func TestPutOperations(t *testing.T) {
	tasks := []repository.Task{
		{ID: 6, Name: "test", Completed: false},
	}
	server, db := newTestServer(tasks)

	t.Run("serve add task", func(t *testing.T) {
		task := repository.Task{
			ID:        6,
			Name:      "Edited",
			Completed: true,
		}
		performUpdateTaskRequest(task, server, t)

		editedTask, _ := db.GetTaskByID(task.ID)
		if editedTask.Name != task.Name ||
			editedTask.Completed != task.Completed {
			t.Errorf("did not get edited task correctly")
		}
	})
}

func newTestServer(tasks []repository.Task) (Api, repository.Repository) {
	db := repository.NewMockDatabase(tasks, nil)
	return NewServer(db), db
}

func performRequest(server Api, r *http.Request, t *testing.T, expectingCode int) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()

	server.ServeHTTP(response, r)
	assertStatus(t, response.Code, expectingCode)

	return response
}

func performFetchRequest[T any](path string, got *T, server Api, t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, path, nil)
	response := performRequest(server, r, t, http.StatusOK)
	assertContentType(t, response, jsonContentType)
	parseResponse(t, response.Body, &got)
}

func performGetTasksRequest(server Api, t *testing.T) (got []repository.Task) {
	performFetchRequest("/tasks", &got, server, t)
	return
}

func performGetTasksByIdRequest(id int64, server Api, t *testing.T) (got repository.Task) {
	performFetchRequest(fmt.Sprintf("/tasks/%v", id), &got, server, t)
	return
}

func performGetTasksByCompletionRequest(completion bool, server Api, t *testing.T) (got []repository.Task) {
	performFetchRequest(fmt.Sprintf("/tasks?completed=%v", completion), &got, server, t)
	return
}

func performAddTaskRequest(server Api, t *testing.T) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(http.MethodPost, "/tasks", nil)
	return performRequest(server, request, t, http.StatusOK)
}

func performUpdateTaskRequest(task repository.Task, server Api, t *testing.T) {
	body, _ := json.Marshal(task)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%v", task.ID), bytes.NewBuffer(body))
	performRequest(server, request, t, http.StatusOK)
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func parseResponse[T any](t testing.TB, body io.Reader, r *T) {
	t.Helper()

	err := json.NewDecoder(body).Decode(r)
	if err != nil {
		t.Fatalf("unable to parse response from server %q: %v", body, err)
	}
}
