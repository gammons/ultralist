package ultralist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGroupByContext(t *testing.T) {
	assert := assert.New(t)

	list := setUpTestMemoryTodoList()
	grouper := &Grouper{}
	grouped := grouper.GroupByContext(list.Todos())

	assert.Equal(2, len(grouped.Groups["root"]), "")
	assert.Equal(1, len(grouped.Groups["more"]), "")
}

func TestGroupByProject(t *testing.T) {
	assert := assert.New(t)

	list := setUpTestMemoryTodoList()
	grouper := &Grouper{}
	grouped := grouper.GroupByProject(list.Todos())

	assert.Equal(2, len(grouped.Groups["test1"]), "")
}

func TestGroupByContextWithPriorityFirst(t *testing.T) {
	assert := assert.New(t)

	var list []*Todo
	list = append(list, &Todo{Subject: "a - one", IsPriority: false})
	list = append(list, &Todo{Subject: "b - two", IsPriority: true})

	grouper := &Grouper{}
	grouped := grouper.GroupByContext(list)

	assert.Equal("b - two", grouped.Groups["No contexts"][0].Subject)
}

func TestGroupByContextSortedByDueDate(t *testing.T) {
	assert := assert.New(t)

	var list []*Todo
	list = append(list, &Todo{Subject: "a - one", IsPriority: false, Due: time.Now().Format(DateFormat)})
	list = append(list, &Todo{Subject: "b - two", IsPriority: false, Due: time.Now().AddDate(0, 0, -1).Format(DateFormat)})
	list = append(list, &Todo{Subject: "c - three", IsPriority: false, Due: ""})

	grouper := &Grouper{}
	grouped := grouper.GroupByContext(list)

	assert.Equal("b - two", grouped.Groups["No contexts"][0].Subject)
}

func TestGroupByContextSortedByDueDateWithNoDuePriority(t *testing.T) {
	assert := assert.New(t)

	var list []*Todo
	list = append(list, &Todo{Subject: "a - one", IsPriority: false, Due: time.Now().Format(DateFormat)})
	list = append(list, &Todo{Subject: "b - two", IsPriority: false, Due: time.Now().AddDate(0, 0, -1).Format(DateFormat)})
	list = append(list, &Todo{Subject: "c - three", IsPriority: true, Due: ""})

	grouper := &Grouper{}
	grouped := grouper.GroupByContext(list)

	assert.Equal("c - three", grouped.Groups["No contexts"][0].Subject)
}

func TestGroupByContextSortedByDueDateWithPriority(t *testing.T) {
	assert := assert.New(t)

	var list []*Todo
	list = append(list, &Todo{Subject: "a - one", IsPriority: true, Due: time.Now().Format(DateFormat)})
	list = append(list, &Todo{Subject: "b - two", IsPriority: false, Due: time.Now().AddDate(0, 0, -1).Format(DateFormat)})
	list = append(list, &Todo{Subject: "c - three", IsPriority: false, Due: ""})

	grouper := &Grouper{}
	grouped := grouper.GroupByContext(list)

	assert.Equal("a - one", grouped.Groups["No contexts"][0].Subject)
}

func TestGroupByContextSortedByDueDateWithArchived(t *testing.T) {
	assert := assert.New(t)

	var list []*Todo
	list = append(list, &Todo{Subject: "a - one", IsPriority: true, Archived: true, Due: time.Now().Format(DateFormat)})
	list = append(list, &Todo{Subject: "b - two", IsPriority: false, Archived: true, Due: time.Now().AddDate(0, 0, -1).Format(DateFormat)})
	list = append(list, &Todo{Subject: "c - three", IsPriority: false, Due: ""})

	grouper := &Grouper{}
	grouped := grouper.GroupByContext(list)

	assert.Equal("c - three", grouped.Groups["No contexts"][0].Subject)
}

func setUpTestMemoryTodoList() *TodoList {
	list := &TodoList{}

	todo1 := NewTodo()
	todo1.Subject = "this is the first subject"
	todo1.Projects = []string{"test1"}
	todo1.Contexts = []string{"root"}
	todo1.Due = "2016-04-04"
	todo1.Archive()
	list.Add(todo1)

	todo2 := NewTodo()
	todo2.Subject = "audit userify for 2FA"
	todo2.Projects = []string{"test1"}
	todo2.Contexts = []string{"root", "more"}
	todo2.Complete()
	list.Add(todo2)

	return list
}
