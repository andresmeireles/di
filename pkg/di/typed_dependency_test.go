package di_test

import (
	"testing"

	"github.com/andresmeireles/di/pkg/di"
	"github.com/andresmeireles/di/testdata"
)

func TestTypedDependency(t *testing.T) {
	t.Run("should create a new typed dependency", func(t *testing.T) {
		dep := di.NewTypedDependency[testdata.TestOne](testdata.NewTestOne)

		if dep.Name() != "testdata.TestOne" {

		}
	})
}
