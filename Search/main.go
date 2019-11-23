
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	//"io/ioutil"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"


	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
)

var collection *mongo.Collection

var redisclient *redis.Client

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Home")
}

func updateCache(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var key string = params["city"]

    //get from redis cache
    var props, err = redisclient.Get(key).Result()
    if err != nil {
      fmt.Println("City not found.")
     json.NewEncoder(w).Encode(err)
      return
    }

    //delete entry from redis cache
    err = redisclient.Del(key).Err()
    if err != nil {
      fmt.Println(err)
      json.NewEncoder(w).Encode(err)
      return
    }

    json.NewEncoder(w).Encode(props)
}

func getPropertiesbyCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//get from redis cache
	var key string = params["city"]
	var props, err = redisclient.Get(key).Result()
	var results []Property

	if err != nil {
    	   
	//fetch from db
	cur, err1 := collection.Find(context.Background(), bson.M{"City": params["city"]})
	if err1 != nil {
		 fmt.Println("City not found.")
		return
	}

	for cur.Next(context.Background()) {
		var result Property
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
	results = append(results, result)
	}
	
	//insert into redis cache
    	value := results
	err = redisclient.Set(key, value, 0).Err()
	} else {
		json.Unmarshal([]byte(props), &results)
	}

	json.NewEncoder(w).Encode(results)
}

func initDb() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@10.0.1.147:27017/?replicaSet=cmpe281&connect=direct")
	//const connectionString = "mongodb://admin:admin@10.0.1.147:27017/?replicaSet=cmpe281&connect=direct"

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
	collection = client.Database("Property").Collection("Listings")

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

	initDb()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/updatecache", updateCache).Methods("POST")
	router.HandleFunc("/search", getPropertiesbyCity).Methods("GET")
	if GetRedisServer() {
		log.Fatal(http.ListenAndServe(":8080", router))
	} else {
	fmt.Println("Server startup failed")
	}
}
