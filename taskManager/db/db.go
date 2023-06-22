package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"taskManager/logger"
)

var Client *mongo.Client
var TasksCollection *mongo.Collection
var UsersCollection *mongo.Collection
var log = logger.New()

func Connect(_URI string) {
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(_URI))

	if err != nil {
		log.Error("Failed to connect to the database: %s", err.Error())
		panic(err)
	}
	Client = cli
	TasksCollection = cli.Database("taskManager").Collection("tasks")
	UsersCollection = cli.Database("taskManager").Collection("users")
}

func Disconnect() {
	err := Client.Disconnect(context.TODO())
	if err != nil {
		log.Error("Failed to disconnect from the database: %s", err.Error())
		panic(err)
	}
}

func InsertItem(collection *mongo.Collection, it interface{}) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(context.TODO(), it)
}

func RemoveItemWithID(collection *mongo.Collection, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return collection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
}

func CompleteTaskWithId(collection *mongo.Collection, id primitive.ObjectID) (*mongo.UpdateResult, error) {
	return collection.UpdateOne(context.TODO(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"status", true}}}})
}

func ReturnAllItems(collection *mongo.Collection, f ...bson.D) (*mongo.Cursor, error) {
	if len(f) > 1 {
		return nil, errors.New(`"ReturnAllItems" must be called with 1 filter at most`)
	}

	var filter bson.D

	if len(f) == 1 {
		filter = f[0]
	} else {
		filter = bson.D{{}}
	}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error("You don't deserve the tasks")
		return nil, err
	}

	return cur, nil
}

func GetItemFromId(collection *mongo.Collection, ID primitive.ObjectID) *mongo.SingleResult {
	return collection.FindOne(context.TODO(), bson.D{{"_id", ID}})

}
func SearchUserFromEmail(email string) *mongo.SingleResult {
	return UsersCollection.FindOne(context.TODO(), bson.D{{"email", email}})
}

func SearchUserFromUsername(username string) *mongo.SingleResult {
	return UsersCollection.FindOne(context.TODO(), bson.D{{"username", username}})
}

func ValidateToken(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	return UsersCollection.UpdateOne(context.TODO(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"token_valid", true}}}})
}

func InvalidateToken(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	return UsersCollection.UpdateOne(context.TODO(), bson.D{{"_id", id}},
		bson.D{{"$set", bson.D{{"token_valid", false}}}})
}
