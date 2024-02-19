package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jie1311/Entain/sport/proto/sport"
)

func TestInit(t *testing.T) {
	sportDB, _ := sql.Open("sqlite3", "./sport.db")
	eventsRepo := NewEventsRepo(sportDB)

	if err := eventsRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
	}
}

func TestList(t *testing.T) {
	sportDB, _ := sql.Open("sqlite3", "./sport.db")
	eventsRepo := NewEventsRepo(sportDB)
	if err := eventsRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
		return
	}

	events, err := eventsRepo.List(nil, nil)
	if err != nil {
		t.Errorf("error when listing events, err: %s", err.Error())
		return
	}
	if len(events) != 100 {
		t.Errorf("list func failed to retrive all events, expected events: %d, got: %d", 100, len(events))
		return
	}
	for _, event := range events {
		if (event.AdvertisedStartTime.AsTime().After(time.Now()) && event.Status != sport.Status_OPEN) ||
			(event.AdvertisedStartTime.AsTime().Before(time.Now()) && event.Status != sport.Status_CLOSED) {
			t.Error("list func failed to populate status")
			return
		}
	}
}

func TestListWithFilter(t *testing.T) {
	sportDB, _ := sql.Open("sqlite3", "./sport.db")
	eventsRepo := NewEventsRepo(sportDB)
	if err := eventsRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
		return
	}

	filter := sport.ListEventsRequestFilter{
		VisibleOnly: true,
	}

	events, err := eventsRepo.List(&filter, nil)
	if err != nil {
		t.Errorf("error when listing events, err: %s", err.Error())
		return
	}
	for _, event := range events {
		if event.Visible != true {
			t.Error("list func failed to apply appropriate filters")
			return
		}
	}
}

func TestListWithSorting(t *testing.T) {
	sportDB, _ := sql.Open("sqlite3", "./sport.db")
	eventsRepo := NewEventsRepo(sportDB)
	if err := eventsRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
		return
	}

	sorting := sport.ListEventsRequestSorting{
		SortBy:  1, // sort by ID
		Descend: true,
	}

	events, err := eventsRepo.List(nil, &sorting)
	if err != nil {
		t.Errorf("error when listing events, err: %s", err.Error())
		return
	}
	for i := 0; i < len(events)-1; i++ {
		if events[i].Id < events[i+1].Id {
			t.Error("list func failed to apply appropriate sorting")
			return
		}
	}
}
