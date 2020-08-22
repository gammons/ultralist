package ultralist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputParser(t *testing.T) {
	assert := assert.New(t)

	parser := &InputParser{}

	filter, _ := parser.Parse("do this thing due:tom status:now,next")

	assert.Equal(2, len(filter.Status))
	assert.Equal("now", filter.Status[0])
	assert.Equal("next", filter.Status[1])

	assert.Equal(1, len(filter.Due))
	assert.Equal("tom", filter.Due[0])
	assert.Equal("do this thing", filter.Subject)
}

func TestSubject(t *testing.T) {
	assert := assert.New(t)

	parser := &InputParser{}

	filter, _ := parser.Parse("due:tom here is the subject")
	assert.Equal("here is the subject", filter.Subject)
	assert.Equal("tom", filter.Due[0])
}

func TestProjects(t *testing.T) {
	assert := assert.New(t)

	parser := &InputParser{}

	filter, _ := parser.Parse("due:tom here is the +project with @context1 and @context2")
	assert.Equal("here is the +project with @context1 and @context2", filter.Subject)
	assert.Equal("project", filter.Projects[0])
	assert.Equal([]string{"context1", "context2"}, filter.Contexts)
}

func TestCompleted(t *testing.T) {
	assert := assert.New(t)

	parser := &InputParser{}

	filter, _ := parser.Parse("lunch with bob completed:true")
	assert.Equal("lunch with bob", filter.Subject)
	assert.Equal(true, filter.Completed)
}

func TestEmptyString(t *testing.T) {
	parser := &InputParser{}

	_, err := parser.Parse("")

	if err == nil {
		fmt.Println("expected an error, got none")
		t.Fail()
	}
}

func TestNoSubject(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("priority:true")

	assert.Equal("", filter.Subject)
	assert.Equal(true, filter.IsPriority)
}
