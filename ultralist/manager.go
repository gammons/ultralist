package ultralist

import (
	"fmt"
	"strconv"

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

	table := tview.NewTable().SetBorders(false)

	for todoRow, todo := range todoList.Todos() {
		table.SetCell(todoRow, 0, tview.NewTableCell(strconv.Itoa(todo.ID)))
		table.SetCell(todoRow, 2, tview.NewTableCell("[ ]"))
		table.SetCell(todoRow, 3, tview.NewTableCell(todo.Due))
		table.SetCell(todoRow, 4, tview.NewTableCell(todo.Subject))
	}

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("\n\nUltralist")

	tableHolder := tview.NewFlex().AddItem(table, 0, 1, false)

	// tableGrid := tview.NewGrid().
	// 	SetRows(1,0,1).

	mainGrid := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 0, 1, false).
			AddItem(tableHolder, 0, 5, false), 0, 1, false)

		// SetRows(3, 0, 3).
		// SetColumns(30, 0, 30).
		// SetBorders(false).

	windowApp := tview.NewApplication().SetRoot(mainGrid, true).EnableMouse(true)

	if err := windowApp.Run(); err != nil {
		panic(err)
	}

	fmt.Println("manager")
}
