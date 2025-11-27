package internal

import (
	"context"
	"fmt"
)

type storage interface {
	CreateEvent(ctx context.Context, event CreateEventRequest) (CreateEventResponse, error)
	GetEvents(ctx context.Context) ([]CreateEventResponse, error)
}

type Service struct {
	storage storage
}

func NewService(storage storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) HelloWorld() string {
	return "Hello, World!"
}

func (s *Service) CreateEvent(ctx context.Context, event CreateEventRequest) (CreateEventResponse, error) {
	if event.Title == "" {
		return CreateEventResponse{}, fmt.Errorf("title cannot be empty: %w", ErrInput)
	}

	if event.Description == "" {
		return CreateEventResponse{}, fmt.Errorf("description cannot be empty: %w", ErrInput)
	}

	if event.StartTime.IsZero() || event.EndTime.IsZero() {
		return CreateEventResponse{}, fmt.Errorf("start time and end time should be set: %w", ErrInput)
	}

	if len(event.Title) <= 100 {
		return CreateEventResponse{}, fmt.Errorf("title should have more than 100 words: %w", ErrInput)
	}

	response, err := s.storage.CreateEvent(ctx, event)

	if err != nil {
		return CreateEventResponse{}, fmt.Errorf("error creating event: %w", err)
	}

	return response, nil
}
