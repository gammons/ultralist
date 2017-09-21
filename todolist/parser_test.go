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

func TestParseAddNote(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("add write the test functions")

	b1 := parser.ParseAddNote(todo, "an 1 TestPasrseAddNote")
	b2 := parser.ParseAddNote(todo, "an 1 TestPasrseDeleteNote")
	b3 := parser.ParseAddNote(todo, "an 1 TestPasrseEditNote")

	if !b1 || !b2 || !b3 {
		t.Error("Fail adding notes, expected 3 notes but", len(todo.Notes))
	}
}

func TestParseDeleteNote(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("add buy notebook")

	todo.Notes = append(todo.Notes, "ASUStek")
	todo.Notes = append(todo.Notes, "Apple")
	todo.Notes = append(todo.Notes, "Dell")
	todo.Notes = append(todo.Notes, "Acer")

	b1 := parser.ParseDeleteNote(todo, "dn 1 1")
	b2 := parser.ParseDeleteNote(todo, "dn 1 1")

	if !b1 || !b2 {
		t.Error("Fail deleting notes, expected 2 notes left but", len(todo.Notes))
	}

	if todo.Notes[0] != "ASUStek" || todo.Notes[1] != "Acer" {
		t.Error("Fail deleting notes,", todo.Notes[0], "and", todo.Notes[1], "are left")
	}
}

func TestParseEditNote(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("add record the weather")

	todo.Notes = append(todo.Notes, "Aug 29 Wed")
	todo.Notes = append(todo.Notes, "Cloudy")
	todo.Notes = append(todo.Notes, "40°C")
	todo.Notes = append(todo.Notes, "Tokyo")

	parser.ParseEditNote(todo, "en 1 0 Aug 29 Tue")
	if todo.Notes[0] != "Aug 29 Tue" {
		t.Error("Fail editing notes, note 0 should be \"Aug 29 Tue\" but got", todo.Notes[0])
	}

	parser.ParseEditNote(todo, "en 1 1 Sunny")
	if todo.Notes[1] != "Sunny" {
		t.Error("Fail editing notes, note 1 should be \"Sunny\" but got", todo.Notes[1])
	}

	parser.ParseEditNote(todo, "en 1 2 22°C")
	if todo.Notes[2] != "22°C" {
		t.Error("Fail editing notes, note 2 should be \"22°C\" but got", todo.Notes[2])
	}

	parser.ParseEditNote(todo, "en 1 3 Seoul")
	if todo.Notes[3] != "Seoul" {
		t.Error("Fail editing notes, note 3 should be \"Seoul\" but got", todo.Notes[3])
	}
}

func TestHandleNotes(t *testing.T) {
	parser := &Parser{}
	todo := parser.ParseNewTodo("add search engine survey")

	if !parser.ParseAddNote(todo, "an 1 www.google.com") {
		t.Error("Expected Notes to be added")
	}
	if todo.Notes[0] != "www.google.com" {
		t.Error("Expected note 1 to be 'www.google.com' but got", todo.Notes[0])
	}

	if !parser.ParseEditNote(todo, "en 1 0 www.duckduckgo.com") {
		t.Error("Expected Notes to be editted")
	}
	if todo.Notes[0] != "www.duckduckgo.com" {
		t.Error("Expected note 1 to be 'www.duckduckgo.com' but got", todo.Notes[0])
	}

	if !parser.ParseDeleteNote(todo, "dn 1 0") {
		t.Error("Expected Notes to be deleted")
	}
	if len(todo.Notes) != 0 {
		t.Error("Expected no note")
	}
}

func TestDueToday(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	expectedDate := bod(time.Now()).Format("2006-01-02")

	todo := parser.ParseNewTodo("do this thing with @bob and @mary due today")
	assert.Equal(expectedDate, todo.Due)

	todo = parser.ParseNewTodo("do this thing with @bob and @mary due tod")
	assert.Equal(expectedDate, todo.Due)
}

func TestDueTomorrow(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	expectedDate := bod(time.Now()).AddDate(0, 0, 1).Format("2006-01-02")

	todo := parser.ParseNewTodo("do this thing with @bob and @mary due tomorrow")
	assert.Equal(expectedDate, todo.Due)

	todo = parser.ParseNewTodo("do this thing with @bob and @mary due tom")
	assert.Equal(expectedDate, todo.Due)
}

func TestDueSpecific(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	todo := parser.ParseNewTodo("do this thing with @bob and @mary due jun 1")
	year := strconv.Itoa(time.Now().Year())
	assert.Equal(fmt.Sprintf("%s-06-01", year), todo.Due)
}

func TestDueSpecificEuropeanDate(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	todo := parser.ParseNewTodo("do this thing with @bob and @mary due 1 jun")
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

func TestDueOnSpecificDateEuropeFormat(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	year := strconv.Itoa(time.Now().Year())
	assert.Equal(fmt.Sprintf("%s-05-02", year), parser.Due("due 2 may", time.Now()))
	assert.Equal(fmt.Sprintf("%s-06-01", year), parser.Due("due 1 jun", time.Now()))
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

func TestParseEditTodoJustDate(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	todo := NewTodo()
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	parser.ParseEditTodo(todo, "e 24 due tom")

	assert.Equal(todo.Due, tomorrow)
}

func TestParseEditTodoJustDateDoesNotEditExistingSubject(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	todo := NewTodo()
	todo.Subject = "pick up the trash"
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	parser.ParseEditTodo(todo, "e 24 due tom")

	assert.Equal(todo.Due, tomorrow)
	assert.Equal(todo.Subject, "pick up the trash")
}

func TestParseEditTodoJustSubject(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	todo := &Todo{Subject: "pick up the trash", Due: "2016-11-25"}

	parser.ParseEditTodo(todo, "e 24 changed the todo")

	assert.Equal(todo.Due, "2016-11-25")
	assert.Equal(todo.Subject, "changed the todo")
}

func TestParseEditTodoSubjectUpdatesProjectsAndContexts(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	todo := &Todo{
		Subject:  "pick up the +trash with @dad",
		Due:      "2016-11-25",
		Projects: []string{"trash"},
		Contexts: []string{"dad"},
	}

	parser.ParseEditTodo(todo, "e 24 get the +garbage with @mom")

	assert.Equal(todo.Due, "2016-11-25")
	assert.Equal(todo.Subject, "get the +garbage with @mom")
	assert.Equal(todo.Projects, []string{"garbage"})
	assert.Equal(todo.Contexts, []string{"mom"})
}

func TestParseEditTodoWithSubjectAndDue(t *testing.T) {
	assert := assert.New(t)
	parser := &Parser{}
	todo := &Todo{
		Subject:  "pick up the +trash with @dad",
		Due:      "2016-11-25",
		Projects: []string{"trash"},
		Contexts: []string{"dad"},
	}
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	parser.ParseEditTodo(todo, "e 24 get the +garbage with @mom due tom")

	assert.Equal(todo.Due, tomorrow)
	assert.Equal(todo.Subject, "get the +garbage with @mom")
}
