package ultralist

import (
	"sort"
	"time"
)

// TodoList is the struct of a list with several todos.
type TodoList struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	IsSynced bool
	Data     []*Todo `json:"todo_items_attributes"`
}

// Load is loading a list with several todos.
func (t *TodoList) Load(todos []*Todo) {
	t.Data = todos
}

// Add is adding a single todo to a todo list.
func (t *TodoList) Add(todo *Todo) {
	todo.ID = t.NextID()
	t.Data = append(t.Data, todo)
}

// Delete is deleting multiple todos from a todo list.
func (t *TodoList) Delete(ids ...int) {
	for _, id := range ids {
		todo := t.FindByID(id)
		if todo == nil {
			continue
		}
		i := -1
		for index, todo := range t.Data {
			if todo.ID == id {
				i = index
			}
		}

		t.Data = append(t.Data[:i], t.Data[i+1:]...)
	}
}

// Complete is completing multiple todos in a todo list.
func (t *TodoList) Complete(ids ...int) {
	for _, id := range ids {
		todo := t.FindByID(id)
		if todo == nil {
			continue
		}
		prevStatus := todo.Status
		todo.Complete()
		t.Delete(id)
		t.Data = append(t.Data, todo)

		r := &Recurrence{}
		if r.HasNextRecurringTodo(todo) {
			next := r.NextRecurringTodo(todo, time.Now())
			next.Status = prevStatus
			t.Add(next)
		}
	}
}

// Uncomplete is uncompleting multiple todos from a todo list.
func (t *TodoList) Uncomplete(ids ...int) {
	for _, id := range ids {
		todo := t.FindByID(id)
		if todo == nil {
			continue
		}
		todo.Uncomplete()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

// Archive is archiving multiple todos from a todo list.
func (t *TodoList) Archive(ids ...int) {
	for _, id := range ids {
		todo := t.FindByID(id)
		if todo == nil {
			continue
		}
		todo.Archive()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

// Unarchive is unarchiving multiple todos from a todo list.
func (t *TodoList) Unarchive(ids ...int) {
	for _, id := range ids {
		todo := t.FindByID(id)
		if todo == nil {
			continue
		}
		todo.Unarchive()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

// Prioritize is prioritizing multiple todos from a todo list.
func (t *TodoList) Prioritize(ids ...int) {
	for _, id := range ids {
		todo := t.FindByID(id)
		if todo == nil {
			continue
		}
		todo.Prioritize()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

// Unprioritize is unprioritizing multiple todos from a todo list.
func (t *TodoList) Unprioritize(ids ...int) {
	for _, id := range ids {
		todo := t.FindByID(id)
		if todo == nil {
			continue
		}
		todo.Unprioritize()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

// SetStatus sets the status of a todo
func (t *TodoList) SetStatus(input string, ids ...int) {
	for _, id := range ids {
		todo := t.FindByID(id)
		if todo == nil {
			continue
		}
		todo.Status = input
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

// IndexOf finds the index of a todo.
func (t *TodoList) IndexOf(todoToFind *Todo) int {
	for i, todo := range t.Data {
		if todo.ID == todoToFind.ID {
			return i
		}
	}
	return -1
}

// ByDate is the by date struct of a todo.
type ByDate []*Todo

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	t1Due := a[i].CalculateDueTime()
	t2Due := a[j].CalculateDueTime()
	return t1Due.Before(t2Due)
}

// Todos is a sorted list of todos.
func (t *TodoList) Todos() []*Todo {
	sort.Sort(ByDate(t.Data))
	return t.Data
}

// MaxID returns the maximum human readable ID of all todo items.
func (t *TodoList) MaxID() int {
	maxID := 0
	for _, todo := range t.Data {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}
	return maxID
}

// NextID returns the next human readable ID.
func (t *TodoList) NextID() int {
	var found bool
	maxID := t.MaxID()
	for i := 1; i <= maxID; i++ {
		found = false
		for _, todo := range t.Data {
			if todo.ID == i {
				found = true
				break
			}
		}
		if !found {
			return i
		}
	}
	return maxID + 1
}

// FindByID finds a todo by ID.
func (t *TodoList) FindByID(id int) *Todo {
	for _, todo := range t.Data {
		if todo.ID == id {
			return todo
		}
	}
	return nil
}

// GarbageCollect deletes todos which are archived.
func (t *TodoList) GarbageCollect() {
	var toDelete []*Todo
	for _, todo := range t.Data {
		if todo.Archived {
			toDelete = append(toDelete, todo)
		}
	}
	for _, todo := range toDelete {
		t.Delete(todo.ID)
	}
}
