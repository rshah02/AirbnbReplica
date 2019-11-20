package models

type Property struct {
	PropertyId    string       `json:"propertyid" bson:"propertyid"`
	Username      string 	   `json:"username" bson:"username"`
	UserId	      string 	   `json:"userid" bson:"userid"`	
	Title         string       `json:"title" bson:"title"`
	Description   string 	   `json:"description" bson:"description"`
	StreetAddr    string       `json:"street" bson:"street"`
	City	      string	   `json:"city" bson:"city"`
	Country       string 	   `json:"country" bson:"country"`				
	ZipCode	      string       `json:"zip" bson:"zip"`
	Bedrooms      int          `json:"bed" bson:"bed"`
	Bathrooms     int          `json:"bath" bson:"bath"`
	Accomodates   int          `json:"accomodates" bson:"accomodates"`
	Currency      string       `json:"currency" bson:"currency"`
	Price	      int	   `json:"price" bson:"price"`	
	MinStay       int          `json:"minstay" bson:"minstay"`
	MaxStay       int          `json:"maxstay" bson:"maxstay"`
	StartDate     string	    `json:"start" bson:"start"`
	EndDate       string	    `json:"end" bson:"end"`
 	PropertyType  PropertyType  `json:"ptype" bson:"ptype"`
	Amenities     Amenities	    `json:"amenities" bson:"amenities"`
	Spaces	      Spaces	    `json:"spaces" bson:"spaces"`
	Image	      string 		`json: "image" bson:"image"`
}	  

type PropertyType struct {
	PrivateBed	bool 		`json:"pbed" bson:"pbed"`	
	Whole		bool		`json:"whole" bson:"whole"`
	Shared		bool		`json:"shared" bson:"shared"`
	
}

type Amenities struct {
	
	Ac 	bool 		`json:"ac" bson:"ac"`
	Heater  bool 		`json:"heater" bson:"heater"`
	TV	bool 		`json:"tv" bson:"tv"`
	Wifi	bool 		`json:"wifi" bson:"wifi"`
}



type Spaces struct {
	
	Kitchen  bool 	`json:"kitchen" bson:"kitchen"`
	Closets  bool 	`json:"closets" bson:"closets"`
	Parking  bool 	`json:"parking" bson:"parking"`
	Gym	 bool 	`json:"gym" bson:"gym"`
	Pool	 bool 	`json:"pool" bson:"pool"`
	
}




    
