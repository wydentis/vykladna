package students_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

func (s *StudentsService) GetStudent(ctx context.Context, id uuid.UUID) (core_domains.Student, error) {
	student, err := s.studentsRepository.GetStudent(ctx, id)

	if err != nil {
		return core_domains.Student{}, fmt.Errorf("get student by id: %w", err)
	}

	return student, nil
}
