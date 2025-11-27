package internal

import "time"

type CreateEventRequest struct {
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type CreateEventResponse struct {
	ID          string
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
}
