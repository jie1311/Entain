package db

import (
	"database/sql"
	"testing"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

func TestInit(t *testing.T) {
	racingDB, _ := sql.Open("sqlite3", "./racing.db")
	racesRepo := NewRacesRepo(racingDB)

	if err := racesRepo.Init(); err != nil {
		t.Errorf("error when initiating db, err: %s", err.Error())
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
			t.Error("list func didn't apply appropriate filters")
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
		SortBy:  1, // sort by ID
		Descend: true,
	}

	races, err := racesRepo.List(nil, &sorting)
	if err != nil {
		t.Errorf("error when listing races, err: %s", err.Error())
		return
	}
	for i := 0; i < len(races)-1; i++ {
		if races[i].Id < races[i+1].Id {
			t.Errorf("list func didn't apply appropriate sorting")
			return
		}
	}
}
