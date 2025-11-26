package api

import (
	"github.com/ObiaNzk/LTK-test-manu/cmd/api/handlers"
	"github.com/ObiaNzk/LTK-test-manu/internal"
	"log"
	"net/http"
)

func main() {
	service := internal.NewService()

	handler := handlers.NewHandler(service)

	router := NewRouter(handler)

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
