package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Status     string `json:"status"`          // "success", "error", "warning"
	Message    string `json:"message"`         // user-facing message
	Error      string `json:"error,omitempty"` // optional internal error
	StatusCode int    `json:"status_code"`     // HTTP status code
}

type APIResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

// Send Success
func Success(c *gin.Context, message string, data interface{}) {
	respond(c, "success", message, http.StatusOK, data, nil)
}

// Send Error
func Error(c *gin.Context, message string, statusCode int, err error) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	respond(c, "error", message, statusCode, nil, err)
}

// Send Bad Request
func Warning(c *gin.Context, message string, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusBadRequest
	}
	respond(c, "warning", message, statusCode, nil, nil)
}

// Core builder function
func respond(c *gin.Context, status string, message string, statusCode int, data interface{}, err error) {
	meta := Meta{
		Status:     status,
		Message:    message,
		StatusCode: statusCode,
	}

	if err != nil {
		meta.Error = err.Error()
	}

	c.JSON(statusCode, APIResponse{
		Meta: meta,
		Data: data,
	})
}
