package main

import (
	internal "challange/internal/repository"
	"fmt"
)

func main() {
	repo := internal.NewDatabase()

	allTasks, _ := repo.GetAllTasks()
	for _, task := range allTasks {
		fmt.Printf("Task %v: %q completion is %t.\n", task.ID, task.Name, task.Completed)
	}

	currentTask, _ := repo.GetTaskByID(6)
	fmt.Printf("Current task is: %q.\n", currentTask.Name)

	newTask := internal.Task{
		Name:      "Checkout the real project",
		Completed: false,
	}
	newID, _ := repo.AddTask(newTask)
	fmt.Printf("Added new task at the end, id is %v.\n", newID)

	taskToEdit := internal.Task{
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
