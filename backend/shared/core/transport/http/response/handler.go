package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(rw http.ResponseWriter, log *core_logger.Logger) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.rw.Header().Set("content-type", "application/json")
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write http response", "err", err)
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, "err", err)
	h.errorResponse(err, http.StatusInternalServerError, msg)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...any)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
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

	logFunc(msg, "err", err)
	h.errorResponse(err, statusCode, msg)
}

func (h *HTTPResponseHandler) errorResponse(err error, statusCode int, msg string) {
	var response map[string]string

	if statusCode >= 500 {
		response = map[string]string{
			"message": msg,
			"error":   "internal server error",
		}
	} else {
		response = map[string]string{
			"message": msg,
			"error":   err.Error(),
		}
	}

	h.JSONResponse(response, statusCode)
}
