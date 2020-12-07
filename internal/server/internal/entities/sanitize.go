package entities

type SanitizeRequest interface {
	GetInput() string
}

type SanitizeResponse interface {
	Result() string
}

type sanitizeRequest struct {
	input string
}

func (r *sanitizeRequest) GetInput() string {
	return r.input
}

func NewSanitizeRequest(input string) SanitizeRequest {
	return &sanitizeRequest{
		input: input,
	}
}
