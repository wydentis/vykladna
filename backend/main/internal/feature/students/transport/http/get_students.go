package students_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

type GetStudentsResponse []StudentDTO

func (h *StudentsHTTPTransport) GetStudents(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

	userInfo := core_auth.FromContext(ctx)

	var studentDomains []core_domains.Student
	var err error

	switch userInfo.Role {
	case core_domains.UserRoleTeacher:
		studentDomains, err = h.studentsService.GetStudentsByOwnerID(ctx, userInfo.ID)
	case core_domains.UserRoleAdmin:
		studentDomains, err = h.studentsService.GetStudents(ctx)
	default:
		studentDomains, err = []core_domains.Student{}, fmt.Errorf("invalid role: %q: %w", userInfo.Role, core_errors.ErrForbidden)
	}
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get students")
		return
	}

	response := GetStudentsResponse(studentDTOsFromDomains(studentDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}
