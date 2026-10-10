package students_transport_http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_http_middleware "github.com/wydentis/vykladna/shared/core/transport_http/middleware"
	core_http_server "github.com/wydentis/vykladna/shared/core/transport_http/server"
)

var (
	idPathValueKey = "id"
)

type StudentsHTTPTransport struct {
	studentsService StudentsService
	authManager     *core_auth.Manager
}

type StudentsService interface {
	GetStudent(ctx context.Context, id uuid.UUID) (core_domains.Student, error)
	GetStudents(ctx context.Context) ([]core_domains.Student, error)
	GetStudentsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]core_domains.Student, error)
	CreateStudent(ctx context.Context, student core_domains.Student) (core_domains.Student, error)
	PatchStudent(ctx context.Context, id uuid.UUID, patch core_domains.StudentPatch) (core_domains.Student, error)
	PatchStudentTID(ctx context.Context, id uuid.UUID, telegramID int64) (core_domains.Student, error)
}

func NewUsersHTTPTransport(usersService StudentsService, authManager *core_auth.Manager) *StudentsHTTPTransport {
	return &StudentsHTTPTransport{
		studentsService: usersService,
		authManager:     authManager,
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
			Method:  http.MethodGet,
			Path:    "/students",
			Handler: h.GetStudents,
			Middlewares: []core_http_middleware.Middleware{
				core_http_middleware.Auth(*h.authManager),
			},
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
