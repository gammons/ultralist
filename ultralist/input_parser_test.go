package ultralist

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInputParser(t *testing.T) {
	assert := assert.New(t)

	parser := &InputParser{}

	filter, _ := parser.Parse("do this thing due:tom status:now,next")

	assert.Equal(2, len(filter.Status))
	assert.Equal("now", filter.Status[0])
	assert.Equal("next", filter.Status[1])

	tomorrow := time.Now().AddDate(0, 0, 1).Format(DATE_FORMAT)
	assert.Equal(tomorrow, filter.Due)
	assert.Equal("do this thing", filter.Subject)
}

func TestSubject(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("due:tom here is the subject")

	assert.Equal("here is the subject", filter.Subject)
	tomorrow := time.Now().AddDate(0, 0, 1).Format(DATE_FORMAT)
	assert.Equal(tomorrow, filter.Due)
}

func TestDueAgenda(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("due:agenda blah blah")

	tomorrow := time.Now().AddDate(0, 0, 1).Format(DATE_FORMAT)
	assert.Equal(tomorrow, filter.DueBefore)
	assert.Equal(true, filter.HasDueBefore)
	assert.Equal(false, filter.HasDue)
}

func TestProjectsInSubject(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("due:tom here is the +project with @context1 and @context2")

	assert.Equal("here is the +project with @context1 and @context2", filter.Subject)
	assert.Equal("project", filter.Projects[0])
	assert.Equal([]string{"context1", "context2"}, filter.Contexts)
}

func TestProjectsAsFilter(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("subject project:project1,-project2")

	assert.Equal([]string{"project1"}, filter.Projects)
	assert.Equal([]string{"project2"}, filter.ExcludeProjects)

	// assert that project filters override subject projects
	filter, _ = parser.Parse("subject +subjectProject project:project1,-project2")

	assert.Equal([]string{"project1"}, filter.Projects)
	assert.Equal([]string{"project2"}, filter.ExcludeProjects)
}

func TestCompleted(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("lunch with bob completed:true")

	assert.Equal("lunch with bob", filter.Subject)
	assert.Equal(true, filter.Completed)
}

func TestStatus(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, err := parser.Parse("status:one,-two")

	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	assert.Equal(true, filter.HasStatus)
	assert.Equal([]string{"one"}, filter.Status)
	assert.Equal([]string{"two"}, filter.ExcludeStatus)

	filter, err = parser.Parse("this is the subject")

	assert.Equal(false, filter.HasStatus)
}

func TestNoSubject(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("priority:true")

	assert.Equal("", filter.Subject)
	assert.Equal(true, filter.IsPriority)
}

func TestInvalidRecurrence(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	_, err := parser.Parse("recur:blah")

	if err == nil {
		fmt.Println("expected an error")
		t.Fail()
	}

	assert.Equal(err.Error(), "I could not understand the recurrence you gave me: 'blah'")
}

func TestNoneRecurrence(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("recur:none")

	assert.Equal(true, filter.HasRecur)
	assert.Equal("", filter.Recur)
}
