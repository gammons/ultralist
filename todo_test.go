package main

import (
	"os"
	"testing"

	"github.com/gammons/todolist/todolist"
)

const testFile = ".todos.json"

func newTestApp() *todolist.App {
	os.Remove(testFile)

	store := todolist.FileStore{FileLocation: testFile}
	store.Initialize()

	return &todolist.App{
		TodoList:  &todolist.TodoList{},
		TodoStore: &store,
	}
}

func TestAddTodo(t *testing.T) {
	app := newTestApp()
	err := routeInput(app, "a", "a foobar")
	if err != nil {
		t.Fatal(err)
	}

	for _, todo := range app.TodoList.Todos() {
		if todo.Subject == "foobar" {
			return
		}
	}

	t.Fatal("Todo not found")
}

func TestEditSubject(t *testing.T) {
	app := newTestApp()

	err := routeInput(app, "a", "a foobar")
	if err != nil {
		t.Fatal(err)
	}

	err = routeInput(app, "es", "es 1 foobilicious")
	if err != nil {
		t.Fatal(err)
	}

	for _, todo := range app.TodoList.Todos() {
		if todo.Subject == "foobilicious" {
			return
		}
	}

	t.Fatal("Todo not found")
}
