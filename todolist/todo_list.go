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

func (t *TodoList) Delete(id int) {
	i := -1
	for index, todo := range t.Data {
		if todo.Id == id {
			i = index
		}
	}

	t.Data = append(t.Data[:i], t.Data[i+1:]...)
}

func (t *TodoList) Complete(id int) {
	todo := t.FindById(id)
	todo.Completed = true
	t.Delete(id)
	t.Data = append(t.Data, todo)
}

func (t *TodoList) Uncomplete(id int) {
	todo := t.FindById(id)
	todo.Completed = false
	t.Delete(id)
	t.Data = append(t.Data, todo)
}

func (t *TodoList) Archive(id int) {
	todo := t.FindById(id)
	todo.Archived = true
	t.Delete(id)
	t.Data = append(t.Data, todo)
}

func (t *TodoList) Unarchive(id int) {
	todo := t.FindById(id)
	todo.Archived = false
	t.Delete(id)
	t.Data = append(t.Data, todo)
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

func (t *TodoList) NextId() int {
	maxId := 0
	for _, todo := range t.Data {
		if todo.Id > maxId {
			maxId = todo.Id
		}
	}
	return maxId + 1
}

func (t *TodoList) FindById(id int) *Todo {
	for _, todo := range t.Data {
		if todo.Id == id {
			return todo
		}
	}
	return nil
}
