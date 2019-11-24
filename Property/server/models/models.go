package models

type Property struct {
	PropertyId    string       `json:"PropertyId" bson:"PropertyId"`
	Username      string 	   `json:"Username" bson:"Username"`
	UserId	      string 	   `json:"UserId" bson:"UserId"`	
	Title         string       `json:"Title" bson:"Title"`
	Description   string 	   `json:"Description" bson:"Description"`
	StreetAddr    string       `json:"StreetAddr" bson:"StreetAddr"`
	City	      string	   `json:"City" bson:"City"`
	Country       string 	   `json:"Country" bson:"Country"`				
	ZipCode	      string       `json:"ZipCode" bson:"ZipCode"`
	Bedrooms      string          `json:"Bedrooms" bson:"Bedrooms"`
	Bathrooms     string          `json:"Bathrooms" bson:"Bathrooms"`
	Accomodates   string       `json:"Accomodates" bson:"Accomodates"`
	Currency      string       `json:"Currency" bson:"Currency"`
	Price	      string	   `json:"Price" bson:"Price"`	
	MinStay       string       `json:"MinStay" bson:"MinStay"`
	MaxStay       string       `json:"MaxStay" bson:"MaxStay"`
	StartDate     string	    `json:"StartDate" bson:"StartDate"`
	EndDate       string	    `json:"EndDate" bson:"EndDate"`
 	PropertyType  PropertyType  `json:"PropertyType" bson:"PropertyType"`
	Amenities     Amenities	    `json:"Amenities" bson:"Amenities"`
	Spaces	      Spaces	    `json:"Spaces" bson:"Spaces"`
	Image	      string 		`json: "Image" bson:"Image"`
}	  

type PropertyType struct {
	PrivateBed	bool 		`json:"PrivateBed" bson:"PrivateBed"`	
	Whole		bool		`json:"Whole" bson:"Whole"`
	Shared		bool		`json:"Shared" bson:"Shared"`
	
}

type Amenities struct {
	
	Ac 	bool 		`json:"Ac" bson:"Ac"`
	Heater  bool 		`json:"Heater" bson:"Heater"`
	TV	bool 		`json:"TV" bson:"TV"`
	Wifi	bool 		`json:"Wifi" bson:"Wifi"`
}



type Spaces struct {
	
	Kitchen  bool 	`json:"Kitchen" bson:"Kitchen"`
	Closets  bool 	`json:"Closets" bson:"Closets"`
	Parking  bool 	`json:"Parking" bson:"Parking"`
	Gym	 bool 	`json:"Gym" bson:"Gym"`
	Pool	 bool 	`json:"Pool" bson:"Pool"`
	
}




    
