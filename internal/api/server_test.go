package api

import (
	"challange/internal/repository"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchOperations(t *testing.T) {
	tasks := []repository.Task{
		{ID: 0, Name: "a", Completed: true},
		{ID: 1, Name: "b", Completed: false},
		{ID: 2, Name: "c", Completed: true},
		{ID: 3, Name: "d", Completed: true},
		{ID: 4, Name: "e", Completed: false},
	}
	db := repository.NewMockDatabase(tasks, nil)
	server := NewServer(db)

	t.Run("serve all tasks", func(t *testing.T) {
		request := newGetAllTasksRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		got := getTasksFromResponse(t, response.Body)
		if len(got) != len(tasks) {
			t.Errorf("did not get correct tasks count")
		}
	})
}

func newGetAllTasksRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	return req
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

func getTasksFromResponse(t testing.TB, body io.Reader) (tasks []repository.Task) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&tasks)
	if err != nil {
		t.Fatalf("unable to parse response from server %q: %v", body, err)
	}

	return
}
