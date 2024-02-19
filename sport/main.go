package main

import (
	"database/sql"
	"flag"
	"log"
	"net"

	"github.com/jie1311/Entain/sport/db"
	"github.com/jie1311/Entain/sport/proto/sport"
	"github.com/jie1311/Entain/sport/service"
	"google.golang.org/grpc"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9001", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s\n", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9001")
	if err != nil {
		return err
	}

	sportDB, err := sql.Open("sqlite3", "./db/sport.db")
	if err != nil {
		return err
	}

	eventsRepo := db.NewEventsRepo(sportDB)
	if err := eventsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	sport.RegisterSportServer(
		grpcServer,
		service.NewsportService(
			eventsRepo,
		),
	)

	log.Printf("gRPC server listening on: %s\n", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
