package core_domains

import "github.com/google/uuid"

// TODO: add validation on all layers

type Student struct {
	ID      uuid.UUID
	Version int

	Name        string
	Surname     string
	PhoneNumber string
	TelegramID  *string
}

func NewStudent(
	ID uuid.UUID,
	Version int,
	Name string,
	Surname string,
	PhoneNumber string,
	TelegramID *string,
) Student {
	return Student{
		ID:          ID,
		Version:     Version,
		Name:        Name,
		Surname:     Surname,
		PhoneNumber: PhoneNumber,
		TelegramID:  TelegramID,
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
		nil,
	)
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

func (s *Student) ApplyPatch(patch StudentPatch) {
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
}
