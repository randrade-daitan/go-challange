package repository

type Task struct {
	ID        int64
	Name      string
	Completed bool
}

func (Task) TableName() string {
	return "task"
}
