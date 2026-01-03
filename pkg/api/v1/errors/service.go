package errors

type APIError struct {
	OriginalError error
	ErrorMessage  string
}

func NewAPIError(originalError error, errorMessage string) *APIError {
	return &APIError{
		OriginalError: originalError,
		ErrorMessage:  errorMessage,
	}
}
