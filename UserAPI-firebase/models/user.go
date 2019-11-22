package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	UserId string `json:"userId" bson:"userId" bson:"_id"`
	Email string `json:"email" bson:"email"`
    Password string	`bson:"password" json:"password"`
	FirstName string  `json:"firstName" bson:"firstName"`
	LastName string  `json:"lastName" bson:"lastName"`
}

type UserToken struct {
	
}