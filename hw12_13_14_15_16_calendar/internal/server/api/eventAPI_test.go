package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/davilov/hw12_13_14_15_calendar/api/go"
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

var uuidCode = uuid.New()

func TestAPI(t *testing.T) {

	t.Run("Test API", func(t *testing.T) {
		CreateCalendarEvent(t)
		UpdateCalendarEvent(t)
		DeleteCalendarEvent(t)
	})

}

func GetCalendarEvent(expected []api.Event) ([]api.Event, error) {

	fmt.Print("Run GetCalendarEvent API...  ")
	resp, err := http.Get("http://localhost:8080/calendar")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	events := make([]api.Event, 0)
	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		return nil, err
	}

	if expected != nil {
		if len(events) != len(expected) {
			return nil, fmt.Errorf("incorrect response length")
		}
	}

	return events, nil

}

func DeleteCalendarEvent(t *testing.T) {

	event, err := GetCalendarEvent(nil)
	fmt.Print("Run DeleteCalendarEvent API...  ")
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/calendar/"+uuidCode.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	_, err = GetCalendarEvent(event[:len(event)-1])
	if err != nil {
		t.Fatal(err)
	}

}

func CreateCalendarEvent(t *testing.T) {

	eventCreate := api.Event{
		EventId:         uuidCode.String(),
		Title:           "test title",
		DateStart:       time.Now(),
		DateEnd:         time.Now().Add(15 * time.Minute),
		Description:     "test description",
		UserId:          1,
		EventNotifyTime: 15,
	}

	event, err := GetCalendarEvent(nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print("Run CreateCalendarEvent API...  ")

	jsonData, err := json.Marshal(eventCreate)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post("http://localhost:8080/calendar", "application/json", bytes.NewBuffer(jsonData))
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		t.Fatal(fmt.Errorf("bad status: %s", resp.Status))
	}
	event = append(event, eventCreate)

	eventAfterCreate, err := GetCalendarEvent(event)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(eventAfterCreate); i++ {
		j := 0
		for ; j < len(event); j++ {
			if event[j].EventId == eventAfterCreate[i].EventId {
				break
			}
		}
		if j < len(event) {
			uuidCode, _ = uuid.Parse(eventAfterCreate[i].EventId)
		}
	}

}

func UpdateCalendarEvent(t *testing.T) {

	eventUpdate := api.Event{
		EventId:         uuidCode.String(),
		Title:           "test title(updated)",
		DateStart:       time.Now(),
		DateEnd:         time.Now().Add(15 * time.Minute),
		Description:     "test description(updated)",
		UserId:          1,
		EventNotifyTime: 15,
	}

	event, err := GetCalendarEvent(nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print("Run UpdateCalendarEvent API...  ")
	jsonData, err := json.Marshal(eventUpdate)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/calendar/"+uuidCode.String(), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Status Code:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		t.Fatal(fmt.Errorf("bad status: %s", resp.Status))
	}
	_, err = GetCalendarEvent(event)
	if err != nil {
		t.Fatal(err)
	}

}
