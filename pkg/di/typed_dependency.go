package di

import "reflect"

type typedDependency[T any] struct {
	builder any
}

func NewTypedDependency[T any](builder any) typedDependency[T] {
	return typedDependency[T]{builder}
}

func (td typedDependency[T]) Builder() any {
	return td.builder
}

func (td typedDependency[T]) Name() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}
