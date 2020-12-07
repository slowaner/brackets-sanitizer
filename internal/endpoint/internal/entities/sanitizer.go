package entities

type ValidateRequest interface {
	GetInput() string
}

type ValidateResponse interface {
	IsValid() bool
}

type validateResponse struct {
	valid bool
}

func (v *validateResponse) IsValid() bool {
	return v.valid
}

func NewValidateResponse(valid bool) ValidateResponse {
	return &validateResponse{
		valid: valid,
	}
}
