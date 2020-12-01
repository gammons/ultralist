package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/ultralist/ultralist/ultralist"
)

// ManagerState represents the state of the tui.
type ManagerState string

const (
	// ModeFocus is "focus mode", where the todos are listed without anything else on the screen.
	ModeFocus ManagerState = "focus_mode"

	// ModeTodoManaging is when the user is actively managing their todos.  There is a highlighted todo,
	// and there are text views at the bottom the represent the commands someone can do on the highlighted todo.
	ModeTodoManaging ManagerState = "todo_managing"

	// ModeEditing is when an input box is focused.  This is mainly used to disable global commands like
	// quitting the app with "q" and switching to other modes.
	ModeEditing ManagerState = "todo_editing"

	// ModeListFiltering is when the user is actively filtering the list.
	// It can be thought of as a subset of ModeEditing.
	ModeListFiltering ManagerState = "list_filtering"

	// ModeSearching is when a user is searching
	ModeSearching ManagerState = "searching"
)

// Manager represents the core of the TUI
type Manager struct {
	App          *tview.Application
	MainArea     *tview.Grid
	TodoTextView *TodoTextView
	FilterArea   *tview.Flex
	CommandsArea *tview.Flex
	DebugArea    *tview.TextView

	State ManagerState

	TodoList *ultralist.TodoList

	// TodoIDs         []int
	// SelectedTodoIdx int

	Grouping   ultralist.Grouping
	TodoFilter *ultralist.Filter
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

	CmdSearch = buildTextView("/:Search")
	CmdGroup  = buildTextView("g:group")
)

// Todo editing controls
var (
	StatusInput = tview.NewInputField()
	DueInput    = tview.NewInputField()
)

// Todo list commands
var (
	GroupSelect     = tview.NewDropDown()
	CompletedSelect = tview.NewDropDown()
	SearchInput     = tview.NewInputField()
)

func NewManager(todoList *ultralist.TodoList) *Manager {
	textView := NewTodoTextView(todoList)

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
		textView.TextView,
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

	app := tview.NewApplication().SetRoot(mainArea, true).EnableMouse(true)

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	manager := &Manager{
		App:          app,
		TodoList:     todoList,
		TodoTextView: textView,
		CommandsArea: commandsArea,
		FilterArea:   filterArea,
		MainArea:     mainArea,
		DebugArea:    debugArea,
		TodoFilter:   &ultralist.Filter{},
		Grouping:     ultralist.GroupByNone,
	}

	manager.App.SetInputCapture(manager.inputCapture)

	// set up the inputs
	manager.setupSearchInput()
	manager.setupGroupSelect()
	manager.setupCompletedSelect()
	manager.setupStatusInput()
	manager.setupDueInput()

	// set up initial state for the app.
	manager.switchStateToModeTodoManaging()
	manager.drawTodos()

	// temporary to debug.
	manager.CommandsArea.AddItem(CmdDebug, 0, 1, false)

	return manager
}

func (m *Manager) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch m.State {
	case ModeTodoManaging:
		m.todoEventsInputCapture(event)
	case ModeFocus:
		m.focusModeInputCapture(event)
	}

	m.drawTodos()

	return event
}

func (m *Manager) focusModeInputCapture(event *tcell.EventKey) {
	if event.Rune() == 'j' || event.Key() == tcell.KeyDown {
		m.TodoTextView.HighlightNextTodo()
	}
	if event.Rune() == 'k' || event.Key() == tcell.KeyUp {
		m.TodoTextView.HighlightPrevTodo()
		m.switchStateToModeTodoManaging()
	}

	if event.Key() == tcell.KeyEnter ||
		event.Rune() == ' ' {
		m.switchStateToModeTodoManaging()
	}

	if event.Rune() == '/' {
		m.switchStateToModeListFiltering()
		m.App.SetFocus(SearchInput)
	}

	if event.Rune() == 'g' {
		m.switchStateToModeListFiltering()
		m.App.SetFocus(GroupSelect)
	}

	// quit the app
	if event.Rune() == 'q' {
		m.App.Stop()
	}
}

