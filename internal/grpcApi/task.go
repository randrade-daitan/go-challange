package grpcApi

import "challange/internal/repository"

func convertDataTask(t *Task) repository.Task {
	return repository.Task{
		ID:        t.Id,
		Name:      t.Name,
		Completed: t.Completed,
	}
}

func convertRepoTask(t repository.Task) *Task {
	return &Task{
		Id:        t.ID,
		Name:      t.Name,
		Completed: t.Completed,
	}
}

func convertRepoTasks(t []repository.Task) *Tasks {
	var tasks []*Task
	for _, rt := range t {
		newTask := convertRepoTask(rt)
		tasks = append(tasks, newTask)
	}
	return &Tasks{
		Task: tasks,
	}
}
