package core_domains

type Nullable[T any] struct {
	Value *T
	Set   bool
}

func NewNullable[T any](value *T, set bool) Nullable[T] {
	return Nullable[T]{
		Value: value,
		Set:   set,
	}
}

func NewNullNullable[T any]() Nullable[T] {
	return NewNullable[T](nil, false)
}

func NewSetNullable[T any](value *T) Nullable[T] {
	return NewNullable(value, true)
}
