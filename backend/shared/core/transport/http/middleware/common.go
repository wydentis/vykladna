package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
)

var (
	requestIDHeader  = "X-Request-ID"
	statusCodeHeader = "Status-Code"
)

func RequestID() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				id, err := uuid.NewV7()
				if err != nil {
					panic(err)
				}

				requestID := id.String()

				r.Header.Set(requestIDHeader, requestID)
				w.Header().Set(requestIDHeader, requestID)

				h.ServeHTTP(w, r)
			},
		)
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				requestID := r.Header.Get(requestIDHeader)

				l := log.With(
					"request_id", requestID,
					"method", r.Method,
					"path", r.URL.Path,
				)

				ctx := l.ToContext(r.Context())

				h.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

func Trace() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				log := core_logger.FromContext(r.Context())

				before := time.Now()

				log.Debug(
					"-> incoming http request",
					"timestamp", before.UTC(),
				)

				h.ServeHTTP(w, r)

				now := time.Now()

				log.Debug(
					"<- done http request",
					"timestamp", now.UTC(),
					"latency", now.Sub(before),
					"status_code", w.Header().Get(statusCodeHeader),
				)
			},
		)
	}
}

func Panic() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				responseHandler := core_http_response.NewHTTPResponseHandler(w, log)

				defer func() {
					if p := recover(); p != nil {
						responseHandler.PanicResponse(p, "panic occurred")
					}
				}()

				h.ServeHTTP(w, r)
			},
		)
	}
}
