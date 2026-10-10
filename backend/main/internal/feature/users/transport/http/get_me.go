package users_transport_http

import (
	"net/http"

	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
)

type GetMeResponse UserDTO

func (h *UsersHTTPTransport) GetMe(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

	userInfo := core_auth.FromContext(ctx)

	user, err := h.usersService.GetUser(ctx, userInfo.ID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	response := UserDTO(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}
