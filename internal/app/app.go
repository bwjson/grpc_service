package app

import (
	grpcapp "github.com/bwjson/grpc_server/internal/app/grpc"
	"github.com/bwjson/grpc_server/internal/db/postgres"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort, db_port, db_pass, db_host, db_name, db_user string) *App {
	db, err := postgres.New(db_host, db_port, db_user, db_name, db_pass)
	if err != nil {
		panic(err)
	}

	// STUDENTS REPO TO GRPC APP
	grpcApp := grpcapp.New(log, grpcPort, db)

	return &App{GRPCServer: grpcApp}
}
