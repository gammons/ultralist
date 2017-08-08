package todolist

import "sort"

type TodoList struct {
	Data []*Todo
}

func (t *TodoList) Load(todos []*Todo) {
	t.Data = todos
}

func (t *TodoList) Add(todo *Todo) {
	todo.Id = t.NextId()
	t.Data = append(t.Data, todo)
}

func (t *TodoList) Delete(ids ...int) {
	for _, id := range ids {
		todo := t.FindById(id)
		if todo == nil {
			continue
		}
		i := -1
		for index, todo := range t.Data {
			if todo.Id == id {
				i = index
			}
		}

		t.Data = append(t.Data[:i], t.Data[i+1:]...)
	}
}

func (t *TodoList) Complete(ids ...int) {
	for _, id := range ids {
		todo := t.FindById(id)
		if todo == nil {
			continue
		}
		todo.Complete()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

func (t *TodoList) Uncomplete(ids ...int) {
	for _, id := range ids {
		todo := t.FindById(id)
		if todo == nil {
			continue
		}
		todo.Uncomplete()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

func (t *TodoList) Archive(ids ...int) {
	for _, id := range ids {
		todo := t.FindById(id)
		if todo == nil {
			continue
		}
		todo.Archive()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

func (t *TodoList) Unarchive(ids ...int) {
	for _, id := range ids {
		todo := t.FindById(id)
		if todo == nil {
			continue
		}
		todo.Unarchive()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

func (t *TodoList) Prioritize(ids ...int) {
	for _, id := range ids {
		todo := t.FindById(id)
		if todo == nil {
			continue
		}
		todo.Prioritize()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

func (t *TodoList) Unprioritize(ids ...int) {
	for _, id := range ids {
		todo := t.FindById(id)
		if todo == nil {
			continue
		}
		todo.Unprioritize()
		t.Delete(id)
		t.Data = append(t.Data, todo)
	}
}

func (t *TodoList) IndexOf(todoToFind *Todo) int {
	for i, todo := range t.Data {
		if todo.Id == todoToFind.Id {
			return i
		}
	}
	return -1
}

type ByDate []*Todo

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	t1Due := a[i].CalculateDueTime()
	t2Due := a[j].CalculateDueTime()
	return t1Due.Before(t2Due)
}

func (t *TodoList) Todos() []*Todo {
	sort.Sort(ByDate(t.Data))
	return t.Data
}

func (t *TodoList) MaxId() int {
	maxId := 0
	for _, todo := range t.Data {
		if todo.Id > maxId {
			maxId = todo.Id
		}
	}
	return maxId
}

func (t *TodoList) NextId() int {
	var found bool
	maxID := t.MaxId()
	for i := 1; i <= maxID; i++ {
		found = false
		for _, todo := range t.Data {
			if todo.Id == i {
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

func (t *TodoList) FindById(id int) *Todo {
	for _, todo := range t.Data {
		if todo.Id == id {
			return todo
		}
	}
	return nil
}

func (t *TodoList) GarbageCollect() {
	var toDelete []*Todo
	for _, todo := range t.Data {
		if todo.Archived {
			toDelete = append(toDelete, todo)
		}
	}
	for _, todo := range toDelete {
		t.Delete(todo.Id)
	}
}
