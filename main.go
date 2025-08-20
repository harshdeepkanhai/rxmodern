package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=rs0"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Println(err)
		}
	}()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to MongoDB")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "RxModern is running")
	})
	log.Println("listening on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
