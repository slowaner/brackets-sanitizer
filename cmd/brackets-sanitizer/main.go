package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/slowaner/brackets-sanitizer/internal/endpoint"
	"github.com/slowaner/brackets-sanitizer/internal/server"
	"github.com/slowaner/brackets-sanitizer/internal/server/logging"
	"github.com/slowaner/brackets-sanitizer/internal/server/metrics"
	"github.com/slowaner/brackets-sanitizer/internal/service"
	grpcserver "github.com/slowaner/grpc-server/server"
)

//go:generate protoc -I ../../pkg/transport/grpc/src --go_out=plugins=grpc:../../pkg/transport/grpc ../../pkg/transport/grpc/src/sanitizer.proto

const (
	metricsNamespace    = "slowaner"
	metricsSubsystem    = "brackets_sanitizer"
	metricsNameCount    = "request_count"
	metricsHelpCount    = "Request count"
	metricsNameDuration = "request_duration"
	metricsHelpDuration = "Request duration"
)

var (
	methodError = []string{"method", "error"}
)

func main() {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.NewSyncLogger(logger)
		logger = level.Info(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = kitlog.With(logger,
			"svc", "brackets-sanitizer",
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

	/* Prometheus */
	counter := kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystem,
		Name:      metricsNameCount,
		Help:      metricsHelpCount,
	}, methodError)
	histogram := kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystem,
		Name:      metricsNameDuration,
		Help:      metricsHelpDuration,
	}, methodError)

	svc := service.NewService()
	endpoints := endpoint.NewEndpoints(svc)
	sanitizerServer := server.NewRegistrar(ctx, endpoints)
	err := sanitizerServer.WrapServer(logging.NewSanitizerServerWrapper(logger))
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	err = sanitizerServer.WrapServer(metrics.NewSanitizerServerWrapper(counter, histogram))
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

	r := http.NewServeMux()
	r.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer close(signalChan)
		<-signalChan

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			_ = level.Error(logger).Log(
				"msg", "can't stop http server",
				"err", err,
			)
		}
	}()

	go func() {
		err = httpServer.ListenAndServe()
		if err != nil {
			_ = level.Error(logger).Log(
				"msg", "http server error",
				"err", err,
			)
		}
	}()

	for err := range errChan {
		if err != nil {
			_ = level.Error(logger).Log("err", err)
		}
	}

	_ = level.Info(logger).Log("msg", "service ended")
}
