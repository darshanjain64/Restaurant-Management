package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://darshanr94dj:Hello%402024@cluster0.o9oelrd.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
  
	// Create a new client and connect to the server
	client, err := mongo.NewClient(opts)
	if err != nil {
	  log.Fatal(err)
	}
  
    ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)

    defer cancel()

	err=client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	  }
   


	
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}


  var Client *mongo.Client = DBinstance()

  func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
    var collection *mongo.Collection = client.Database("restauarnt").Collection(collectionName)

	return collection
  }