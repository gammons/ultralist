package cli

import (
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/fatih/color"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	blue        = color.New(0, color.FgBlue)
	blueBold    = color.New(color.Bold, color.FgBlue)
	green       = color.New(0, color.FgGreen)
	greenBold   = color.New(color.Bold, color.FgGreen)
	cyan        = color.New(0, color.FgCyan)
	cyanBold    = color.New(color.Bold, color.FgCyan)
	magenta     = color.New(0, color.FgMagenta)
	magentaBold = color.New(color.Bold, color.FgMagenta)
	red         = color.New(0, color.FgRed)
	redBold     = color.New(color.Bold, color.FgRed)
	white       = color.New(0, color.FgWhite)
	whiteBold   = color.New(color.Bold, color.FgWhite)
	yellow      = color.New(0, color.FgYellow)
	yellowBold  = color.New(color.Bold, color.FgYellow)
)

// ScreenPrinter is the default struct of this file
type ScreenPrinter struct {
	Writer         *io.Writer
	UnicodeSupport bool
	ListNotes      bool
	ShowStatus     bool
}

// NewScreenPrinter creates a new screeen printer.
func NewScreenPrinter(unicodeSupport bool, listNotes bool, showStatus bool) *ScreenPrinter {
	w := new(io.Writer)
	formatter := &ScreenPrinter{
		Writer:         w,
		UnicodeSupport: unicodeSupport,
		ListNotes:      listNotes,
		ShowStatus:     showStatus,
	}
	return formatter
}

// Print prints the output of ultralist to the terminal screen.
func (f *ScreenPrinter) Print(groupedTodos *ultralist.GroupedTodos) {
	var keys []string
	for key := range groupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	tabby := tabby.NewCustom(tabwriter.NewWriter(color.Output, 0, 0, 2, ' ', 0))
	tabby.AddLine()
	for _, key := range keys {
		tabby.AddLine(cyan.Sprint(key))
		for _, todo := range groupedTodos.Groups[key] {
			f.printTodo(tabby, todo)
		}
		tabby.AddLine()
	}
	tabby.Print()
}

func (f *ScreenPrinter) printTodo(tabby *tabby.Tabby, todo *ultralist.Todo) {
	if f.ShowStatus {
		tabby.AddLine(
			f.formatID(todo.ID, todo.IsPriority),
			f.formatCompleted(todo.Completed),
			f.formatDue(todo),
			f.formatStatus(todo.Status, todo.IsPriority),
			f.formatSubject(todo.Subject, todo.IsPriority))
	} else {
		tabby.AddLine(
			f.formatID(todo.ID, todo.IsPriority),
			f.formatCompleted(todo.Completed),
			f.formatDue(todo),
			f.formatStatus(todo.Status, todo.IsPriority),
			f.formatSubject(todo.Subject, todo.IsPriority))
	}

	if f.ListNotes {
		for nid, note := range todo.Notes {
			tabby.AddLine(
				"  "+cyan.Sprint(strconv.Itoa(nid)),
				white.Sprint(""),
				white.Sprint(""),
				white.Sprint(""),
				white.Sprint(""),
				white.Sprint(note))
		}
	}
}

func (f *ScreenPrinter) formatID(ID int, isPriority bool) string {
	if isPriority {
		return yellowBold.Sprint(strconv.Itoa(ID))
	}
	return yellow.Sprint(strconv.Itoa(ID))
}

func (f *ScreenPrinter) formatCompleted(completed bool) string {
	if completed {
		if f.UnicodeSupport {
			return white.Sprint("[✔]")
		}
		return white.Sprint("[x]")
	}
	return white.Sprint("[ ]")
}

func (f *ScreenPrinter) formatDue(todo *ultralist.Todo) string {
	if todo.Due == "" {
		return white.Sprint("          ")
	}

	if todo.IsPriority {
		return f.printPriorityDue(todo)
	}
	return f.printDue(todo)
}

func (f *ScreenPrinter) formatStatus(status string, isPriority bool) string {
	if status == "" {
		return green.Sprint("          ")
	}

	if len(status) < 10 {
		for x := len(status); x <= 10; x++ {
			status += " "
		}
	}

	statusRune := []rune(status)

	if isPriority {
		return greenBold.Sprintf("%-10v", string(statusRune[0:10]))
	}
	return green.Sprintf("%-10s", string(statusRune[0:10]))
}

func (f *ScreenPrinter) formatInformation(todo *ultralist.Todo) string {
	var information []string
	if todo.IsPriority {
		information = append(information, "*")
	} else {
		information = append(information, " ")
	}
	if todo.HasNotes() {
		information = append(information, "N")
	} else {
		information = append(information, " ")
	}

	return white.Sprint(strings.Join(information, ""))
}

func (f *ScreenPrinter) printDue(todo *ultralist.Todo) string {
	dueTime, _ := time.Parse(ultralist.DateFormat, todo.Due)

	if todo.DueToday() {
		return blue.Sprint("today     ")
	} else if todo.DueTomorrow() {
		return blue.Sprint("tomorrow  ")
	} else if todo.PastDue() && !todo.Completed {
		return red.Sprint(dueTime.Format("Mon Jan 02"))
	}
	return blue.Sprint(dueTime.Format("Mon Jan 02"))
}

func (f *ScreenPrinter) printPriorityDue(todo *ultralist.Todo) string {
	dueTime, _ := time.Parse(ultralist.DateFormat, todo.Due)

	if todo.DueToday() {
		return blueBold.Sprint("today     ")
	} else if todo.DueTomorrow() {
		return blueBold.Sprint("tomorrow  ")
	} else if todo.PastDue() && !todo.Completed {
		return redBold.Sprint(dueTime.Format("Mon Jan 02"))
	}
	return blueBold.Sprint(dueTime.Format("Mon Jan 02"))
}

func (f *ScreenPrinter) formatSubject(subject string, isPriority bool) string {
	splitted := strings.Split(subject, " ")

	if isPriority {
		return f.printPrioritySubject(splitted)
	}
	return f.printSubject(splitted)
}

func (f *ScreenPrinter) printPrioritySubject(splitted []string) string {
	coloredWords := []string{}
	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, magentaBold.Sprint(word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, redBold.Sprint(word))
		} else {
			coloredWords = append(coloredWords, whiteBold.Sprint(word))
		}
	}
	return strings.Join(coloredWords, " ")
}

func (f *ScreenPrinter) printSubject(splitted []string) string {
	coloredWords := []string{}
	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, magenta.Sprint(word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, red.Sprint(word))
		} else {
			coloredWords = append(coloredWords, white.Sprint(word))
		}
	}
	return strings.Join(coloredWords, " ")
}
