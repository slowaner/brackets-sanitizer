package entities

type SanitizeRequest interface {
	GetInput() string
}

type SanitizeResponse interface {
	Result() string
}

type sanitizeResponse struct {
	result string
}

func (r *sanitizeResponse) Result() string {
	return r.result
}

func NewSanitizeResponse(result string) SanitizeResponse {
	return &sanitizeResponse{
		result: result,
	}
}
