package router

import (
	"log"
	"net/http"
	"../middleware"
	"github.com/gorilla/mux"
	h "github.com/gorilla/handlers"
)


// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()
	headersOk := h.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := h.AllowedOrigins([]string{"*"})
	methodsOk := h.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	router.HandleFunc("/booking", middleware.GetBookings).Methods("GET")
	//Returns all the bookings of the passed userId
	router.HandleFunc("/{userId}/myBookings", middleware.GetUserBookings).Methods("GET")
	//Creates a new booking for that property (user id is passed in the body)
	router.HandleFunc("/property/book", middleware.DoBooking).Methods("POST")
	//updates a booking based on the passed values
	router.HandleFunc("/{userId}/updateBooking/{bookingId}", middleware.UpdateBooking).Methods("PUT")
	//deletes a booking of that user which is passed and that passed booking
	router.HandleFunc("/{userId}/myBookings/remove/{bookingId}", middleware.DeleteBooking).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3000", h.CORS(headersOk, methodsOk, originsOk)(router)))
	return router
}