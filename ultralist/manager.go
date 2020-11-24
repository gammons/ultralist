package ultralist

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Manager struct {
	TodoTextView *tview.TextView
	CommandsArea *tview.Flex
	MainArea     *tview.Grid
	DebugArea    *tview.TextView
	App          *tview.Application
	Commands     map[string]*tview.TextView
	TodoList     *TodoList

	TodoIDs         []int
	SelectedTodoID  int
	SelectedTodoIdx int
}

const (
	ColorBackground = 0x151515 // background color
	ColorForeground = "#d0d0d0"

	ColorBlue   = "#6A9FB5"
	ColorRed    = "#AC4142"
	ColorPurple = "#AA759F"
	ColorGreen  = "#90A959"
	ColorYellow = "#F4BF75"
)

func NewManager(todoList *TodoList) *Manager {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetBorder(false)
	textView.SetRegions(true)

	mainArea := tview.NewGrid()
	mainArea.SetBackgroundColor(tcell.NewHexColor(ColorBackground))
	mainArea.SetBorder(false)
	mainArea.SetRows(3, -1, 3)

	commandsArea := tview.NewFlex()
	commandsArea.SetBorder(false)
	commandsArea.SetBackgroundColor(tcell.NewHexColor(ColorBackground))

	debugArea := tview.NewTextView()

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

	manager := &Manager{
		TodoList:     todoList,
		TodoTextView: textView,
		CommandsArea: commandsArea,
		MainArea:     mainArea,
		DebugArea:    debugArea,
	}

	var commands map[string]*tview.TextView
	commands = make(map[string]*tview.TextView)
	commands["debug"] = manager.buildTextView("")
	manager.CommandsArea.AddItem(commands["debug"], 0, 1, false)
	commands["complete"] = manager.buildTextView("c:complete")
	manager.CommandsArea.AddItem(commands["complete"], 0, 1, false)

	commands["prioritize"] = manager.buildTextView("p:prioritize")
	commands["archive"] = manager.buildTextView("a:archive")
	manager.Commands = commands

	return manager
}

func (m *Manager) RunManager() {
	m.App = tview.NewApplication().SetRoot(m.MainArea, true).EnableMouse(true)
	m.SelectedTodoIdx = 0

	m.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			if m.SelectedTodoIdx < len(m.TodoIDs)-1 {
				m.SelectedTodoIdx += 1
			}
		}
		if event.Rune() == 'k' {
			if m.SelectedTodoIdx > 0 {
				m.SelectedTodoIdx -= 1
			}
		}

		if event.Rune() == 'c' {
			todo := m.TodoList.FindByID(m.TodoIDs[m.SelectedTodoIdx])
			if todo.Completed {
				m.TodoList.Uncomplete(m.TodoIDs[m.SelectedTodoIdx])
			} else {
				m.TodoList.Complete(m.TodoIDs[m.SelectedTodoIdx])
			}
		}

		if event.Key() == tcell.KeyTab {
			fmt.Println("tab")
		}

		m.drawTodos()
		m.TodoTextView.Highlight(strconv.Itoa(m.SelectedTodoIdx))

		return event
	})

	m.drawTodos()
	if err := m.App.Run(); err != nil {
		panic(err)
	}
}

func (m *Manager) drawTodos() {
	grouper := &Grouper{}
	groupedTodos := grouper.GroupByProject(m.TodoList.Todos())

	var todoIDs []int
	viewPrinter := &ViewPrinter{}

	var keys []string
	for key := range groupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	m.TodoTextView.Clear()

	totalDisplayedTodos := 0
	for _, key := range keys {
		fmt.Fprintf(m.TodoTextView, "\n[%s]%s[%s]\n", ColorBlue, key, ColorForeground)

		for _, todo := range groupedTodos.Groups[key] {
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

			if totalDisplayedTodos == m.SelectedTodoIdx {
				m.buildCommandsMenu(todo)
			}

			totalDisplayedTodos++
		}
	}
	m.TodoIDs = todoIDs
}

func (m *Manager) buildTextView(label string) *tview.TextView {
	view := tview.NewTextView()
	view.SetBackgroundColor(tcell.NewHexColor(ColorBackground))
	view.SetText(label)
	return view
}

func (m *Manager) buildCommandsMenu(todo *Todo) {
	m.Commands["debug"].SetText(fmt.Sprintf("todoIDs: %v, id:%v", m.TodoIDs, m.SelectedTodoIdx))
	if todo.Completed {
		m.Commands["complete"].SetText("c:uncomplete")
	} else {
		m.Commands["complete"].SetText("c:complete")
	}
}
