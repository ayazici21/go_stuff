package main

import (
	"fmt"
	"taskManager/db"
	"taskManager/manager"
)

const _URI = "mongodb://localhost:27017"

func main() {
	db.Connect(_URI)

	defer db.Disconnect()

	for true {
		choice := manager.DisplayMenu()
		switch choice {
		case 0:
			manager.AddTask()
			break
		case 1:
			manager.ViewTasks()
			break
		case 2:
			manager.CompleteTask()
			break
		case 3:
			manager.DeleteTask()
			break
		case 4:
			manager.UseFilter()
			break
		case 5:
			fmt.Println("Thank you for using the Task Manager. Goodbye!")
			return
		}
	}
}
