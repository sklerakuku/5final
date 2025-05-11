// calculation/parser_test.go
package calculation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected float64
		wantErr  bool
	}{
		{"Simple addition", "2 + 2", 4, false},
		{"Operator precedence", "2 + 2 * 2", 6, false},
		{"Parentheses", "(2 + 2) * 2", 8, false},
		{"Division", "10 / 2", 5, false},
		{"Division by zero", "10 / 0", 0, true},
		{"Unary minus", "-5 + 10", 5, false},
		{"Invalid syntax", "2 + * 2", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.expr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
