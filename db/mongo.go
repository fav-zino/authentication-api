package db

import (
	"authentication_api/config"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection

func ConnectToDB() error {

	clientOptions := options.Client().ApplyURI(config.AppConfig.DatabaseURL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	UserCollection = client.Database("test").Collection("users")

	index := mongo.IndexModel{
		Keys: bson.M{
			"email": 1, // Create an index on the email field for faster lookups(arrange in ascending order)
		},
		Options: options.Index().SetUnique(true), //make field unique
	}
	_, err = UserCollection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return err
	}

	return nil

}
