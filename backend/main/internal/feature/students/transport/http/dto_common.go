package students_transport_http

import (
	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type StudentDTO struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`

	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	PhoneNumber string  `json:"phone_number"`
	TelegramID  *string `json:"telegram_id"`
}

func studentDTOFromDomain(student core_domains.Student) StudentDTO {
	return StudentDTO{
		ID:          student.ID,
		Version:     student.Version,
		Name:        student.Name,
		Surname:     student.Surname,
		PhoneNumber: student.PhoneNumber,
		TelegramID:  student.TelegramID,
	}
}
