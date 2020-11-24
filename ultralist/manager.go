package ultralist

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type ManagerState string

const (
	ModeFocus         ManagerState = "focus_mode"
	ModeTodoEditing   ManagerState = "todo_editing"
	ModeListFiltering ManagerState = "list_filtering"
	ModeStatusEditing ManagerState = "status_editing"
	ModeSearching     ManagerState = "searching"
)

type Manager struct {
	App          *tview.Application
	MainArea     *tview.Grid
	TodoTextView *tview.TextView
	FilterArea   *tview.Flex
	CommandsArea *tview.Flex
	DebugArea    *tview.TextView

	State ManagerState

	TodoList     *TodoList
	GroupedTodos *GroupedTodos

	TodoIDs         []int
	SelectedTodoIdx int

	SearchTerm string
}

const (
	ColorBackground = 0x151515 // background color
	ColorForeground = "#d0d0d0"

	ColorBlue   = "#6A9FB5"
	ColorRed    = "#AC4142"
	ColorPurple = "#AA759F"
	ColorGreen  = "#90A959"
	ColorYellow = "#F4BF75"
	ColorGray   = "#606060"
)

// Todo commands
var (
	CmdDebug      = buildTextView("")
	CmdComplete   = buildTextView("c:complete")
	CmdUncomplete = buildTextView("c:uncomplete")

	CmdPrioritize   = buildTextView("p:prioritize")
	CmdUnprioritize = buildTextView("p:unprioritize")

	CmdArchive   = buildTextView("a:archive")
	CmdUnarchive = buildTextView("a:unarchive")

	CmdStatus = buildTextView("s:status")
	CmdDue    = buildTextView("d:due")
	CmdDelete = buildTextView("x:delete")
	CmdSearch = buildTextView("Search: ")
)

// Todo list commands
var ()

func NewManager(todoList *TodoList) *Manager {
	textView := tview.NewTextView()
	textView.SetBackgroundColor(tcell.NewHexColor(0x101010))
	textView.SetDynamicColors(true)
	textView.SetBorder(false)
	textView.SetRegions(true)

	mainArea := tview.NewGrid()
	mainArea.SetBackgroundColor(tcell.NewHexColor(ColorBackground))
	mainArea.SetBorder(false)
	mainArea.SetRows(2, -1, 2)

	filterArea := tview.NewFlex()
	filterArea.SetBorder(false)
	filterArea.SetBackgroundColor(tcell.NewHexColor(ColorBackground))

	commandsArea := tview.NewFlex()
	commandsArea.SetBorder(false)
	commandsArea.SetBackgroundColor(tcell.NewHexColor(ColorBackground))

	debugArea := tview.NewTextView()

	mainArea.AddItem(
		filterArea,
		0,     // row
		0,     // column
		1,     // rowSpan
		1,     // colSpan
		0,     // minGridHeight
		0,     // minGridWidth
		false) // focus

	mainArea.AddItem(
		textView,
		1,    // row
		0,    // column
		1,    // rowSpan
		1,    // colSpan
		0,    // minGridHeight
		0,    // minGridWidth
		true) // focus

	mainArea.AddItem(
		commandsArea,
		2,     // row
		0,     // column
		1,     // rowSpan
		1,     // colSpan
		0,     // minGridHeight
		0,     // minGridWidth
		false) // focus

	grouper := &Grouper{}
	groupedTodos := grouper.GroupByProject(todoList.Todos())

	app := tview.NewApplication().SetRoot(mainArea, true).EnableMouse(true)

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	manager := &Manager{
		App:             app,
		TodoList:        todoList,
		GroupedTodos:    groupedTodos,
		TodoTextView:    textView,
		CommandsArea:    commandsArea,
		FilterArea:      filterArea,
		MainArea:        mainArea,
		DebugArea:       debugArea,
		SelectedTodoIdx: 0,
		SearchTerm:      "",
	}

	manager.App.SetInputCapture(manager.inputCapture)
	manager.switchStateToModeTodoEditing()
	manager.drawTodos()

	manager.CommandsArea.AddItem(CmdDebug, 0, 1, false)

	return manager
}

