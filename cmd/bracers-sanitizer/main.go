package main

import (
	"context"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/slowaner/bracers-sanitizer/internal/server/logging"

	"github.com/slowaner/bracers-sanitizer/internal/endpoint"
	"github.com/slowaner/bracers-sanitizer/internal/server"
	"github.com/slowaner/bracers-sanitizer/internal/service"
	grpcserver "github.com/slowaner/grpc-server/server"
)

//go:generate protoc -I ../../pkg/transport/grpc/src --go_out=plugins=grpc:../../pkg/transport/grpc ../../pkg/transport/grpc/src/sanitizer.proto

func main() {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.NewSyncLogger(logger)
		logger = level.Info(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = kitlog.With(logger,
			"svc", "bracers-sanitizer",
			"ts", kitlog.DefaultTimestampUTC,
			"caller", kitlog.DefaultCaller,
		)
	}

	_ = level.Info(logger).Log("msg", "service started")

	grpcConfig := grpcserver.Config{
		GrpcEnable:    true,
		GrpcHost:      "127.0.0.1",
		GrpcPort:      24560,
		GrpcEnableTls: false,
	}

	ctx := context.Background()

	svc := service.NewService()
	endpoints := endpoint.NewEndpoints(svc)
	sanitizerServer := server.NewRegistrar(ctx, endpoints)
	err := sanitizerServer.WrapServer(logging.NewSanitizerServerWrapper(logger))
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	srv, err := grpcserver.NewGrpcForServers(
		ctx,
		grpcConfig,
		[]grpcserver.Registrar{
			sanitizerServer,
		},
	)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	errChan := srv.Serve()

	for err := range errChan {
		if err != nil {
			_ = level.Error(logger).Log("err", err)
		}
	}

	_ = level.Info(logger).Log("msg", "service ended")
}
