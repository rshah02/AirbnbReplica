package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"net/http"
	"../models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3" 
        "github.com/aws/aws-sdk-go/service/s3/s3manager"   
)

// DB connection string for localhost mongoDB

const connectionString = "mongodb+srv://admin:admin@cluster0-uvain.mongodb.net/test?retryWrites=true&w=majority"

// Database Name
const dbName = "test"

// Collection name
const collName = "listing"

// collection object/instance
var collection *mongo.Collection

//s3 bucket name
const bucket = "listingpcs"


// create connection with mongo db
func init() {

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
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

	collection = client.Database(dbName).Collection(collName)

	fmt.Println("Collection instance created!")
}

// CreateProperty create property route
func CreateProperty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var task models.Property
	_ = json.NewDecoder(r.Body).Decode(&task)

	fmt.Println("I have decoded the request", task.Image)
	file, err := os.Open(task.Image)
	if err != nil {
	    fmt.Printf("Unable to open file %q, %v", err)
	}
	defer file.Close()

	insertOneListing(task)		//to insert the record into MongoDB
	//upload the image to S3 bucket "listingpcs"
	sess, err := session.NewSession(&aws.Config{
	Region: aws.String("us-east-1")},
	)
	uploader := s3manager.NewUploader(sess)
	
	_, err = uploader.Upload(&s3manager.UploadInput{
	Bucket: aws.String(bucket),
    	Key: aws.String(task.Image),
    	Body: file,
	})
	if err != nil {
	    // 	Print the error and exit.
	    fmt.Printf("Unable to upload %q to %q, %v", task.Image, bucket, err)
	}

	
	fmt.Printf("Successfully uploaded %q to %q\n", task.Image, bucket)
	
	json.NewEncoder(w).Encode(task)
}

func insertOneListing(task models.Property) {

	insertResult, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
}


// GetAllProperty get all the property route
func GetAllProperty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getAllProperty()
	json.NewEncoder(w).Encode(payload)
}

func getAllProperty() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
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
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)
		fmt.Println("\n \n")
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

//Update a single property


func UpdateProperty(w http.ResponseWriter, r *http.Request) {
	personID := mux.Vars(r)["id"]
	var listing models.Property
	_ = json.NewDecoder(r.Body).Decode(&listing)
	update(listing, personID)

}

func update(listing models.Property, personID string) {
	filter := bson.M{"propertyid": personID}
	update := bson.M{"$set": bson.M{"title": listing.Title, "price" : listing.Price,"description" : listing.Description},}
	result := collection.FindOneAndUpdate(context.Background(), filter, update)
	
	fmt.Println(result)
}


//Delete a single property

func DeleteProperty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	params := mux.Vars(r)["id"]
	var listing1 models.Property
	var task models.Property
	task = listing1
	filter1 := bson.M{"propertyid": params}
	err := collection.FindOne(context.Background(), filter1).Decode(&listing1)
	
//	fmt.Println("Task image is : ", task.Image)
	
//	fmt.Printf("value of listing1.Image is : %
	sess, err := session.NewSession(&aws.Config{
	Region: aws.String("us-east-1")},
	)

	svc := s3.New(sess)
	fmt.Println("Task image is : ", task.Image)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
	    Bucket: aws.String(bucket),
	    Key: aws.String(listing1.Image),
	})
	if err != nil {
    	// Print the error and exit.
    	fmt.Printf("Unable to delete %q from %q, %v", listing1.Image, bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
	    Bucket: aws.String(bucket),
	    Key:    aws.String(listing1.Image),
	})

	fmt.Printf("Successfully deleted %q from %q\n", listing1.Image, bucket)

	_ = json.NewDecoder(r.Body).Decode(&listing1)
	deleteOneTask(listing1, params)	
	
	/*file, err := os.Open(listing1.Image)
	if err != nil {
	    fmt.Printf("Unable to open file %q, %v", err)
	}
	defer file.Close()
*/

}

// delete one task from the DB, delete by ID
func deleteOneTask(listing1 models.Property, params string) {
	
	filter := bson.M{"propertyid": params}
	result, err := collection.DeleteOne(context.Background(), filter)
	
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}


//Get a single record

func GetProperty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	
	params := mux.Vars(r)["id"]

	var task models.Property

//	_ = json.NewDecoder(r.Body).Decode(&task)

//	fmt.Println("get id is : %q ", getid)
	filter := bson.M{"propertyid": params}

	err := collection.FindOne(context.Background(), filter).Decode(&task)
	fmt.Println("Task image is : %q ", task.Image)
	if err != nil {
		fmt.Println(err)
		return
	} else {
	// Print out data from the document result
		json.NewEncoder(w).Encode(task)
	}

}

func GetManyProperty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


	params := mux.Vars(r)["user"]

	var listing models.Property
		_ = json.NewDecoder(r.Body).Decode(&listing)

	filter := bson.M{"username": params}

	// find all documents
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// iterate through all documents
	for cursor.Next(context.Background()) {
	    var p models.Property
	    // decode the document
	    if err := cursor.Decode(&p); err != nil {
	    	log.Fatal(err)
	    }
	    json.NewEncoder(w).Encode(p)
	}

	// check if the cursor encountered any errors while iterating 
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	

}


