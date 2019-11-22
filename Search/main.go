
package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"


	"github.com/gorilla/mux"
)

var collection *mongo.Collection

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Home")
}

func updateCache(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(req)
    var key string = params["city"]

    //get from redis cache
    var props, err = client.Get(key).Result()
    if err != nil {
      fmt.Println("City not found.")
      formatter.JSON(w, http.StatusNotFound, err)
      return
    }

    //delete entry from redis cache
    err = client.Del(key).Err()
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusInternalServerError, err)
      return
    }

     formatter.JSON(w, http.StatusOK)
    }
}

func getPropertiesbyCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//get from redis cache
	var key string = params["city"]
	var props, err = client.Get(key).Result()
	var results []models.Property

	if err != nil {
    	   
	//fetch from db
	err1 := collection.Find(context.Background(), bson.M{"City": params["city"]})
	if err1 != nil {
		 fmt.Println("City not found.")
		return
	}

	for cur.Next(context.Background()) {
		var result models.Property
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
	results = append(results, result)
	}
	
	//insert into redis cache
    	value, _ := results
	err = redisclient.Set(key, value, 0).Err()
	}
	else {
		results := props
	}

	json.NewEncoder(w).Encode(results)
}

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	collection := client.Database("Property").Collection("Listings")

}

func GetRedisServer() bool {
	fmt.Println("Connecting to Redis server..")
  	redisclient = NewRedisServer()

  	fmt.Println("PING")
  	pong, err := redisclient.Ping().Result()
  	if err != nil {
    	fmt.Println("Could not connect to Redis server:", err)
    	return false
  	}

  	fmt.Println(pong)
  	return true
}


func main() {

	init()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/updatecache", updateCache).Methods("POST")
	router.HandleFunc("/search", getPropertiesbyCity).Methods("GET")
	if GetRedisServer() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}
	else {
	fmt.Println("Server startup failed")
	}
}
