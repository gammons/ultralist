package todolist

import "testing"

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
