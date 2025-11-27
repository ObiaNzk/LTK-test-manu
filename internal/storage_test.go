package internal_test

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ObiaNzk/LTK-test-manu/internal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	suite.Suite
	db      *sql.DB
	mock    sqlmock.Sqlmock
	storage *internal.Storage
}

func (s *StorageTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)

	s.db = db
	s.mock = mock
	s.storage = internal.NewStorage(db)
}

func (s *StorageTestSuite) TearDownTest() {
	s.db.Close()

	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *StorageTestSuite) TestCreateEvent_Success() {
	ctx := context.Background()
	now := time.Now()
	title := strings.Repeat("a", 101)

	request := internal.CreateEventRequest{
		Title:       title,
		Description: "Test Description",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec("INSERT INTO events \\(id,title, description, start_time, end_time, created_at\\) VALUES \\(\\$1,\\$2, \\$3, \\$4, \\$5,\\$6\\)").
		WithArgs(
			sqlmock.AnyArg(),
			title,
			"Test Description",
			now,
			now.Add(time.Hour),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	result, err := s.storage.CreateEvent(ctx, request)

	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), result.ID)
	require.Equal(s.T(), title, result.Title)
	require.Equal(s.T(), "Test Description", result.Description)
	require.Equal(s.T(), now, result.StartTime)
	require.Equal(s.T(), now.Add(time.Hour), result.EndTime)
	require.NotZero(s.T(), result.CreatedAt)
}

func (s *StorageTestSuite) TestCreateEvent_BeginTxError() {
	ctx := context.Background()
	now := time.Now()
	title := strings.Repeat("a", 101)

	request := internal.CreateEventRequest{
		Title:       title,
		Description: "Test Description",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
	}

	s.mock.ExpectBegin().WillReturnError(errors.New("begin transaction failed"))

	_, err := s.storage.CreateEvent(ctx, request)

	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "creating transaction")
}

func (s *StorageTestSuite) TestCreateEvent_ExecError() {
	ctx := context.Background()
	now := time.Now()
	title := strings.Repeat("a", 101)

	request := internal.CreateEventRequest{
		Title:       title,
		Description: "Test Description",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec("INSERT INTO events \\(id,title, description, start_time, end_time, created_at\\) VALUES \\(\\$1,\\$2, \\$3, \\$4, \\$5,\\$6\\)").
		WithArgs(
			sqlmock.AnyArg(),
			title,
			"Test Description",
			now,
			now.Add(time.Hour),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("insert failed"))

	s.mock.ExpectRollback()

	_, err := s.storage.CreateEvent(ctx, request)

	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "creating event")
}

func (s *StorageTestSuite) TestCreateEvent_CommitError() {
	ctx := context.Background()
	now := time.Now()
	title := strings.Repeat("a", 101)

	request := internal.CreateEventRequest{
		Title:       title,
		Description: "Test Description",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec("INSERT INTO events \\(id,title, description, start_time, end_time, created_at\\) VALUES \\(\\$1,\\$2, \\$3, \\$4, \\$5,\\$6\\)").
		WithArgs(
			sqlmock.AnyArg(),
			title,
			"Test Description",
			now,
			now.Add(time.Hour),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

	_, err := s.storage.CreateEvent(ctx, request)

	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "commit failed")
}

func (s *StorageTestSuite) TestGetEvents_Success_MultipleRows() {
	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "description",
		"start_time", "end_time", "created_at",
	}).
		AddRow(
			"id-1",
			strings.Repeat("a", 101),
			"Description 1",
			now,
			now.Add(time.Hour),
			now,
		).
		AddRow(
			"id-2",
			strings.Repeat("b", 101),
			"Description 2",
			now.Add(2*time.Hour),
			now.Add(3*time.Hour),
			now,
		)

	s.mock.ExpectQuery("SELECT id,title, description, start_time, end_time, created_at FROM events ORDER BY start_time ASC").
		WillReturnRows(rows)

	results, err := s.storage.GetEvents(ctx)

	require.NoError(s.T(), err)
	require.Len(s.T(), results, 2)
	require.Equal(s.T(), "id-1", results[0].ID)
	require.Equal(s.T(), strings.Repeat("a", 101), results[0].Title)
	require.Equal(s.T(), "Description 1", results[0].Description)
	require.Equal(s.T(), "id-2", results[1].ID)
	require.Equal(s.T(), strings.Repeat("b", 101), results[1].Title)
	require.Equal(s.T(), "Description 2", results[1].Description)
}

func (s *StorageTestSuite) TestGetEvents_Success_EmptyTable() {
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{
		"id", "title", "description",
		"start_time", "end_time", "created_at",
	})

	s.mock.ExpectQuery("SELECT id,title, description, start_time, end_time, created_at FROM events ORDER BY start_time ASC").
		WillReturnRows(rows)

	results, err := s.storage.GetEvents(ctx)

	require.NoError(s.T(), err)
	require.Empty(s.T(), results)
}

func (s *StorageTestSuite) TestGetEvents_QueryError() {
	ctx := context.Background()

	s.mock.ExpectQuery("SELECT id,title, description, start_time, end_time, created_at FROM events ORDER BY start_time ASC").
		WillReturnError(errors.New("database connection lost"))

	_, err := s.storage.GetEvents(ctx)

	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "creating event")
}

func (s *StorageTestSuite) TestGetEventByID_Success() {
	ctx := context.Background()
	eventID := "test-id-123"
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "description",
		"start_time", "end_time", "created_at",
	}).AddRow(
		eventID,
		strings.Repeat("a", 101),
		"Test Description",
		now,
		now.Add(time.Hour),
		now,
	)

	s.mock.ExpectQuery("SELECT id, title, description, start_time, end_time, created_at FROM events WHERE id = \\$1").
		WithArgs(eventID).
		WillReturnRows(rows)

	result, err := s.storage.GetEventByID(ctx, eventID)

	require.NoError(s.T(), err)
	require.Equal(s.T(), eventID, result.ID)
	require.Equal(s.T(), strings.Repeat("a", 101), result.Title)
	require.Equal(s.T(), "Test Description", result.Description)
	require.Equal(s.T(), now, result.StartTime)
	require.Equal(s.T(), now.Add(time.Hour), result.EndTime)
	require.Equal(s.T(), now, result.CreatedAt)
}

func (s *StorageTestSuite) TestGetEventByID_NotFound() {
	ctx := context.Background()
	eventID := "nonexistent-id"

	s.mock.ExpectQuery("SELECT id, title, description, start_time, end_time, created_at FROM events WHERE id = \\$1").
		WithArgs(eventID).
		WillReturnError(sql.ErrNoRows)

	_, err := s.storage.GetEventByID(ctx, eventID)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, internal.ErrNotFound)
	require.Contains(s.T(), err.Error(), "event not found")
}

func (s *StorageTestSuite) TestGetEventByID_QueryError() {
	ctx := context.Background()
	eventID := "test-id"

	s.mock.ExpectQuery("SELECT id, title, description, start_time, end_time, created_at FROM events WHERE id = \\$1").
		WithArgs(eventID).
		WillReturnError(errors.New("database connection error"))

	_, err := s.storage.GetEventByID(ctx, eventID)

	require.Error(s.T(), err)
	require.NotErrorIs(s.T(), err, internal.ErrNotFound)
	require.Contains(s.T(), err.Error(), "getting event")
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
