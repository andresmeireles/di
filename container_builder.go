package di

import (
	"context"
	"time"
)

type buildResult struct {
	hasError     bool
	errorMessage string
}

func fail(message string) buildResult {
	return buildResult{hasError: true, errorMessage: message}
}

func ok() buildResult {
	return buildResult{hasError: false, errorMessage: ""}
}

type containerBuilder struct {
	dependencies []Dependency
	timeout      *float64
	recursion    *int
}

// Create new container builder struct
//
// Parameters:
//   - dependency ([]Dependency) slice of dependencies
//   - timeout (*float64) time in milliseconds to define a deadline to container build, if is nil, it will be 2 minutes
//   - recursion (*int) max number of times that builder function must execute before panic, if is nil, it will be 1000 times
func NewContainerBuilder(
	dependency []Dependency,
	timeout *float64,
	recursion *int,
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

	containerChannel := make(chan *Container, 1)
	errorChannel := make(chan buildResult, 1)
	go db.resolveDependencies(ctx, recursion, deps, errorChannel, containerChannel)
	build := <-errorChannel

	if build.hasError {
		panic(build.errorMessage)
	}

	return <-containerChannel
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
	errorChannel chan buildResult,
	containerChannel chan *Container,
) {
	container := new(Container)
	numberOfRecursions := 0
	for _, dep := range dependencies {
		select {
		case <-ctx.Done():
			containerChannel <- nil
			errorChannel <- fail("timeout")
			return
		default:
			if numberOfRecursions >= maxRecursion {
				containerChannel <- nil
				errorChannel <- fail("max recursion reached")
				return
			}

			container.Add(dep)
			numberOfRecursions++
		}
	}

	errorChannel <- ok()
	containerChannel <- container
}
