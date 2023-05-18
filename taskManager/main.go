package main

import (
	"fmt"
	"manager"
)

func main() {
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
			fmt.Println("Thank you for using the Task Manager. Goodbye!")
			return
		}
	}
}
