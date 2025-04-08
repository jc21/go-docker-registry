//go:build unit
// +build unit

package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoSomething(t *testing.T) {
	result := DoSomething()
	assert.Equal(t, "abc123", result)
}
