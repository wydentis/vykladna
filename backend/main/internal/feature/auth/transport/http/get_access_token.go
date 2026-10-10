package auth_transport_http

import (
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

type GetAccessTokenResponse TokenDTO

func (h *AuthTransportHTTP) GetAccessToken(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

	c, err := r.Cookie(refreshCookieName)
	if err != nil {
		err = fmt.Errorf("get refresh cookie: %v: %w", err, core_errors.ErrUnauthorized)
		responseHandler.ErrorResponse(err, "failed to get access token")
		return
	}

	access, err := h.authService.GetAccessToken(ctx, c.Value)
	if err != nil {
		if errors.Is(err, core_errors.ErrUnauthorized) {
			expiredCookie := newRefreshCookie("", -1)
			responseHandler.SetCookie(expiredCookie)
		}

		responseHandler.ErrorResponse(err, "failed to get access token")
		return
	}

	responseHandler.NoStore()

	response := tokenDTOFromAccess(access)
	responseHandler.JSONResponse(response, http.StatusOK)
}
