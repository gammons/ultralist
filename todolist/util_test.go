package todolist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluralize(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("todo", pluralize(1, "todo", "todos"))
	assert.Equal("todos", pluralize(2, "todo", "todos"))
}
