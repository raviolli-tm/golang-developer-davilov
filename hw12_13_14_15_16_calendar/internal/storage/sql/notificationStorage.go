package sqlstorage

import (
	"context"
	"fmt"

	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
)

type NotificationStorage struct {
	s *Storage
}

func (s *Storage) NewNotificationStorage() *NotificationStorage {
	return &NotificationStorage{s: s}
}

func (s *NotificationStorage) CreateNotification(
	ctx context.Context,
	n storage.Notification,
) (storage.Notification, error) {
	query := `INSERT INTO calendar.notification
		(notification_event_title, notification_event_date, notification_user_id)
		VALUES ($1, $2, $3)
		RETURNING notification_id, notification_date`

	err := s.s.db.QueryRowContext(ctx, query,
		n.NotificationEventTitle,
		n.NotificationEventDate,
		n.NotificationUserId,
	).Scan(&n.NotificationId, &n.NotificationDate)
	if err != nil {
		return storage.Notification{}, fmt.Errorf("failed to insert notification: %w", err)
	}

	return n, nil
}
