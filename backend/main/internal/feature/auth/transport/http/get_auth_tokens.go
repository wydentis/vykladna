package auth_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_request "github.com/wydentis/vykladna/shared/core/transport_http/request"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	utils_validation "github.com/wydentis/vykladna/shared/utils/validation"
)

type GetAuthTokensRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *GetAuthTokensRequest) Validate() error {
	if err := utils_validation.ValidateUsername(r.Username); err != nil {
		return fmt.Errorf("'username' validation failed: %w", err)
	} else if err := utils_validation.ValidatePassword(r.Password); err != nil {
		return fmt.Errorf("'password' validation failed: %w", err)
	}

	return nil
}

type GetAuthTokensResponse TokenDTO

func (h *AuthTransportHTTP) GetAuthTokens(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

	var request GetAuthTokensRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	tokens, err := h.authService.GetAuthTokens(ctx, request.Username, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get auth tokens")
		return
	}

	c := newRefreshCookie(tokens.Refresh, 7*24*60*60)
	responseHandler.SetCookie(c)
	responseHandler.NoStore()

	response := GetAuthTokensResponse(tokenDTOFromAccess(tokens.Access))
	responseHandler.JSONResponse(response, http.StatusOK)
}
