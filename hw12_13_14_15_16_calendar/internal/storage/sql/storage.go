package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
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
	s.db, err = sql.Open("pgx", dsn)
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
