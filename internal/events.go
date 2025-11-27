package internal

import (
	"context"
	"fmt"
)

type storage interface {
	CreateEvent(ctx context.Context, event CreateEventRequest) (CreateEventResponse, error)
	GetEvents(ctx context.Context) ([]CreateEventResponse, error)
	GetEventByID(ctx context.Context, id string) (CreateEventResponse, error)
}

type Service struct {
	storage storage
}

func NewService(storage storage) *Service {
	return &Service{
		storage: storage,
	}
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
		return CreateEventResponse{}, fmt.Errorf("creating event: %w", err)
	}

	return response, nil
}

func (s *Service) GetEventByID(ctx context.Context, id string) (CreateEventResponse, error) {
	if id == "" {
		return CreateEventResponse{}, fmt.Errorf("empty id: %w", ErrInput)
	}

	event, err := s.storage.GetEventByID(ctx, id)
	if err != nil {
		return CreateEventResponse{}, fmt.Errorf("getting event: %w", err)
	}

	return event, nil
}

func (s *Service) GetEvents(ctx context.Context) ([]CreateEventResponse, error) {
	events, err := s.storage.GetEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting events: %w", err)
	}

	return events, nil
}
