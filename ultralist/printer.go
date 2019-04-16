package ultralist

// Printer is an interface for printing grouped todos.
type Printer interface {
	Print(*GroupedTodos, int, bool)
}
