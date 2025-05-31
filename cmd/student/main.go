package main

import (
	"github.com/bwjson/grpc_server/internal/app"
	"github.com/bwjson/grpc_server/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// config
	cfg := config.ParseConfig()

	// logger
	log := setupLogger(cfg.Env)
	log.Info("Logger started")

	// app
	application := app.New(
		log,
		cfg.GRPC.Port,
		cfg.Postgres.Port,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Name,
		cfg.Postgres.User,
	)

	// gRPC server
	go func() {
		if err := application.GRPCServer.Run(); err != nil {
			log.Error("Cannot start gRPC server")
			panic(err)
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()
	log.Info("GRPC server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	}

	return log
}
