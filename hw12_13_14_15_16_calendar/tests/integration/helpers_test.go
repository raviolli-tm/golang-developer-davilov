//go:build integration

package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

type event struct {
	EventID         string    `json:"event_id"`
	Title           string    `json:"title"`
	DateStart       time.Time `json:"date_start"`
	DateEnd         time.Time `json:"date_end"`
	Description     string    `json:"description"`
	UserID          int32     `json:"user_id"`
	EventNotifyTime int32     `json:"event_notify_time"`
}

var (
	baseURL     string
	db          *sql.DB
	testTimeout time.Duration
)

func TestMain(m *testing.M) {
	baseURL = envOr("CALENDAR_URL", "http://localhost:8888")
	dsn := envOr("POSTGRES_DSN", "postgres://postgres_user:postgres_password@localhost:5432/postgres_db?sslmode=disable")
	testTimeout = parseDuration(envOr("INTEGRATION_TEST_TIMEOUT", "180s"), 180*time.Second)

	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	deadline := time.Now().Add(testTimeout)
	if err := waitForAPI(deadline); err != nil {
		log.Fatalf("calendar API is not available: %v", err)
	}

	if err := waitForDatabase(deadline); err != nil {
		log.Fatalf("database is not available: %v", err)
	}

	os.Exit(m.Run())
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func parseDuration(value string, fallback time.Duration) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return fallback
	}
	return duration
}

func waitForAPI(deadline time.Time) error {
	client := &http.Client{Timeout: 3 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/hello-world")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		time.Sleep(time.Second)
	}

	return fmt.Errorf("timeout waiting for %s", baseURL)
}

func waitForDatabase(deadline time.Time) error {
	for time.Now().Before(deadline) {
		if err := db.Ping(); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("timeout waiting for database")
}

func cleanupDatabase(t *testing.T) {
	t.Helper()

	if _, err := db.Exec(`DELETE FROM calendar.notification`); err != nil {
		t.Fatalf("clean notifications: %v", err)
	}
	if _, err := db.Exec(`DELETE FROM calendar.calendar_event`); err != nil {
		t.Fatalf("clean events: %v", err)
	}
}

func createEvent(t *testing.T, title string, start, end time.Time) event {
	t.Helper()

	return createEventWithNotify(t, title, start, end, 10080)
}

func createEventWithNotify(t *testing.T, title string, start, end time.Time, notifyMinutes int32) event {
	t.Helper()

	payload := event{
		Title:           title,
		DateStart:       start,
		DateEnd:         end,
		Description:     "integration test event",
		UserID:          100,
		EventNotifyTime: notifyMinutes,
	}

	status, body := doJSON(t, http.MethodPost, "/calendar", payload)
	if status != http.StatusOK {
		t.Fatalf("create event %q: status %d, body %s", title, status, body)
	}

	var created event
	if err := json.Unmarshal(body, &created); err != nil {
		t.Fatalf("decode created event: %v", err)
	}
	if created.EventID == "" {
		t.Fatalf("created event has no id")
	}

	return created
}

func listEvents(t *testing.T, path string) []event {
	t.Helper()

	status, body := doJSON(t, http.MethodGet, path, nil)
	if status != http.StatusOK {
		t.Fatalf("GET %s: status %d, body %s", path, status, body)
	}

	var events []event
	if err := json.Unmarshal(body, &events); err != nil {
		t.Fatalf("decode events: %v", err)
	}

	return events
}

func doJSON(t *testing.T, method, path string, payload interface{}) (int, []byte) {
	t.Helper()

	var reader io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal payload: %v", err)
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, baseURL+path, reader)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}

	return resp.StatusCode, body
}

func containsEvent(events []event, id string) bool {
	for _, item := range events {
		if item.EventID == id {
			return true
		}
	}
	return false
}

func uniqueTitle(prefix string) string {
	return fmt.Sprintf("%s %d", prefix, time.Now().UnixNano())
}

func notificationCount(t *testing.T, title string) int {
	t.Helper()

	var count int
	err := db.QueryRow(`SELECT count(*) FROM calendar.notification WHERE notification_event_title = $1`, title).Scan(&count)
	if err != nil {
		t.Fatalf("query notifications: %v", err)
	}

	return count
}

func formatDate(date time.Time) string {
	return date.UTC().Format("2006-01-02")
}
