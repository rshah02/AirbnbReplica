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
	ProfileId     bson.ObjectId `json:"ProfileId" bson:"_id"`
	//ProfileId     string  `json:"profileId" bson:"profileId"`
	UserId        string  `json:"UserId" bson:"UserId"`
	FirstName 	  string  `json:"FirstName" bson:"FirstName"`
	LastName 	  string  `json:"LastName" bson:"LastName"`
	Email         string  `json:"Email" bson:"Email"`
	City          string  `json:"City" bson:"City"`
	Description   string  `json:"Description" bson:"Description"`
	Gender        string  `json:"Gender" bson:"Gender"`
	BirthDate     time.Time `json:"BirthDate" bson:"BirthDate"`
	Languages	  string   `json:"Languages" bson:"Languages"`
	Photo	      string   `json:"Photo" bson:"Photo"`
	Phone         string  	`json:"Phone" bson:"Phone"`
}


var profiles map[string] Profile
