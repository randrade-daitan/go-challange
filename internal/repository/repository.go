package repository

type Repository interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id int64) (Task, error)
	GetTasksByCompletion(isCompleted bool) ([]Task, error)

	AddTask(task Task) (int64, error)
	EditTask(task Task) error
}
