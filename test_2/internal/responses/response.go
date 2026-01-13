package responses

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error   string       `json:"error"`
	Code    int          `json:"code"`
	Errors  []string     `json:"errors,omitempty"`  // Array of validation errors
	Details interface{}  `json:"details,omitempty"` // Additional details if needed
}

// ValidationErrorDetail represents a single validation error
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrorResponse represents a validation error response
type ValidationErrorResponse struct {
	Error  string                  `json:"error"`
	Code   int                     `json:"code"`
	Errors []ValidationErrorDetail `json:"errors"`
}

// NewSuccess creates a new success response
func NewSuccess(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
	}
}

// NewSuccessWithMeta creates a new success response with metadata
func NewSuccessWithMeta(message string, data interface{}, meta interface{}) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

// NewError creates a new error response
func NewError(errorMsg string, code int) ErrorResponse {
	return ErrorResponse{
		Error: errorMsg,
		Code:  code,
	}
}

// NewValidationError creates a new validation error response
func NewValidationError(errorMsg string, errors []ValidationErrorDetail) ValidationErrorResponse {
	return ValidationErrorResponse{
		Error:  errorMsg,
		Code:   422, // Unprocessable Entity
		Errors: errors,
	}
}