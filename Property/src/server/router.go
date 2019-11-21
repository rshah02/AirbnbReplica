package main

import (
	
	"github.com/gorilla/mux"
)


// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
    router.HandleFunc("/ping",pingHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/property",GetAllProperty).Methods("GET", "OPTIONS")
	router.HandleFunc("/property/addProperty", CreateProperty).Methods("POST", "OPTIONS")
	router.HandleFunc("/property/{id}", GetProperty).Methods("GET", "OPTIONS")
	router.HandleFunc("/property1/{user}", GetManyProperty).Methods("GET", "OPTIONS")
	router.HandleFunc("/property/{id}", UpdateProperty).Methods("PUT", "OPTIONS")
	router.HandleFunc("/property/{id}", DeleteProperty).Methods("DELETE", "OPTIONS")
	return router
}
