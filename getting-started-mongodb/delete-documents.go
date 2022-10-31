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

	result, err := podcastsCollection.DeleteOne(ctx, bson.M{"title": "The Polyglot Developer Podcast"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
	
	result, err = episodesCollection.DeleteMany(ctx, bson.M{"duration": 25})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)

	if err = podcastsCollection.Drop(ctx); err != nil {
		log.Fatal(err)
	}
}