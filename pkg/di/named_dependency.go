package di

type namedDependency struct {
	name string
	builder any
}

func NewNamedDependency(name string, implementation any) namedDependency {
	return namedDependency{builder: implementation, name: name}
}

func (di namedDependency) Builder() any {
	return di.builder
}

func (di namedDependency) Name() string {
	return di.name
}

