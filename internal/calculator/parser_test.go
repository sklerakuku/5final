package calculation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// internal/calculator/parser_test.go
func TestParser(t *testing.T) {
	result, _ := Parse("2 + 2 * 2")
	assert.Equal(t, 6.0, result)
}
