package db

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

func TestInit(t *testing.T) {
	racingDB, _ := sql.Open("sqlite3", "./racing.db")
	racesRepo := NewRacesRepo(racingDB)

	if err := racesRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
	}
}

func TestList(t *testing.T) {
	racingDB, _ := sql.Open("sqlite3", "./racing.db")
	racesRepo := NewRacesRepo(racingDB)
	if err := racesRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
		return
	}

	races, err := racesRepo.List(nil, nil)
	if err != nil {
		t.Errorf("error when listing races, err: %s", err.Error())
		return
	}
	if len(races) != 100 {
		t.Errorf("list func failed to retrive all races, expected races: %d, got: %d", 100, len(races))
		return
	}
	for _, race := range races {
		if (race.AdvertisedStartTime.AsTime().After(time.Now()) && race.Status != racing.Status_OPEN) ||
			(race.AdvertisedStartTime.AsTime().Before(time.Now()) && race.Status != racing.Status_CLOSED) {
			t.Error("list func failed to populate status")
			return
		}
	}
}

func TestListWithFilter(t *testing.T) {
	racingDB, _ := sql.Open("sqlite3", "./racing.db")
	racesRepo := NewRacesRepo(racingDB)
	if err := racesRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
		return
	}

	filter := racing.ListRacesRequestFilter{
		MeetingIds:  []int64{10},
		VisibleOnly: true,
	}

	races, err := racesRepo.List(&filter, nil)
	if err != nil {
		t.Errorf("error when listing races, err: %s", err.Error())
		return
	}
	for _, race := range races {
		if race.MeetingId != 10 || race.Visible != true {
			t.Error("list func failed to apply appropriate filters")
			return
		}
	}
}

func TestListWithSorting(t *testing.T) {
	racingDB, _ := sql.Open("sqlite3", "./racing.db")
	racesRepo := NewRacesRepo(racingDB)
	if err := racesRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
		return
	}

	sorting := racing.ListRacesRequestSorting{
		SortBy:  racing.SortBy_ID,
		Descend: true,
	}

	races, err := racesRepo.List(nil, &sorting)
	if err != nil {
		t.Errorf("error when listing races, err: %s", err.Error())
		return
	}
	for i := 0; i < len(races)-1; i++ {
		if races[i].Id < races[i+1].Id {
			t.Error("list func failed to apply appropriate sorting")
			return
		}
	}
}

func TestGet(t *testing.T) {
	racingDB, _ := sql.Open("sqlite3", "./racing.db")
	racesRepo := NewRacesRepo(racingDB)
	if err := racesRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
		return
	}
	var testId int64 = 96

	race, err := racesRepo.Get(testId)
	if err != nil {
		t.Errorf("error when getting a race, err: %s", err.Error())
		return
	}

	if race == nil {
		fmt.Printf("no race with this id (%d) found \n", testId)
		return
	}

	if race.Id != testId {
		t.Errorf("get func failed to retrive a race by ID, expected id: %d, got: %d", testId, race.Id)
		return
	}
}
