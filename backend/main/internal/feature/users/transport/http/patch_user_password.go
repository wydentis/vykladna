package users_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_request "github.com/wydentis/vykladna/shared/core/transport_http/request"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	utils_validation "github.com/wydentis/vykladna/shared/utils/validation"
)

type PatchUserPasswordRequest struct {
	Password string `json:"password"`
}

func (r *PatchUserPasswordRequest) Validate() error {
	if err := utils_validation.ValidatePassword(r.Password); err != nil {
		return fmt.Errorf("'password' validation failed: %w", err)
	}

	return nil
}

func (h *UsersHTTPTransport) PatchUserPassword(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

	var request PatchUserPasswordRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	userInfo := core_auth.FromContext(ctx)

	err := h.usersService.PatchUserPassword(ctx, userInfo.ID, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	responseHandler.NoContentResponse()
}
