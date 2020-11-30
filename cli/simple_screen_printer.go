package cli

import (
	"fmt"
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

// ScreenPrinter is the default struct of this file
type SimpleScreenPrinter struct {
	Writer         *io.Writer
	UnicodeSupport bool
	ListNotes      bool
	ShowStatus     bool
}

// NewScreenPrinter creates a new screeen printer.
func NewSimpleScreenPrinter(unicodeSupport bool, listNotes bool, showStatus bool) *SimpleScreenPrinter {
	w := new(io.Writer)
	formatter := &SimpleScreenPrinter{
		Writer:         w,
		UnicodeSupport: unicodeSupport,
		ListNotes:      listNotes,
		ShowStatus:     showStatus,
	}
	return formatter
}

// Print prints the output of ultralist to the terminal screen.
func (f *SimpleScreenPrinter) Print(groupedTodos *ultralist.GroupedTodos) {
	var keys []string
	for key := range groupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	tabby := tabby.NewCustom(tabwriter.NewWriter(color.Output, 0, 0, 2, ' ', 0))
	tabby.AddLine()
	for _, key := range keys {
		tabby.AddLine(fmt.Sprint(key))
		for _, todo := range groupedTodos.Groups[key] {
			f.printTodo(tabby, todo)
		}
		tabby.AddLine()
	}
	tabby.Print()
}

func (f *SimpleScreenPrinter) printTodo(tabby *tabby.Tabby, todo *ultralist.Todo) {
	if f.ShowStatus {
		tabby.AddLine(
			f.formatID(todo.ID, todo.IsPriority),
			f.formatCompleted(todo.Completed),
			f.formatInformation(todo),
			f.formatDue(todo),
			f.formatSubject(todo.Subject, todo.IsPriority))
	} else {
		tabby.AddLine(
			f.formatID(todo.ID, todo.IsPriority),
			f.formatCompleted(todo.Completed),
			f.formatDue(todo),
			f.formatSubject(todo.Subject, todo.IsPriority))
	}
	if f.ListNotes {
		for nid, note := range todo.Notes {
			tabby.AddLine(
				"  "+fmt.Sprint(strconv.Itoa(nid)),
				fmt.Sprint(""),
				fmt.Sprint(""),
				fmt.Sprint(""),
				fmt.Sprint(note))
		}
	}
}

func (f *SimpleScreenPrinter) formatID(ID int, isPriority bool) string {
	if isPriority {
		return fmt.Sprint(strconv.Itoa(ID))
	}
	return fmt.Sprint(strconv.Itoa(ID))
}

func (f *SimpleScreenPrinter) formatCompleted(completed bool) string {
	if completed {
		if f.UnicodeSupport {
			return fmt.Sprint("[✔]")
		} else {
			return fmt.Sprint("[x]")
		}
	}
	return fmt.Sprint("[ ]")
}

func (f *SimpleScreenPrinter) formatDue(todo *ultralist.Todo) string {
	if todo.Due == "" {
		return fmt.Sprint("          ")
	}

	if todo.IsPriority {
		return f.printPriorityDue(todo)
	}
	return f.printDue(todo)
}

func (f *SimpleScreenPrinter) formatInformation(todo *ultralist.Todo) string {
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
	return fmt.Sprint(strings.Join(information, ""))
}

func (f *SimpleScreenPrinter) printDue(todo *ultralist.Todo) string {
	dueTime, _ := time.Parse(ultralist.DateFormat, todo.Due)

	if todo.DueToday() {
		return fmt.Sprint("today     ")
	} else if todo.DueTomorrow() {
		return fmt.Sprint("tomorrow  ")
	} else if todo.PastDue() && !todo.Completed {
		return fmt.Sprint(dueTime.Format("Mon Jan 02"))
	}
	return fmt.Sprint(dueTime.Format("Mon Jan 02"))
}

func (f *SimpleScreenPrinter) printPriorityDue(todo *ultralist.Todo) string {
	dueTime, _ := time.Parse(ultralist.DateFormat, todo.Due)

	if todo.DueToday() {
		return fmt.Sprint("today     ")
	} else if todo.DueTomorrow() {
		return fmt.Sprint("tomorrow  ")
	} else if todo.PastDue() {
		return fmt.Sprint(dueTime.Format("Mon Jan 02"))
	}
	return fmt.Sprint(dueTime.Format("Mon Jan 02"))
}

func (f *SimpleScreenPrinter) formatSubject(subject string, isPriority bool) string {
	splitted := strings.Split(subject, " ")

	if isPriority {
		return f.printPrioritySubject(splitted)
	}
	return f.printSubject(splitted)
}

func (f *SimpleScreenPrinter) printPrioritySubject(splitted []string) string {
	coloredWords := []string{}
	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprint(word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprint(word))
		} else {
			coloredWords = append(coloredWords, fmt.Sprint(word))
		}
	}
	return strings.Join(coloredWords, " ")
}

func (f *SimpleScreenPrinter) printSubject(splitted []string) string {
	coloredWords := []string{}
	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprint(word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprint(word))
		} else {
			coloredWords = append(coloredWords, fmt.Sprint(word))
		}
	}
	return strings.Join(coloredWords, " ")
}
