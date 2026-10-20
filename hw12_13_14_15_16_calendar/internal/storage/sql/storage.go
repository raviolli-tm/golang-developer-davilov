package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

type DatabaseConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
	SSLMode  string `yaml:"sslmode"`
}

type Storage struct {
	db     *sql.DB
	config DatabaseConf
}

// New connects to the database and returns a ready-to-use storage.
func New(conf DatabaseConf, ctx context.Context) (*Storage, error) {
	strg := &Storage{config: conf}

	if err := strg.Connect(ctx); err != nil {
		return nil, err
	}

	return strg, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	if s.db != nil {
		return nil
	}

	sslMode := s.config.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.config.Host,
		s.config.Port,
		s.config.Username,
		s.config.Password,
		s.config.Database,
		sslMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.db == nil {
		return nil
	}

	err := s.db.Close()
	s.db = nil

	if err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}

	return nil
}
