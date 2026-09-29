package core_http_response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	core_errors "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/errors"
	core_logger "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  *gin.Context
}

func NewHTTPResponseHandler(
	log *core_logger.Logger,
	rw *gin.Context,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrForbidden):
		statusCode = http.StatusForbidden
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}

	h.rw.JSON(statusCode, response)
}
