package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/appErrors"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

const eventColumns = `event_id, title, date_start, date_end, event_description, user_id, event_notify_time`

type EventStorage struct {
	s *Storage
}

func (s *Storage) NewEventStorage() *EventStorage {
	return &EventStorage{s: s}
}

func (s *EventStorage) CreateEvent(e storage.Event, ctx context.Context) (storage.Event, error) {
	if err := e.Validate(); err != nil {
		return storage.Event{}, err
	}

	query := `INSERT INTO calendar.calendar_event
		(title, date_start, date_end, event_description, user_id, event_notify_time)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING event_id`

	var id uuid.UUID
	err := s.s.db.QueryRowContext(ctx, query,
		e.Title,
		e.DateStart,
		e.DateEnd,
		e.Description,
		e.UserId,
		e.EventNotifyTime,
	).Scan(&id)
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to insert event: %w", err)
	}

	e.ID = id

	return e, nil
}

func (s *EventStorage) ReadEvents(ctx context.Context) ([]storage.Event, error) {
	query := `SELECT ` + eventColumns + ` FROM calendar.calendar_event ORDER BY date_start`

	return s.queryEvents(ctx, query)
}

// ReadEventsByPeriod returns events whose start time belongs to [from, to).
func (s *EventStorage) ReadEventsByPeriod(from, to time.Time, ctx context.Context) ([]storage.Event, error) {
	query := `SELECT ` + eventColumns + ` FROM calendar.calendar_event
		WHERE date_start >= $1 AND date_start < $2
		ORDER BY date_start`

	return s.queryEvents(ctx, query, from, to)
}

func (s *EventStorage) queryEvents(ctx context.Context, query string, args ...interface{}) ([]storage.Event, error) {
	rows, err := s.s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}
	defer rows.Close()

	result := make([]storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.DateStart,
			&event.DateEnd,
			&event.Description,
			&event.UserId,
			&event.EventNotifyTime,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		result = append(result, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}

	return result, nil
}

func (s *EventStorage) UpdateEvent(id uuid.UUID, e storage.Event, ctx context.Context) (storage.Event, error) {
	if err := e.Validate(); err != nil {
		return storage.Event{}, err
	}

	query := `UPDATE calendar.calendar_event
		SET title = $1, date_start = $2, date_end = $3, event_description = $4,
		    user_id = $5, event_notify_time = $6
		WHERE event_id = $7`

	result, err := s.s.db.ExecContext(ctx, query,
		e.Title,
		e.DateStart,
		e.DateEnd,
		e.Description,
		e.UserId,
		e.EventNotifyTime,
		id,
	)
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to update event: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return storage.Event{}, appErrors.ErrIdDoesNotExist
	}

	e.ID = id

	return e, nil
}

func (s *EventStorage) DeleteEvent(id uuid.UUID, ctx context.Context) (uuid.UUID, error) {
	result, err := s.s.db.ExecContext(ctx, `DELETE FROM calendar.calendar_event WHERE event_id = $1`, id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete event: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return uuid.Nil, appErrors.ErrIdDoesNotExist
	}

	return id, nil
}

// Close closes the underlying database connection.
func (s *EventStorage) Close(ctx context.Context) error {
	return s.s.Close(ctx)
}
