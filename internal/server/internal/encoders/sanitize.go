package encoders

import (
	"context"

	"github.com/pkg/errors"
	"github.com/slowaner/brackets-sanitizer/internal/server/internal/entities"
	sgrpc "github.com/slowaner/brackets-sanitizer/pkg/transport/grpc"
)

func EncodeSanitizeResponse(_ context.Context, r interface{}) (resp interface{}, err error) {
	res, ok := r.(entities.SanitizeResponse)
	if !ok {
		err = errors.New("invalid response type")
		return
	}

	resp = &sgrpc.SanitizeResponse{
		Result: res.Result(),
	}
	return
}
func DecodeSanitizeRequest(_ context.Context, r interface{}) (request interface{}, err error) {
	req, ok := r.(*sgrpc.SanitizeRequest)
	if !ok {
		err = errors.New("invalid request type")
		return
	}

	request = entities.NewSanitizeRequest(req.Input)
	return
}
