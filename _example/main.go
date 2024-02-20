package main

import (
	"fmt"

	"github.com/andresmeireles/di"
)

type DependencyInterface interface {
	GetSum(num1, num2 int) int
}

type Person struct {
	Name string
	Age  int
}

func NewPerson() Person {
	return Person{"John", 25}
}

type DependencyImplementation struct{}

func (d DependencyImplementation) GetSum(num1, num2 int) int {
	return num1 + num2
}

type ComplexStruct struct {
	Imp1   DependencyInterface
	Person Person
}

func NewComplexStruct(imp1 DependencyInterface, person Person) ComplexStruct {
	return ComplexStruct{imp1, person}
}

func main() {
	deps := []di.Dependency{
		di.NewTypedDependency[ComplexStruct](NewComplexStruct),
		di.NewTypedDependency[DependencyInterface](func() DependencyImplementation { return DependencyImplementation{} }),
		di.NewTypedDependency[Person](NewPerson),
	}

	builder := di.NewContainerBuilder(deps, nil, nil, true)
	container := builder.Build()

	complexStruct, _ := di.Get[ComplexStruct](*container)

	fmt.Println("Sum:", complexStruct.Imp1.GetSum(6, 8), "Says", complexStruct.Person.Name, "with", complexStruct.Person.Age, "years old")
}
