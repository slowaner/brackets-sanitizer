package server

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/pkg/errors"
	"github.com/slowaner/brackets-sanitizer/internal/server/internal/encoders"
	sgrpc "github.com/slowaner/brackets-sanitizer/pkg/transport/grpc"
	"google.golang.org/grpc"
)

type endpoints interface {
	GetValidate() endpoint.Endpoint
}

type wrapper func(server sgrpc.BracketsSanitizerServer) (srv sgrpc.BracketsSanitizerServer, err error)

type Registrar interface {
	Register(ctx context.Context, grpcServer *grpc.Server) (err error)
	WrapServer(wrapper wrapper) (err error)
}

var _ sgrpc.BracketsSanitizerServer = (*sanitizerServer)(nil)

type sanitizerServer struct {
	validate gt.Handler
}

var _ Registrar = (*registrar)(nil)

type registrar struct {
	srv sgrpc.BracketsSanitizerServer
}

func (r *registrar) WrapServer(wrapper wrapper) (err error) {
	srv, err := wrapper(r.srv)
	if err != nil {
		return
	}

	r.srv = srv
	return
}

func (srv *sanitizerServer) Validate(ctx context.Context, request *sgrpc.ValidateRequest) (response *sgrpc.ValidateResponse, err error) {
	_, resp, err := srv.validate.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	response, ok := resp.(*sgrpc.ValidateResponse)
	if !ok {
		err = errors.New("invalid request type")
		return
	}
	return
}

func (r *registrar) Register(
	ctx context.Context,
	grpcServer *grpc.Server,
) (err error) {
	sgrpc.RegisterBracketsSanitizerServer(grpcServer, r.srv)
	return
}

func (r *registrar) GetServer() (srv sgrpc.BracketsSanitizerServer) {
	return r.srv
}

func New(ctx context.Context, endpoints endpoints) sgrpc.BracketsSanitizerServer {
	return &sanitizerServer{
		validate: gt.NewServer(
			endpoints.GetValidate(),
			encoders.DecodeValidateRequest,
			encoders.EncodeValidateResponse,
		),
	}
}

func NewRegistrar(ctx context.Context, endpoints endpoints) Registrar {
	return &registrar{
		srv: New(ctx, endpoints),
	}
}
