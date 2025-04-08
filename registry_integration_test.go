//go:build integration
// +build integration

package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoSomethingModule(t *testing.T) {
	result := DoSomething()
	assert.Equal(t, "abc123", result)
}
