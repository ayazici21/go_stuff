package manager

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
		logger.Error("Body parser: %s", err.Error())
		return c.Status(400).SendString("Could not parse body")
	}

	tsk := NewTask(body.Title, body.Description)
	_, err = db.InsertItem(tsk)

	if err != nil {
		logger.Error("Insert item: %s", err)
		return c.Status(400).SendString("Could not insert the item.")
	}

	logger.Info("Task '%s' insert success.", tsk.ID)
	return c.Status(200).JSON(tsk) // send task so it can be easily added in the UI
}

func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		logger.Error("Invalid ID: %s", id)
		return c.Status(400).SendString("Received an invalid ObjectID")
	}

	_, err = db.RemoveItemWithID(objID)

	if err != nil {
		logger.Error("Remove item: %s", err.Error())
		return c.Status(400).SendString("Could not delete the item.")
	}

	logger.Info("Task '%s' remove success.", objID.Hex())
	return c.Status(200).SendString("Task successfully deleted")
}

func CompleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Invalid ID: %s", id)
		return c.Status(400).SendString("Received an invalid ObjectID")
	}
	_, err = db.CompleteTaskWithId(objID)

	if err == nil {
		logger.Info("Task '%s' mark success.", objID.Hex())
		return c.Status(200).SendString("Task successfully marked completed.")
	}

	logger.Error("Complete task: %s", err.Error())
	return c.Status(400).SendString("Could not mark the task completed.")
}

func ViewTasks(c *fiber.Ctx) error {
	cur, err := db.ReturnAllTasks(filter)
	if err != nil {
		logger.Error("Return tasks: %s", err.Error())
		return c.Status(400).SendString("An error occurred.")
	}

	tasks := make([]Task, 0)
	for cur.Next(context.TODO()) {
		var task Task
		err := cur.Decode(&task)
		if err != nil {
			logger.Error("Decode: %s", err.Error())
			return c.Status(400).SendString("An error occurred.")
		}

		tasks = append(tasks, task)
	}
	return c.Status(200).JSON(tasks)
}

func UseFilter(c *fiber.Ctx) error {
	f, err := strconv.Atoi(c.Params("filter"))

	if err != nil {
		logger.Error("An error occurred: %s", err.Error())
		return c.Status(400).SendString("An error occurred.")
	}

	filter = f
	logger.Info("Switched filter to %d", filter)
	return c.Status(200).SendString("Successfully switched filter.")
}

func GetFilter(c *fiber.Ctx) error {
	return c.Status(200).JSON(bson.M{"filter": filter})
}
