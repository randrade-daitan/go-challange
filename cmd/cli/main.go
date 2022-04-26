package main

import (
	"challange/internal/repository"
	"fmt"
)

func main() {
	repo := repository.NewRepository()

	allTasks, _ := repo.GetAllTasks()
	for _, t := range allTasks {
		fmt.Printf("Task %v: %q completion is %t.\n", t.ID, t.Name, t.Completed)
	}

	currentTask, _ := repo.GetTaskByID(6)
	fmt.Printf("Current task is: %q.\n", currentTask.Name)

	newTask := repository.Task{
		Name:      "Checkout the real project",
		Completed: false,
	}
	newID, _ := repo.AddTask(newTask)
	fmt.Printf("Added new task at the end, id is %v.\n", newID)

	taskToEdit := repository.Task{
		ID:        newID,
		Name:      "Celebrate",
		Completed: false,
	}
	repo.EditTask(taskToEdit)
	editedTask, _ := repo.GetTaskByID(newID)
	fmt.Printf("Edited task %v with new name: %q.\n", editedTask.ID, editedTask.Name)

	completedTasks, _ := repo.GetTasksByCompletion(true)
	fmt.Printf("Already done %v tasks!.\n", len(completedTasks))
}
