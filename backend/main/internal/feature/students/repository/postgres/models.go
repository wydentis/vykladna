package students_postgres_repository

import (
	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type StudentModel struct {
	ID      uuid.UUID
	Version int

	Name        string
	Surname     string
	PhoneNumber string
	TelegramID  *string
}

func studentDomainFromModel(model StudentModel) core_domains.Student {
	return core_domains.Student{
		ID:          model.ID,
		Version:     model.Version,
		Name:        model.Name,
		Surname:     model.Surname,
		PhoneNumber: model.PhoneNumber,
		TelegramID:  model.TelegramID,
	}
}
