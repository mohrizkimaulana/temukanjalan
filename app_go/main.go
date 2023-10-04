package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Creating a new router instance
	router := mux.NewRouter()

	// Creating a new controller to handle response
	// for '/' endpoint
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"response": "This is from Golang Backend, Hi Pintu",
		}

		// Setting the response
		json.NewEncoder(rw).Encode(response)
	})

	// Logging on console
	log.Println("[bootup]: Server is running at port: 3000")

	http.ListenAndServe(":3000", router) // starting the server on port 5000
}
