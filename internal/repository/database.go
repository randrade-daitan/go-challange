package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

func NewDatabase() Repository {
	cfg := mysql.Config{
		User:   DBUser(),
		Passwd: DBPass(),
		Net:    DBProtocol,
		Addr:   DBURL(),
		DBName: DBName(),
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return &Database{db}
}

func (db *Database) GetAllTasks() ([]Task, error) {
	return db.queryForTasks("SELECT * FROM task")
}

func (db *Database) GetTaskByID(id int64) (Task, error) {
	return db.queryRowForTask("SELECT * FROM task WHERE id = ?", id)
}

func (db *Database) GetTasksByCompletion(isCompleted bool) ([]Task, error) {
	return db.queryForTasks("SELECT * FROM task WHERE completed = ?", isCompleted)
}

func (db *Database) AddTask(t Task) (int64, error) {
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

func (db *Database) EditTask(t Task) error {
	_, err := db.executeQuery("UPDATE task SET name = ?, completed = ? WHERE id = ?", t.Name, t.Completed, t.ID)
	return err
}

func (db *Database) queryForTasks(query string, args ...any) ([]Task, error) {
	tasks := []Task{}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not fetch tasks rows: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var task Task
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

func (db *Database) queryRowForTask(query string, args ...any) (Task, error) {
	var task Task

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
