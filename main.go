package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://monsterr:0JYrw9FGAM6ybvNk@cluster0.5cgyh.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	quickstartDatabase := client.Database("quickstart")
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	episodesCollection := quickstartDatabase.Collection("episodes")

	podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: "The Polyglot Developer Podcast"},
		{Key: "author", Value: "Nic Raboy"},
		{Key: "tags", Value: bson.A{"development", "programming", "coding"}},
	})

	episodeResult, err := episodesCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title", "GraphQL for API Development"},
			{"description", "Learn about GraphQL from the co-creator of GraphQL, Lee Byron."},
			{"duration", 25},
		},
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title", "Progressive Web Application Development"},
			{"description", "Learn about PWA development with Tara Manicsic."},
			{"duration", 32},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted %v documents into episode collection!\n", len(episodeResult.InsertedIDs))
}