package di

import (
	"fmt"
	"reflect"
)

// Get dependency by generic definition, with generic with any type
//
// Parameter:
//   - container (di.Container) container
func Get[T any](container Container) (T, error) {
	name := reflect.TypeOf((*T)(nil)).Elem().String()
	dep, err := container.Get(name)

	if err != nil {
		return *new(T), fmt.Errorf("dependency %s not exists", name)
	}

	castedDep, ok := dep.(T)
	if !ok {
		return *new(T), fmt.Errorf("dependency %s is not %s", name, reflect.TypeOf((*T)(nil)).Elem().String())
	}

	return castedDep, nil
}
