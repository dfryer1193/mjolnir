package optional

type Option[T any] struct {
	value T
	empty bool
}

func New[T any](value T) Option[T] {
	return Option[T]{value: value, empty: false}
}

func Empty[T any]() Option[T] {
	return Option[T]{empty: true}
}

func Some[T any](value T) Option[T] {
	return Option[T]{value: value, empty: false}
}

func (o Option[T]) IsEmpty() bool {
	return o.empty
}

func (o Option[T]) Get() T {
	return o.value
}
