package students_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

func (s *StudentsService) PatchStudent(ctx context.Context, id uuid.UUID, patch core_domains.StudentPatch) (core_domains.Student, error) {
	student, err := s.studentsRepository.GetStudent(ctx, id)
	if err != nil {
		return core_domains.Student{}, fmt.Errorf("get user: %w", err)
	}

	student.ApplyPatch(patch)

	student, err = s.studentsRepository.PatchStudent(ctx, student)
	if err != nil {
		return core_domains.Student{}, fmt.Errorf("patch user: %w", err)
	}

	return student, nil
}
