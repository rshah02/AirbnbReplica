package models

//import "go.mongodb.org/mongo-driver/bson/primitive"


type Booking struct {
	//Gen       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BookingID string  `json:"BookingID" bson:"BookingID"`
	PropertyID string `json:"PropertyID" bson:"PropertyID"`
	UserId string `json:"UserId" bson:"UserId"`
	Title string `json:"Title" bson:"Title"`
	Guests string `json:"Guests" bson:"Guests"`
	CheckInDate string `json:"CheckInDate" bson:"CheckInDate"`
	CheckOutDate string `json:"CheckOutDate" bson:"CheckOutDate"`
	Message string `json:"Message" bson:"Message"`
	Amount string  `json:"Amount" bson:"Amount"`
}

