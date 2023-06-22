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

var log = logger.New()

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      bool               `bson:"status" json:"status"`
	Owner       primitive.ObjectID `bson:"owner" json:"owner"`
}

type taskBody struct {
	Title       string `json:"title" xml:"title" form:"title"`
	Description string `json:"description" xml:"description" form:"description"`
}

func NewTask(title string, desc string, ownerID primitive.ObjectID) *Task {
	return &Task{primitive.NewObjectID(), title, desc, false, ownerID}
}

func AddTask(c *fiber.Ctx) error {
	var body taskBody
	err := c.BodyParser(&body)
	if err != nil {
		log.Error("Body parser: %s", err)
		return c.Status(400).JSON(fiber.Map{"message": "Invalid JSON"})
	}
	log.Info(c.Get("ownerID"))
	id := c.Locals("ownerID").(primitive.ObjectID)
	if err != nil {
		log.Error("%s", err)
	}

	tsk := NewTask(body.Title, body.Description, id)
	_, err = db.InsertItem(db.TasksCollection, tsk)

	if err != nil {
		log.Error("Insert item: %s", err)
		return c.Status(400).JSON(fiber.Map{"message": "Could not insert the item."})
	}

	log.Info("Task '%s' insert success.", tsk.ID)
	return c.Status(200).JSON(tsk) // send task so it can be easily added in the UI
}

func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Error("Invalid ID: %s", id)
		return c.Status(400).JSON(fiber.Map{"message": "Received an invalid ObjectID"})
	}

	res, err := validateOwner(objID, c.Locals("ownerID").(primitive.ObjectID))

	if err != nil {
		log.Error("%s", err)
		return c.Status(500).JSON(fiber.Map{"message": "An error occurred"})
	}

	if !res {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden"})
	}

	_, err = db.RemoveItemWithID(db.TasksCollection, objID)

	if err != nil {
		log.Error("Remove item: %s", err)
		return c.Status(400).JSON(fiber.Map{"message": "Could not delete the item."})
	}

	log.Info("Task '%s' removed.", objID.Hex())
	return c.Status(200).JSON(fiber.Map{"message": "Task successfully deleted"})
}

func CompleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error("Invalid ID: %s", id)
		return c.Status(400).JSON(fiber.Map{"message": "Received an invalid ObjectID"})
	}

	res, err := validateOwner(objID, c.Locals("ownerID").(primitive.ObjectID))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "An error occurred"})
	}

	if !res {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden"})
	}

	_, err = db.CompleteTaskWithId(db.TasksCollection, objID)

	if err == nil {
		log.Info("Task '%s' mark success.", objID.Hex())
		return c.Status(200).JSON(fiber.Map{"message": "Task successfully marked completed."})
	}

	log.Error("Complete task: %s", err)
	return c.Status(400).JSON(fiber.Map{"message": "Could not mark the task completed."})
}

func ViewTasks(c *fiber.Ctx) error {
	filter, err := strconv.Atoi(c.Query("filter"))
	if err != nil {
		log.Error("Filter: %s", c.Query("filter"))
		return c.Status(406).JSON(fiber.Map{"message": "Invalid query"})
	}

	id := c.Locals("ownerID")
	var f bson.D

	switch filter {
	case 0:
		f = bson.D{{"owner", id}}
		break
	case 1:
		f = bson.D{{"status", true}, {"owner", id}}
		break
	case 2:
		f = bson.D{{"status", false}, {"owner", id}}
	}

	cur, err := db.ReturnAllItems(db.TasksCollection, f)

	if err != nil {
		log.Error("Return tasks: %s", err)
		return c.Status(500).JSON(fiber.Map{"message": "An error occurred."})
	}

	tasks := make([]Task, 0)
	for cur.Next(context.TODO()) {
		var task Task
		err = cur.Decode(&task)
		if err != nil {
			log.Error("Decode: %s", err)
			return c.Status(500).JSON(fiber.Map{"message": "An error occurred."})
		}
		tasks = append(tasks, task)
	}
	return c.Status(200).JSON(tasks)
}

func validateOwner(id primitive.ObjectID, ownerID primitive.ObjectID) (bool, error) {
	res := db.GetItemFromId(db.TasksCollection, id)
	if res.Err() != nil {
		log.Error("%s", res.Err())
		return false, res.Err()
	}

	var task Task
	err := res.Decode(&task)

	if err != nil {
		log.Error("%s", err)
		return false, err
	} else if task.Owner.Hex() != ownerID.Hex() {
		return false, nil
	}
	return true, nil
}
