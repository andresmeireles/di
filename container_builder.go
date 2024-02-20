package di

import (
	"context"
	"fmt"
	"time"
)

type buildResult struct {
	container    *Container
	hasError     bool
	errorMessage string
}

func fail(message string) buildResult {
	return buildResult{hasError: true, errorMessage: message, container: nil}
}

func ok(container *Container) buildResult {
	return buildResult{container: container, hasError: false, errorMessage: ""}
}

type containerBuilder struct {
	dependencies []Dependency
	timeout      *float64
	recursion    *int
	debug        bool
}

// Create new container builder struct
//
// Parameters:
//   - dependency ([]Dependency) slice of dependencies
//   - timeout (*float64) time in milliseconds to define a deadline to container build, if is nil, it will be 2 minutes
//   - recursion (*int) max number of times that builder function must execute before panic, if is nil, it will be 1000 times
//   - debug (bool) if is true, it will print debug messages
func NewContainerBuilder(
	dependency []Dependency,
	timeout *float64,
	recursion *int,
	debug bool,
) containerBuilder {
	rec := 1000
	if recursion != nil {
		rec = *recursion
	}

	deadline := 120000.0
	if timeout != nil {
		deadline = *timeout
	}

	return containerBuilder{
		dependencies: dependency,
		timeout:      &deadline,
		recursion:    &rec,
		debug:        debug,
	}
}

// Build container and return pointer to di.Container
func (db containerBuilder) Build() *Container {
	deps := db.dependencies
	recursion := *db.recursion
	duration := time.Duration(*db.timeout) * time.Millisecond
	deadline := time.Now().Add(duration)
	ctx, cancelFunc := context.WithDeadline(context.Background(), deadline)

	defer cancelFunc()

	resultChannel := make(chan buildResult, 1)

	go db.resolveDependencies(ctx, recursion, deps, resultChannel)

	build := <-resultChannel

	if build.hasError {
		panic(build.errorMessage)
	}

	return build.container
}

// Resolve multiple dependencies
//
// Parameters:
//   - ctx (context.Context) context
//   - maxRecursion (int) number of total recursions when resolve dependencies
//   - dependencies ([]di.Dependency) slice of dependencies to be resolved
//   - errorChannel (chan buildResult) channel with result struct
//   - containerChannel (chan *di.Container) channel with pointer to di.Container
func (db containerBuilder) resolveDependencies(
	ctx context.Context,
	maxRecursion int,
	dependencies []Dependency,
	resultChannel chan buildResult,
) {
	container := new(Container)
	numberOfRecursions := 0
	deps := dependencies
	index := 0
	maxIndex := len(dependencies)

	for len(deps) != 0 {
		if maxIndex == 1 {
			if err := container.Add(deps[0]); err != nil {
				resultChannel <- fail(fmt.Sprintf("Cannot build container. Last dependency cannot be build, err %s", err))

				return
			}
		}

		// reset index when reach last position and not been completed
		if index > maxIndex-1 {
			index = 0
		}

		select {
		case <-ctx.Done():
			resultChannel <- fail("timeout")

			return
		default:
			if db.debug {
				fmt.Printf("Resolving %s\n", deps[index].Name())
			}

			if numberOfRecursions >= maxRecursion {
				resultChannel <- fail("max recursion reached")

				return
			}

			err := container.Add(deps[index])

			if err != nil && db.debug {
				fmt.Printf("Cannot build %s, err %s\n", deps[index].Name(), err)
			}

			numberOfRecursions++

			// remove existing index, when succeeded
			if err == nil {
				if db.debug {
					fmt.Printf("Added %s\n", deps[index].Name())
				}

				deps = append(deps[:index], deps[index+1:]...)
				maxIndex--
				index = 0

				continue
			}

			if db.debug {
				for _, d := range deps {
					fmt.Println("Must resolve:", d.Name())
				}
			}

			index++
		}
	}

	resultChannel <- ok(container)
}
