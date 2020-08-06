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

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("\n\nUltralist").SetTitle("Header").SetBorder(true)

	//todoListHolder := tview.NewFlex().SetDirection(tview.FlexRow)
	todoListHolder := tview.NewGrid()
	todoListHolder.SetSize(len(todoList.Todos()), 1, 4, 62)
	// todoListHolder.SetMinSize(1, 60)
	// todoListHolder.SetRows(4, 4)
	//todoListHolder.SetColumns(62)
	todoListHolder.SetBorder(true)
	todoListHolder.SetTitle("TodoListHolder")

	//printer := NewTviewPrinter(true)

	for idx, todo := range todoList.Todos() {
		// todoHolder := tview.NewBox().SetBorder(true).
		// 	//AddItem(tview.NewTextView().SetDynamicColors(true).SetText(printer.FormatSubject(todo.Subject, todo.IsPriority)), 0, 10, false)
		// 	SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		// 		tview.Print(screen, todo.Subject, x+2, y+1, width-2, tview.AlignLeft, tcell.ColorWhite)
		// 		return x, y, width, height
		// 	})

		todoHolder := tview.NewGrid()
		todoHolder.SetTitle("TodoHolder").SetBorder(true)
		todoHolder.SetColumns(5, 55)
		todoHolder.SetRows(1, 1)
		//todoHolder.SetSize(2, 2, 3, 60).SetBorder(true)
		//subjectView := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText(printer.FormatSubject(todo.Subject, todo.IsPriority))
		subjectView := tview.NewTextView().SetDynamicColors(true)
		subjectView.SetText(todo.Subject)
		//subjectView.SetBorder(true)
		todoHolder.AddItem(subjectView, 0, 1, 1, 1, 0, 0, true)

		dueView := tview.NewTextView().SetDynamicColors(true)
		dueView.SetText("[ ]")
		todoHolder.AddItem(dueView, 0, 0, 1, 1, 0, 0, true)

		// todoHolder := tview.NewFlex()
		// todoHolder.SetTitle("TodoHolder").SetBorder(true)
		//
		// dueView := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
		// dueView.SetText("[ ]")
		// todoHolder.AddItem(dueView, 5, 0, true)
		//
		// subjectView := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
		// subjectView.SetText(todo.Subject)
		// todoHolder.AddItem(subjectView, 55, 0, true)

		todoListHolder.AddItem(todoHolder, idx, 0, 1, 1, 0, 0, true)
	}

	//AddItem(tview.NewTextView().SetDynamicColors(true).SetText(printer.FormatID(todo.ID, todo.IsPriority)), 0, 1, false).
	//AddItem(tview.NewTextView().SetText(todo.Due), 0, 1, false).

	//todoHolder.SetBorder(true)
	//todoListHolder.AddItem(todoHolder, 0, 1, false)
	// 	table.SetCell(todoRow, 0, tview.NewTableCell(strconv.Itoa(todo.ID)))
	// 	table.SetCell(todoRow, 2, tview.NewTableCell("[ ]"))
	// 	table.SetCell(todoRow, 3, tview.NewTableCell(todo.Due))
	// 	table.SetCell(todoRow, 4, tview.NewTableCell(todo.Subject))
	//}

	todoListArea := tview.NewFrame(todoListHolder).SetBorders(2, 2, 2, 2, 20, 20) //.AddItem(table, 0, 1, false)
	todoListArea.SetTitle("TodoListArea")

	// tableGrid := tview.NewGrid().
	// 	SetRows(1,0,1).

	mainArea := tview.NewFlex()
	mainArea.SetBorder(true).SetTitle("MainArea")
	mainArea.SetDirection(tview.FlexRow)

	gridArea := tview.NewGrid().SetRows(1, 0, 1)
	gridArea.SetBorders(true).SetTitle("GridArea")

	mainArea.AddItem(header, 5, 0, false)
	mainArea.AddItem(todoListHolder, 60, 0, false)

	// mainArea.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
	// 	AddItem(header, 5, 0, false).
	// 	AddItem(todoListArea, 60, 0, false), 0, 1, false)

	// SetRows(3, 0, 3).
	// SetColumns(30, 0, 30).
	// SetBorders(false).

	windowApp := tview.NewApplication().SetRoot(mainArea, true).EnableMouse(true)

	if err := windowApp.Run(); err != nil {
		panic(err)
	}

	fmt.Println("manager")
}
