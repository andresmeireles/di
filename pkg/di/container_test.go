package di_test

import (
	"testing"

	"github.com/andresmeireles/di/pkg/di"
	"github.com/andresmeireles/di/testdata"
)

func TestContainer(t *testing.T) {
	t.Run("should create container", func(t *testing.T) {
		d1 := testdata.NewTestOne
		d2 := testdata.NewTestTwo

		de1 := di.NewTypedDependency[testdata.TestOne](d1)
		de2 := di.NewNamedDependency("testdata.TestTwo", d2)

		container := di.Container{}
		container.Add(de1)
		container.Add(de2)

		dep1, err := container.Get("testdata.TestOne")
		if err != nil {
			t.Fatal(err)
		}

		dep2, err := container.Get("testdata.TestTwo")
		if err != nil {
			t.Fatal(err)
		}

		cast1 := dep1.(testdata.TestOne)
		cast2 := dep2.(testdata.TestTwo)

		if cast1.Name != "testone" {
			t.Fatalf("Dep message must be %s, received %s", "testone", cast1.Name)
		}

		if cast2.Name != "test two" {
			t.Fatalf("Dep message must be %s, received %s", "test two", cast2.Name)
		}
	})

	t.Run("should create a dependecy with depdencies", func(t *testing.T) {
		d1 := testdata.NewTestOne
		d2 := testdata.NewTestTwo
		d3 := testdata.NewTestThree

		de1 := di.NewTypedDependency[testdata.TestOne](d1)
		de2 := di.NewNamedDependency("testdata.TestTwo", d2)
		de3 := di.NewTypedDependency[testdata.TestThree](d3)

		container := new(di.Container)
		container.Add(de1)
		container.Add(de2)
		container.Add(de3)

		dep3, err := container.Get("testdata.TestThree")
		if err != nil {
			t.Fatal(err)
		}

		cast3, _ := dep3.(testdata.TestThree)

		if cast3.One.Name != "testone" {
			t.Fatalf("Want %s, Received %s", "testone", cast3.One.Name)
		}
	})

	t.Run("should break when dependecy not exists", func(t *testing.T) {
		container := new(di.Container)
		_, err := container.Get("dep")

		if err == nil {
			t.Fatal("error must no be nil")
		}

		if err.Error() != "dependency dep not exists" {
			t.Fatalf("expected %s. received %s", "a", err.Error())
		}
	})

	t.Run("should panic when implementation is not a function", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				if err != "Implementation of name is not a function." {
					t.Fatalf("received %s. expected %s", err, "Implementation of name is not a function.")
				}
			}
		}()

		container := new(di.Container)
		de := di.NewNamedDependency("name", 2)
		container.Add(de)
	})

	t.Run("should break when name is empty", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				if err != "name cannot be empty" {
					t.Fatalf("received %s. expected %s", err, "name cannot be empty")
				}
			}
		}()

		container := new(di.Container)
		dep := di.NewNamedDependency("", func() int { return 1 })
		container.Add(dep)
	})

	t.Run("should break when name is long empty", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				if err != "name cannot be empty" {
					t.Fatalf("received %s. expected %s", err, "name cannot be empty")
				}
			}
		}()

		container := new(di.Container)
		dep := di.NewNamedDependency("                 ", func() int { return 1 })
		container.Add(dep)
	})
}
