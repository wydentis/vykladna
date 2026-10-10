package students_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

func (s *StudentsService) GetStudentsByOwnerID(ctx context.Context, userID uuid.UUID) ([]core_domains.Student, error) {
	students, err := s.studentsRepository.GetStudentsByOwnerID(ctx, userID)
	if err != nil {
		return []core_domains.Student{}, fmt.Errorf("get students: %w", err)
	}

	return students, nil
}
