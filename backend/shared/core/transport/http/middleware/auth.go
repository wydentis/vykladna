package core_http_middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

var (
	authHeader   = "Authorization"
	bearerScheme = "Bearer"
)

func bearerToken(r *http.Request) (string, bool) {
	scheme, token, ok := strings.Cut(r.Header.Get(authHeader), " ")
	if !ok || !strings.EqualFold(scheme, bearerScheme) {
		return "", false
	}

	return token, true
}

func Auth(authManager core_auth.Manager) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

				token, ok := bearerToken(r)
				if !ok {
					responseHandler.ErrorResponse(fmt.Errorf("%s: %w", "invalid bearer token", core_errors.ErrUnauthorized), "invalid token format")
				}

				claims, err := authManager.ParseAccess(token)
				if err != nil {
					responseHandler.ErrorResponse(fmt.Errorf("%v: %w", err, core_errors.ErrUnauthorized), "unauthorized")
					return
				}

				userID, err := uuid.Parse(claims.Subject)
				if err != nil {
					responseHandler.ErrorResponse(fmt.Errorf("%v: %w", err, core_errors.ErrUnauthorized), "invalid user id")
					return
				}
				userInfo := core_auth.NewUserInfo(userID, core_domains.UserRole(claims.Role))

				ctx = userInfo.ToContext(r.Context())

				h.ServeHTTP(rw, r.WithContext(ctx))
			},
		)
	}
}
