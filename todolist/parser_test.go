package todolist

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseSubject(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("do this thing")
	if todo.Subject != "do this thing" {
		t.Error("Expected todo.Subject to equal 'do this thing'")
	}
}

func TestParseSubjectWithDue(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("do this thing due tomorrow")
	if todo.Subject != "do this thing" {
		t.Error("Expected todo.Subject to equal 'do this thing', got ", todo.Subject)
	}
}

func TestParseExpandProjects(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	correctFormat := parser.ExpandProject("ex 113 +meeting: figures, slides, coffee, suger")
	assert.Equal("+meeting", correctFormat)
	wrongFormat1 := parser.ExpandProject("ex 114 +meeting figures, slides, coffee, suger")
	assert.Equal("", wrongFormat1)
	wrongFormat2 := parser.ExpandProject("ex 115 meeting: figures, slides, coffee, suger")
	assert.Equal("", wrongFormat2)
	wrongFormat3 := parser.ExpandProject("ex 116 meeting figures, slides, coffee, suger")
	assert.Equal("", wrongFormat3)
	wrongFormat4 := parser.ExpandProject("ex 117 +重要な會議: 図, コーヒー, 砂糖")
	assert.Equal("+重要な會議", wrongFormat4)
}

func TestParseProjects(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("do this thing +proj1 +proj2 +專案3 +proj-name due tomorrow")
	if len(todo.Projects) != 4 {
		t.Error("Expected Projects length to be 3")
	}
	if todo.Projects[0] != "proj1" {
		t.Error("todo.Projects[0] should equal 'proj1' but got", todo.Projects[0])
	}
	if todo.Projects[1] != "proj2" {
		t.Error("todo.Projects[1] should equal 'proj2' but got", todo.Projects[1])
	}
	if todo.Projects[2] != "專案3" {
		t.Error("todo.Projects[2] should equal '專案3' but got", todo.Projects[2])
	}
	if todo.Projects[3] != "proj-name" {
		t.Error("todo.Projects[3] should equal 'proj-name' but got", todo.Projects[3])
	}
}

func TestParseContexts(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("do this thing with @bob and @mary due tomorrow")
	if len(todo.Contexts) != 2 {
		t.Error("Expected Projects length to be 2")
	}
	if todo.Contexts[0] != "bob" {
		t.Error("todo.Contexts[0] should equal 'mary' but got", todo.Contexts[0])
	}
	if todo.Contexts[1] != "mary" {
		t.Error("todo.Contexts[1] should equal 'mary' but got", todo.Contexts[1])
	}
}

func TestDueToday(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("do this thing with @bob and @mary due today")
	if todo.Due != bod(time.Now()).Format("2006-01-02") {
		fmt.Println("Date is different", todo.Due, time.Now())
	}
	todo = parser.ParseNewTodo("do this thing with @bob and @mary due tod")
	if todo.Due != bod(time.Now()).Format("2006-01-02") {
		fmt.Println("Date is different", todo.Due, time.Now())
	}
}

func TestDueTomorrow(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("do this thing with @bob and @mary due tomorrow")
	if todo.Due != bod(time.Now()).AddDate(0, 0, 1).Format("2006-01-02") {
		fmt.Println("Date is different", todo.Due, time.Now())
	}
	todo = parser.ParseNewTodo("do this thing with @bob and @mary due tom")
	if todo.Due != bod(time.Now()).AddDate(0, 0, 1).Format("2006-01-02") {
		fmt.Println("Date is different", todo.Due, time.Now())
	}
}

func TestDueSpecific(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	todo := parser.ParseNewTodo("do this thing with @bob and @mary due jun 1")
	year := strconv.Itoa(time.Now().Year())
	assert.Equal(fmt.Sprintf("%s-06-01", year), todo.Due)
}

func TestMondayOnSunday(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	now, _ := time.Parse("2006-01-02", "2016-04-24")
	assert.Equal("2016-04-25", parser.monday(now))
}

func TestMondayOnMonday(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	now, _ := time.Parse("2006-01-02", "2016-04-25")
	assert.Equal("2016-04-25", parser.monday(now))
}

func TestMondayOnTuesday(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	now, _ := time.Parse("2006-01-02", "2016-04-26")
	assert.Equal("2016-05-02", parser.monday(now))
}

func TestTuesdayOnMonday(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	now, _ := time.Parse("2006-01-02", "2016-04-25")
	assert.Equal("2016-04-26", parser.tuesday(now))
}

func TestTuesdayOnWednesday(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	now, _ := time.Parse("2006-01-02", "2016-04-27")
	assert.Equal("2016-05-03", parser.tuesday(now))
}

func TestDueOnSpecificDate(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	year := strconv.Itoa(time.Now().Year())
	assert.Equal(fmt.Sprintf("%s-05-02", year), parser.Due("due may 2", time.Now()))
	assert.Equal(fmt.Sprintf("%s-06-01", year), parser.Due("due jun 1", time.Now()))
}

func TestDueOnSpecificDateEuropean(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	year := strconv.Itoa(time.Now().Year())
	assert.Equal(fmt.Sprintf("%s-05-02", year), parser.Due("due 2 may", time.Now()))
}

func TestDueIntelligentlyChoosesCorrectYear(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	marchTime, _ := time.Parse("2006-01-02", "2016-03-25")
	januaryTime, _ := time.Parse("2006-01-02", "2016-01-05")
	septemberTime, _ := time.Parse("2006-01-02", "2016-09-25")
	decemberTime, _ := time.Parse("2006-01-02", "2016-12-25")

	assert.Equal("2016-01-10", parser.parseArbitraryDate("jan 10", januaryTime))
	assert.Equal("2016-01-10", parser.parseArbitraryDate("jan 10", marchTime))
	assert.Equal("2017-01-10", parser.parseArbitraryDate("jan 10", septemberTime))
	assert.Equal("2017-01-10", parser.parseArbitraryDate("jan 10", decemberTime))
}

func TestParseCommandIdSubject(t *testing.T) {
	assert := assert.New(t)
	parser := Parser{"es 24 a new subject"}
	command, id, subject := parser.Parse()

	assert.Equal("es", command)
	assert.Equal(24, id)
	assert.Equal("a new subject", subject)
}

func TestParseInvalidCommandIdSubject(t *testing.T) {
	assert := assert.New(t)
	input := "es a new project"
	parser := Parser{input}
	command, id, subject := parser.Parse()

	assert.Equal("", command)
	assert.Equal(-1, id)
	assert.Equal(input, subject)
}
