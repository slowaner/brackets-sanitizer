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

func (r *validateResponse) IsValid() bool {
	return r.valid
}

func NewValidateResponse(valid bool) ValidateResponse {
	return &validateResponse{
		valid: valid,
	}
}
