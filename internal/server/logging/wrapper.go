package logging

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	sgrpc "github.com/slowaner/brackets-sanitizer/pkg/transport/grpc"
)

type loggingWrapper struct {
	logger log.Logger
	srv    sgrpc.BracketsSanitizerServer
}

func (l *loggingWrapper) Validate(ctx context.Context, request *sgrpc.ValidateRequest) (resp *sgrpc.ValidateResponse, err error) {
	defer func(from time.Time) {
		_ = l.wrappedLogger(err).Log(
			"method", "Validate",
			"request", request,
			"response", resp,
			"executionTime", time.Since(from),
			"err", err,
		)
	}(time.Now())
	return l.srv.Validate(ctx, request)
}

func (l *loggingWrapper) wrappedLogger(err error) log.Logger {
	if err != nil {
		return level.Error(l.logger)
	}
	return l.logger
}

func NewSanitizerServerWrapper(
	logger log.Logger,
) func(server sgrpc.BracketsSanitizerServer) (srv sgrpc.BracketsSanitizerServer, err error) {
	return func(server sgrpc.BracketsSanitizerServer) (srv sgrpc.BracketsSanitizerServer, err error) {
		srv = &loggingWrapper{
			logger: logger,
			srv:    server,
		}
		return
	}
}
