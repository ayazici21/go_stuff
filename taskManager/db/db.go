package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection
var Counter int

type DB struct {
	Client     *mongo.Client
	Database   *mongo.Database
	EntryCount int
}

func Connect(_URI string) {
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(_URI))

	if err != nil {
		fmt.Println("Failed to connect to the database:", err.Error())
		panic(err)
	}
	Client = cli
	Collection = cli.Database("taskManager").Collection("tasks")

	pipeline := []bson.M{ // stolen fr fr
		{
			"$group": bson.M{
				"_id": "",
				"max": bson.M{"$max": "$task_id"},
			},
		},
	}

	cur, err := Collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		fmt.Println("I'm too tired to handle errors.")
		panic(err)
	}
	var res []bson.M
	err = cur.All(context.TODO(), &res)
	if err != nil {
		fmt.Println("I'm too tired to handle errors.")
		panic(err)
	}
	if len(res) == 0 {
		Counter = 0
	} else {
		Counter = int(res[0]["max"].(int32)) // this is int32 not int 🤓👆
	}
}

func Disconnect() {
	err := Client.Disconnect(context.TODO())
	if err != nil {
		fmt.Println("Failed to disconnect from the database:", err.Error())
		panic(err)
	}
}

func InsertItem(doc interface{}) *mongo.InsertOneResult {
	res, err := Collection.InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Println("Could not insert item to the database.")
		return nil
	}
	return res
}

func RemoveItemWithID(id int) (*mongo.DeleteResult, error) {
	return Collection.DeleteOne(context.TODO(), bson.D{{"task_id", id}})
}

func CompleteTaskWithId(id int) (*mongo.UpdateResult, error) {
	return Collection.UpdateOne(context.TODO(), bson.D{{"task_id", id}}, bson.D{{"$set", bson.D{{"status", true}}}})
}

func ReturnAllTasks(f int) *mongo.Cursor {
	var filter bson.D
	switch f {
	case 1:
		filter = bson.D{{}}
		break
	case 2:
		filter = bson.D{{"status", false}}
		break
	case 3:
		filter = bson.D{{"status", true}}
		break

	}

	cur, err := Collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("You don't deserve the tasks")
		panic(err)
	}

	return cur
}
