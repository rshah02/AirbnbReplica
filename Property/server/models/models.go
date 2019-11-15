package models

import "go.mongodb.org/mongo-driver/bson/primitive"


type Property struct {
	//Id        string `json:"id"`
	RecId       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`	
	PropertyId  string       `json:"propertyid" bson:"propertyid"`
	UserId      string        `json:"userid" bson:"userid"`
	Title       string        `json:"title" bson:"title"`
	IsAvailable bool 	  `json:"avail" bson:"avail"`

/*	Type []     struct {
		Whole bool `json:"whole" bson:"whole"`
		Private bool `json:"private" bson:"private"`
		Shared 	bool `json:"shared" bson:"shared"` 	  
	} `json:"type" bson:"type"`
*/
	Guests        int         `json:"guests" bson:"guests"`	 
	Description   string 	  `json:"description" bson:"description"`
	Location      string 	  `json:"location" bson:"location"`
	URltoimage    string      `json:"urltoimage" bson:"urltoimage"`
	Price	      int	  `json:"price" bson:"price"`
	
}

