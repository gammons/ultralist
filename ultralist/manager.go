package ultralist

import (
	"fmt"

	"github.com/rivo/tview"
)

type Manager struct{}

func (m *Manager) RunManager(todoList *TodoList) {

	// var buf bytes.Buffer
	// printer := &ScreenPrinter{Writer: &buf}
	// grouper := &Grouper{}
	// printer.Print(grouper.GroupByNothing(todoList.Todos()), false)
	//
	// textView := tview.NewTextView().SetText(buf.String())

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("\n\nUltralist")

	todoListHolder := tview.NewFlex().SetDirection(tview.FlexRow)
	todoListHolder.SetBorder(false)

	printer := NewTviewPrinter(true)

	for _, todo := range todoList.Todos() {
		todoHolder := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewTextView().SetDynamicColors(true).SetText(printer.FormatID(todo.ID, todo.IsPriority)), 0, 1, false).
			AddItem(tview.NewTextView().SetText(todo.Due), 0, 1, false).
			AddItem(tview.NewTextView().SetDynamicColors(true).SetText(printer.FormatSubject(todo.Subject, todo.IsPriority)), 0, 10, false)

		todoHolder.SetBorder(true)
		todoListHolder.AddItem(todoHolder, 0, 1, false)
		// 	table.SetCell(todoRow, 0, tview.NewTableCell(strconv.Itoa(todo.ID)))
		// 	table.SetCell(todoRow, 2, tview.NewTableCell("[ ]"))
		// 	table.SetCell(todoRow, 3, tview.NewTableCell(todo.Due))
		// 	table.SetCell(todoRow, 4, tview.NewTableCell(todo.Subject))
	}

	todoListArea := tview.NewFrame(todoListHolder).SetBorders(2, 2, 2, 2, 20, 20) //.AddItem(table, 0, 1, false)

	// tableGrid := tview.NewGrid().
	// 	SetRows(1,0,1).

	mainGrid := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 0, 1, false).
			AddItem(todoListArea, 0, 5, false), 0, 1, false)

		// SetRows(3, 0, 3).
		// SetColumns(30, 0, 30).
		// SetBorders(false).

	windowApp := tview.NewApplication().SetRoot(mainGrid, true).EnableMouse(true)

	if err := windowApp.Run(); err != nil {
		panic(err)
	}

	fmt.Println("manager")
}
