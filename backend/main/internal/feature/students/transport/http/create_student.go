package students_transport_http

import (
	"fmt"
	"net/http"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_request "github.com/wydentis/vykladna/shared/core/transport_http/request"
	core_http_response "github.com/wydentis/vykladna/shared/core/transport_http/response"
	utils_validation "github.com/wydentis/vykladna/shared/utils/validation"
)

type CreateStudentRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone_number"`
}

func (r *CreateStudentRequest) Validate() error {
	if err := utils_validation.ValidateName(r.Name); err != nil {
		return fmt.Errorf("'name' validation failed: %w", err)
	}
	if err := utils_validation.ValidateSurname(r.Surname); err != nil {
		return fmt.Errorf("'surname' validation failed: %w", err)
	}
	if err := utils_validation.ValidatePhoneNumber(r.PhoneNumber); err != nil {
		return fmt.Errorf("'phone_number' validation failed: %w", err)
	}

	return nil
}

type CreateStudentResponse StudentDTO

func (h *StudentsHTTPTransport) CreateStudent(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, log)

	var request CreateStudentRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	studentDomain := core_domains.NewStudentUninitialized(
		request.Name,
		request.Surname,
		request.PhoneNumber,
	)

	studentDomain, err := h.studentsService.CreateStudent(ctx, studentDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create student")
		return
	}

	response := CreateStudentResponse(studentDTOFromDomain(studentDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}
