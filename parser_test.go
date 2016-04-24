package todolist

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/now"
)

func TestParseSubject(t *testing.T) {
	parser := &Parser{}
	todo := parser.Parse("do this thing")
	if todo.Subject != "do this thing" {
		t.Error("Expected todo.Subject to equal 'do this thing'")
	}
}

func TestParseSubjectWithDue(t *testing.T) {
	parser := &Parser{}
	todo := parser.Parse("do this thing due tomorrow")
	if todo.Subject != "do this thing" {
		t.Error("Expected todo.Subject to equal 'do this thing', got ", todo.Subject)
	}
}

func TestParseProjects(t *testing.T) {
	parser := &Parser{}
	todo := parser.Parse("do this thing +proj1 +proj2 due tomorrow")
	if len(todo.Projects) != 2 {
		t.Error("Expected Projects length to be 2")
	}
	if todo.Projects[0] != "proj1" {
		t.Error("todo.Projects[0] should equal 'proj1' but got", todo.Projects[0])
	}
	if todo.Projects[1] != "proj2" {
		t.Error("todo.Projects[1] should equal 'proj2' but got", todo.Projects[1])
	}
}

func TestParseContexts(t *testing.T) {
	parser := &Parser{}
	todo := parser.Parse("do this thing with @bob and @mary due tomorrow")
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
	todo := parser.Parse("do this thing with @bob and @mary due today")
	if todo.Due != now.BeginningOfDay() {
		fmt.Println("Date is different", todo.Due, time.Now())
	}
}

func TestDueTomorrow(t *testing.T) {
	parser := &Parser{}
	todo := parser.Parse("do this thing with @bob and @mary due tomorrow")
	if todo.Due != now.BeginningOfDay().AddDate(0, 0, 1) {
		fmt.Println("Date is different", todo.Due, time.Now())
	}
}

//func TestDueNextWeek(t *testing.T) {
//	parser := &Parser{}
//
//	fmt.Println("about to parse")
//	todo := parser.Parse("do this thing with @bob and @mary due next week")
//	fmt.Println(todo.Due)
//}

func TestDueMonday(t *testing.T) {
	parser := &Parser{}
	todo := parser.Parse("do this thing with @bob and @mary due mon")
	fmt.Println(todo.Due)
}
