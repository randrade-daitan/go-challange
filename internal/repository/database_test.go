package internal

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFetchOperations(t *testing.T) {
	database, mock := newMockDatabase(t)
	defer database.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "completed"}).
		AddRow(1, "one", false).
		AddRow(2, "two", true)
	mock.ExpectQuery("SELECT * FROM task").WillReturnRows(rows)

	t.Run("get all tasks", func(t *testing.T) {
		tasks, err := database.GetAllTasks()
		if err != nil {
			t.Errorf("could not fetch all tasks: %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("unexpected number of tasks")
		}
	})
}

func newMockDatabase(t *testing.T) (*Database, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return &Database{db}, mock
}
