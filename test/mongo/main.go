package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/rtntubmt97/profiler/internal/defines"
)

func main() {
	var err error
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	testDb := client.Database("test")
	profileCollection := testDb.Collection("profile")

	var profile defines.Profile
	err = profileCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&profile)

	fmt.Println(profile)
}
