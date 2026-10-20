package api

import (
	"net/http"
	"strings"
	"time"

	genApi "github.com/davilov/hw12_13_14_15_calendar/api/go"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
)

const dateFormat = "2006-01-02"

// PeriodRouter adds the day/week/month listing endpoints from the technical
// specification to the generated REST API. It implements genApi.Router, so it
// can be passed to genApi.NewRouter together with the generated controller.
type PeriodRouter struct {
	service *EventAPIService
}

func NewPeriodRouter(service *EventAPIService) genApi.Router {
	return &PeriodRouter{service: service}
}

func (p *PeriodRouter) Routes() genApi.Routes {
	return genApi.Routes{
		"GetCalendarEventsByDay": genApi.Route{
			Name:        "GetCalendarEventsByDay",
			Method:      http.MethodGet,
			Pattern:     "/calendar/day",
			HandlerFunc: p.handle(storage.PeriodDay),
		},
		"GetCalendarEventsByWeek": genApi.Route{
			Name:        "GetCalendarEventsByWeek",
			Method:      http.MethodGet,
			Pattern:     "/calendar/week",
			HandlerFunc: p.handle(storage.PeriodWeek),
		},
		"GetCalendarEventsByMonth": genApi.Route{
			Name:        "GetCalendarEventsByMonth",
			Method:      http.MethodGet,
			Pattern:     "/calendar/month",
			HandlerFunc: p.handle(storage.PeriodMonth),
		},
	}
}

func (p *PeriodRouter) OrderedRoutes() []genApi.Route {
	routes := p.Routes()

	return []genApi.Route{
		routes["GetCalendarEventsByDay"],
		routes["GetCalendarEventsByWeek"],
		routes["GetCalendarEventsByMonth"],
	}
}

func (p *PeriodRouter) handle(period storage.Period) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := parseDate(r.URL.Query().Get("date"))
		if err != nil {
			_ = genApi.EncodeJSONResponse(genApi.Error{
				Code:    "invalid_date",
				Message: err.Error(),
			}, statusPtr(http.StatusUnprocessableEntity), w)
			return
		}

		result, err := p.service.GetCalendarEventsByPeriod(r.Context(), date, period)
		if err != nil {
			result = errorResponse(err)
		}

		_ = genApi.EncodeJSONResponse(result.Body, &result.Code, w)
	}
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	}

	trimmed := strings.TrimSpace(value)
	if date, err := time.Parse(dateFormat, trimmed); err == nil {
		return date, nil
	}

	date, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return time.Time{}, err
	}

	return date.UTC(), nil
}

func statusPtr(status int) *int {
	return &status
}
