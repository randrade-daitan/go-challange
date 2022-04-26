package repository

// A storage for tasks
type Repository interface {

	// Fetches all tasks stored.
	GetAllTasks() ([]Task, error)

	// Try to fetch a task given an id.
	GetTaskByID(id int64) (Task, error)

	// Fetches all tasks by its completion status.
	GetTasksByCompletion(isCompleted bool) ([]Task, error)

	// Adds a new task with value passed through the task parameter.
	// Returns the id of the newly created task.
	AddTask(task Task) (int64, error)

	// Edits a task from the value passed through the task parameter.
	EditTask(task Task) error
}