func (m *Manager) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch m.State {
	case ModeTodoEditing:
		m.todoEventsInputCapture(event)
	case ModeFocus:
		m.focusModeInputCapture(event)
	}

	// handle global events
	m.globalEventsInputCapture(event)

	m.drawTodos()

	return event
}

func (m *Manager) focusModeInputCapture(event *tcell.EventKey) {
	if event.Rune() == 'j' || event.Key() == tcell.KeyDown {
		if m.SelectedTodoIdx < len(m.TodoIDs)-1 {
			m.SelectedTodoIdx += 1
		}
		m.switchStateToModeTodoEditing()
	}
	if event.Rune() == 'k' || event.Key() == tcell.KeyUp {
		if m.SelectedTodoIdx > 0 {
			m.SelectedTodoIdx -= 1
		}
		m.switchStateToModeTodoEditing()
	}
}

func (m *Manager) todoEventsInputCapture(event *tcell.EventKey) {
	if len(m.TodoIDs) == 0 {
		return
	}

	todo := m.TodoList.FindByID(m.TodoIDs[m.SelectedTodoIdx])

	if event.Rune() == 'j' || event.Key() == tcell.KeyDown {
		if m.SelectedTodoIdx < len(m.TodoIDs)-1 {
			m.SelectedTodoIdx += 1
		}
	}
	if event.Rune() == 'k' || event.Key() == tcell.KeyUp {
		if m.SelectedTodoIdx > 0 {
			m.SelectedTodoIdx -= 1
		}
	}

	// complete
	if event.Rune() == 'c' {
		if todo.Completed {
			m.TodoList.Uncomplete(m.TodoIDs[m.SelectedTodoIdx])
		} else {
			m.TodoList.Complete(m.TodoIDs[m.SelectedTodoIdx])
		}
	}

	// prioritize
	if event.Rune() == 'p' {
		if todo.IsPriority {
			m.TodoList.Unprioritize(m.TodoIDs[m.SelectedTodoIdx])
		} else {
			m.TodoList.Prioritize(m.TodoIDs[m.SelectedTodoIdx])
		}
	}

	// archive
	if event.Rune() == 'a' {
		if todo.Archived {
			m.TodoList.Unarchive(m.TodoIDs[m.SelectedTodoIdx])
		} else {
			m.TodoList.Archive(m.TodoIDs[m.SelectedTodoIdx])
		}
	}

	// set status
	if event.Rune() == 's' {
		m.switchStateToModeStatusEditing()
	}
}

func (m *Manager) globalEventsInputCapture(event *tcell.EventKey) {
	if event.Rune() == '/' {
		m.switchStateToSearching()
	}

	// switch the app to a different context
	if event.Rune() == ' ' {
		if m.State == ModeTodoEditing {
			m.switchStateToModeListFiltering()
		} else {
			m.switchStateToModeTodoEditing()
		}
	}

	// switch to focus mode
	if event.Key() == tcell.KeyEsc {
		if m.State != ModeFocus {
			m.switchStateToModeFocus()
		} else {
			m.switchStateToModeTodoEditing()
		}
	}

	// quit the app
	if event.Rune() == 'q' {
		m.App.Stop()
	}

}

func (m *Manager) switchStateToModeListFiltering() {
	m.State = ModeListFiltering

	m.CommandsArea.Clear()
	m.FilterArea.Clear()
	m.FilterArea.AddItem(tview.NewTextView().SetText("List filtering"), 0, 1, false)
}

func (m *Manager) switchStateToModeTodoEditing() {
	m.State = ModeTodoEditing
	m.drawTodos()
}

func (m *Manager) switchStateToModeFocus() {
	m.State = ModeFocus
	m.CommandsArea.Clear()
	m.drawTodos()
}

