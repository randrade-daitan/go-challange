package internal

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFetchOperations(t *testing.T) {
	t.Run("get all tasks", func(t *testing.T) {
		database, mock, rows := newMockDatabase(t)
		defer database.Close()
		rows.
			AddRow(1, "one", false).
			AddRow(2, "two", true)

		mock.ExpectQuery("SELECT * FROM task").WillReturnRows(rows)
		tasks, err := database.GetAllTasks()

		if err != nil {
			t.Errorf("could not fetch all tasks: %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("unexpected number of tasks")
		}
	})

	t.Run("get task by id", func(t *testing.T) {
		database, mock, rows := newMockDatabase(t)

		rows.AddRow(15, "task", false)
		mock.
			ExpectQuery("SELECT * FROM task WHERE id = ?").
			WillReturnRows(rows)

		task, err := database.GetTaskByID(15)

		if err != nil {
			t.Errorf("could not fetch task by id: %v", err)
		}
		if task.ID != 15 {
			t.Errorf("unexpected task id: %v", task.ID)
		}
	})

	t.Run("add new task", func(t *testing.T) {
		task := Task{10, "Title", false}
		database, mock, _ := newMockDatabase(t)

		mock.
			ExpectExec("INSERT INTO task (name, completed) VALUES (?, ?)").
			WithArgs(task.Name, task.Completed).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := database.AddTask(task)

		if err != nil {
			t.Errorf("could not add task: %v", err)
		}
		if id != 1 {
			t.Errorf("unexpected task id: %v", id)
		}
	})

	t.Run("edit existing task", func(t *testing.T) {
		task := Task{2, "edited", true}
		database, mock, rows := newMockDatabase(t)

		rows.AddRow(2, "task", false)
		mock.
			ExpectExec("UPDATE task SET name = ?, completed = ? WHERE id = ?").
			WithArgs(task.Name, task.Completed, task.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := database.EditTask(task)
		if err != nil {
			t.Errorf("could not add task: %v", err)
		}
	})
}

func newMockDatabase(t *testing.T) (*Database, sqlmock.Sqlmock, *sqlmock.Rows) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "completed"})
	return &Database{db}, mock, rows
}
