package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/appErrors"
	"github.com/google/uuid"
)

type Event struct {
	ID              uuid.UUID `json:"event_id"`
	Title           string    `json:"title"`
	DateStart       time.Time `json:"date_start"`
	DateEnd         time.Time `json:"date_end"`
	Description     string    `json:"description"`
	UserId          int       `json:"user_id"`
	EventNotifyTime int       `json:"event_notify_time"`
}

// Validate checks the event business rules. The returned error wraps
// appErrors.ErrValidation so that the transport layer can map it to HTTP 422.
func (e *Event) Validate() error {
	var errstrings []string

	if e.Title == "" {
		errstrings = append(errstrings, appErrors.ErrTitleIsNotSet.Error())
	}

	if e.DateStart.IsZero() {
		errstrings = append(errstrings, appErrors.ErrDateStartIsNotSet.Error())
	}

	if e.DateEnd.IsZero() {
		errstrings = append(errstrings, appErrors.ErrDateEndIsNotSet.Error())
	}

	if !e.DateStart.IsZero() && !e.DateEnd.IsZero() && !e.DateStart.Before(e.DateEnd) {
		errstrings = append(errstrings, appErrors.ErrDateStartBelowDateEnd.Error())
	}

	if e.UserId == 0 {
		errstrings = append(errstrings, appErrors.ErrUserIdIsNotSet.Error())
	}

	if e.EventNotifyTime < 0 {
		errstrings = append(errstrings, appErrors.ErrNotifyTimeIsNegative.Error())
	}

	if len(errstrings) > 0 {
		return fmt.Errorf("%w: %s", appErrors.ErrValidation, strings.Join(errstrings, "; "))
	}

	return nil
}

// ErrInvalidPeriod is returned by PeriodRange for an unknown period.
var ErrInvalidPeriod = errors.New("invalid period")

// Period defines a supported calendar listing interval.
type Period string

const (
	PeriodDay   Period = "day"
	PeriodWeek  Period = "week"
	PeriodMonth Period = "month"
)

// PeriodRange returns the half-open interval [from, to) for the period that
// contains date. Week starts on Monday, month starts on the 1st. All bounds
// are calculated in UTC to keep the behaviour independent from the server
// timezone.
func PeriodRange(date time.Time, period Period) (time.Time, time.Time, error) {
	d := date.UTC()
	start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)

	switch period {
	case PeriodDay:
		return start, start.AddDate(0, 0, 1), nil
	case PeriodWeek:
		offset := (int(start.Weekday()) + 6) % 7 // Monday == 0
		start = start.AddDate(0, 0, -offset)
		return start, start.AddDate(0, 0, 7), nil
	case PeriodMonth:
		start = time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)
		return start, start.AddDate(0, 1, 0), nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("%w: %q", ErrInvalidPeriod, period)
	}
}
