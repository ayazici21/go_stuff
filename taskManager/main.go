package main

import (
	"fmt"
	"taskManager/db"
	"taskManager/manager"
	"taskManager/util"
)

const _URI = "mongodb://localhost:27017"

func main() {
	db.Connect(_URI)

	defer db.Disconnect()

	for true {
		choice := displayMenu()

		switch choice {
		case 1:
			manager.AddTask()
			break
		case 2:
			manager.ViewTasks()
			break
		case 3:
			manager.CompleteTask()
			break
		case 4:
			manager.DeleteTask()
			break
		case 5:
			manager.UseFilter()
			break
		case 6:
			fmt.Println("Thank you for using the Task Manager. Goodbye!")
			return
		}
	}
}

func displayMenu() int {
	for {
		num, err := util.GetInt(`Menu:
1. Add a Task
2. View Tasks
3. Mark a Task as Completed
4. Delete a Task
5. Filter Tasks by Status
6. Exit

Enter your choice: `)

		if err == nil {
			if num >= 1 && num <= 6 {
				return num
			}
		}
		fmt.Println("Please enter one of the options.")
	}
}
