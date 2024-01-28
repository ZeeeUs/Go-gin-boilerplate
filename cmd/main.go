package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZeeeUs/Go-gin-boilerplate/internal/config"
	"github.com/ZeeeUs/Go-gin-boilerplate/internal/service"
	"github.com/ZeeeUs/Go-gin-boilerplate/internal/storage"
	"github.com/ZeeeUs/Go-gin-boilerplate/internal/transport"
	"github.com/ZeeeUs/Go-gin-boilerplate/internal/transport/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		panic(err)
	}
	logger := cfg.Logger()

	pgConn, err := pgxpool.NewWithConfig(context.Background(), cfg.PgxConnConfig())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	strg := storage.New(pgConn)
	svc := service.New(logger, strg)
	hands := handlers.New(logger, svc)
	srv := transport.New(logger, cfg.Server.Host, hands)

	go func() {
		srv.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
