package core_domains

import (
	"fmt"

	"github.com/google/uuid"
	utils_validation "github.com/wydentis/vykladna/shared/utils/validation"
)

type Student struct {
	ID      uuid.UUID
	Version int64

	Name           string
	Surname        string
	PhoneNumber    string
	TelegramSynced bool
	TelegramID     *int64
}

func NewStudent(
	ID uuid.UUID,
	Version int64,
	Name string,
	Surname string,
	PhoneNumber string,
	TelegramSynced bool,
	TelegramID *int64,
) Student {
	return Student{
		ID:             ID,
		Version:        Version,
		Name:           Name,
		Surname:        Surname,
		PhoneNumber:    PhoneNumber,
		TelegramSynced: TelegramSynced,
		TelegramID:     TelegramID,
	}
}

func NewStudentUninitialized(
	Name string,
	Surname string,
	PhoneNumber string,
) Student {
	return NewStudent(
		UninitializedID,
		UninitializedVersion,
		Name,
		Surname,
		PhoneNumber,
		false,
		nil,
	)
}

func (s *Student) Validate() error {
	var subErr error

	if err := utils_validation.ValidateName(s.Name); err != nil {
		subErr = err
	}
	if err := utils_validation.ValidateSurname(s.Surname); err != nil {
		subErr = err
	}
	if err := utils_validation.ValidatePhoneNumber(s.PhoneNumber); err != nil {
		subErr = err
	}

	if subErr != nil {
		return fmt.Errorf("'student' validation: %w", subErr)
	}

	return nil
}

type StudentPatch struct {
	Name        Nullable[string]
	Surname     Nullable[string]
	PhoneNumber Nullable[string]
}

func NewStudentPatch(
	Name Nullable[string],
	Surname Nullable[string],
	PhoneNumber Nullable[string],
) StudentPatch {
	return StudentPatch{
		Name:        Name,
		Surname:     Surname,
		PhoneNumber: PhoneNumber,
	}
}

func (s *StudentPatch) Validate() error {
	if s.Name.Set && s.Name.Value == nil {
		return fmt.Errorf("'name' name cannot be null")
	}
	if s.Surname.Set && s.Surname.Value == nil {
		return fmt.Errorf("'surname' name cannot be null")
	}
	if s.PhoneNumber.Set && s.PhoneNumber.Value == nil {
		return fmt.Errorf("'phone number' name cannot be null")
	}

	return nil
}

func (s *Student) ApplyPatch(patch StudentPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate student patch: %w", err)
	}

	tmp := *s

	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}
	if patch.Surname.Set {
		tmp.Surname = *patch.Surname.Value
	}
	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = *patch.PhoneNumber.Value
	}

	*s = tmp

	return s.Validate()
}
