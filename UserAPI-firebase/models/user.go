package models

type User struct {
	UserId string `json:"UserId" bson:"UserId" bson:"_id"`
	Email string `json:"Email" bson:"Email"`
    Password string	`bson:"Password" json:"Password"`
	FirstName string  `json:"FirstName" bson:"FirstName"`
	LastName string  `json:"LastName" bson:"LastName"`
}

type UserToken struct {
	
}