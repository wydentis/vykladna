package students_transport_http

import (
	"fmt"
	"net/http"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_request "github.com/wydentis/vykladna/shared/core/transport_http/request"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	core_http_types "github.com/wydentis/vykladna/shared/core/transport_http/types"
	utils_validation "github.com/wydentis/vykladna/shared/utils/validation"
)

type PathStudentRequest struct {
	Name        core_http_types.Nullable[string] `json:"name"`
	Surname     core_http_types.Nullable[string] `json:"surname"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PathStudentRequest) Validate() error {
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
	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value == nil {
			return fmt.Errorf("'phone_number' cannot be null")
		}
		if err := utils_validation.ValidatePhoneNumber(*r.PhoneNumber.Value); err != nil {
			return fmt.Errorf("'phone_number' validation failed: %w", err)
		}
	}

	return nil
}

type PatchStudentResponse StudentDTO

func (h *StudentsHTTPTransport) PatchStudent(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, logger)

	id, err := core_http_request.GetUUIDPathValue(r, idPathValueKey)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id from path")
		return
	}

	var request PathStudentRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	patch := userPatchFromRequest(request)

	student, err := h.studentsService.PatchStudent(ctx, id, patch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchStudentResponse(studentDTOFromDomain(student))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PathStudentRequest) core_domains.StudentPatch {
	return core_domains.NewStudentPatch(
		request.Name.ToDomain(),
		request.Surname.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
