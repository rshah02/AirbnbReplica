package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
//	"time"
	"os"
	"net/http"
	"../models"
	"mime/multipart"
	"bytes"
	"path/filepath"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3" 
       // "github.com/aws/aws-sdk-go/service/s3/s3manager" 
	"github.com/aws/aws-sdk-go/aws/awserr"  
	"github.com/nu7hatch/gouuid"
//	"github.com/aws/aws-sdk-go/aws/credentials"
	//"math/big"
	"math/rand"
	"strconv"

)	

// DB connection string for localhost mongoDB

const connectionString = "mongodb+srv://admin:admin@cluster0-f9fkk.mongodb.net/test?retryWrites=true&w=majority"
//const connectionString = "mongo -u admin -p mongo1234 --authenticationDatabase admin mongodb://primary:27017,secondary1:27017,secondary2:27017/?replicaSet=cmpe281"
//const connectionString = "mongodb://admin:mongo1234@10.0.1.144:27017/?replicaSet=cmpe281&connect=direct"

// Database Name
const dbName = "Property"

// Collection name
const collName = "Listings"

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


func Phandler(w http.ResponseWriter, r *http.Request) {
	
	file, fileHeader, err := r.FormFile("task.Image")
	if err != nil {
		log.Println("err is ", err)
		fmt.Fprintf(w, "Could not get uploaded file")
		return
	}
	defer file.Close()

	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	fmt.Println("i am here")
	if err != nil {
		fmt.Fprintf(w, "Could not upload file")
	}

	fileName, err := UploadFileToS3(s, file, fileHeader)

	fmt.Println("fileName is : ", fileName)
	params := mux.Vars(r)["id"]

	var task models.Property
	err1 := collection.FindOne(context.Background(), params).Decode(&task)
	fmt.Println("PropertyId is : ", task.PropertyId)
	fmt.Println("params is : ", params)

//	_ = json.NewDecoder(r.Body).Decode(&task)

//	fmt.Println("get id is : %q ", getid)
//	filter := bson.M{"propertyid": params}


	fmt.Println("Task image is : %q ", task.Image)
	if err1 != nil {
		fmt.Println(err1)
		return
	} else {
	// Print out data from the document result
		task.Image = fileName		
		json.NewEncoder(w).Encode(task)
	}
	
	
//  _ = json.NewDecoder(r.Body).Decode(&listing)  // config settings: this is where you choose the bucket,
	
	//fmt.Println("err is : \n",err)

	fmt.Fprintf(w, "Image uploaded successfully: %v", fileName)
}



func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
  // get the file size and read
  // the file content into a buffer
  size := fileHeader.Size
  buffer := make([]byte, size)
  file.Read(buffer)

const max = 6555555000
 n := rand.Intn(max)
 str:= strconv.Itoa(n)
  // create a unique file name for the file
  tempFileName := "pictures/" + str +  filepath.Ext(fileHeader.Filename)
fmt.Println("tempFilename is : " , tempFileName)
  // filename, content-type and storage class of the file
  // you're uploading
  _, err := s3.New(s).PutObject(&s3.PutObjectInput{
     Bucket:               aws.String("listingpcs"),
     Key:                  aws.String(tempFileName),
     ACL:                  aws.String("public-read"),// could be private if you want it to be access by only authorized users
     Body:                 bytes.NewReader(buffer),
     ContentLength:        aws.Int64(int64(size)),
     ContentType:        aws.String(http.DetectContentType(buffer)),
     ContentDisposition:   aws.String("attachment"),
     ServerSideEncryption: aws.String("AES256"),
     StorageClass:         aws.String("INTELLIGENT_TIERING"),
  })
  if err != nil {
     return "", err
  }

  return tempFileName, err
}

// CreateProperty create property route
func CreateProperty(w http.ResponseWriter, r *http.Request) {
	/*w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type","*")*/
	u, err := uuid.NewV4()
	//params := mux.Vars(r)
	if err != nil {
		log.Fatal(err)
	}
	var propid=u.String()
	var task models.Property
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.PropertyId = propid
	
	fmt.Println("I have decoded the request")


	insertOneListing(task)		//to insert the record into MongoDB
	
	json.NewEncoder(w).Encode(task)
}

func insertOneListing(task models.Property) {

	insertResult, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
}


//GetAllProperty get all the property route
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
	temp := mux.Vars(r)["id"]
	fmt.Println("my temp is : ",temp)
	var listing models.Property
	_ = json.NewDecoder(r.Body).Decode(&listing)
	update(listing, temp)

}

func update(listing models.Property, temp string) {
	filter := bson.M{"propertyid": temp}
	update := bson.M{"$set": bson.M{ "title": listing.Title, "description" : listing.Description,"street" : listing.StreetAddr,"city" : listing.City,                 						"country" : listing.Country,"zip" : listing.ZipCode,"bed" : listing.Bedrooms,"bath" : listing.Bathrooms,
					"accomodates" : listing.Accomodates,"currency" : listing.Currency,"price" : listing.Price, "minstay" : listing.MinStay,
					"maxstay" : listing.MaxStay, "start" : listing.StartDate,"end" : listing.EndDate, "ptype.pbed" : listing.PropertyType.PrivateBed, 						"ptype.whole" : listing.PropertyType.Whole, "ptype.shared" : listing.PropertyType.Shared, "amenities.ac" : listing.Amenities.Ac, 						"amenities.heater" : listing.Amenities.Heater, "amenities.tv" : listing.Amenities.TV, "amenities.wifi" : listing.Amenities.Wifi, 						"spaces.kitchen" : listing.Spaces.Kitchen, "spaces.closets" : listing.Spaces.Closets, "spaces.parking" : listing.Spaces.Parking, 						"spaces.gym" : listing.Spaces.Gym, "spaces.pool" : listing.Spaces.Pool },}
	result := collection.FindOneAndUpdate(context.Background(), filter, update)
	fmt.Println("hi" , listing.PropertyType.Whole)	
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
	
	fmt.Println("dEleted from db")
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
	
	sess, err := session.NewSession(&aws.Config{
	Region: aws.String("us-east-1")},
	)
	svc := s3.New(sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(task.Image),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		
	}

	fmt.Println(result)
}

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}

func GetManyProperty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


	params := mux.Vars(r)["user"]

	var listing models.Property
		_ = json.NewDecoder(r.Body).Decode(&listing)

	filter := bson.M{"userid": params}

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
	    	    
	    sess, err := session.NewSession(&aws.Config{
            Region: aws.String("us-east-1")},
        	)
	    svc := s3.New(sess)
	    
	    input := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(p.Image),
		}
    	
	     result, err := svc.GetObject(input)
	     
	     json.NewEncoder(w).Encode(p)
	     if err != nil {
		fmt.Println(err.Error())
	     }
	     
	    // fmt.Println("Printing result")
	     fmt.Println(result)
	     }
	     
	     
	     	    
		// check if the cursor encountered any errors while iterating 
	     if err := cursor.Err(); err != nil {
			log.Fatal(err)
	     }
	

}


