package storage

import (
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	NotificationId         uuid.UUID `json:"notification_id"`
	NotificationEventTitle string    `json:"notification_event_title"`
	NotificationEventDate  time.Time `json:"notification_event_date"`
	NotificationDate       time.Time `json:"notification_date"`
	NotificationUserId     int       `json:"notification_user_id"`
}
