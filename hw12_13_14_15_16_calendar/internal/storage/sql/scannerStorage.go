package sqlstorage

import (
	"context"
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"time"
)

type ScannerStorage struct {
	s *Storage
}

func (s *Storage) NewScannerStorage() *ScannerStorage {
	scannerStorage := ScannerStorage{s}
	return &scannerStorage
}

func (s *ScannerStorage) ReadEventToNotify(ctx context.Context, dateTime time.Time, timeRange int) ([]storage.Event, error) {

	dateTime = dateTime.Truncate(time.Minute)

	tx, err := s.s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin ReadEventToNotify transaction: %w", err)
	}
	defer tx.Rollback()

	query := `SELECT event_id, user_id, title from calendar.calendar_event ca
	where ca.date_start - INTERVAL '1 minutes'*ca.event_notify_time > $1::timestamp and
	ca.date_start - INTERVAL '1 minutes'*ca.event_notify_time <= $2::timestamp`

	rows, err := tx.QueryContext(ctx, query, dateTime, dateTime.Add(time.Duration(timeRange)*time.Minute))

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
		return nil, fmt.Errorf("failed to commit ReadEventToNotify transaction: %w", err)
	}

	return result, nil
}

func (s *ScannerStorage) ReadEventToDelete(ctx context.Context) ([]uuid.UUID, error) {

	tx, err := s.s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin ReadEventToDelete transaction: %w", err)
	}
	defer tx.Rollback()

	query := `select event_id from calendar.calendar_event where date_end < $1::timestamp;`

	rows, err := tx.QueryContext(ctx, query, time.Now().Add(-1*time.Hour*24*365))
	if err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}
	defer rows.Close()

	result := make([]uuid.UUID, 0)
	for rows.Next() {
		event := uuid.UUID{}
		if err := rows.Scan(&event); err != nil {
			continue
		}
		result = append(result, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit ReadEventToDelete transaction: %w", err)
	}

	return result, nil
}

func (s *ScannerStorage) DeleteEventByIDs(ctx context.Context, ids []uuid.UUID) ([]uuid.UUID, error) {

	tx, err := s.s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin DeleteEventByIDs transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`DELETE FROM calendar.calendar_event where event_id = $1`)
	if err != nil {
		return nil, err
	}

	result := make([]uuid.UUID, 0)

	for _, id := range ids {
		_, err := stmt.Exec(id)
		if err == nil {
			result = append(result, id)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit DeleteEventByIDs transaction: %w", err)
	}

	return result, nil
}
