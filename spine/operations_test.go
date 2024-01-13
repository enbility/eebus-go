package spine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperations(t *testing.T) {
	operations := NewOperations(true, false)
	assert.NotNil(t, operations)

	text := operations.String()
	assert.NotEqual(t, 0, len(text))

	data := operations.Information()
	assert.NotNil(t, data)

	operations2 := NewOperations(true, true)
	assert.NotNil(t, operations2)

	text = operations2.String()
	assert.NotEqual(t, 0, len(text))

	data = operations2.Information()
	assert.NotNil(t, data)

	operations3 := NewOperations(false, false)
	assert.NotNil(t, operations3)

	text = operations3.String()
	assert.NotEqual(t, 0, len(text))

	data = operations3.Information()
	assert.NotNil(t, data)
}
