package todolist

type Printer interface {
	Print(*GroupedTodos, bool)
}
