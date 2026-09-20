package students_transport_http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_http_server "github.com/wydentis/vykladna/shared/core/transport_http/server"
)

var (
	idPathValueKey = "id"
)

type StudentsHTTPTransport struct {
	studentsService StudentsService
}

type StudentsService interface {
	GetStudent(ctx context.Context, id uuid.UUID) (core_domains.Student, error)
	CreateStudent(ctx context.Context, student core_domains.Student) (core_domains.Student, error)
	PatchStudent(ctx context.Context, id uuid.UUID, patch core_domains.StudentPatch) (core_domains.Student, error)
	PatchStudentTID(ctx context.Context, id uuid.UUID, telegramID int64) (core_domains.Student, error)
}

func NewUsersHTTPTransport(usersService StudentsService) *StudentsHTTPTransport {
	return &StudentsHTTPTransport{
		studentsService: usersService,
	}
}

func (h *StudentsHTTPTransport) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    fmt.Sprintf("/students/{%s}", idPathValueKey),
			Handler: h.GetStudent,
		},
		{
			Method:  http.MethodPost,
			Path:    "/students",
			Handler: h.CreateStudent,
		},
		{
			Method:  http.MethodPatch,
			Path:    fmt.Sprintf("/students/{%s}", idPathValueKey),
			Handler: h.PatchStudent,
		},
		{
			Method:  http.MethodPatch,
			Path:    fmt.Sprintf("/students/{%s}/telegram-id", idPathValueKey),
			Handler: h.PatchStudentTID,
		},
	}
}
