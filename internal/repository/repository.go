package repository

type Repository interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id int64) (Task, error)
	GetTasksByCompletion(isCompleted bool) ([]Task, error)

	AddTask(t Task) (int64, error)
	EditTask(t Task) error
}
