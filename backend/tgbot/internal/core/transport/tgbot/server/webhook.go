package core_tgbot_server

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_middleware "github.com/wydentis/vykladna/shared/core/transport_http/middleware"
	core_http_server "github.com/wydentis/vykladna/shared/core/transport_http/server"
	core_tgbot_middleware "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/middleware"
	core_tgbot_types "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/types"
)

const (
	webhookPath          = "/webhook"
	maxBodySize          = 1 << 20
	telegramSecretHeader = "X-Telegram-Bot-Api-Secret-Token"
)

type submitter interface {
	Submit(core_tgbot_types.Update) error
}

type router struct {
	pool  *pool
	final map[reflect.Type]core_tgbot_middleware.HandlerFunc
}

func (r router) Submit(u core_tgbot_types.Update) error {
	if _, ok := r.final[u.Type()]; !ok {
		return nil
	}

	return r.pool.Submit(u)
}

func webhookRoute(cfg Config, s submitter) core_http_server.Route {
	route := core_http_server.Route{
		Method:  http.MethodPost,
		Path:    webhookPath,
		Handler: webhookHandler(s),
	}
	if cfg.Secret != "" {
		route.Middlewares = []core_http_middleware.Middleware{
			secretToken(cfg.Secret),
		}
	}

	return route
}

func webhookHandler(s submitter) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		log := core_logger.FromContext(r.Context())

		body, err := io.ReadAll(http.MaxBytesReader(rw, r.Body, maxBodySize))
		if err != nil {
			http.Error(rw, "bad request", http.StatusBadRequest)
			return
		}

		var u core_tgbot_types.Update
		if err := json.Unmarshal(body, &u); err != nil {
			log.Error("failed to decode telegram update", "err", err, "body", truncate(body, 2048))
			rw.WriteHeader(http.StatusOK)
			return
		}

		if u.Payload != nil {
			if err := s.Submit(u); err != nil {
				log.Error("submit telegram update", "update_id", u.UpdateID, "err", err)
				rw.WriteHeader(http.StatusServiceUnavailable)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
	}
}

func truncate(b []byte, n int) string {
	if len(b) > n {
		return string(b[:n]) + "..."
	}

	return string(b)
}

func secretToken(secret string) core_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			got := r.Header.Get(telegramSecretHeader)
			if subtle.ConstantTimeCompare([]byte(got), []byte(secret)) != 1 {
				rw.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(rw, r)
		})
	}
}
