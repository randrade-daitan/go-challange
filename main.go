package main

import (
	internal "challange/internal/repository"
	"fmt"
)

func main() {
	repo := internal.NewDatabase()

	allTasks, _ := repo.GetAllTasks()
	for _, task := range allTasks {
		fmt.Printf("Task %q completion is %t", task.Name, task.Completed)
	}
}
