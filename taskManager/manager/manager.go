package manager

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"taskManager/db"
	"taskManager/util"
)

var filter = 1

type Task struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ID          int                `bson:"task_id" json:"task_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      bool               `bson:"status" json:"status"`
}

func NewTask(title string, desc string) *Task {
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

func AddTask() {
	title, err := util.GetString("Enter the title of the task: ")

	if err != nil {
		fmt.Println(err)
		return
	}

	desc, err := util.GetString("Enter the description of the task: ")

	if err != nil {
		fmt.Println(err)
		return
	}

	tsk := NewTask(title, desc)
	_, err = db.InsertItem(tsk)

	if err != nil {
		fmt.Println("Something happened.")
		return
	}

	fmt.Println("Task added successfully.")
}

func DeleteTask() {
	taskId, err := util.GetInt("Which task do you want to delete? ")

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
	taskId, err := util.GetInt("Which task do you want to mark completed? ")

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
	cur := db.ReturnAllTasks(filter)

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

func UseFilter() {
	for {
		f, err := util.GetInt(`Filter Tasks by Status:
1: No filter
2: Pending
3: Completed

Enter your choice: `)

		if err != nil {
			fmt.Println("You won't break me.")
			continue
		}

		filter = f
		break
	}

}
