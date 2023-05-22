package manager

import (
	"bufio"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
	"strings"
	"taskManager/db"
)

var input = bufio.NewReader(os.Stdin)

type Task struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ID          int                `bson:"task_id" json:"task_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      bool               `bson:"status" json:"status"`
}

func NewTask(title string, desc string) (tsk *Task) {
	db.Counter++
	return &Task{primitive.NewObjectID(), db.Counter, title, desc, false}
}

func (tsk *Task) printTask() {
	status := "Pending"
	if tsk.Status {
		status = "Completed"
	}

	fmt.Printf(
		"Task ID: %d\n"+
			"Title: %s\n"+
			"Description: %s\n"+
			"Status: %s\n", tsk.ID, tsk.Title, tsk.Description, status,
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

	tsk := NewTask(strings.TrimSpace(title), strings.TrimSpace(desc))
	db.InsertItem(tsk)

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

	_, err = db.RemoveItemWithID(taskId)

	if err != nil {
		fmt.Printf("Could not delete item with id %d from the database", taskId)
		return
	}

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
	_, err = db.CompleteTaskWithId(taskId)

	if err == nil {
		fmt.Println("Task marked completed successfully!")
		return
	}

	fmt.Println("An error occurred because of your incompetence.")
}

func ViewTasks() {
	cur := db.ReturnAllTasks()

	for cur.Next(context.TODO()) {
		var task Task
		err := cur.Decode(&task)
		if err != nil {
			fmt.Println("You don't deserve this.")
			panic(err)
		}

		task.printTask()
		fmt.Println()
	}
}
