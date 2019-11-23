package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	h "github.com/gorilla/handlers"
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
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")        

	params := mux.Vars(r)["city"]	

	filter := bson.M{"city": params}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
	results = append(results, result)
}
	json.NewEncoder(w).Encode(results)


}
 /*       
        //var p Property
        //_ = json.NewDecoder(r.Body).Decode(&p)
//        var prop []Property
        
        //get from redis cache
      //  var key string = params["city"]
        
        //log.Print("cty: %v\n", key)
        
  //      var props, err = redisclient.Get(key).Result()
        
    //    log.Print(props)
        //var results []Property
        
       // if err != nil {
           
        filter := bson.M{"City": params["city"]}           
        //fetch from db
        cur, err1 := collection.Find(context.Background(), filter)
        


        log.Print(cur)
        //var prop Property
        //err1 := collection.Find(context.Background(), filter).All(&prop)
        //log.Print("prop taken: %v\n", prop.City)
        if err1 != nil {
                 fmt.Println(err1)
                return
        }
       var results []primitive.M
 
        for cur.Next(context.Background()) {
                var result bson.M
                e := cur.Decode(&result)
                if e != nil {
                        log.Fatal(e)
                }
        results = append(results, result)
        }
        json.NewEncoder(w).Encode(results)
        //insert into redis cache
//        value := results
  //      err = redisclient.Set(key, value, 0).Err()*/
    //    }
        // else {
                //json.Unmarshal([]byte(props), &results)
//	}


func initDb() {
	// Set client options
	//clientOptions := options.Client().ApplyURI("mongodb://admin:admin@10.0.1.147:27017/?replicaSet=cmpe281&connect=direct")
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:admin@cluster0-f9fkk.mongodb.net/test?retryWrites=true&w=majority")
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

/*func GetRedisServer() bool {
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
}*/

func pingHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("ping")
	mapD := map[string]string{"message": "API Working"}
	mapB, _ := json.Marshal(mapD)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(mapB)
}

func main() {

	initDb()
	headersOk := h.AllowedHeaders([]string{"X-Requested-With", "content-type", "Authorization", "Accept-Encoding", "Accept","X-CSRF-Token","Content-Length"})
	originsOk := h.AllowedOrigins([]string{"*"})
	methodsOk := h.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT","OPTIONS"})

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	//router.HandleFunc("/updatecache", updateCache).Methods("POST")
	router.HandleFunc("/search/get/{city}", getPropertiesbyCity).Methods("GET")
	router.HandleFunc("/search/ping", pingHandler).Methods("GET")
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", h.CORS(headersOk, methodsOk, originsOk)(router)))
	/*if GetRedisServer() {
		log.Fatal(http.ListenAndServe(":8080", router))
	} else {
	fmt.Println("Server startup failed")
	}*/
}
