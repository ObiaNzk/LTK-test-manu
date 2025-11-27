package internal

import (
	"context"
	"errors"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) HelloWorld() string {
	return "Hello, World!"
}

func (s *Service) CreateEvent(ctx context.Context, event CreateEventRequest) (CreateEventResponse, error) {

	return errors.New("not impleneted")
}
