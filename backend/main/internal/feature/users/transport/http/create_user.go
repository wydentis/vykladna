package users_transport_http

import (
	"fmt"
	"net/http"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_request "github.com/wydentis/vykladna/shared/core/transport_http/request"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	utils_validation "github.com/wydentis/vykladna/shared/utils/validation"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

func (r *CreateUserRequest) Validate() error {
	if err := utils_validation.ValidateUsername(r.Username); err != nil {
		return fmt.Errorf("'username' validation failed: %w", err)
	}
	if err := utils_validation.ValidateName(r.Name); err != nil {
		return fmt.Errorf("'name' validation failed: %w", err)
	}
	if err := utils_validation.ValidateSurname(r.Surname); err != nil {
		return fmt.Errorf("'surname' validation failed: %w", err)
	}
	if err := utils_validation.ValidatePassword(r.Password); err != nil {
		return fmt.Errorf("'password' validation failed: %w", err)
	}

	return nil
}

type CreateUserResponse UserDTO

func (h *UsersHTTPTransport) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

	var req CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	userDomain := core_domains.NewUserUninitialized(req.Username, req.Name, req.Surname)

	user, err := h.usersService.CreateUser(ctx, userDomain, req.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