func (m *Manager) switchStateToSearching() {
	m.SearchTerm = ""
	m.State = ModeSearching
	m.FilterArea.Clear()

	input := tview.NewInputField().SetLabel("Search: ").SetFieldWidth(10)
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		m.SearchTerm = input.GetText()
		m.drawTodos()
		return event
	})

	input.SetDoneFunc(func(key tcell.Key) {
		m.SelectedTodoIdx = 0
		m.switchStateToModeTodoEditing()
		m.App.SetFocus(m.MainArea)
	})

	m.FilterArea.AddItem(input, 0, 1, false)
	m.App.SetFocus(input)
	m.drawTodos()
}

func (m *Manager) switchStateToModeStatusEditing() {
	m.State = ModeStatusEditing
	m.CommandsArea.Clear()
	todo := m.TodoList.FindByID(m.TodoIDs[m.SelectedTodoIdx])

	input := tview.NewInputField().SetLabel("Set status: ").SetFieldWidth(10)
	input.SetText(todo.Status)

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			todo.Status = input.GetText()
		}
		m.switchStateToModeTodoEditing()
	})

	m.CommandsArea.AddItem(input, 0, 1, false)
	m.App.SetFocus(input)
	m.drawTodos()
}

func (m *Manager) RunManager() {
	if err := m.App.Run(); err != nil {
		panic(err)
	}
}

func (m *Manager) drawTodos() {

	var todoIDs []int
	viewPrinter := &ViewPrinter{}

	var keys []string
	for key := range m.GroupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	m.TodoTextView.Clear()
	m.TodoTextView.Highlight("-1")

	// searchRegex := regexp.MustCompile(fmt.Sprintf("%s", m.SearchTerm))

	totalDisplayedTodos := 0
	for _, key := range keys {
		var filteredTodos []*Todo
		for _, todo := range m.GroupedTodos.Groups[key] {
			if m.SearchTerm != "" {
				match, _ := regexp.MatchString(m.SearchTerm, todo.Subject)
				if !match {
					break
				}
			}
			filteredTodos = append(filteredTodos, todo)
		}

		if len(filteredTodos) > 0 {
			fmt.Fprintf(m.TodoTextView, "\n[%s]%s[%s]\n", ColorBlue, key, ColorForeground)
		}

		for _, todo := range filteredTodos {
			fmt.Fprintf(
				m.TodoTextView,
				"[\"%v\"]%s  %s  %s  %s  %s[\"\"]\n",
				totalDisplayedTodos,
				viewPrinter.FormatID(todo),
				viewPrinter.FormatCompleted(todo),
				viewPrinter.FormatDue(todo),
				viewPrinter.FormatStatus(todo),
				viewPrinter.FormatSubject(todo),
			)

			todoIDs = append(todoIDs, todo.ID)

			if totalDisplayedTodos == m.SelectedTodoIdx && m.State == ModeTodoEditing {
				m.buildTodoCommandsMenu(todo)
				m.TodoTextView.Highlight(strconv.Itoa(m.SelectedTodoIdx))
			}

			totalDisplayedTodos++
		}
	}
	m.TodoIDs = todoIDs
}

func (m *Manager) buildTodoCommandsMenu(todo *Todo) {
	m.CommandsArea.Clear()

	if todo.Completed {
		m.CommandsArea.AddItem(CmdUncomplete, 0, 1, false)
	} else {
		m.CommandsArea.AddItem(CmdComplete, 0, 1, false)
	}

	if todo.IsPriority {
		m.CommandsArea.AddItem(CmdUnprioritize, 0, 1, false)
	} else {
		m.CommandsArea.AddItem(CmdPrioritize, 0, 1, false)
	}

	if todo.Archived {
		m.CommandsArea.AddItem(CmdUnarchive, 0, 1, false)
	} else {
		m.CommandsArea.AddItem(CmdArchive, 0, 1, false)
	}

	m.CommandsArea.AddItem(CmdStatus, 0, 1, false)
	m.CommandsArea.AddItem(CmdDue, 0, 1, false)
	m.CommandsArea.AddItem(CmdDelete, 0, 1, false)
}

func buildTextView(label string) *tview.TextView {
	view := tview.NewTextView()
	view.SetBackgroundColor(tcell.NewHexColor(ColorBackground))
	view.SetText(label)
	return view
}
