package errors

import (
	"codim/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error" example:"Error message describing what went wrong"`
	Code  string `json:"code,omitempty" example:"VALIDATION_ERROR"`
}

func HandleError(ctx *gin.Context, log *logger.Logger, err *APIError, statusCode int) {
	if err.OriginalError != nil {
		log.WithFields(map[string]interface{}{
			"error":       err.ErrorMessage,
			"status_code": statusCode,
		}).Error(err.OriginalError)

		if ctx != nil {
			ctx.JSON(statusCode, ErrorResponse{
				Error: err.ErrorMessage,
			})
		}
	} else {
		log.WithFields(map[string]interface{}{
			"status_code": statusCode,
		}).Error(err.OriginalError)

		if ctx != nil {
			ctx.JSON(statusCode, ErrorResponse{
				Error: err.ErrorMessage,
			})
		}
	}
}
