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
	// TODO b.ii.2: Implement get task by id
	var task Task
	return task, nil
}

func (db *Database) GetTaskByCompletion(isCompleted bool) (Task, error) {
	// TODO b.ii.5: Implement get task by completion
	var task Task
	return task, nil
}

func (db *Database) AddTask(t Task) (int64, error) {
	// TODO b.ii.3: Add a task
	return 0, nil
}

func (db *Database) EditTask(t Task) error {
	// TODO b.ii.4: Update a task
	return nil
}
