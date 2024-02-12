package main

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
}
