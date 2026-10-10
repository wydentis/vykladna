package students_service

import (
	"context"
	"fmt"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

func (s *StudentsService) GetStudents(ctx context.Context) ([]core_domains.Student, error) {
	students, err := s.studentsRepository.GetStudents(ctx)
	if err != nil {
		return []core_domains.Student{}, fmt.Errorf("get students: %w", err)
	}

	return students, nil
}
