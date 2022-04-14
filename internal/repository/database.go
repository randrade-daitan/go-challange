package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

func NewDatabase() Repository {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "challenge",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return &Database{db}
}

func (db *Database) GetAllTasks() ([]Task, error) {
	var tasks []Task

	rows, err := db.Query("SELECT * FROM task")
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

func (db *Database) GetTaskByID(id int64) (Task, error) {
	var task Task

	row := db.QueryRow("SELECT * FROM task WHERE id = ?", id)
	if err := row.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("GetTaskByID %d: no such task", id)
		}
		return task, fmt.Errorf("GetTaskByID %d: %v", id, err)
	}

	return task, nil
}

func (db *Database) GetTaskByCompletion(isCompleted bool) (Task, error) {
	// TODO b.ii.5: Implement get task by completion
	var task Task
	return task, nil
}

func (db *Database) AddTask(t Task) (int64, error) {
	result, err := db.Exec("INSERT INTO task (name, completed) VALUES (?, ?)", t.Name, t.Completed)
	if err != nil {
		return 0, fmt.Errorf("addTask: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addTask: %v", err)
	}
	return id, nil
}

func (db *Database) EditTask(t Task) error {
	// TODO b.ii.4: Update a task
	return nil
}
