package tui

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/ultralist/ultralist/ultralist"
)

type TodoTextView struct {
	TextView        *tview.TextView
	SelectedTodoIdx int

	TodoList *ultralist.TodoList

	// the current todo IDs that are listed in the TextView.
	TodoIDs []int

	// a flag to determine whether the selected todo should be highlighted within the text view.
	HighlightSelectedTodo bool
}

//func NewTodoTextView(todoList *ultralist.TodoList, filter *ultralist.Filter) *TodoTextView {
func NewTodoTextView(todoList *ultralist.TodoList) *TodoTextView {
	textView := tview.NewTextView()
	textView.SetWrap(false)
	textView.SetBackgroundColor(tcell.NewHexColor(0x101010))
	textView.SetDynamicColors(true)
	textView.SetBorder(false)
	textView.SetRegions(true)

	return &TodoTextView{
		TextView:              textView,
		SelectedTodoIdx:       0,
		HighlightSelectedTodo: true,
		TodoList:              todoList,
		// Filter:   filter,
		// TodoList: todoList,
	}
}

func (t *TodoTextView) DrawTodos(filter *ultralist.Filter, grouping ultralist.Grouping) {
	var todoIDs []int
	var keys []string
	viewPrinter := &ViewPrinter{}

	todoFilter := &ultralist.TodoFilter{
		Filter: filter,
		Todos:  t.TodoList.Todos(),
	}

	grouper := &ultralist.Grouper{}
	var groups *ultralist.GroupedTodos

	switch grouping {
	case ultralist.GroupByNone:
		groups = grouper.GroupByNothing(todoFilter.ApplyFilter())
	case ultralist.GroupByProject:
		groups = grouper.GroupByProject(todoFilter.ApplyFilter())
	case ultralist.GroupByContext:
		groups = grouper.GroupByContext(todoFilter.ApplyFilter())
	case ultralist.GroupByStatus:
		groups = grouper.GroupByStatus(todoFilter.ApplyFilter())
	}

	for key := range groups.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	t.TextView.Clear()
	// not sure if I need this
	// t.TextView.Highlight("-1")

	totalDisplayedTodos := 0
	for _, key := range keys {
		if len(groups.Groups[key]) > 0 {
			fmt.Fprintf(t.TextView, "\n[%s]%s[%s]\n", ColorBlue, key, ColorForeground)
		}

		for _, todo := range groups.Groups[key] {
			fmt.Fprintf(
				t.TextView,
				"[\"%v\"]%s  %s  %s  %s  %s[\"\"]\n",
				totalDisplayedTodos,
				viewPrinter.FormatID(todo),
				viewPrinter.FormatCompleted(todo),
				viewPrinter.FormatDue(todo),
				viewPrinter.FormatStatus(todo),
				viewPrinter.FormatSubject(todo),
			)

			todoIDs = append(todoIDs, todo.ID)

			if totalDisplayedTodos == t.SelectedTodoIdx && t.HighlightSelectedTodo {
				// TODO: implement this in Manager
				// m.buildTodoCommandsMenu(todo)
				t.TextView.Highlight(strconv.Itoa(t.SelectedTodoIdx))
			}

			totalDisplayedTodos++
		}
	}
	t.TextView.ScrollTo(t.SelectedTodoIdx, 0)
	t.TodoIDs = todoIDs
}

func (t *TodoTextView) SelectedTodo() *ultralist.Todo {
	if len(t.TodoIDs) == 0 {
		return nil
	}

	return t.TodoList.FindByID(t.TodoIDs[t.SelectedTodoIdx])
}

func (t *TodoTextView) HighlightNextTodo() {
	if t.SelectedTodoIdx < len(t.TodoIDs)-1 {
		t.SelectedTodoIdx++
	}
}

func (t *TodoTextView) HighlightPrevTodo() {
	if t.SelectedTodoIdx > 0 {
		t.SelectedTodoIdx--
	}
}

func (t *TodoTextView) ResetSelectedTodoIdx() {
	t.SelectedTodoIdx = 0
}

// func (m *Manager) handleTodoScrollLocation(selectedTodoLineLocation int) {
// 	_, _, _, height := m.TodoTextView.GetRect()
//
// 	// handle the top of the list
// 	// handle when scrolling in the middle
//
// 	// hight is 5
// 	// selectedTodoLineLocation is 5
// 	// todo list length is 10
// 	// keep highlighted todo in middle
//
// }
//
