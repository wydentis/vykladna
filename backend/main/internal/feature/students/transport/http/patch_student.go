package students_transport_http

import (
	"net/http"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_request "github.com/wydentis/vykladna/shared/core/transport_http/request"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	core_http_types "github.com/wydentis/vykladna/shared/core/transport_http/types"
)

type PathStudentRequest struct {
	Name        core_http_types.Nullable[string] `json:"name"`
	Surname     core_http_types.Nullable[string] `json:"surname"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
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
