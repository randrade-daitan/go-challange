package api

import (
	"challange/internal/api/proto"
	"challange/internal/repository"
)

func convertDataTask(t *proto.Task) repository.Task {
	return repository.Task{
		ID:        t.Id,
		Name:      t.Name,
		Completed: t.Completed,
	}
}

func convertRepoTask(t repository.Task) *proto.Task {
	return &proto.Task{
		Id:        t.ID,
		Name:      t.Name,
		Completed: t.Completed,
	}
}

func convertRepoTasks(t []repository.Task) *proto.Tasks {
	var tasks []*proto.Task
	for _, rt := range t {
		newTask := convertRepoTask(rt)
		tasks = append(tasks, newTask)
	}
	return &proto.Tasks{
		Task: tasks,
	}
}
