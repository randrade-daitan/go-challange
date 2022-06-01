package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMySqlGetAllTasks(t *testing.T) {
	repo := RepositoryForTesting(testableDatabase, t)
	repo.TestRepositoryGetTasks("SELECT * FROM task", t)
}

func TestMySqlGetTaskByID(t *testing.T) {
	repo := RepositoryForTesting(testableDatabase, t)
	repo.TestRepositoryGetTaskById("SELECT * FROM task WHERE id = ?", t)
}

func TestMySqlAddTask(t *testing.T) {
	task := Task{10, "DB Test", false}
	repo := RepositoryForTesting(testableDatabase, t)

	repo.Mock.
		ExpectExec("INSERT INTO task (name, completed) VALUES (?, ?)").
		WithArgs(task.Name, task.Completed).
		WillReturnResult(sqlmock.NewResult(task.ID, 1))

	repo.TestRepositoryAddTask(task, t)
}

func TestMySqlEditTask(t *testing.T) {
	task := Task{2, "edited", true}
	repo := RepositoryForTesting(testableDatabase, t)

	repo.Mock.
		ExpectExec("UPDATE task SET name = ?, completed = ? WHERE id = ?").
		WithArgs(task.Name, task.Completed, task.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo.TestRepositoryEditTask(task, t)
}

func TestMySqlGetTasksByCompletion(t *testing.T) {
	q := "SELECT * FROM task WHERE completed = ?"

	t.Run("fetch completed tasks", func(t *testing.T) {
		repo := RepositoryForTesting(testableDatabase, t)
		repo.TestRepositoryGetTasksByCompletion(q, t, true)
	})

	t.Run("fetch uncompleted tasks", func(t *testing.T) {
		repo := RepositoryForTesting(testableDatabase, t)
		repo.TestRepositoryGetTasksByCompletion(q, t, false)
	})
}

func testableDatabase(db *sql.DB, t *testing.T) Repository {
	t.Helper()
	return &Database{db}
}
