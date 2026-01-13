package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    int         `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// HandleError sends a JSON error response with the appropriate status code
func HandleError(c *gin.Context, statusCode int, errorMsg string, details ...interface{}) {
	resp := ErrorResponse{
		Error: errorMsg,
		Code:  statusCode,
	}
	
	if len(details) > 0 {
		resp.Details = details[0]
	}
	
	c.JSON(statusCode, resp)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, errorMsg string, details ...interface{}) {
	HandleError(c, http.StatusBadRequest, errorMsg, details...)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, errorMsg string, details ...interface{}) {
	HandleError(c, http.StatusUnauthorized, errorMsg, details...)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, errorMsg string, details ...interface{}) {
	HandleError(c, http.StatusNotFound, errorMsg, details...)
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, errorMsg string, details ...interface{}) {
	HandleError(c, http.StatusConflict, errorMsg, details...)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, errorMsg string, details ...interface{}) {
	HandleError(c, http.StatusInternalServerError, errorMsg, details...)
}

// ValidationError sends a 422 Unprocessable Entity response for validation errors
func ValidationError(c *gin.Context, errorMsg string, details ...interface{}) {
	HandleError(c, http.StatusUnprocessableEntity, errorMsg, details...)
}