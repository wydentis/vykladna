package students_transport_http

import (
	"net/http"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_request "github.com/wydentis/vykladna/shared/core/transport_http/request"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
)

type PatchStudentTIDRequest struct {
	TelegramID *string `json:"telegram_id"`
}

type PatchStudentTIDResponse StudentDTO

func (h *StudentsHTTPTransport) PatchStudentTID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, logger)

	id, err := core_http_request.GetUUIDPathValue(r, idPathValueKey)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id from path")
		return
	}

	var request PatchStudentTIDRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	student, err := h.studentsService.PatchStudentTID(ctx, id, request.TelegramID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch student telegram id")
		return
	}

	response := PatchStudentResponse(studentDTOFromDomain(student))
	responseHandler.JSONResponse(response, http.StatusOK)
}
