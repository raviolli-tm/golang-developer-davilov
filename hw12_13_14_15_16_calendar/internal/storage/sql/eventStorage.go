package sqlstorage

import (
	"context"
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type EventStorage struct {
	s *Storage
}

func (s *Storage) NewEventStorage() *EventStorage {
	eventStorage := EventStorage{s}
	return &eventStorage
}

func (s *EventStorage) CreateEvent(e storage.Event, ctx context.Context) (storage.Event, error) {

	tx, err := s.s.db.BeginTx(ctx, nil)
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO calendar.calendar_event (title, date_start, date_end, event_description, user_id, event_notify_time) VALUES ($1, $2, $3, $4, $5, $6) RETURNING event_id`
	rows, err := tx.QueryContext(ctx, query,
		e.Title,
		e.DateStart,
		e.DateEnd,
		e.Description,
		e.UserId,
		e.EventNotifyTime)
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to insert event: %w", err)
	}
	var id uuid.UUID
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return storage.Event{}, fmt.Errorf("failed to scan ID: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	e.ID = id
	return e, nil
}
func (s *EventStorage) ReadEvents(ctx context.Context) ([]storage.Event, error) {

	tx, err := s.s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `select event_id, title, date_start, date_end, event_description, user_id, event_notify_time from calendar.calendar_event`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}
	defer rows.Close()

	result := make([]storage.Event, 0)
	for rows.Next() {
		event := storage.Event{}
		if err := rows.Scan(&event.ID,
			&event.Title,
			&event.DateStart,
			&event.DateEnd,
			&event.Description,
			&event.UserId,
			&event.EventNotifyTime); err != nil {
			continue
		}
		result = append(result, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}
func (s *EventStorage) UpdateEvent(id uuid.UUID, e storage.Event, ctx context.Context) (storage.Event, error) {

	tx, err := s.s.db.BeginTx(ctx, nil)
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `update calendar.calendar_event set title = $1, date_start = $2, date_end = $3, event_description = $4, user_id = $5, event_notify_time = $6 where event_id = $7`

	_, err = tx.QueryContext(ctx, query,
		e.Title,
		e.DateStart,
		e.DateEnd,
		e.Description,
		e.UserId,
		e.EventNotifyTime,
		id.String())
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to update event: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	e.ID = id
	return e, nil

}
func (s *EventStorage) DeleteEvent(id uuid.UUID, ctx context.Context) (uuid.UUID, error) {
	tx, err := s.s.db.BeginTx(ctx, nil) // *sql.Tx
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `delete from calendar.calendar_event where event_id = $1`
	_, err = tx.ExecContext(ctx, query,
		id.String())
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete event: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return id, nil
}
