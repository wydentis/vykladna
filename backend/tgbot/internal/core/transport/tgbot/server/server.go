package core_tgbot_server

import (
	"context"
	"reflect"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_middleware "github.com/wydentis/vykladna/shared/core/transport_http/middleware"
	core_http_server "github.com/wydentis/vykladna/shared/core/transport_http/server"
	core_tgbot_middleware "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/middleware"
	core_tgbot_types "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/types"
)

const APIVersionTelegram core_http_server.APIVersion = "telegram"

type Server struct {
	cfg         Config
	httpCfg     core_http_server.Config
	logger      *core_logger.Logger
	middlewares []core_tgbot_middleware.Middleware
	handlers    map[reflect.Type]core_tgbot_middleware.HandlerFunc
}

func NewServer(
	cfg Config,
	httpCfg core_http_server.Config,
	logger *core_logger.Logger,
	middlewares ...core_tgbot_middleware.Middleware,
) *Server {
	return &Server{
		cfg:         cfg,
		httpCfg:     httpCfg,
		logger:      logger,
		middlewares: middlewares,
		handlers:    map[reflect.Type]core_tgbot_middleware.HandlerFunc{},
	}
}

func (s *Server) RegisterRouters(routers ...*Router) {
	for _, r := range routers {
		for t, h := range r.handlers {
			if _, dup := s.handlers[t]; dup {
				panic("tgbot server: duplicate route for " + t.String())
			}
			s.handlers[t] = core_tgbot_middleware.Chain(h, r.middlewares...)
		}
	}
}

func (s *Server) Run(ctx context.Context) error {
	final := make(map[reflect.Type]core_tgbot_middleware.HandlerFunc, len(s.handlers))
	for t, h := range s.handlers {
		final[t] = core_tgbot_middleware.Chain(h, s.middlewares...)
	}

	dispatch := func(ctx context.Context, u core_tgbot_types.Update) error {
		if h, ok := final[u.Type()]; ok {
			return h(ctx, u)
		}

		return nil
	}

	p := newPool(dispatch, s.logger, s.cfg.Workers, s.cfg.QueueSize)
	p.start()

	apiRouter := core_http_server.NewAPIVersionRouter(APIVersionTelegram)
	apiRouter.RegisterRoutes(webhookRoute(s.cfg, router{pool: p, final: final}))

	httpServer := core_http_server.NewHTTPServer(
		s.httpCfg,
		s.logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(s.logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	httpServer.RegisterAPIVersionRouters(apiRouter)

	err := httpServer.Run(ctx)
	p.stop(s.httpCfg.ShutdownTimeout)
	return err
}
