package encoders

import (
	"context"

	"github.com/pkg/errors"
	"github.com/slowaner/brackets-sanitizer/internal/server/internal/entities"
	sgrpc "github.com/slowaner/brackets-sanitizer/pkg/transport/grpc"
)

func EncodeValidateResponse(_ context.Context, r interface{}) (resp interface{}, err error) {
	res, ok := r.(entities.ValidateResponse)
	if !ok {
		err = errors.New("invalid response type")
		return
	}

	resp = &sgrpc.ValidateResponse{
		Valid: res.IsValid(),
	}
	return
}
func DecodeValidateRequest(_ context.Context, r interface{}) (request interface{}, err error) {
	req, ok := r.(*sgrpc.ValidateRequest)
	if !ok {
		err = errors.New("invalid request type")
		return
	}

	request = entities.NewValidateRequest(req.Input)
	return
}
