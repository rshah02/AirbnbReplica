package main

import (
	"context"
	"encoding/json"
	"log"
	http "net/http"
	"./models"
	
	h "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var mongodb_database = "airbnbClonedb"
var mongodb_collection = "User"

//var app *App

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
	var user models.User
	_ = json.NewDecoder(req.Body).Decode(&user)
	
	unqueId := uuid.Must(uuid.NewV4())
	user.UserId = unqueId.String()
	
	log.Print("attempting to connect to db")

	session, err := mgo.Dial("mongodb://admin:cmpe281@10.0.1.209:27017/?replicaSet=cmpe281&connect=direct")
	
	if err != nil {
		log.Print(err)
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
	var person models.User
	_ = json.NewDecoder(req.Body).Decode(&person)

	session, err := mgo.Dial("mongodb://admin:cmpe281@10.0.1.209:27017/?replicaSet=cmpe281&connect=direct")
	if err != nil {
		panic(err)
		ErrorWithJSON(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"Email": person.Email,
		"Password": person.Password}
	
	var user models.User
	err = c.Find(query).One(&user)
	if err == mgo.ErrNotFound {
		ErrorWithJSON(w, "Login Failed", http.StatusUnauthorized)
		return
	}

	// initialize firebase sdk with service account
	opt := option.WithCredentialsFile("airbnb-clone.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
        log.Fatal("error initializing app: %v\n", err)
	}
	log.Print(app)
		
	// create tokens
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal("error getting Auth client: %v\n", err)
	}

	token, err := client.CustomToken(context.Background(), user.UserId)
	if err != nil {
		log.Fatal("error minting custom token: %v\n", err)
	}
	log.Printf("Got custom token: %v\n", token)	
	
	userData := bson.M{
	"Email":   user.Email,
	"FirstName":    user.FirstName,
	"LastName": user.LastName,
	"UserId":      user.UserId,
	"Token": token}
	
	respBody, err := json.MarshalIndent(userData, "", "  ")
	ResponseWithJSON(w, respBody, http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)["id"]

	session, err := mgo.Dial("mongodb://admin:cmpe281@10.0.1.209:27017/?replicaSet=cmpe281&connect=direct")
	if err != nil {
		panic(err)
		ErrorWithJSON(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"UserId": params}
	
	log.Print("parameter: %v\n", params)
	
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

/*func authenticationMiddleware(au http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//verify token
	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
        log.Fatal("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", token)
	})
}*/

func initfirebase() {

	
}
	

func main() {
	log.Print("hello")		
	
	router := mux.NewRouter()
	headersOk := h.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := h.AllowedOrigins([]string{"*"})
	methodsOk := h.AllowedMethods([]string{"POST", "DELETE"})

	router.HandleFunc("/users/signup", RegisterUser).Methods("POST")
	router.HandleFunc("/users/login", UserSignIn).Methods("POST")
	router.HandleFunc("/users/delete/{id}", DeleteUser).Methods("DELETE")
	
	//auth middleware test
	//router.HandleFunc("/users/{id}", authenticationMiddleware(DeleteUser)).Methods("DELETE")
	// testing
    router.HandleFunc("/users/ping", pingHandler).Methods("GET")

	//log.Print("initializing firebase")
	//initfirebase()	
	
	log.Fatal(http.ListenAndServe(":3000", h.CORS(headersOk, methodsOk, originsOk)(router)))
}
