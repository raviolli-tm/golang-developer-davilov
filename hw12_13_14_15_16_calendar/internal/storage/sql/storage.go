package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"log"
	"sync"
)

type DatabaseConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
}

type Storage struct {
	db     *sql.DB
	mu     sync.Mutex
	config DatabaseConf
}

func New(conf DatabaseConf, ctx context.Context) *Storage {
	strg := Storage{config: conf, mu: sync.Mutex{}}
	err := strg.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &strg
}

func (s *Storage) Connect(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		return nil
	}

	var err error
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		s.config.Host,
		s.config.Port,
		s.config.Username,
		s.config.Password,
		s.config.Database)
	s.db, err = sql.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}

	err = s.db.PingContext(ctx)
	if err != nil {
		_ = s.db.Close()
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return nil
	}

	done := make(chan error, 1)

	go func() {
		done <- s.db.Close()
	}()

	select {
	case err := <-done:
		s.db = nil
		if err != nil {
			return fmt.Errorf("failed to close db: %w", err)
		}
		return nil

	case <-ctx.Done():
		return fmt.Errorf("context canceled: %w", ctx.Err())

	}
}

func (s *Storage) CreateEvent(e storage.Event, ctx context.Context) error {

	tx, err := s.db.BeginTx(ctx, nil) // *sql.Tx
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	query := `INSERT INTO calendar.calendar_event (title, date_start, date_end, event_description, user_id, event_notify_time) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, query,
		e.Title,
		e.DateStart,
		e.DateEnd,
		e.Description,
		e.UserId,
		e.EventNotifyTime)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
func (s *Storage) ReadEvents(ctx context.Context) ([]storage.Event, error) {

	tx, err := s.db.BeginTx(ctx, nil) // *sql.Tx
	if err != nil {
		log.Fatal(err)
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
		err = fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}
func (s *Storage) UpdateEvent(id uuid.UUID, e storage.Event, ctx context.Context) error {

	tx, err := s.db.BeginTx(ctx, nil) // *sql.Tx
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	query := `update calendar.calendar_event set title = $1, date_start = $2, date_end = $3, event_description = $4, user_id = $5, event_notify_time = $6 where event_id = $7`
	_, err = tx.ExecContext(ctx, query,
		e.Title,
		e.DateStart,
		e.DateEnd,
		e.Description,
		e.UserId,
		e.EventNotifyTime,
		id.String())
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	err = tx.Commit() // или tx.Rollback()
	if err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil

}
func (s *Storage) DeleteEvent(id uuid.UUID, ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil) // *sql.Tx
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	query := `delete from calendar.calendar_event where event_id = $1`
	_, err = tx.ExecContext(ctx, query,
		id.String())
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	err = tx.Commit() // или tx.Rollback()
	if err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
