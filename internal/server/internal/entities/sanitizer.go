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

func (v *validateRequest) GetInput() string {
	return v.input
}

func NewValidateRequest(input string) ValidateRequest {
	return &validateRequest{
		input: input,
	}
}
