package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"github.com/slowaner/brackets-sanitizer/internal/endpoint/internal/entities"
)

type service interface {
	Validate(ctx context.Context, input string) (valid bool, err error)
	Sanitize(ctx context.Context, input string) (result string, err error)
}

type Endpoints interface {
	GetValidate() endpoint.Endpoint
	GetSanitize() endpoint.Endpoint
}

type endpoints struct {
	validate endpoint.Endpoint
	sanitize endpoint.Endpoint
}

func (e endpoints) GetValidate() endpoint.Endpoint {
	return e.validate
}

func (e endpoints) GetSanitize() endpoint.Endpoint {
	return e.sanitize
}

func makeValidateEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		req, ok := request.(entities.ValidateRequest)
		if !ok {
			err = errors.New("invalid request type")
			return
		}

		valid, err := s.Validate(ctx, req.GetInput())
		if err != nil {
			return
		}

		resp = entities.NewValidateResponse(valid)
		return
	}
}

func makeSanitizeEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		req, ok := request.(entities.SanitizeRequest)
		if !ok {
			err = errors.New("invalid request type")
			return
		}

		result, err := s.Sanitize(ctx, req.GetInput())
		if err != nil {
			return
		}

		resp = entities.NewSanitizeResponse(result)
		return
	}
}

func NewEndpoints(s service) Endpoints {
	return &endpoints{
		validate: makeValidateEndpoint(s),
		sanitize: makeSanitizeEndpoint(s),
	}
}
