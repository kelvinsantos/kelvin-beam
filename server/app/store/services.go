package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	DatabaseName      = "beam"
	ScooterCollection = "scooter"
)

func getDB(client *mongo.Client) *mongo.Database {
	return client.Database(DatabaseName)
}

func getNow() primitive.DateTime {
	return primitive.DateTime(time.Now().UnixNano() / 1e6)
}

func GetScooterListCommand(client *mongo.Client) ([]Scooter, error) {

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var dEntries []Scooter
	query := bson.D{}
	cursor, err := getDB(client).Collection(ScooterCollection).Find(ctx, query)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var entry Scooter
		_ = cursor.Decode(&entry)
		dEntries = append(dEntries, entry)
	}

	return dEntries, nil
}
