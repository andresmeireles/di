package di

import (
	"fmt"
	"reflect"
	"strings"
)

// Container
type Container struct {
	container map[string]any
}

// Set function on container.
//
// Parameters:
//   - name (string) name of dependency
//   - implementation (any) initialized struct or value to add on container
func (c *Container) set(name string, implementation any) {
	if c.container == nil {
		c.container = map[string]any{}
	}
	c.container[name] = implementation
}

// Get dependency by name
//
// Parameter:
//   - name (string) name of dependency
func (c Container) Get(name string) (any, error) {
	value, ok := c.container[name]
	if !ok {
		return nil, fmt.Errorf("dependency %s not exists", name)
	}

	return value, nil
}

// Add dependency
//
// Parameter:
//   - dependency (di.Dependency) Dependency interface
func (c *Container) Add(dependency Dependency) error {
	implementation := dependency.Builder()
	name := dependency.Name()

	if len(strings.Trim(name, "")) == 0 {
		panic("name cannot be empty")
	}

	if reflect.TypeOf(implementation).Kind() != reflect.Func {
		panic(fmt.Sprintf("Implementation of %s is not a function.", name))
	}

	params := c.funcHasParams(implementation)

	if params == 0 {
		resolve, _ := c.resolve(implementation, 0)
		c.set(name, resolve)

		return nil
	}

	buildDep, err := c.resolve(implementation, params)
	if err != nil {
		return err
	}

	c.set(name, buildDep)

	return nil
}

// Check if implementation has params and return the number of parameters
//
// Parameter:
//   - implementation (any) implementation function
func (c Container) funcHasParams(implementation any) int {
	return reflect.TypeOf(implementation).NumIn()
}

// Resolve the implementation function
//
// Parameters:
//   - implementation (any) implementation function
//   - numberOfParams (int) number of parameters in implementation function
func (c *Container) resolve(implementation any, numberOfParams int) (any, error) {
	reference := reflect.TypeOf(implementation)
	paramValues := make([]reflect.Value, numberOfParams)

	for i := 0; i < numberOfParams; i++ {
		paramName := reference.In(i).String()
		param, err := c.Get(paramName)

		if err != nil {
			return nil, err
		}

		paramValues[i] = reflect.ValueOf(param)
	}

	resolvedImplementation := reflect.ValueOf(implementation).Call(paramValues)[0].Interface()

	return resolvedImplementation, nil
}
