package ultralist

type Printer interface {
	Print(*GroupedTodos, bool)
}
