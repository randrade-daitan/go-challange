package repository

type repository interface {
	getAllTasks() ([]task, error)
	getTaskByID(id int64) (task, error)

	addTask(t task) (int64, error)
	editTask(t task) error
	deleteTask(t task) error
}
