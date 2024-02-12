package di_test

import (
	"testing"

	"github.com/andresmeireles/di/pkg/di"
)

func TestNamedDependency(t *testing.T) {
	t.Run("should create a new named dependency", func(t *testing.T) {
		dep := di.NewNamedDependency("dep", func() int { return 1 })

		if dep.Name() != "dep" {
			t.Errorf("Name of dependency must be %s, received %s", "dep", dep.Name())
		}
	})
}
