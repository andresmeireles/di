package di_test

import (
	"testing"

	"github.com/andresmeireles/di"
	"github.com/andresmeireles/di/testdata"
)

func TestContainerBuilder(t *testing.T) {
	deps := []di.Dependency{
		di.NewTypedDependency[testdata.TestOne](testdata.NewTestOne),
		di.NewTypedDependency[testdata.TestTwo](testdata.NewTestTwo),
		di.NewTypedDependency[testdata.TestThree](testdata.NewTestThree),
	}

	t.Run("create container", func(t *testing.T) {
		builder := di.NewContainerBuilder(deps, nil, nil)
		container := builder.Build()

		d1, err := container.Get("testdata.TestOne")
		if err != nil {
			t.Fatal(err)
		}

		cast, _ := d1.(testdata.TestOne)
		if cast.Name != "testone" {
			t.Fatalf("expected %s. received %s", "testone", cast.Name)
		}
	})

	t.Run("should break by max recursion", func(t *testing.T) {
		defer func() {
			p := recover()

			if p == nil {
				t.Fatal("should panic")
			}

			if p != "max recursion reached" {
				t.Fatalf("expected %s. received %s", "max recursion reached", p)
			}
		}()

		recursion := 1
		builder := di.NewContainerBuilder(deps, nil, &recursion)
		builder.Build()
	})

	t.Run("should stop by timeout", func(t *testing.T) {
		defer func() {
			p := recover()

			if p == nil {
				t.Fatal("should panic")
			}

			if p != "timeout" {
				t.Fatalf("expected: %s. received: %s", "timeout", p)
			}
		}()

		timeoutMilliseconds := 0.00001
		builder := di.NewContainerBuilder(deps, &timeoutMilliseconds, nil)
		builder.Build()
	})

	t.Run("should create even in mixed order", func(t *testing.T) {
		d := []di.Dependency{
			di.NewTypedDependency[testdata.TestFour](testdata.NewTestFour),
			di.NewTypedDependency[testdata.TestThree](testdata.NewTestThree),
			di.NewNamedDependency("testdata.TestTwo", testdata.NewTestTwo),
			di.NewTypedDependency[testdata.TestOne](testdata.NewTestOne),
		}

		builder := di.NewContainerBuilder(d, nil, nil)
		container := builder.Build()

		t3, _ := di.Get[testdata.TestThree](*container)

		if t3.One.Name != "testone" {
			t.Fatalf("expected %s. received %s", "testone", t3.One.Name)
		}
	})

	t.Run("should never create container", func(t *testing.T) {
		defer func() {
			p := recover()

			if p == nil {
				t.Fatal("should panic")
			}

			if p != "Cannot build container. Last dependency cannot be build, err dependency testdata.TestThree not exists" {
				t.Fatalf("expected %s. received %s", "Cannot build container. Last dependency cannot be build, err dependency testdata.TestThree not exists", p)
			}
		}()

		d := []di.Dependency{
			di.NewTypedDependency[testdata.TestFour](testdata.NewTestFour),
			di.NewTypedDependency[testdata.TestOne](testdata.NewTestOne),
		}

		builder := di.NewContainerBuilder(d, nil, nil)
		builder.Build()
	})
}
