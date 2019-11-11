package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	//"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Booking struct{
	BookingID string 
	Date string 
	Name string 
}

func insert(w http.ResponseWriter, r *http.Request){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
    log.Fatal(err)
		}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
    log.Fatal(err)
		}

	fmt.Println("Connected to MongoDB!")
	collection := client.Database("test").Collection("bookings")
	a := Booking{"B1", "22-12-2019", "A"}
	insertResult, err := collection.InsertOne(context.TODO(), a) //To insert single docum,ent

	//insertResult, err := collection.InsertMany(context.TODO(), bookings)
	if err != nil {
	    log.Fatal(err)
			}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	fmt.Fprintf(w, "Inserted! check db")
	err = client.Disconnect(context.TODO())

	if err != nil {
    log.Fatal(err)
		}
	fmt.Println("Connection to MongoDB closed.")
}

func welcome(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to Booking! To insert a record change url to localhost:3000/insert and hit :)")
}

func handleRequests(){
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/welcome", welcome)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main(){
	handleRequests()
}