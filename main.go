package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/astoyanov87/matches-service/docs" // Update with the path to your generated docs folder
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Initialize Redis connection
	rdb := initRedis()

	// Define routes
	http.HandleFunc("/api/v1/matches", func(w http.ResponseWriter, r *http.Request) {
		// Call the handler to get matches by status
		getMatchesByStatusHandler(rdb, w, r)
	})

	// Define routes
	http.HandleFunc("/api/v1/match/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/match endpoint requested")
		// Call the handler to get matches by status
		getMatchByIdHandler(rdb, w, r)
	})

	http.HandleFunc("/api/v1/swagger/", httpSwagger.WrapHandler)
	// Start the HTTP server
	log.Println("Server is running on port 8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
