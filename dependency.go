package di

type Dependency interface {
	Builder() any
	Name() string
}
