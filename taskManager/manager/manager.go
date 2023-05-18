package manager

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var input = bufio.NewReader(os.Stdin)
var curId = 0

type task struct {
	id          int
	title       string
	description string
	status      bool
}

var tasks = make(map[int]task)

func printTask(tsk task) {
	status := "Pending"
	if tsk.status {
		status = "Completed"
	}

	fmt.Printf(
		"Task ID: %d\n"+
			"Title: %s\n"+
			"Description: %s\n"+
			"Status: %s\n", tsk.id, tsk.title, tsk.description, status,
	)
}

func DisplayMenu() int {
	fmt.Print("Menu:\n" +
		"1. Add a Task\n" +
		"2. View Tasks\n" +
		"3. Mark a Task as Completed\n" +
		"4. Delete a Task\n" +
		"5. Exit\n\n" +
		"Enter your choice: ",
	)

	line, err := input.ReadString('\n')
	if err == nil {
		line = strings.TrimSpace(line)
		if len(line) != 1 {
			fmt.Println("Please enter one of the options.")
			return DisplayMenu()
		}
		res := strings.Index("12345", line)
		if res != -1 {
			return res
		}
		fmt.Println("Please enter one of the options.")
		return DisplayMenu()
	}

	fmt.Println("You cannot escape your fate. Please enter one of the options.")
	return DisplayMenu()
}

func AddTask() {
	fmt.Print("Enter the title of the task: ")
	title, err := input.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Enter the description of the task: ")
	desc, err := input.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	curId++
	tasks[curId] = task{curId, strings.TrimSpace(title),
		strings.TrimSpace(desc), false}

	fmt.Println("Task added successfully.\n")
}

func DeleteTask() {
	fmt.Println("Which task do you want to delete? ")
	s_taskId, err := input.ReadString('\n')
	if err != nil {
		fmt.Println("Die")
		return
	}
	taskId, err := strconv.Atoi(strings.TrimSpace(s_taskId))
	if err != nil {
		fmt.Println("Die")
		return
	}
	_, ok := tasks[taskId]
	if !ok {
		fmt.Println("Die")
		return
	}

	delete(tasks, taskId)
	fmt.Println("Task deleted successfully!")
}

func CompleteTask() {
	fmt.Println("Which task do you want to complete? ")
	s_taskId, err := input.ReadString('\n')
	if err != nil {
		fmt.Println("Die")
		return
	}
	taskId, err := strconv.Atoi(strings.TrimSpace(s_taskId))
	if err != nil {
		fmt.Println("Die")
		return
	}

	task, ok := tasks[taskId]
	if ok {
		task.status = true
		tasks[taskId] = task
		fmt.Println("Task marked completed successfully!")
		return
	}

	fmt.Println("An error occurred because of your incompetence.")
}

func ViewTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found!")
		return
	}

	fmt.Println("Tasks:\n")
	for _, task := range tasks {
		printTask(task)
	}
	fmt.Println()
}
