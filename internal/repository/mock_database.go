package repository

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
	// TODO Will be implemented on c.iv.3
	var task Task
	return task, nil
}

func (db *MockDatabase) GetTasksByCompletion(isCompleted bool) ([]Task, error) {
	// TODO Will be implemented on c.iv.2
	return nil, nil
}

func (db *MockDatabase) AddTask(t Task) (int64, error) {
	// TODO Will be implemented on c.iv.4
	return 0, nil
}

func (db *MockDatabase) EditTask(t Task) error {
	// TODO Will be implemented on c.iv.5
	return nil
}
