package ultralist

// MemoryPrinter is the main struct of this file.
type MemoryPrinter struct {
	Groups *GroupedTodos
}

// Print is printing grouped todos.
func (m *MemoryPrinter) Print(groupedTodos *GroupedTodos, maxTodoID int, printNotes bool) {
	m.Groups = groupedTodos
}
