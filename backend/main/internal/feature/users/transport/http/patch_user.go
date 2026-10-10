package users_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_request "github.com/wydentis/vykladna/shared/core/transport_http/request"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	core_http_types "github.com/wydentis/vykladna/shared/core/transport_http/types"
	utils_validation "github.com/wydentis/vykladna/shared/utils/validation"
)

type PatchUserRequest struct {
	Name    core_http_types.Nullable[string] `json:"name"`
	Surname core_http_types.Nullable[string] `json:"surname"`
}

func (r *PatchUserRequest) Validate() error {
	if r.Name.Set {
		if r.Name.Value == nil {
			return fmt.Errorf("'name' cannot be null")
		}
		if err := utils_validation.ValidateName(*r.Name.Value); err != nil {
			return fmt.Errorf("'name' validation failed: %w", err)
		}
	}
	if r.Surname.Set {
		if r.Surname.Value == nil {
			return fmt.Errorf("'surname' cannot be null")
		}
		if err := utils_validation.ValidateSurname(*r.Surname.Value); err != nil {
			return fmt.Errorf("'surname' validation failed: %w", err)
		}
	}

	return nil
}

type PatchUserResponse UserDTO

func (h *UsersHTTPTransport) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	userInfo := core_auth.FromContext(ctx)
	patch := userPatchFromRequest(request)

	user, err := h.usersService.PatchUser(ctx, userInfo.ID, patch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) core_domains.UserPatch {
	return core_domains.NewUserPatch(
		request.Name.ToDomain(),
		request.Surname.ToDomain(),
	)
}
