package main

import (
	"auth-microserivce/internal/app"
	"auth-microserivce/internal/config"
	"auth-microserivce/pkg/logger"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	time.Sleep(10 * time.Second)
	log := logger.SetupLogger(logger.Local, "./logs.log")
	log.Info("Starting Application")
	configPath := flag.String("c", "", "Path to the configuration file")
	flag.Parse()

	if *configPath != "" {
		log.Info("Trying to load configuration from", "file", *configPath)
	}
	cfg := config.MustLoad(*configPath)

	application := app.New(log, cfg.GRPCConfig, cfg.StorageConfig, cfg.TokenTTL, cfg.Secret)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop
	log.Info("Got signal", slog.String("signal", sig.String()))

	application.GRPCServer.Stop()
	log.Info("App stopped")
}
