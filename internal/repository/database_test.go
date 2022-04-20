package repository

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAllTasks(t *testing.T) {
	database, mock, rows := newMockDatabase(t)
	defer database.Close()

	rowCount := addRandomRows(rows, Task{0, "", true})
	mock.
		ExpectQuery("SELECT * FROM task").
		WillReturnRows(rows)

	tasks, err := database.GetAllTasks()

	if err != nil {
		t.Errorf("could not fetch all tasks: %v", err)
	}
	if len(tasks) != rowCount {
		t.Errorf("unexpected number of tasks")
	}
}

func TestGetTaskByID(t *testing.T) {
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
}

func TestAddTask(t *testing.T) {
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
}

func TestEditTask(t *testing.T) {
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
}

func TestGetTasksByCompletion(t *testing.T) {
	t.Run("fetch completed tasks", func(t *testing.T) {
		testFetchTestsByCompletion(t, true)
	})

	t.Run("fetch incompleted tasks", func(t *testing.T) {
		testFetchTestsByCompletion(t, false)
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

func testFetchTestsByCompletion(t *testing.T, isCompleted bool) {
	t.Helper()
	database, mock, rows := newMockDatabase(t)

	rowCount := addRandomRows(rows, Task{0, "", isCompleted})
	mock.
		ExpectQuery("SELECT * FROM task WHERE completed = ?").
		WithArgs(isCompleted).
		WillReturnRows(rows)

	tasks, err := database.GetTasksByCompletion(isCompleted)
	if err != nil {
		t.Errorf("could not fetch tasks by completion: %v", err)
	}

	fetchCount := len(tasks)
	if fetchCount != rowCount {
		t.Errorf("expected %v tasks by completion, but got %v", rowCount, fetchCount)
	}
}

func addRandomRows(rows *sqlmock.Rows, base Task) int {
	min := 10
	max := 30
	rand.Seed(time.Now().UnixNano())
	count := rand.Intn(max-min+1) + min

	for i := 0; i < count; i++ {
		rows.AddRow(i, base.Name+fmt.Sprint(i), base.Completed)
	}

	return count
}
