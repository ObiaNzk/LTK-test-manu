package internal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Storage struct {
	db sql.DB
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(ctx context.Context, event CreateEventRequest) (CreateEventResponse, error) {
	trx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return CreateEventResponse{}, fmt.Errorf("creating transaction :%w", err)
	}

	defer trx.Rollback()

	// this could be created from the DB
	id := uuid.NewString()
	createdAt := time.Now().UTC()

	query := "INSERT INTO events (id,title, description, start_time, end_time, created_at) VALUES ($1,$2, $3, $4, $5,$6)"

	if _, err := s.db.ExecContext(ctx, query, id, event.Title, event.Description, event.StartTime, event.EndTime, createdAt); err != nil {
		return CreateEventResponse{}, fmt.Errorf("creating event: %w", err)
	}

	result := CreateEventResponse{
		ID:          id,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		CreatedAt:   createdAt,
	}

	return result, trx.Commit()
}

func (s *Storage) GetEvents(ctx context.Context) ([]CreateEventResponse, error) {
	query := "SELECT id,title, description, start_time, end_time, created_at FROM events ORDER BY start_time ASC"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return []CreateEventResponse{}, fmt.Errorf("creating event: %w", err)
	}

	defer rows.Close()

	var results []CreateEventResponse

	for rows.Next() {
		var event CreateEventResponse
		if err := rows.Scan(&event.ID, &event.CreatedAt); err != nil {
			return []CreateEventResponse{}, fmt.Errorf("scanning event: %w", err)
		}

		results = append(results, event)
	}

	return results, nil
}
