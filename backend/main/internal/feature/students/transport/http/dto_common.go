package students_transport_http

import (
	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type StudentDTO struct {
	ID      uuid.UUID `json:"id"`
	Version int64     `json:"version"`

	Name           string `json:"name"`
	Surname        string `json:"surname"`
	PhoneNumber    string `json:"phone_number"`
	TelegramSynced bool   `json:"telegram_synced"`
	TelegramID     *int64 `json:"telegram_id"`

	OwnerUserID uuid.UUID `json:"owner_user_id"`
}

func studentDTOFromDomain(student core_domains.Student) StudentDTO {
	return StudentDTO{
		ID:             student.ID,
		Version:        student.Version,
		Name:           student.Name,
		Surname:        student.Surname,
		PhoneNumber:    student.PhoneNumber,
		TelegramSynced: student.TelegramSynced,
		TelegramID:     student.TelegramID,
		OwnerUserID:    student.OwnerUserID,
	}
}

func studentDTOsFromDomains(students []core_domains.Student) []StudentDTO {
	dtos := make([]StudentDTO, len(students))
	for i, v := range students {
		dtos[i] = studentDTOFromDomain(v)
	}

	return dtos
}
