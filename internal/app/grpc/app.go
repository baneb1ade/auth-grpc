package grpcapp

import (
	"auth-microserivce/internal/config"
	auth2 "auth-microserivce/internal/domain/auth"
	psql2 "auth-microserivce/internal/domain/auth/db"
	"auth-microserivce/internal/grpc/auth"
	"auth-microserivce/pkg/client/psql"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"time"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func NewApp(log *slog.Logger, grpcCfg config.GRPCConfig, storageCfg config.StorageConfig, secret string, tokenTTL time.Duration) *App {
	op := "app.grpc.app.NewAPp"
	logger := log.With(slog.String("op", op))

	logger.Info("Initializing GRPC server")
	gRPCServer := grpc.NewServer()
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		storageCfg.StorageUser,
		storageCfg.StoragePass,
		storageCfg.StorageHost,
		storageCfg.StoragePort,
		storageCfg.StorageDatabase,
	)
	client, err := psql.NewClient(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	storage := psql2.NewStorage(client, log)
	service := auth2.NewService(log, storage, tokenTTL, secret)
	auth.Register(gRPCServer, service)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       grpcCfg.Port,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port))

	lis, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("address", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
