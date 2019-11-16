package api

type response struct {
	Message string      `json:"message,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

// template for success response
func SuccessResponse(content interface{}) *response {
	return &response{
		"Success!",
		content,
	}
}

// template for error response
func ErrorResponse(e string) *response {
	return &response{
		e,
		nil,
	}
}
