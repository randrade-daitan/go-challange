package mysqlrepo

import (
	"challange/internal/repository"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

// Creates a new repository using the vanilla implementation.
func NewRepository() repository.Repository {
	cfg := mysql.Config{
		User:   repository.DBUser(),
		Passwd: repository.DBPass(),
		Net:    repository.DBProtocol,
		Addr:   repository.DBURL(),
		DBName: repository.DBName(),
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return &Database{db}
}

func (db *Database) GetAllTasks() ([]repository.Task, error) {
	return db.queryForTasks("SELECT * FROM task")
}

func (db *Database) GetTaskByID(id int64) (repository.Task, error) {
	return db.queryRowForTask("SELECT * FROM task WHERE id = ?", id)
}

func (db *Database) GetTasksByCompletion(isCompleted bool) ([]repository.Task, error) {
	return db.queryForTasks("SELECT * FROM task WHERE completed = ?", isCompleted)
}

func (db *Database) AddTask(t repository.Task) (int64, error) {
	result, err := db.executeQuery("INSERT INTO task (name, completed) VALUES (?, ?)", t.Name, t.Completed)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addTask: %v", err)
	}

	return id, nil
}

func (db *Database) EditTask(t repository.Task) error {
	_, err := db.executeQuery("UPDATE task SET name = ?, completed = ? WHERE id = ?", t.Name, t.Completed, t.ID)
	return err
}

func (db *Database) queryForTasks(query string, args ...any) ([]repository.Task, error) {
	tasks := []repository.Task{}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not fetch tasks rows: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var task repository.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			return nil, fmt.Errorf("could not fetch tasks next row: %v", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not fetch tasks row: %v", err)
	}

	return tasks, nil
}

func (db *Database) queryRowForTask(query string, args ...any) (repository.Task, error) {
	var task repository.Task

	row := db.QueryRow(query, args...)
	if err := row.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
		return task, fmt.Errorf("error scanning task row: %v", err)
	}

	return task, nil
}

func (db *Database) executeQuery(query string, args ...any) (sql.Result, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %v", err)
	}
	return result, nil
}
