package api

import (
	"challange/internal/repository"
	"testing"
)

func TestRecordInsertAndGetTasks(t *testing.T) {
	tasks := []repository.Task{}
	server, _ := newTestServer(tasks)

	performAddTaskRequest(server, t)
	performAddTaskRequest(server, t)
	performAddTaskRequest(server, t)

	t.Run("get tasks", func(t *testing.T) {
		got := performGetTasksRequest(server, t)
		if len(got) != 3 {
			t.Errorf("did not get correct tasks count")
		}
	})

	t.Run("update and check tasks", func(t *testing.T) {
		task := repository.Task{
			ID:        1,
			Name:      "Edited",
			Completed: true,
		}
		performUpdateTaskRequest(task, server, t)

		complete := performGetTasksByCompletionRequest(true, server, t)
		if len(complete) != 1 {
			t.Errorf("did not get correct complete tasks count")
		}

		incomplete := performGetTasksByCompletionRequest(false, server, t)
		if len(incomplete) != 2 {
			t.Errorf("did not get correct incomplete tasks count")
		}
	})

	t.Run("add and check tasks", func(t *testing.T) {
		performAddTaskRequest(server, t)
		got := performGetTasksRequest(server, t)
		if len(got) != 4 {
			t.Errorf("did not get correct tasks count after adding")
		}
	})
}
