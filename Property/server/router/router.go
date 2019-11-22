package router

import (
	"../middleware"
	"github.com/gorilla/mux"
)


// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/property", middleware.GetAllProperty).Methods("GET", "OPTIONS")
	router.HandleFunc("/property/addProperty", middleware.CreateProperty).Methods("POST", "OPTIONS")
	router.HandleFunc("/property/{id}", middleware.GetProperty).Methods("GET", "OPTIONS")
	router.HandleFunc("/property1/{user}", middleware.GetManyProperty).Methods("GET", "OPTIONS")
	router.HandleFunc("/property/{id}", middleware.UpdateProperty).Methods("PUT", "OPTIONS")
	router.HandleFunc("/property/{id}", middleware.DeleteProperty).Methods("DELETE", "OPTIONS")
	return router
}
