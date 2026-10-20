//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEventCRUD(t *testing.T) {
	cleanupDatabase(t)

	start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	title := uniqueTitle("crud")

	created := createEvent(t, title, start, start.Add(time.Hour))

	events := listEvents(t, "/calendar")
	if !containsEvent(events, created.EventID) {
		t.Fatalf("created event %s is not present in the listing", created.EventID)
	}

	update := event{
		Title:           title + " updated",
		DateStart:       start,
		DateEnd:         start.Add(2 * time.Hour),
		Description:     "updated description",
		UserID:          100,
		EventNotifyTime: 10080,
	}

	status, body := doJSON(t, http.MethodPut, "/calendar/"+created.EventID, update)
	if status != http.StatusOK {
		t.Fatalf("update event: status %d, body %s", status, body)
	}

	events = listEvents(t, "/calendar")
	found := false
	for _, item := range events {
		if item.EventID == created.EventID {
			found = true
			if item.Title != update.Title {
				t.Fatalf("event title was not updated: got %q, want %q", item.Title, update.Title)
			}
		}
	}
	if !found {
		t.Fatalf("updated event %s is not present in the listing", created.EventID)
	}

	status, body = doJSON(t, http.MethodDelete, "/calendar/"+created.EventID, nil)
	if status != http.StatusOK {
		t.Fatalf("delete event: status %d, body %s", status, body)
	}

	events = listEvents(t, "/calendar")
	if containsEvent(events, created.EventID) {
		t.Fatalf("deleted event %s is still present in the listing", created.EventID)
	}

	status, body = doJSON(t, http.MethodDelete, "/calendar/"+created.EventID, nil)
	if status != http.StatusNotFound {
		t.Fatalf("delete missing event: status %d, body %s", status, body)
	}

	status, body = doJSON(t, http.MethodPut, "/calendar/"+uuid.New().String(), update)
	if status != http.StatusNotFound {
		t.Fatalf("update missing event: status %d, body %s", status, body)
	}
}

func TestEventValidation(t *testing.T) {
	cleanupDatabase(t)

	start := time.Now().UTC().Add(time.Hour).Truncate(time.Second)

	tests := []struct {
		name    string
		payload event
	}{
		{
			name: "empty title",
			payload: event{
				DateStart: start,
				DateEnd:   start.Add(time.Hour),
				UserID:    1,
			},
		},
		{
			name: "date end before date start",
			payload: event{
				Title:     "invalid period",
				DateStart: start,
				DateEnd:   start.Add(-time.Hour),
				UserID:    1,
			},
		},
		{
			name: "user is not set",
			payload: event{
				Title:     "no user",
				DateStart: start,
				DateEnd:   start.Add(time.Hour),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, body := doJSON(t, http.MethodPost, "/calendar", test.payload)
			if status != http.StatusUnprocessableEntity {
				t.Fatalf("expected 422, got %d, body %s", status, body)
			}

			var apiError struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}
			if err := json.Unmarshal(body, &apiError); err != nil {
				t.Fatalf("decode error response: %v, body %s", err, body)
			}
			if apiError.Code != "validation_error" {
				t.Fatalf("unexpected error code %q", apiError.Code)
			}
		})
	}

	status, body := doJSON(t, http.MethodPut, "/calendar/not-a-uuid", event{})
	if status != http.StatusUnprocessableEntity {
		t.Fatalf("invalid id: expected 422, got %d, body %s", status, body)
	}
}

func TestEventListingByPeriod(t *testing.T) {
	cleanupDatabase(t)

	now := time.Now().UTC().Truncate(time.Second)
	today := createEvent(t, uniqueTitle("today"), now.Add(time.Hour), now.Add(2*time.Hour))
	inTwoDays := createEvent(t, uniqueTitle("two-days"), now.AddDate(0, 0, 2), now.AddDate(0, 0, 2).Add(time.Hour))
	inTenDays := createEvent(t, uniqueTitle("ten-days"), now.AddDate(0, 0, 10), now.AddDate(0, 0, 10).Add(time.Hour))
	inFortyDays := createEvent(t, uniqueTitle("forty-days"), now.AddDate(0, 0, 40), now.AddDate(0, 0, 40).Add(time.Hour))

	dayEvents := listEvents(t, "/calendar/day?date="+formatDate(now))
	if !containsEvent(dayEvents, today.EventID) {
		t.Fatalf("day listing does not contain today's event")
	}
	if containsEvent(dayEvents, inTwoDays.EventID) {
		t.Fatalf("day listing contains an event from another day")
	}

	weekEvents := listEvents(t, "/calendar/week?date="+formatDate(now))
	if !containsEvent(weekEvents, today.EventID) || !containsEvent(weekEvents, inTwoDays.EventID) {
		t.Fatalf("week listing does not contain events of the current week")
	}
	if containsEvent(weekEvents, inTenDays.EventID) {
		t.Fatalf("week listing contains an event from the next week")
	}

	monthEvents := listEvents(t, "/calendar/month?date="+formatDate(now))
	if !containsEvent(monthEvents, today.EventID) {
		t.Fatalf("month listing does not contain today's event")
	}
	if containsEvent(monthEvents, inFortyDays.EventID) {
		t.Fatalf("month listing contains an event from another month")
	}

	status, body := doJSON(t, http.MethodGet, "/calendar/day?date=not-a-date", nil)
	if status != http.StatusUnprocessableEntity {
		t.Fatalf("invalid date: expected 422, got %d, body %s", status, body)
	}
}

func TestStorerSavesNotification(t *testing.T) {
	cleanupDatabase(t)

	title := uniqueTitle("notification")
	start := time.Now().UTC().Add(5 * time.Second).Truncate(time.Second)

	createEventWithNotify(t, title, start, start.Add(10*time.Minute), 0)

	deadline := time.Now().Add(testTimeout)
	for time.Now().Before(deadline) {
		if notificationCount(t, title) > 0 {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("notification %q was not saved by the storer within %s", title, testTimeout)
}
