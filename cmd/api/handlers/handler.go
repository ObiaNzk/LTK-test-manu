package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ObiaNzk/LTK-test-manu/internal"
	"github.com/go-chi/chi/v5"
)

//go:generate mockgen -destination=mocks/mock_events_service.go -package=mocks github.com/ObiaNzk/LTK-test-manu/cmd/api/handlers eventsService

type eventsService interface {
	CreateEvent(ctx context.Context, event internal.CreateEventRequest) (internal.CreateEventResponse, error)
	GetEventByID(ctx context.Context, id string) (internal.CreateEventResponse, error)
}

type Handler struct {
	eventsService eventsService
}

func NewHandler(service eventsService) *Handler {
	return &Handler{
		eventsService: service,
	}
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

func (h *Handler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "empty event", http.StatusBadRequest)
		return
	}

	event, err := h.eventsService.GetEventByID(ctx, id)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("error getting event: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		CreatedAt   time.Time `json:"created_at"`
	}{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		CreatedAt:   event.CreatedAt,
	}

	jsonResult, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error creating json response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)
}

func decodeBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
