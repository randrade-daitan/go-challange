package repository

import "errors"

type MockDatabase struct {
	tasks []Task
	err   error
}

func NewMockDatabase(tasks []Task, err error) Repository {
	return &MockDatabase{tasks, err}
}

func (db *MockDatabase) GetAllTasks() ([]Task, error) {
	return db.tasks, db.err
}

func (db *MockDatabase) GetTaskByID(id int64) (Task, error) {
	var task Task
	for _, task := range db.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return task, errors.New("could not find task")
}

func (db *MockDatabase) GetTasksByCompletion(isCompleted bool) ([]Task, error) {
	tasks := []Task{}

	for _, task := range db.tasks {
		if task.Completed == isCompleted {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (db *MockDatabase) AddTask(t Task) (int64, error) {
	newId := int64(len(db.tasks))
	t.ID = newId
	db.tasks = append(db.tasks, t)
	return newId, nil
}

func (db *MockDatabase) EditTask(t Task) error {
	for i, task := range db.tasks {
		if task.ID == t.ID {
			task.Name = t.Name
			task.Completed = t.Completed
			db.tasks[i] = task
			return nil
		}
	}

	return errors.New("could not find task to edit")
}