func (m *Manager) todoEventsInputCapture(event *tcell.EventKey) {
	todo := m.TodoTextView.SelectedTodo()
	if todo == nil {
		return
	}

	if event.Rune() == 'j' || event.Key() == tcell.KeyDown {
		m.TodoTextView.HighlightNextTodo()
	}
	if event.Rune() == 'k' || event.Key() == tcell.KeyUp {
		m.TodoTextView.HighlightPrevTodo()
	}

	// complete
	if event.Rune() == 'c' {
		if todo.Completed {
			todo.Uncomplete()
		} else {
			todo.Complete()
		}
	}

	// prioritize
	if event.Rune() == 'p' {
		if todo.IsPriority {
			todo.Unprioritize()
		} else {
			todo.Prioritize()
		}
	}

	// archive
	if event.Rune() == 'a' {
		if todo.Archived {
			todo.Unarchive()
		} else {
			todo.Archive()
		}
	}

	// set status
	if event.Rune() == 's' {
		m.editTodoStatus()
	}

	// set due
	if event.Rune() == 'd' {
		m.editTodoDue()
	}

	if event.Rune() == '/' {
		m.switchStateToModeListFiltering()
		m.App.SetFocus(SearchInput)
	}

	if event.Rune() == 'g' {
		m.switchStateToModeListFiltering()
		m.App.SetFocus(GroupSelect)
	}

	// quit the app
	if event.Rune() == 'q' {
		m.App.Stop()
	}

	if event.Key() == tcell.KeyEsc {
		m.switchStateToModeFocus()
	}

	if event.Key() == tcell.KeyTab {
		m.switchStateToModeListFiltering()
		m.App.SetFocus(SearchInput)
	}
}

func (m *Manager) switchStateToModeListFiltering() {
	m.State = ModeListFiltering
	m.FilterArea.Clear()

	m.CommandsArea.Clear()
	m.CommandsArea.AddItem(CmdSearch, 0, 1, false)
	m.CommandsArea.AddItem(CmdGroup, 0, 1, false)

	m.FilterArea.AddItem(SearchInput, 0, 1, false)
	m.FilterArea.AddItem(GroupSelect, 0, 1, false)
	m.FilterArea.AddItem(CompletedSelect, 0, 1, false)

	m.drawTodos()
}

func (m *Manager) switchStateToModeTodoManaging() {
	m.State = ModeTodoManaging
	m.FilterArea.Clear()
	m.drawTodos()
}

func (m *Manager) switchStateToModeFocus() {
	m.State = ModeFocus
	m.CommandsArea.Clear()
	m.drawTodos()
}

func (m *Manager) switchStateToSearching() {
	m.TodoFilter.Subject = ""
	m.TodoFilter.HasSubject = false
	m.State = ModeSearching
	m.FilterArea.Clear()

	input := tview.NewInputField().SetLabel("Search: ").SetFieldWidth(10)
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		m.TodoFilter.Subject = input.GetText()
		m.TodoFilter.HasSubject = true
		m.drawTodos()
		return event
	})

	input.SetDoneFunc(func(key tcell.Key) {
		m.TodoTextView.ResetSelectedTodoIdx()
		m.switchStateToModeTodoManaging()
		m.App.SetFocus(m.MainArea)
	})

	m.FilterArea.AddItem(input, 0, 1, false)
	m.App.SetFocus(input)
	m.drawTodos()
}

func (m *Manager) editTodoStatus() {
	m.State = ModeEditing
	m.CommandsArea.Clear()

	todo := m.TodoTextView.SelectedTodo()
	if todo == nil {
		return
	}

	StatusInput.SetText(todo.Status)

	StatusInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			todo.Status = StatusInput.GetText()
		}
		m.switchStateToModeTodoManaging()
	})

	m.CommandsArea.AddItem(StatusInput, 0, 1, false)
	m.App.SetFocus(StatusInput)
	m.drawTodos()
}

func (m *Manager) editTodoDue() {
	m.State = ModeEditing
	m.CommandsArea.Clear()

	todo := m.TodoTextView.SelectedTodo()
	if todo == nil {
		return
	}

	DueInput.SetText(todo.Due)

	DueInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			todo.Due = DueInput.GetText()
		}
		m.switchStateToModeTodoManaging()
	})

	m.CommandsArea.AddItem(DueInput, 0, 1, false)
	m.App.SetFocus(DueInput)
	m.drawTodos()
}

func (m *Manager) RunManager() {
	if err := m.App.Run(); err != nil {
		panic(err)
	}
}

