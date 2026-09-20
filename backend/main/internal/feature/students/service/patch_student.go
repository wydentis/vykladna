package students_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (s *StudentsService) PatchStudent(ctx context.Context, id uuid.UUID, patch core_domains.StudentPatch) (core_domains.Student, error) {
	student, err := s.studentsRepository.GetStudent(ctx, id)
	if err != nil {
		return core_domains.Student{}, fmt.Errorf("get user: %w", err)
	}

	if err := student.ApplyPatch(patch); err != nil {
		return core_domains.Student{}, fmt.Errorf("apply student patch: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	student, err = s.studentsRepository.PatchStudent(ctx, student)
	if err != nil {
		return core_domains.Student{}, fmt.Errorf("patch user: %w", err)
	}

	return student, nil
}
