package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	r := mux.NewRouter()

	r.HandleFunc("/person", createPerson).Methods("POST")
	r.HandleFunc("/person/{id}", getPerson).Methods("GET")
	r.HandleFunc("/people", getPeople).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	personCollection := client.Database("MyFirstDatabase").Collection("People")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, _ := personCollection.InsertOne(ctx, person)
	json.NewEncoder(w).Encode(result)

}

func getPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var people []Person
	personCollection := client.Database("MyFirstDatabase").Collection("People")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, cancel := personCollection.Find(ctx, bson.M{})
	if cancel != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + cancel.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + cancel.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(people)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	personCollection := client.Database("MyFirstDatabase").Collection("People")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := personCollection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(person)
}