func (m *Manager) buildTodoCommandsMenu(todo *ultralist.Todo) {
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

func (m *Manager) setupSearchInput() {
	SearchInput := SearchInput.SetLabel("Search: ").SetFieldWidth(10)
	SearchInput.SetText(m.TodoFilter.Subject)
	SearchInput.SetFieldBackgroundColor(tcell.NewHexColor(0x505050))
	SearchInput.SetLabelColor(tcell.NewHexColor(0xd0d0d0))

	SearchInput.SetChangedFunc(func(term string) {
		m.TodoFilter.Subject = term
		m.TodoFilter.HasSubject = true
		CmdSearch.SetText(fmt.Sprintf("'%s'", term))
		m.drawTodos()
	})

	SearchInput.SetDoneFunc(func(key tcell.Key) {
		m.TodoTextView.ResetSelectedTodoIdx()
		if key == tcell.KeyTab || key == tcell.KeyBacktab {
			m.App.SetFocus(GroupSelect)
		} else {
			m.switchStateToModeTodoManaging()
			m.App.SetFocus(m.MainArea)
		}
	})
}

func (m *Manager) setupCompletedSelect() {
	CompletedSelect.SetLabel("Completed: ")
	CompletedSelect.SetFieldBackgroundColor(tcell.NewHexColor(0x505050))
	CompletedSelect.SetLabelColor(tcell.NewHexColor(0xd0d0d0))

	CompletedSelect.AddOption("All", func() {
		m.TodoFilter.HasCompleted = false
		m.drawTodos()
	})
	CompletedSelect.AddOption("true", func() {
		m.TodoFilter.HasCompleted = true
		m.TodoFilter.Completed = true
		m.drawTodos()
	})
	CompletedSelect.AddOption("false", func() {
		m.TodoFilter.HasCompleted = true
		m.TodoFilter.Completed = false
		m.drawTodos()
	})

	CompletedSelect.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyBacktab {
			m.App.SetFocus(GroupSelect)
		}
		if key == tcell.KeyTab {
			m.switchStateToModeTodoManaging()
			m.App.SetFocus(m.MainArea)
		}
	})
}

func (m *Manager) setupGroupSelect() {
	GroupSelect.SetLabel("Group: ").SetFieldWidth(10)
	GroupSelect.SetFieldBackgroundColor(tcell.NewHexColor(0x505050))
	GroupSelect.SetLabelColor(tcell.NewHexColor(0xd0d0d0))

	GroupSelect.AddOption("None", func() {
		m.Grouping = ultralist.GroupByNone
		m.drawTodos()
		m.switchStateToModeTodoManaging()
		m.App.SetFocus(m.MainArea)

	})
	GroupSelect.AddOption("Project", func() {
		m.Grouping = ultralist.GroupByProject
		m.drawTodos()
		m.switchStateToModeTodoManaging()
		m.App.SetFocus(m.MainArea)
	})
	GroupSelect.AddOption("Context", func() {
		m.Grouping = ultralist.GroupByContext
		m.drawTodos()
		m.switchStateToModeTodoManaging()
		m.App.SetFocus(m.MainArea)
	})
	GroupSelect.AddOption("Status", func() {
		m.Grouping = ultralist.GroupByStatus
		m.drawTodos()
		m.switchStateToModeTodoManaging()
		m.App.SetFocus(m.MainArea)
	})

	GroupSelect.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyBacktab {
			m.App.SetFocus(SearchInput)
		}
		if key == tcell.KeyTab {
			m.App.SetFocus(CompletedSelect)
		}
	})
}

func (m *Manager) setupStatusInput() {
	StatusInput.SetLabel("Set status: ").SetFieldWidth(10)
	StatusInput.SetFieldBackgroundColor(tcell.NewHexColor(0x505050))
	StatusInput.SetLabelColor(tcell.NewHexColor(0xd0d0d0))
}

func (m *Manager) setupDueInput() {
	DueInput.SetLabel("Set due: ").SetFieldWidth(10)
	DueInput.SetFieldBackgroundColor(tcell.NewHexColor(0x505050))
	DueInput.SetLabelColor(tcell.NewHexColor(0xd0d0d0))
}

func (m *Manager) drawTodos() {
	m.TodoTextView.DrawTodos(m.TodoFilter, m.Grouping)
}

func buildTextView(label string) *tview.TextView {
	view := tview.NewTextView()
	view.SetBackgroundColor(tcell.NewHexColor(ColorBackground))
	view.SetText(label)
	return view
}
