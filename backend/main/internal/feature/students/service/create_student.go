package students_service

import (
	"context"
	"fmt"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (s *StudentsService) CreateStudent(ctx context.Context, student core_domains.Student) (core_domains.Student, error) {
	if err := student.Validate(); err != nil {
		return core_domains.Student{}, fmt.Errorf("validate student: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	student, err := s.studentsRepository.CreateStudent(ctx, student)

	if err != nil {
		return core_domains.Student{}, fmt.Errorf("create student: %w", err)
	}

	return student, nil
}
