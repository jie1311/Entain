package service

import (
	"github.com/jie1311/Entain/sport/db"
	"github.com/jie1311/Entain/sport/proto/sport"
	"golang.org/x/net/context"
)

type Sport interface {
	// ListEvents will return a collection of events.
	ListEvents(ctx context.Context, in *sport.ListEventsRequest) (*sport.ListEventsResponse, error)
}

// sportService implements the sport interface.
type sportService struct {
	eventsRepo db.EventsRepo
}

// NewsportService instantiates and returns a new sportService.
func NewsportService(eventsRepo db.EventsRepo) Sport {
	return &sportService{eventsRepo}
}

func (s *sportService) ListEvents(ctx context.Context, in *sport.ListEventsRequest) (*sport.ListEventsResponse, error) {
	events, err := s.eventsRepo.List(in.Filter, in.Sorting)
	if err != nil {
		return nil, err
	}

	return &sport.ListEventsResponse{Events: events}, nil
}
