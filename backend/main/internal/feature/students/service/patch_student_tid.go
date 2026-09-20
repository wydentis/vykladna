package students_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (s *StudentsService) PatchStudentTID(ctx context.Context, id uuid.UUID, telegramID int64) (core_domains.Student, error) {
	student, err := s.studentsRepository.GetStudent(ctx, id)
	if err != nil {
		return core_domains.Student{}, fmt.Errorf("get user: %w", err)
	}

	if student.TelegramSynced {
		return core_domains.Student{}, fmt.Errorf("telegram is already synced: %w", core_errors.ErrConflict)
	}

	student.TelegramSynced = true
	student.TelegramID = &telegramID

	student, err = s.studentsRepository.PatchStudent(ctx, student)
	if err != nil {
		return core_domains.Student{}, fmt.Errorf("patch user: %w", err)
	}

	return student, nil
}
