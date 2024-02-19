package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, sort *racing.ListRacesRequestSorting) ([]*racing.Race, error)

	// GetRace will return a race.
	Get(id int64) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, sorting *racing.ListRacesRequestSorting) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	query = r.applySort(query, sorting)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) Get(id int64) (*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]
	query, args = r.matchID(query, id)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	races, err := r.scanRaces(rows)
	if err != nil {
		return nil, err
	}

	if len(races) > 1 {
		return nil, fmt.Errorf("sql: primary key search should not return more than one results, got %d results", len(races))
	}

	if len(races) == 0 {
		return nil, nil
	}

	return races[0], nil

}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	if filter.VisibleOnly {
		clauses = append(clauses, "visible IS true")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *racesRepo) applySort(query string, sorting *racing.ListRacesRequestSorting) string {
	// default behavior when sorting is not provided or sort_by is unspecified
	if sorting == nil || sorting.SortBy == 0 {
		sorting = &racing.ListRacesRequestSorting{
			SortBy: 6, // set default to ADVERTISED_START_TIME
		}
	}

	query += " ORDER BY " + sorting.SortBy.String()

	if sorting.Descend {
		query += " DESC "
	}

	return query
}

func (r *racesRepo) matchID(query string, id int64) (string, []interface{}) {
	var args []interface{}

	query += " WHERE id IS ?"

	args = append(args, id)

	return query, args
}

func (m *racesRepo) scanRaces(rows *sql.Rows) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		if ts.AsTime().After(time.Now()) {
			race.Status = 1 // OPEN
		} else {
			race.Status = 2 // CLOSED
		}

		races = append(races, &race)
	}

	return races, nil
}
