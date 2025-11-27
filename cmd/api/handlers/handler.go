package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ObiaNzk/LTK-test-manu/internal"
)

type eventsService interface {
	HelloWorld() string
	CreateEvent(ctx context.Context, event internal.CreateEventRequest) (internal.CreateEventResponse, error)
}

type Handler struct {
	eventsService eventsService
}

func NewHandler(service eventsService) *Handler {
	return &Handler{
		eventsService: service,
	}
}

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Message string `json:"message" validate:"required"`
	}

	var payload body

	if err := decodeBody(r, &payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %s", err.Error()), http.StatusBadRequest)

		return
	}

	defer r.Body.Close()

	message := h.eventsService.HelloWorld()
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(message))
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type body struct {
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time" validate:"required"`
		EndTime     time.Time `json:"end_time" validate:"required"`
	}

	var payload body

	if err := decodeBody(r, &payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %s", err.Error()), http.StatusBadRequest)

		return
	}

	if payload.Title == "" {
		http.Error(w, "empty title", http.StatusBadRequest)
	}

	if payload.StartTime.IsZero() || payload.EndTime.IsZero() {
		http.Error(w, "start time and end time should be set", http.StatusBadRequest)
	}

	if len(payload.Title) <= 100 {
		http.Error(w, "title should have more than 100 words", http.StatusBadRequest)
	}

	if payload.StartTime.After(payload.EndTime) {
		http.Error(w, "start time should be before end time", http.StatusBadRequest)
	}

	defer r.Body.Close()

	event := internal.CreateEventRequest{
		Title:       payload.Title,
		Description: payload.Description,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
	}

	result, err := h.eventsService.CreateEvent(ctx, event)

	if err != nil {
		if errors.Is(err, internal.ErrPepito) {
			http.Error(w, "pepito", http.StatusConflict)
		}

		message := fmt.Sprintf("Error creating event: %s", err.Error())
		http.Error(w, message, http.StatusInternalServerError)
	}

	response := struct {
		ID          string    `json:"id" validate:"required"`
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time" validate:"required"`
		EndTime     time.Time `json:"end_time" validate:"required"`
		CreatedAt   time.Time `json:"created_at" validate:"required"`
	}{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		StartTime:   result.StartTime,
		EndTime:     result.EndTime,
		CreatedAt:   result.CreatedAt,
	}

	jsonResult, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "creating json response", http.StatusInternalServerError)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(jsonResult)
}

func decodeBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
