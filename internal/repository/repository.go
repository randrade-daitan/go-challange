package internal

type Repository interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id int64) (Task, error)
	GetTaskByCompletion(isCompleted bool) (Task, error)

	AddTask(t Task) (int64, error)
	EditTask(t Task) error
}
