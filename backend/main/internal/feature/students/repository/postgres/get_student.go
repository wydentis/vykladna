package students_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	core_postgres_pool "github.com/wydentis/vykladna/shared/core/db/postgres/pool"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (r *StudentsRepository) GetStudent(ctx context.Context, id uuid.UUID) (core_domains.Student, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, name, surname, phone_number, telegram_synced, telegram_id
		FROM students.accounts
		WHERE id=$1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var model StudentModel

	err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Name,
		&model.Surname,
		&model.PhoneNumber,
		&model.TelegramSynced,
		&model.TelegramID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return core_domains.Student{}, fmt.Errorf("student with id='%s': %w", id, core_errors.ErrNotFound)
		}
		return core_domains.Student{}, fmt.Errorf("scan error: %w", err)
	}

	return studentDomainFromModel(model), nil
}
