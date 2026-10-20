package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type ScannerStorage struct {
	s *Storage
}

func (s *Storage) NewScannerStorage() *ScannerStorage {
	return &ScannerStorage{s: s}
}

// ReadEventToNotify returns events whose notification time belongs to the
// half-open interval (from, to].
func (s *ScannerStorage) ReadEventToNotify(
	ctx context.Context,
	from, to time.Time,
) ([]storage.Event, error) {
	query := `SELECT ` + eventColumns + ` FROM calendar.calendar_event
		WHERE date_start - INTERVAL '1 minute' * event_notify_time > $1::timestamp
		  AND date_start - INTERVAL '1 minute' * event_notify_time <= $2::timestamp
		ORDER BY date_start`

	rows, err := s.s.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to read events to notify: %w", err)
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
		return nil, fmt.Errorf("failed to read events to notify: %w", err)
	}

	return result, nil
}

// ReadEventToDelete returns IDs of events that ended more than one year ago.
func (s *ScannerStorage) ReadEventToDelete(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := s.s.db.QueryContext(ctx,
		`SELECT event_id FROM calendar.calendar_event WHERE date_end < $1::timestamp`,
		time.Now().Add(-365*24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read events to delete: %w", err)
	}
	defer rows.Close()

	result := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan event id: %w", err)
		}
		result = append(result, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read events to delete: %w", err)
	}

	return result, nil
}

// DeleteEventByIDs deletes the given events and returns the IDs that were
// actually removed.
func (s *ScannerStorage) DeleteEventByIDs(ctx context.Context, ids []uuid.UUID) ([]uuid.UUID, error) {
	tx, err := s.s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin delete transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	deleted := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		result, err := tx.ExecContext(ctx, `DELETE FROM calendar.calendar_event WHERE event_id = $1`, id)
		if err != nil {
			return nil, fmt.Errorf("failed to delete event %s: %w", id, err)
		}
		if affected, err := result.RowsAffected(); err == nil && affected > 0 {
			deleted = append(deleted, id)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit delete transaction: %w", err)
	}

	return deleted, nil
}
