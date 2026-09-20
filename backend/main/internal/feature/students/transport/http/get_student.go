package students_transport_http

import (
	"net/http"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_request "github.com/wydentis/vykladna/shared/core/transport_http/request"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
)

type GetStudentResponse StudentDTO

func (h *StudentsHTTPTransport) GetStudent(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, logger)

	id, err := core_http_request.GetUUIDPathValue(r, idPathValueKey)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id from path")
		return
	}

	student, err := h.studentsService.GetStudent(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get student")
		return
	}

	response := GetStudentResponse(studentDTOFromDomain(student))

	responseHandler.JSONResponse(response, http.StatusOK)
}
