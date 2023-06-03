package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"taskManager/logger"
)

var Client *mongo.Client
var Collection *mongo.Collection

func Connect(_URI string) {
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(_URI))

	if err != nil {
		fmt.Println("Failed to connect to the database:", err.Error())
		panic(err)
	}
	Client = cli
	Collection = cli.Database("taskManager").Collection("tasks")
}

func Disconnect() {
	err := Client.Disconnect(context.TODO())
	if err != nil {
		logger.Error("Failed to disconnect from the database: %s", err.Error())
		panic(err)
	}
}

func InsertItem(it interface{}) (*mongo.InsertOneResult, error) {
	return Collection.InsertOne(context.TODO(), it)
}

func RemoveItemWithID(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return Collection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
}

func CompleteTaskWithId(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	return Collection.UpdateOne(context.TODO(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"status", true}}}})
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
	}

	cur, err := Collection.Find(context.TODO(), filter)
	if err != nil {
		logger.Error("You don't deserve the tasks")
		return nil
	}

	return cur
}
