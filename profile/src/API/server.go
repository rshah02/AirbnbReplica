/*ll API in Go (Version 3)
	Uses MongoDB and RabbitMQ 
	(For use with Kong API Key)
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// MongoDB Config
// var mongodb_server = os.Getenv("SERVER")
// var mongodb_database = os.Getenv("DATABASE")
// var mongodb_collection = os.Getenv("COLLECTION")
// var mongo_admin_database = os.Getenv("ADMIN_DATABASE")
// var mongo_username = os.Getenv("USERNAME")
// var mongo_password = os.Getenv("PASSWORD")
var mongodb_server = "127.0.0.1"
var mongodb_database = "admin"
var mongodb_collection = "profile"
var mongo_admin_database = "admin"
var mongo_username = "admin"
var mongo_password = "9210"

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping",pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/profile", profileHandler(formatter)).Methods("GET")
	mx.HandleFunc("/profile", profileUpdateHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/profile", createProfileHandler(formatter)).Methods("POST")
	//mx.HandleFunc("/profile/{id}", getProfileByIdHandler(formatter)).Methods("GET")
	mx.HandleFunc("/profile", deleteProfileHandler(formatter)).Methods("DELETE")
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API version 1.0 alive!"})
	}
}

// API Create New Profile
func createProfileHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var profileObject Profile
		_ = json.NewDecoder(req.Body).Decode(&profileObject)
		profileObject.ProfileId=bson.NewObjectId()
		fmt.Println("profile:", profileObject.FirstName )
		info := &mgo.DialInfo{
			Addrs:    []string{mongodb_server},
			Timeout:  60 * time.Second,
			Database: mongodb_database,
			Username: mongo_username,
			Password: mongo_password,
		}
		session, err := mgo.DialWithInfo(info)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		return
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)
		// var profile Profile
		var x Profile
	 	err =c.Find(bson.M{"userId":profileObject.UserId}).One(&x)
		 if err !=nil{
		
						err = c.Insert(profileObject)
						if err != nil {
							if mgo.IsDup(err) {
							formatter.JSON(w,http.StatusBadRequest,"User with this ID already exists")
							return
						}
						formatter.JSON(w, http.StatusInternalServerError, err.Error())
						return
					}
		
			if err != nil {
			log.Fatal(err)
			}
	
			formatter.JSON(w, http.StatusOK, profileObject)
		}else {
			formatter.JSON(w,http.StatusBadRequest,"user exists")
		}

	}
}



// GET profile Handler
func profileHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var profile Profile
		_ = json.NewDecoder(req.Body).Decode(&profile)
		fmt.Println("profile:", profile.ProfileId)
		info := &mgo.DialInfo{
			Addrs:    []string{mongodb_server},
			Timeout:  60 * time.Second,
			Database: mongodb_database,
			Username: mongo_username,
			Password: mongo_password,
		}	
		session, err := mgo.DialWithInfo(info)
        if err != nil {
                panic(err)
        }
        defer session.Close()
        session.SetMode(mgo.Monotonic, true)
        c := session.DB(mongodb_database).C(mongodb_collection)
		var result Profile
		query := bson.M{"_id": profile.ProfileId}
		
		if err =c.Find(query).One(&result);err != nil {
			
				formatter.JSON(w, http.StatusInternalServerError, err.Error())}
        fmt.Println("profile:", result )
		formatter.JSON(w, http.StatusOK, result)
	}
}

// API Update profile
func profileUpdateHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		info := &mgo.DialInfo{
			Addrs:    []string{mongodb_server},
			Timeout:  60 * time.Second,
			Database: mongodb_database,
			Username: mongo_username,
			Password: mongo_password,
		}
		session, err := mgo.DialWithInfo(info)
		
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var newProfile Profile
		if err := json.NewDecoder(req.Body).Decode(&newProfile); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request")
			return
		}

		if err := c.UpdateId(newProfile.ProfileId, &newProfile); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var updatedProfile Profile
		if err := c.FindId(newProfile.ProfileId).One(&updatedProfile); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, updatedProfile)

	}
}


//API delete profile
func deleteProfileHandler  (formatter *render.Render) http.HandlerFunc {
return func	(w http.ResponseWriter, req *http.Request) {

	info := &mgo.DialInfo{
		Addrs:    []string{mongodb_server},
		Timeout:  60 * time.Second,
		Database: mongodb_database,
		Username: mongo_username,
		Password: mongo_password,
	}

	session, err := mgo.DialWithInfo(info)
	defer session.Close()
	if err != nil {
		formatter.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	var profileObject Profile
	_ = json.NewDecoder(req.Body).Decode(&profileObject)
	fmt.Println("Profile Id: ", profileObject.ProfileId)
	err = c.Remove(bson.M{"_id": profileObject.ProfileId})
   fmt.Println("err",err)
	if err != nil {
		fmt.Println("profile not found")
		formatter.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	formatter.JSON(w, http.StatusOK, "Deleted")

}
}


