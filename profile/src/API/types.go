/*
	Gumball API in Go (Version 3)
	Uses MongoDB and RabbitMQ 
	(For use with Kong API Key)
*/

	
package main
import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Profile struct {
	ProfileId     bson.ObjectId `json:"profileId" bson:"_id"`
	//ProfileId     string  `json:"profileId" bson:"profileId"`
	UserId        string  `json:"userId" bson:"userId"`
	FirstName 	  string  `json:"firstName" bson:"firstName"`
	LastName 	  string  `json:"lastName" bson:"lastName"`
	Email         string  `json:"email" bson:"email"`
	City          string  `json:"city" bson:"city"`
	Description   string  `json:"description" bson:"description"`
	Gender        string  `json:"gender" bson:"gender"`
	BirthDate     time.Time `json:"birthDate" bson:"birthDate"`
	Languages	  string   `json:"languages" bson:"languages"`
	Photo	      string   `json:"photo" bson:"photo"`
}


var profiles map[string] Profile
