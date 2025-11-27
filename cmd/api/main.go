package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ObiaNzk/LTK-test-manu/cmd/api/handlers"
	"github.com/ObiaNzk/LTK-test-manu/internal"
	"github.com/ObiaNzk/LTK-test-manu/internal/platform"
)

func main() {
	cfg := newLocalConfig()

	db, err := platform.NewDB(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	storage := internal.NewStorage(db)
	service := internal.NewService(storage)
	handler := handlers.NewHandler(service)

	router := NewRouter(handler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
