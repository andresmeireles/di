package di

import "reflect"

type typedDependency[T any] struct {
	builder any
}

// Create new dependency with generic type
//
// Parameter:
//   - builder (any) implementation function
func NewTypedDependency[T any](builder any) typedDependency[T] {
	return typedDependency[T]{builder}
}

// Get dependency build implementation
func (td typedDependency[T]) Builder() any {
	return td.builder
}

// Get dependency name
func (td typedDependency[T]) Name() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}
