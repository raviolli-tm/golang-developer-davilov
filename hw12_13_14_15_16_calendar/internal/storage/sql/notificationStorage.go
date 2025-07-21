package sqlstorage

import (
	"context"
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"time"
)

type NotificationStorage struct {
	s *Storage
}

func (s *Storage) NewNotificationStorage() *NotificationStorage {
	notificationStorage := NotificationStorage{s}
	return &notificationStorage
}

func (s *NotificationStorage) CreateNotification(ctx context.Context, n storage.Notification) (storage.Notification, error) {

	tx, err := s.s.db.BeginTx(ctx, nil)
	if err != nil {
		return storage.Notification{}, fmt.Errorf("failed to begin create transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO calendar.notification (title, date_event, user_id) VALUES ($1, $2, $3) RETURNING notification_id;`
	rows, err := tx.QueryContext(ctx, query,
		n.NotificationEventTitle,
		n.NotificationEventDate,
		n.NotificationUserId,
	)
	if err != nil {
		return storage.Notification{}, fmt.Errorf("failed to insert notification: %w", err)
	}
	var id uuid.UUID
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return storage.Notification{}, fmt.Errorf("failed to scan notification ID: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return storage.Notification{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	n.NotificationId = id
	n.NotificationDate = time.Now()

	return n, nil
}
