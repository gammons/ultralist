package ultralist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEventLogsWithAddingTodo(t *testing.T) {
	assert := assert.New(t)
	todo := &Todo{Subject: "testing", Completed: false, Archived: false}
	list := &TodoList{}
	list.Add(todo)

	logger := NewEventLogger(list, &MemoryStore{})
	todo2 := NewTodo()
	list.Add(todo2)
	logger.CreateEventLogs()

	assert.Equal(1, len(logger.Events))
	assert.Equal(AddEvent, logger.Events[0].EventType)
	assert.Equal(todo2.ID, logger.Events[0].Object.ID)
}

func TestCreateEventLogsWithAddingMultipleTodos(t *testing.T) {
	assert := assert.New(t)
	todo := &Todo{Subject: "testing", Completed: false, Archived: false}
	list := &TodoList{}
	list.Add(todo)

	logger := NewEventLogger(list, &MemoryStore{})
	todo2 := NewTodo()
	list.Add(todo2)
	todo3 := NewTodo()
	list.Add(todo3)
	logger.CreateEventLogs()

	assert.Equal(2, len(logger.Events))
	assert.Equal(AddEvent, logger.Events[0].EventType)
	assert.Equal(todo2.ID, logger.Events[0].Object.ID)
	assert.Equal(AddEvent, logger.Events[1].EventType)
	assert.Equal(todo3.ID, logger.Events[1].Object.ID)
}

func TestUpdateEvent(t *testing.T) {
	assert := assert.New(t)
	todo := &Todo{Subject: "testing", Completed: false, Archived: false}
	list := &TodoList{}
	list.Add(todo)

	logger := NewEventLogger(list, &MemoryStore{})
	todo.Subject = "testing2"
	logger.CreateEventLogs()

	assert.Equal(1, len(logger.Events))
	assert.Equal(UpdateEvent, logger.Events[0].EventType)
	assert.Equal(todo.ID, logger.Events[0].Object.ID)
}

func TestDeleteEvent(t *testing.T) {
	assert := assert.New(t)
	todo := &Todo{Subject: "testing", Completed: false, Archived: false}
	todo2 := &Todo{Subject: "testing", Completed: false, Archived: false}
	list := &TodoList{}
	list.Add(todo)
	list.Add(todo2)

	logger := NewEventLogger(list, &MemoryStore{})
	list.Delete(todo.ID)
	logger.CreateEventLogs()

	assert.Equal(1, len(logger.Events))
	assert.Equal(DeleteEvent, logger.Events[0].EventType)
	assert.Equal(todo.ID, logger.Events[0].Object.ID)
}
