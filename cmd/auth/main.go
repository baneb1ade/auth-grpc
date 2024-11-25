package main

import (
	"auth-microserivce/internal/app"
	"auth-microserivce/internal/config"
	"auth-microserivce/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := logger.SetupLogger(logger.Local, "./logs.log")
	log.Info("Starting Application")
	cfg := config.MustLoad()

	application := app.New(log, cfg.GRPCConfig, cfg.StorageConfig, cfg.TokenTTL, cfg.Secret)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop
	log.Info("Got signal", slog.String("signal", sig.String()))

	application.GRPCServer.Stop()
	log.Info("App stopped")
}