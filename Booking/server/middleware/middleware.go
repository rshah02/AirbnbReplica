package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	//"flag"

	"../models"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sns"
    "github.com/aws/aws-sdk-go/aws"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//const connectionString = "mongodb://localhost:27017"
const connectionString = "mongodb://admin:admin@primary:27017/admin?replicaSet=cmpe281"

const dbName = "admin"

const collName = "bookings"

var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection = client.Database(dbName).Collection(collName)
	fmt.Println("Collection instance created!")
}

func GetBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	payload := getAllBookings()
	json.NewEncoder(w).Encode(payload)
}

func getAllBookings() []primitive.M {
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
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func GetUserBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println("getAllUserBookings "+params["userId"])
	filter := bson.M{"UserId" : params["userId"]}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var results []models.Booking
	for cur.Next(context.Background()) {
		var result models.Booking
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
	results = append(results, result)
}
json.NewEncoder(w).Encode(results)
}

func DoBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	u, err := uuid.NewV4()
	params := mux.Vars(r)
	if err != nil {
                log.Fatal(err)
        } 
	var bookID = u.String()
	var b models.Booking
	_ = json.NewDecoder(r.Body).Decode(&b)
	b.BookingID = bookID
	b.PropertyID = params["proprtyId"]
	 fmt.Println(b.BookingID, r.Body)
	doOneBooking(b)
	sess, err := session.NewSession(&aws.Config{
    Region: aws.String("us-west-2")},
    )
    svc := sns.New(sess)
/*	topicPtr := flag.String("t", "arn:aws:sns:us-west-2:574372214394:MyTopic", "The ARN of the topic to which the user subscribes")
    flag.Parse()
    topicArn := *topicPtr
        fmt.Println("You must supply a message and topic ARN")
        fmt.Println("Usage: go run SnsPublish.go -m MESSAGE -t TOPIC-ARN")
        os.Exit(1)
    }*/
 sub, err := svc.Subscribe(&sns.SubscribeInput{
    	Endpoint:              aws.String(b.EMail),
        Protocol:              aws.String("email"),
        ReturnSubscriptionArn: aws.Bool(true), 
        TopicArn:              aws.String("arn:aws:sns:us-west-2:574372214394:MyTopic"),
    })
     if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    fmt.Println(*sub.SubscriptionArn)
    result, err := svc.Publish(&sns.PublishInput{
        Message:  aws.String(b.BookingID),
        TopicArn: aws.String("arn:aws:sns:us-west-2:574372214394:MyTopic"),
    })

    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    fmt.Println(*result.MessageId)
    json.NewEncoder(w).Encode(b)
}

func doOneBooking(obj models.Booking) {
	insertResult, err := collection.InsertOne(context.Background(), obj)
	fmt.Println("In Insert Block")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
}

func UpdateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var b models.Booking
	_ = json.NewDecoder(r.Body).Decode(&b)
	 fmt.Println(b, r.Body)
	fmt.Println("In Update Block")
	updateOneBooking(params["userId"], params["bookingId"], b)
		sess, err := session.NewSession(&aws.Config{
    Region: aws.String("us-west-2")},
    )
    svc := sns.New(sess)
/*	topicPtr := flag.String("t", "arn:aws:sns:us-west-2:574372214394:MyTopic", "The ARN of the topic to which the user subscribes")
    flag.Parse()
    topicArn := *topicPtr
        fmt.Println("You must supply a message and topic ARN")
        fmt.Println("Usage: go run SnsPublish.go -m MESSAGE -t TOPIC-ARN")
        os.Exit(1)
    }*/
 sub, err := svc.Subscribe(&sns.SubscribeInput{
    	Endpoint:              aws.String(b.EMail),
        Protocol:              aws.String("email"),
        ReturnSubscriptionArn: aws.Bool(true), 
        TopicArn:              aws.String("arn:aws:sns:us-west-2:574372214394:MyTopic"),
    })
     if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    fmt.Println(*sub.SubscriptionArn)
    result, err := svc.Publish(&sns.PublishInput{
        Message:  aws.String(b.BookingID),
        TopicArn: aws.String("arn:aws:sns:us-west-2:574372214394:MyTopic"),
    })

    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    fmt.Println(*result.MessageId)
	json.NewEncoder(w).Encode(b)
}

func updateOneBooking(uId string, bId string, upBooking models.Booking) {
	//fmt.Println(obj)
	fmt.Println("In Update Block")
	//id, _ := primitive.ObjectIDFromHex(obj)
	filter := bson.M{"BookingID": bId, "UserId": uId}
	update := bson.M{"$set": bson.M{"BookingID": upBooking.BookingID, "PropertyID": upBooking.PropertyID, "UserId": upBooking.UserId, "Title": upBooking.Title, "Guests": upBooking.Guests, "CheckInDate": upBooking.CheckInDate, "CheckOutDate": upBooking.CheckOutDate, "Message": upBooking.Message, "Amount": upBooking.Amount, "EMail": upBooking.EMail}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

func DeleteBooking(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	filter := bson.M{"BookingID" : params["bookingId"]}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var mail string
	for cur.Next(context.Background()) {
		var result models.Booking
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		fmt.Println("Email",result.EMail)
		mail = result.EMail
	}
	sess, err := session.NewSession(&aws.Config{
    Region: aws.String("us-west-2")},
    )
    svc := sns.New(sess)
sub, err := svc.Subscribe(&sns.SubscribeInput{
    	Endpoint:              aws.String(mail),
        Protocol:              aws.String("email"),
        ReturnSubscriptionArn: aws.Bool(true), 
        TopicArn:              aws.String("arn:aws:sns:us-west-2:574372214394:MyTopic"),
    })
     if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    fmt.Println(*sub.SubscriptionArn)
    result, err := svc.Publish(&sns.PublishInput{
        Message:  aws.String(params["bookingId"]),
        TopicArn: aws.String("arn:aws:sns:us-west-2:574372214394:MyTopic"),
    })

    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    fmt.Println(*result.MessageId)
	deleteOneBooking(params["userId"], params["bookingId"])
	json.NewEncoder(w).Encode(struct{ Message string }{"Deleted"})

}

func deleteOneBooking(uId string, bId string) {
	fmt.Println(uId)
	fmt.Println("In Delete Block")
	filter := bson.M{"UserId": uId, "BookingID" : bId}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted count: ", result.DeletedCount)
}

