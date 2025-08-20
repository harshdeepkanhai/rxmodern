package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Prescription struct {
	ID       string `bson:"_id"`
	Patient  string `bson:"patient"`
	Drug     string `bson:"drug"`
	Quantity int    `bson:"quantity"`
	Version  int    `bson:"version"`
}

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
	prescription := client.Database("rxmodern").Collection("prescriptions")
	page, err := template.ParseFiles("templates/prescription.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var rx Prescription
		err := prescription.FindOne(r.Context(), bson.M{"_id": "rx-1001"}).Decode(&rx)
		if err != nil {
			log.Println(err)
			http.Error(w, "could not load prescription", http.StatusInternalServerError)
			return
		}
		err = page.Execute(w, rx)
		if err != nil {
			log.Println(err)
		}
	})
	log.Println("listening on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
