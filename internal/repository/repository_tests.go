package repository

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type TestableRepository struct {
	db   *sql.DB
	Repo Repository
	Mock sqlmock.Sqlmock
	Rows *sqlmock.Rows
}

func RepositoryForTesting(fc func(*sql.DB, *testing.T) Repository, t *testing.T) TestableRepository {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating new mock: %v", err)
	}

	repo := fc(db, t)
	rows := sqlmock.NewRows([]string{"id", "name", "completed"})

	return TestableRepository{
		db:   db,
		Repo: repo,
		Mock: mock,
		Rows: rows,
	}
}

func (tr TestableRepository) TestRepositoryGetTasks(query string, t *testing.T) {
	t.Helper()
	defer tr.db.Close()

	rowCount := addRandomRows(tr.Rows, Task{0, "", true})
	tr.Mock.
		ExpectQuery(query).
		WillReturnRows(tr.Rows)

	tasks, err := tr.Repo.GetAllTasks()

	if err != nil {
		t.Errorf("could not fetch all tasks: %v", err)
	}
	if len(tasks) != rowCount {
		t.Errorf("unexpected number of tasks")
	}
}

func (tr TestableRepository) TestRepositoryGetTaskById(query string, t *testing.T) {
	t.Helper()
	defer tr.db.Close()

	tr.Rows.AddRow(15, "task", false)
	tr.Mock.
		ExpectQuery(query).
		WillReturnRows(tr.Rows)

	task, err := tr.Repo.GetTaskByID(15)

	if err != nil {
		t.Errorf("could not fetch task by id: %v", err)
	}
	if task.ID != 15 {
		t.Errorf("unexpected task id: %v", task.ID)
	}
}

func (tr TestableRepository) TestRepositoryGetTasksByCompletion(query string, t *testing.T, isCompleted bool) {
	t.Helper()
	defer tr.db.Close()

	rowCount := addRandomRows(tr.Rows, Task{0, "", isCompleted})
	tr.Mock.
		ExpectQuery(query).
		WithArgs(isCompleted).
		WillReturnRows(tr.Rows)

	tasks, err := tr.Repo.GetTasksByCompletion(isCompleted)
	if err != nil {
		t.Errorf("could not fetch tasks by completion: %v", err)
	}

	fetchCount := len(tasks)
	if fetchCount != rowCount {
		t.Errorf("expected %v tasks by completion, but got %v", rowCount, fetchCount)
	}
}

func (tr TestableRepository) TestRepositoryAddTask(task Task, t *testing.T) {
	t.Helper()
	defer tr.db.Close()

	id, err := tr.Repo.AddTask(task)
	if err != nil {
		t.Errorf("could not add task: %v", err)
	}
	if id != task.ID {
		t.Errorf("unexpected task id: %v", id)
	}
}

func (tr TestableRepository) TestRepositoryEditTask(task Task, t *testing.T) {
	t.Helper()
	defer tr.db.Close()

	tr.Rows.AddRow(2, "task", false)
	err := tr.Repo.EditTask(task)
	if err != nil {
		t.Errorf("could not add task: %v", err)
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
