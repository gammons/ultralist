package ultralist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextID(t *testing.T) {
	assert := assert.New(t)
	todo := &Todo{Subject: "testing", Completed: false, Archived: false}
	list := &TodoList{}
	assert.Equal(1, list.NextID())
	list.Add(todo)
	assert.Equal(2, list.NextID())
}

func TestNextIDWhenTodoDeleted(t *testing.T) {
	assert := assert.New(t)
	todo1 := &Todo{Subject: "testing 1", Completed: false, Archived: false}
	todo2 := &Todo{Subject: "testing 2", Completed: false, Archived: false}
	todo3 := &Todo{Subject: "testing 3", Completed: false, Archived: false}
	todo4 := &Todo{Subject: "testing 4", Completed: false, Archived: false}
	list := &TodoList{}

	list.Add(todo1)
	list.Add(todo2)
	list.Add(todo3)

	list.Delete(2)
	assert.Equal(4, list.NextID())
	list.Add(todo4)
	assert.Equal(5, list.NextID())
	list.Delete(1)
	assert.Equal(5, list.NextID())
}

func TestMaxID(t *testing.T) {
	assert := assert.New(t)
	todo1 := &Todo{Subject: "testing 1", Completed: false, Archived: false}
	todo2 := &Todo{Subject: "testing 2", Completed: false, Archived: false}
	list := &TodoList{}
	assert.Equal(0, list.MaxID())
	list.Add(todo1)
	assert.Equal(1, list.MaxID())
	list.Add(todo2)
	assert.Equal(2, list.MaxID())
}

func TestIndexOf(t *testing.T) {
	assert := assert.New(t)
	todo := &Todo{Subject: "Grant"}
	list := SetUpTestMemoryTodoList()
	assert.Equal(-1, list.IndexOf(todo))
	assert.Equal(0, list.IndexOf(list.Data[0]))
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	list := SetUpTestMemoryTodoList()
	assert.Equal(2, len(list.Data))
	list.Delete(1)
	assert.Equal(1, len(list.Data))
}

func TestComplete(t *testing.T) {
	assert := assert.New(t)
	list := SetUpTestMemoryTodoList()
	assert.Equal(false, list.FindByID(1).Completed)
	list.Complete(1)
	assert.Equal(true, list.FindByID(1).Completed)
}

func TestArchive(t *testing.T) {
	assert := assert.New(t)
	list := SetUpTestMemoryTodoList()
	assert.Equal(false, list.FindByID(2).Archived)
	list.Archive(2)
	assert.Equal(true, list.FindByID(2).Archived)
}
func TestUnarchive(t *testing.T) {
	assert := assert.New(t)
	list := SetUpTestMemoryTodoList()
	assert.Equal(true, list.FindByID(1).Archived)
	list.Unarchive(1)
	assert.Equal(false, list.FindByID(1).Archived)
}

func TestUncomplete(t *testing.T) {
	assert := assert.New(t)
	list := SetUpTestMemoryTodoList()
	assert.Equal(true, list.FindByID(2).Completed)
	list.Uncomplete(2)
	assert.Equal(false, list.FindByID(2).Completed)
}

func TestGarbageCollect(t *testing.T) {
	assert := assert.New(t)
	list := &TodoList{}
	todo1 := &Todo{Subject: "testing 1", Completed: false, Archived: true}
	todo2 := &Todo{Subject: "testing 2", Completed: false, Archived: false}
	todo3 := &Todo{Subject: "testing 3", Completed: false, Archived: true}
	list.Add(todo1)
	list.Add(todo2)
	list.Add(todo3)

	list.GarbageCollect()

	assert.Equal(len(list.Data), 1)
	assert.Equal(2, list.MaxID())
	assert.Equal(3, list.NextID())
}

func TestPrioritizeNotInTodosJson(t *testing.T) {
	assert := assert.New(t)
	list := SetUpTestMemoryTodoList()
	assert.Equal(false, list.FindByID(2).IsPriority)
}

func TestPrioritizeTodo(t *testing.T) {
	assert := assert.New(t)
	list := &TodoList{}
	todo := &Todo{Archived: false, Completed: false, Subject: "testing", IsPriority: false}
	list.Add(todo)
	list.Prioritize(1)
	assert.Equal(true, list.FindByID(1).IsPriority)
	list.Unprioritize(1)
	assert.Equal(false, list.FindByID(1).IsPriority)
}
