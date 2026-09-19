package students_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_postgres_pool "github.com/wydentis/vykladna/shared/core/db/postgres/pool"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (r *StudentsRepository) PatchStudent(ctx context.Context, student core_domains.Student) (core_domains.Student, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE students.accounts
		SET name = $1, surname = $2, phone_number = $3, telegram_id = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING id, version, name, surname, phone_number, telegram_id;
	`

	var studentModel StudentModel

	row := r.pool.QueryRow(
		ctx, query,
		student.Name,
		student.Surname,
		student.PhoneNumber,
		student.TelegramID,
		student.ID,
		student.Version,
	)

	if err := row.Scan(
		&studentModel.ID,
		&studentModel.Version,
		&studentModel.Name,
		&studentModel.Surname,
		&studentModel.PhoneNumber,
		&studentModel.TelegramID,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return core_domains.Student{}, fmt.Errorf("user with id = %s concurrently accessed: %w", &student.ID, core_errors.ErrConflict)
		}

		return core_domains.Student{}, fmt.Errorf("scan error: %w", err)
	}

	studentDomain := studentDomainFromModel(studentModel)
	return studentDomain, nil
}
