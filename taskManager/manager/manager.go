package manager

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"taskManager/db"
	"taskManager/logger"
)

var filter = 1

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      bool               `bson:"status" json:"status"`
}

type taskBody struct {
	Title       string `json:"title" xml:"title" form:"title"`
	Description string `json:"description" xml:"description" form:"description"`
}

func NewTask(title string, desc string) *Task {
	return &Task{primitive.NewObjectID(), title, desc, false}
}

func AddTask(c *fiber.Ctx) error {
	var body taskBody
	err := c.BodyParser(&body)
	if err != nil {
		logger.Error("Body parser returned an error: %s", err.Error())
		return c.Status(400).SendString("Could not parse body")
	}

	tsk := NewTask(body.Title, body.Description)
	_, err = db.InsertItem(tsk)

	if err != nil {
		fmt.Println("Something happened.")
		return c.Status(401).SendString("Could not insert the item.")
	}

	fmt.Println("Task added successfully.")
	return c.Status(200).SendString("Task successfully added")
}

func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("%s is not a valid object ID.", id)
		return c.Status(402).SendString("Received an invalid ObjectID")
	}

	_, err = db.RemoveItemWithID(objID)

	if err != nil {
		logger.Error("Could not delete item with id %s from the database", objID.Hex())
		return c.Status(401).SendString("Could not delete the item.")
	}

	logger.Info("Task with id %s has been deleted.", objID.Hex())
	return c.Status(200).SendString("Task successfully deleted")
}

func CompleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("%s is not a valid object ID.", id)
		return c.Status(402).SendString("Received an invalid ObjectID")
	}
	_, err = db.CompleteTaskWithId(objID)

	if err == nil {
		logger.Info("%s marked completed.", objID.Hex())
		return c.Status(200).SendString("Task successfully marked completed.")
	}

	logger.Error("An error occurred because of your incompetence.")
	return c.Status(401).SendString("Skill issue tbh")
}

func ViewTasks(c *fiber.Ctx) error {
	cur := db.ReturnAllTasks(filter)
	tasks := make([]Task, 0)
	for cur.Next(context.TODO()) {
		var task Task
		err := cur.Decode(&task)
		if err != nil {
			return c.Status(403).SendString("An error occurred")
		}

		tasks = append(tasks, task)
	}
	return c.Status(200).JSON(tasks)
}

func UseFilter(c *fiber.Ctx) error {
	f, err := strconv.Atoi(c.Params("filter"))

	if err != nil {
		logger.Error("An error occurred: %s", err.Error())
		return c.Status(403).SendString("An error occurred")
	}

	filter = f
	return c.Status(200).SendString("Successfully switched filter")
}
