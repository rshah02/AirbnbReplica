package main

import (
	"fmt"
	"log"
	"net/http"
	"./middleware"
	"github.com/gorilla/mux"
	h "github.com/gorilla/handlers"

)

func main() {
	router := mux.NewRouter()
	headersOk := h.AllowedHeaders([]string{"X-Requested-With", "content-type", "Authorization", "Accept-Encoding", "Accept","X-CSRF-Token","Content-Length"})
	originsOk := h.AllowedOrigins([]string{"*"})
	methodsOk := h.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT","OPTIONS"})

	router.HandleFunc("/property", middleware.GetAllProperty).Methods("GET")
	router.HandleFunc("/property/addProperty", middleware.CreateProperty).Methods("POST")
	router.HandleFunc("/property/{id}", middleware.GetProperty).Methods("GET")
	router.HandleFunc("/property1/{user}", middleware.GetManyProperty).Methods("GET")
	router.HandleFunc("/property/{id}", middleware.UpdateProperty).Methods("PUT")
	router.HandleFunc("/property/{id}", middleware.DeleteProperty).Methods("DELETE","OPTIONS")



	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", h.CORS(headersOk, methodsOk, originsOk)(router)))
}
