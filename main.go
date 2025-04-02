package main

import (
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/yourusername/library-api/handlers"
)

func main() {
	r := mux.NewRouter()

	// Initialize routes
	handlers.InitBookRoutes(r)

	// Server configuration
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server started on :8080")
	log.Fatal(srv.ListenAndServe())
}