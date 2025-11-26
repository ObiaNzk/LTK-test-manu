package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type helloWorldService interface {
	HelloWorld() string
}

type Handler struct {
	helloWorldService helloWorldService
}

func NewHandler(service helloWorldService) *Handler {
	return &Handler{
		helloWorldService: service,
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

	message := h.helloWorldService.HelloWorld()
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(message))
}

func decodeBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
