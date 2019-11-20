package router

import (
	"../middleware"
	"github.com/gorilla/mux"
)


// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/booking", middleware.GetBookings).Methods("GET")
	//Returns all the bookings of the passed userId
	router.HandleFunc("/{userId}/mybookings", middleware.GetUserBookings).Methods("GET")
	//Creates a new booking for that property (user id is passed in the body)
	router.HandleFunc("/property/{proprtyId}/book", middleware.DoBooking).Methods("POST")
	//updates a booking based on the passed values
	router.HandleFunc("/{userId}/updateBooking/{bookingId}", middleware.UpdateBooking).Methods("PUT")
	//deletes a booking of that user which is passed and that passed booking
	router.HandleFunc("/{userId}/mybookings/remove/{bookingId}", middleware.DeleteBooking).Methods("DELETE")
	return router
}