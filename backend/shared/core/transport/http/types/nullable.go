package core_http_types

import (
	"encoding/json"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type Nullable[T any] struct {
	core_domains.Nullable[T]
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Set = true

	if string(data) == "null" {
		n.Value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Value = &value
	return nil
}

func (n Nullable[T]) ToDomain() core_domains.Nullable[T] {
	return core_domains.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
