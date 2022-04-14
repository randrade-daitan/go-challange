package main

import (
	internal "challange/internal/repository"
	"fmt"
)

func main() {
	repo := internal.NewDatabase()

	allTasks, _ := repo.GetAllTasks()
	for _, task := range allTasks {
		fmt.Printf("Task %q completion is %t.\n", task.Name, task.Completed)
	}

	currentTask, _ := repo.GetTaskByID(6)
	fmt.Printf("Task %q is the current being done.\n", currentTask.Name)
}
