package students_postgres_repository

import (
	"context"
	"fmt"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

func (r *StudentsRepository) GetStudents(ctx context.Context) ([]core_domains.Student, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, name, surname, phone_number, telegram_synced, telegram_id, owner_user_id
	FROM students.accounts
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select students: %w", err)
	}

	var studentModels []StudentModel
	for rows.Next() {
		var studentModel StudentModel

		if err := rows.Scan(
			&studentModel.ID,
			&studentModel.Version,
			&studentModel.Name,
			&studentModel.Surname,
			&studentModel.PhoneNumber,
			&studentModel.TelegramSynced,
			&studentModel.TelegramID,
			&studentModel.OwnerUserID,
		); err != nil {
			return nil, fmt.Errorf("scan students: %w", err)
		}

		studentModels = append(studentModels, studentModel)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	domains := studentDomainsFromModels(studentModels)
	return domains, nil
}
