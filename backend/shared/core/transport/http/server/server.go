package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_middleware "github.com/wydentis/vykladna/shared/core/transport_http/middleware"
)

type HTTPServer struct {
	mux    *http.ServeMux
	cfg    Config
	logger *core_logger.Logger

	middlewares []core_http_middleware.Middleware
}

func NewHTTPServer(cfg Config, log *core_logger.Logger, middlewares ...core_http_middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		cfg:         cfg,
		logger:      log,
		middlewares: middlewares,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: s.mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.logger.Warn(fmt.Sprintf("start HTTP server on %s:%d", s.cfg.Addr, s.cfg.Port))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		s.logger.Warn("shutdown HTTP server ...")

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.logger.Warn("HTTP server stopped")
	}

	return nil
}
