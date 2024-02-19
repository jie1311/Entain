package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jie1311/Entain/sport/proto/sport"
	_ "github.com/mattn/go-sqlite3"
)

// EventsRepo provides repository access to events.
type EventsRepo interface {
	// Init will initialise our events repository.
	Init() error

	// List will return a list of events.
	List(filter *sport.ListEventsRequestFilter, sort *sport.ListEventsRequestSorting) ([]*sport.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewEventsRepo creates a new events repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

// Init prepares the event repository dummy data.
func (r *eventsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy events.
		err = r.seed()
	})

	return err
}

func (r *eventsRepo) List(filter *sport.ListEventsRequestFilter, sorting *sport.ListEventsRequestSorting) ([]*sport.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()[eventsList]

	query, args = r.applyFilter(query, filter)

	query = r.applySort(query, sorting)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

func (r *eventsRepo) applyFilter(query string, filter *sport.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if filter.VisibleOnly {
		clauses = append(clauses, "visible IS true")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *eventsRepo) applySort(query string, sorting *sport.ListEventsRequestSorting) string {
	// default behavior when sorting is not provided or sort_by is unspecified
	if sorting == nil || sorting.SortBy == 0 {
		sorting = &sport.ListEventsRequestSorting{
			SortBy: 6, // set default to ADVERTISED_START_TIME
		}
	}

	query += " ORDER BY " + sorting.SortBy.String()

	if sorting.Descend {
		query += " DESC "
	}

	return query
}

func (m *eventsRepo) scanEvents(rows *sql.Rows) ([]*sport.Event, error) {
	var events []*sport.Event

	for rows.Next() {
		var event sport.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.Name, &event.Type, &event.Location, &event.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		event.AdvertisedStartTime = ts

		if ts.AsTime().After(time.Now()) {
			event.Status = sport.Status_OPEN
		} else {
			event.Status = sport.Status_CLOSED
		}

		events = append(events, &event)
	}

	return events, nil
}
