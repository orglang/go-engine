package option

type ADT[T any] interface {
	IsEmpty() bool
	Get() T
}

type some[T any] struct {
	val T
}

func (some[T]) IsEmpty() bool {
	return false
}

func (adt some[T]) Get() T {
	return adt.val
}

func Some[T any](val T) ADT[T] {
	return some[T]{val}
}

type none[T any] struct{}

func (none[T]) IsEmpty() bool {
	return true
}

func (adt none[T]) Get() T {
	panic("nothing to return")
}

func None[T any]() ADT[T] {
	return none[T]{}
}
