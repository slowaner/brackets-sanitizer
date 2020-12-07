package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
	sgrpc "github.com/slowaner/brackets-sanitizer/pkg/transport/grpc"
)

type metricsWrapper struct {
	reqCount    metrics.Counter
	reqDuration metrics.Histogram

	srv sgrpc.BracketsSanitizerServer
}

func (w *metricsWrapper) Validate(ctx context.Context, request *sgrpc.ValidateRequest) (resp *sgrpc.ValidateResponse, err error) {
	defer w.recordMetrics("Validate", time.Now(), err)
	return w.srv.Validate(ctx, request)
}

func (w *metricsWrapper) recordMetrics(method string, startTime time.Time, err error) {
	labels := []string{
		"method", method,
		"error", strconv.FormatBool(err != nil),
	}
	w.reqCount.With(labels...).Add(1)
	w.reqDuration.With(labels...).Observe(time.Since(startTime).Seconds())
}

func NewSanitizerServerWrapper(
	reqCount metrics.Counter,
	reqDuration metrics.Histogram,
) func(server sgrpc.BracketsSanitizerServer) (srv sgrpc.BracketsSanitizerServer, err error) {
	return func(server sgrpc.BracketsSanitizerServer) (srv sgrpc.BracketsSanitizerServer, err error) {
		srv = &metricsWrapper{
			reqCount:    reqCount,
			reqDuration: reqDuration,

			srv: server,
		}
		return
	}
}
