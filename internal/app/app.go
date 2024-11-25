package app

import (
	grpcapp "auth-microserivce/internal/app/grpc"
	"auth-microserivce/internal/config"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcCfg config.GRPCConfig, storageCfg config.StorageConfig, tokenTTL time.Duration, secret string) *App {
	grpcApp := grpcapp.NewApp(log, grpcCfg, storageCfg, secret, tokenTTL)

	return &App{
		GRPCServer: grpcApp,
	}
}
