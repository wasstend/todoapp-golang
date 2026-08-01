package core_http_response

type ErrorResponse struct {
	Error   string `json:"error" example:"error"`
	Message string `json:"message" example:"short human-readable error message"`
}
