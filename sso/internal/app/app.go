package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	//TODO: инициализировать хранилище (storege)

	//TODO: init auth server(auth)

	grpcApp := grpcapp.New(log, nil, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
