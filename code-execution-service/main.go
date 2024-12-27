package main

import (
	"log"
	"net/http"
	"os"
	"time"

	pkg "code-execution-service/packages"

	"github.com/gorilla/mux"
)

const DefaultPort = "8080"

func main() {
	// Create a new execution service instance.
	service := pkg.NewExecutionService()

	// Create a new router.
	r := mux.NewRouter()

	// Apply middleware to the router.
	r.Use(LoggingMiddleware)
	r.Use(service.RateLimiter.Limit) // Now correctly using Limit method

	// Define the /execute endpoint with the POST method.
	r.HandleFunc("/execute", service.HandleExecute).Methods("POST")

	// Get the port from the environment or use the default.
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	// Start the server.
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// LoggingMiddleware logs each incoming HTTP request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)

		// Log the request details after it's been handled.
		log.Printf("[%s] %s %s - %v", r.Method, r.URL, r.RemoteAddr, time.Since(startTime))
	})
}
