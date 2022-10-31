package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(""))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// Connect
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	
	quickstartDatabase := client.Database("quickstart")
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	episodesCollection := quickstartDatabase.Collection("episodes")

	id, _ := primitive.ObjectIDFromHex("6101a2583f41ce359013adbf")
	result, err := podcastsCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"author", "Nic Raboy"}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)

	// UpdateMany
	result, err = podcastsCollection.UpdateMany(
		ctx,
		bson.M{"title": "The Polyglot Developer Podcast"},
		bson.D{
			{"$set", bson.D{{"author", "Nicolas Raboy"}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)

	//ReplaceOne
	result, err = podcastsCollection.ReplaceOne(
		ctx,
		bson.M{"author": "Nic Raboy"},
		bson.M{
			"title":  "The Nic Raboy Show",
			"author": "Nicolas Raboy",
		},
	)
	fmt.Printf("Replaced %v Documents!\n", result.ModifiedCount)

	
}