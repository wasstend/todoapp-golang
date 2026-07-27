package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
	core_logger "github.com/wasstend/todoapp-golang/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(
	logger *core_logger.Logger,
	w http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: logger,
		w:   w,
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) JsonResponse(responseBody any, statusCode int) {
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		h.log.Error("write http response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrorConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrorInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrorNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug // Warn
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

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JsonResponse(response, statusCode)
}
