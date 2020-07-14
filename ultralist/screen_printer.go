package ultralist

import (
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/fatih/color"
)

var (
	blue        = color.New(0, color.FgBlue)
	blueBold    = color.New(color.Bold, color.FgBlue)
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
}

// NewScreenPrinter creates a new screeen printer.
func NewScreenPrinter(unicodeSupport bool) *ScreenPrinter {
	w := new(io.Writer)
	formatter := &ScreenPrinter{Writer: w, UnicodeSupport: unicodeSupport}
	return formatter
}

// Print prints the output of ultralist to the terminal screen.
func (f *ScreenPrinter) Print(groupedTodos *GroupedTodos, printNotes bool, showStatus bool) {
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
			f.printTodo(tabby, todo, printNotes, showStatus)
		}
		tabby.AddLine()
	}
	tabby.Print()
}

<<<<<<< HEAD
func (f *ScreenPrinter) printTodo(tabby *tabby.Tabby, todo *Todo, printNotes bool, showStatus bool) {
	if showStatus {
		tabby.AddLine(
			f.formatID(todo.ID, todo.IsPriority),
			f.formatCompleted(todo.Completed),
			f.formatInformation(todo),
			f.formatDue(todo.Due, todo.IsPriority, todo.Completed),
			f.formatSubject(todo.Subject, todo.IsPriority))
	} else {
		tabby.AddLine(
			f.formatID(todo.ID, todo.IsPriority),
			f.formatCompleted(todo.Completed),
			f.formatDue(todo.Due, todo.IsPriority, todo.Completed),
			f.formatSubject(todo.Subject, todo.IsPriority))
	}

=======
func (f *ScreenPrinter) printTodo(tabby *tabby.Tabby, todo *Todo, printNotes bool) {
	tabby.AddLine(
		f.formatID(todo.ID, todo.IsPriority),
		f.formatCompleted(todo.Completed),
		f.formatDue(todo.Due, todo.IsPriority, todo.Completed),
		f.formatInformation(todo),
		f.formatSubject(todo.Subject, todo.IsPriority))
>>>>>>> eea610d... add to simple screen printer
	if printNotes {
		for nid, note := range todo.Notes {
			tabby.AddLine(
				"  "+cyan.Sprint(strconv.Itoa(nid)),
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
		} else {
			return white.Sprint("[x]")
		}
	}
	return white.Sprint("[ ]")
}

func (f *ScreenPrinter) formatDue(due string, isPriority bool, completed bool) string {
	if due == "" {
		return white.Sprint("          ")
	}
	dueTime, _ := time.Parse("2006-01-02", due)

	if isPriority {
		return f.printPriorityDue(dueTime, completed)
	}
	return f.printDue(dueTime, completed)
}

func (f *ScreenPrinter) formatInformation(todo *Todo) string {
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
	if todo.Archived {
		information = append(information, "A")
	} else {
		information = append(information, " ")
	}

	if todo.StartedDate != "" {
		information = append(information, "S")
	} else {
		information = append(information, " ")
	}
	return white.Sprint(strings.Join(information, " "))
}

func (f *ScreenPrinter) printDue(due time.Time, completed bool) string {
	if isToday(due) {
		return blue.Sprint("today     ")
	} else if isTomorrow(due) {
		return blue.Sprint("tomorrow  ")
	} else if isPastDue(due) && !completed {
		return red.Sprint(due.Format("Mon Jan 02"))
	}
	return blue.Sprint(due.Format("Mon Jan 02"))
}

func (f *ScreenPrinter) printPriorityDue(due time.Time, completed bool) string {
	if isToday(due) {
		return blueBold.Sprint("today     ")
	} else if isTomorrow(due) {
		return blueBold.Sprint("tomorrow  ")
	} else if isPastDue(due) && !completed {
		return redBold.Sprint(due.Format("Mon Jan 02"))
	}
	return blueBold.Sprint(due.Format("Mon Jan 02"))
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
