package students_postgres_repository

import (
	"context"
	"fmt"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

func (r *StudentsRepository) CreateStudent(ctx context.Context, student core_domains.Student) (core_domains.Student, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO students.accounts (name, surname, phone_number)
		VALUES ($1, $2, $3)
		RETURNING id, version, name, surname, phone_number, telegram_id
	`

	row := r.pool.QueryRow(ctx, query, student.Name, student.Surname, student.PhoneNumber)

	var studentModel StudentModel
	if err := row.Scan(
		&studentModel.ID,
		&studentModel.Version,
		&studentModel.Name,
		&studentModel.Surname,
		&studentModel.PhoneNumber,
		&studentModel.TelegramID,
	); err != nil {
		return core_domains.Student{}, fmt.Errorf("scan error: %w", err)
	}

	return studentDomainFromModel(studentModel), nil
}
