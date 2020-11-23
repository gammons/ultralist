package ultralist

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Manager struct{}

func (m *Manager) RunManager(todoList *TodoList) {

	grouper := &Grouper{}
	groupedTodos := grouper.GroupByProject(todoList.Todos())

	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetRegions(true)

	viewPrinter := &ViewPrinter{}

	count := 0
	for key := range groupedTodos.Groups {
		fmt.Fprintf(textView, "\n[#6a9fb5]%s[#d0d0d0]\n", key)

		for _, todo := range groupedTodos.Groups[key] {
			sidx := strconv.Itoa(count)

			id := viewPrinter.FormatID(todo)
			completed := viewPrinter.FormatCompleted(todo)
			subject := fmt.Sprintf("[\"%s\"]%s[\"\"]", sidx, viewPrinter.FormatSubject(todo))

			fmt.Fprintf(textView, "%s %s %s\n", id, completed, subject)

			count++
		}
	}

	textView.Highlight("0")
	textView.SetBackgroundColor(tcell.NewHexColor(0x151515))

	mainArea := tview.NewGrid()
	mainArea.SetBorder(true).SetTitle(" Your list ")
	mainArea.SetRows(3, -1, 3)

	// need to set this specifically to the height of the textView
	mainArea.AddItem(
		textView,
		1,    // row
		0,    // column
		1,    // rowSpan
		1,    // colSpan
		0,    // minGridHeight
		0,    // minGridWidth
		true) // focus

	// mainArea.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
	// 	AddItem(header, 5, 0, false).
	// 	AddItem(todoListArea, 60, 0, false), 0, 1, false)

	// SetRows(3, 0, 3).
	// SetColumns(30, 0, 30).
	// SetBorders(false).

	windowApp := tview.NewApplication().SetRoot(mainArea, true).EnableMouse(true)

	highlightedTodoID := 0
	windowApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			if highlightedTodoID < len(todoList.Todos())-1 {
				highlightedTodoID += 1
			}
		}
		if event.Rune() == 'k' {
			if highlightedTodoID > 0 {
				highlightedTodoID -= 1
			}
		}
		if event.Key() == tcell.KeyTab {
			fmt.Println("tab")
		}

		textView.Highlight(strconv.Itoa(highlightedTodoID))

		return event
	})

	if err := windowApp.Run(); err != nil {
		panic(err)
	}
}
