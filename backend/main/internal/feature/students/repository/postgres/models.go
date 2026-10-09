package students_postgres_repository

import (
	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type StudentModel struct {
	ID      uuid.UUID
	Version int64

	Name           string
	Surname        string
	PhoneNumber    string
	TelegramSynced bool
	TelegramID     *int64

	OwnerUserID uuid.UUID
}

func studentDomainFromModel(model StudentModel) core_domains.Student {
	return core_domains.Student{
		ID:             model.ID,
		Version:        model.Version,
		Name:           model.Name,
		Surname:        model.Surname,
		PhoneNumber:    model.PhoneNumber,
		TelegramSynced: model.TelegramSynced,
		TelegramID:     model.TelegramID,
		OwnerUserID:    model.OwnerUserID,
	}
}

func studentDomainsFromModels(models []StudentModel) []core_domains.Student {
	domains := make([]core_domains.Student, len(models))
	for i, v := range models {
		domains[i] = studentDomainFromModel(v)
	}

	return domains
}
