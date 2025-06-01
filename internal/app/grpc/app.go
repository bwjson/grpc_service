package grpcapp

import (
	"fmt"
	"github.com/bwjson/grpc_server/internal/db/postgres"
	studentgrpc "github.com/bwjson/grpc_server/internal/gprc/student"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
	repo       postgres.StudentRepo
}

func New(log *slog.Logger, port string, repo postgres.StudentRepo) *App {
	gRPCServer := grpc.NewServer()

	studentgrpc.Register(gRPCServer, repo)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "gprcapp.Run"

	log := a.log.With(
		slog.String("op", op),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "gprcapp.Stop"

	a.log.With(
		slog.String("op", op),
	)

	a.gRPCServer.GracefulStop()
}
