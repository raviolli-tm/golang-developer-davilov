package storage

import (
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/appErrors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Event struct {
	ID              uuid.UUID `json:"event_id"`
	Title           string    `json:"title"`
	DateStart       time.Time `json:"date_start"`
	DateEnd         time.Time `json:"date_end"`
	Description     string    `json:"description"`
	UserId          int       `json:"user_id"`
	EventNotifyTime int       `json:"event_notify_time"`
	// TODO
}

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

	if e.DateStart.After(e.DateEnd) {
		errstrings = append(errstrings, appErrors.ErrDateStartBelowDateEnd.Error())
	}

	if e.UserId == 0 {
		errstrings = append(errstrings, appErrors.ErrUserIdIsNotSet.Error())
	}

	if len(errstrings) > 0 {
		return fmt.Errorf(strings.Join(errstrings, "\n"))
	}
	return nil

}
