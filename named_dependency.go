package di

// Dependency with name and implementation
type namedDependency struct {
	name    string
	builder any
}

// Create new namedDependency
//
// Parameters:
//   - name (string) name of dependency
//   - implementation (any) implementation function
func NewNamedDependency(name string, implementation any) namedDependency {
	return namedDependency{builder: implementation, name: name}
}

// Get builder definition
func (di namedDependency) Builder() any {
	return di.builder
}

// Get dependency name
func (di namedDependency) Name() string {
	return di.name
}
