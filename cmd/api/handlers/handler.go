package handlers

import (
	"LTK-test-manu/internal"
	"net/http"
)

type Handler struct {
	service *internal.Service
}

func NewHandler(service *internal.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	message := h.service.HelloWorld()
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
