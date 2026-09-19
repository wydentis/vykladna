package students_service

import (
	"context"

	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type StudentsService struct {
	studentsRepository StudentsRepository
}

type StudentsRepository interface {
	GetStudent(ctx context.Context, id uuid.UUID) (core_domains.Student, error)
	CreateStudent(ctx context.Context, student core_domains.Student) (core_domains.Student, error)
	PatchStudent(ctx context.Context, student core_domains.Student) (core_domains.Student, error)
}

func NewStudentsService(studentsRepository StudentsRepository) *StudentsService {
	return &StudentsService{
		studentsRepository: studentsRepository,
	}
}
