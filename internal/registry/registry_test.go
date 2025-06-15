package registry_test

import (
	"testing"

	"github.com/ei-sugimoto/wakemae/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	t.Parallel()

	rg := registry.NewRegistry()
	rg.AddA("10.0.0.2", "redis")
	rg.AddCNAME("cache", "redis")
	ips, ok := rg.Resolve("cache")
	if !ok {
		t.Fatal("failed to resolve cache")
	}
	assert.Len(t, ips, 1)
	assert.Equal(t, ips[0], "10.0.0.2")
}
