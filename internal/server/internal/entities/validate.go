package entities

type ValidateRequest interface {
	GetInput() string
}

type ValidateResponse interface {
	IsValid() bool
}

type validateRequest struct {
	input string
}

func (r *validateRequest) GetInput() string {
	return r.input
}

func NewValidateRequest(input string) ValidateRequest {
	return &validateRequest{
		input: input,
	}
}
