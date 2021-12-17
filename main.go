package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var client *mongo.Client

func main() {
	fmt.Println("Application started successfully!")
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	r := mux.NewRouter()

	r.HandleFunc("/person", createPerson).Methods("POST")

	http.ListenAndServe(":8080", r)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	personCollection := client.Database("MyFirstDatabase").Collection("People")
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	result, _ := personCollection.InsertOne(ctx, person)
	json.NewEncoder(w).Encode(result)

}
