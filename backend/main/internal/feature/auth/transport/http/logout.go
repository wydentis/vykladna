package auth_transport_http

import (
	"net/http"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
)

func (h *AuthTransportHTTP) Logout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

	c := newRefreshCookie("", -1)
	responseHandler.SetCookie(c)

	responseHandler.NoStore()
	responseHandler.NoContentResponse()
}
