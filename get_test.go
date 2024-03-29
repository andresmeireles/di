package di_test

import (
	"testing"

	"github.com/andresmeireles/di"
	"github.com/andresmeireles/di/testdata"
)

func TestGet(t *testing.T) {
	deps := []di.Dependency{
		di.NewTypedDependency[testdata.TestOne](testdata.NewTestOne),
		di.NewTypedDependency[testdata.TestTwo](testdata.NewTestTwo),
		di.NewTypedDependency[testdata.TestThree](testdata.NewTestThree),
	}
	builder := di.NewContainerBuilder(deps, nil, nil, false)
	container := builder.Build()

	t.Run("should get testone by his type", func(t *testing.T) {
		testOne, err := di.Get[testdata.TestOne](*container)

		if err != nil {
			t.Fatalf("Expected nil. Received: %s", err)
		}

		if testOne.Name != "testone" {
			t.Fatalf("Expected: testone. Received: %s", testOne.Name)
		}
	})

	t.Run("should get complex struct by his type", func(t *testing.T) {
		t3, err := di.Get[testdata.TestThree](*container)

		if err != nil {
			t.Fatalf("Expected nil. Received: %s", err)
		}

		if t3.Two.Name != "test two" {
			t.Fatalf("Expected: test two. Received: %s", t3.Two.Name)
		}
	})
}
