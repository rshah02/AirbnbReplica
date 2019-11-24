package main

import (
	"encoding/json"
	"fmt"
	"log"
	http "net/http"
	"os"
	"time"
	
	"github.com/auth0-community/auth0"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	jose "gopkg.in/square/go-jose.v2"
)

var mongodb_server = "127.0.0.1"
var mongodb_database = "admin"
var mongodb_collection = "user"
var mongo_admin_database = "admin"
var mongo_username = "admin"
var mongo_password = "cmpe281"

func pingHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("hello")
	mapD := map[string]string{"message": "API Working"}
	mapB, _ := json.Marshal(mapD)
	ResponseWithJSON(w, mapB, http.StatusOK)
	return
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	mapD := map[string]string{"message": message}
	mapB, _ := json.Marshal(mapD)
	ResponseWithJSON(w, mapB, code)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

//new user sign up
func RegisterUser(w http.ResponseWriter, req *http.Request) {
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	unqueId := uuid.Must(uuid.NewV4())
	user.UserId = unqueId.String()
	info := &mgo.DialInfo{
		Addrs:    []string{mongodb_server},
		Timeout:  60 * time.Second,
		Database: mongodb_database,
		Username: mongo_username,
		Password: mongo_password,
	}

	session, err := mgo.DialWithInfo(info)

	if err != nil {
		ErrorWithJSON(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)

	err = c.Insert(user)
	if err != nil {
		if mgo.IsDup(err) {
			ErrorWithJSON(w, "User with this ID already exists", http.StatusBadRequest)
			return
		}
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		return
	}

	respBody, err := json.MarshalIndent(user, "", "  ")
	ResponseWithJSON(w, respBody, http.StatusOK)
}

func UserSignIn(w http.ResponseWriter, req *http.Request) {
	var person User
	_ = json.NewDecoder(req.Body).Decode(&person)
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
		ErrorWithJSON(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"email": person.Email,
		"password": person.Password}
	var user User

	err = c.Find(query).One(&user)
	if err == mgo.ErrNotFound {
		ErrorWithJSON(w, "Login Failed", http.StatusUnauthorized)
		return
	}
	userData := bson.M{
		"email":   user.Email,
		"FirstName":    user.FirstName,
		"LastName": user.LastName,
		"UserId":      user.UserId}

	respBody, err := json.MarshalIndent(userData, "", "  ")
	ResponseWithJSON(w, respBody, http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
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
		ErrorWithJSON(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"_id": params["UserId"]}
	err = c.Remove(query)

	if err != nil {
		switch err {
		default:
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed delete user: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(w, "User not found", http.StatusNotFound)
			return
		}
	}

	respBody, err := json.MarshalIndent("User deleted", "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func authenticationMiddleware(au http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	
		SigningSecret := []byte("JV3XRl71yLojhYcvDxY51IWNmVge36CM");
		audience := []string{"airbnb-clone"}		
		
		configuration := auth0.NewConfiguration(auth0.NewKeyProvider(SigningSecret), audience, "https://dev-xr8-x4am.auth0.com/" , jose.HS256)
		validator := auth0.NewValidator(configuration)
		
		token, err := validator.ValidateRequest(r)
	
		if err != nil {
            log.Print(err)
            log.Print("Token is not valid:", token)
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
        } else {
            au.ServeHTTP(w, r)
        }
		
	})
}

func main() {
	log.Print("hello")
	router := mux.NewRouter()
	router.HandleFunc("/users/signup", RegisterUser).Methods("POST")
	router.HandleFunc("/users/signin", UserSignIn).Methods("POST")
	router.HandleFunc("/users/{id}", authenticationMiddleware(DeleteUser)).Methods("DELETE")
	// testing
    router.HandleFunc("/users/ping", pingHandler).Methods("GET")	
	log.Fatal(http.ListenAndServe(":3000", router))
}
