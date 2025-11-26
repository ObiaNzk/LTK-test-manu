package api

import (
	"LTK-test-manu/cmd/api/handlers"
	"LTK-test-manu/internal"
	"log"
	"net/http"
)

func main() {
	// Initialize service
	service := internal.NewService()

	// Initialize handler
	handler := handlers.NewHandler(service)

	// Initialize router
	router := NewRouter(handler)

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
