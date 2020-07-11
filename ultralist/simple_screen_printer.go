package ultralist

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
)

// ScreenPrinter is the default struct of this file
type SimpleScreenPrinter struct {
	Writer         *io.Writer
	UnicodeSupport bool
}

// NewScreenPrinter creates a new screeen printer.
func NewSimpleScreenPrinter(unicodeSupport bool) *SimpleScreenPrinter {
	w := new(io.Writer)
	formatter := &SimpleScreenPrinter{Writer: w, UnicodeSupport: unicodeSupport}
	return formatter
}

// Print prints the output of ultralist to the terminal screen.
func (f *SimpleScreenPrinter) Print(groupedTodos *GroupedTodos, printNotes bool) {
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
			f.printTodo(tabby, todo, printNotes)
		}
		tabby.AddLine()
	}
	tabby.Print()
}

func (f *SimpleScreenPrinter) printTodo(tabby *tabby.Tabby, todo *Todo, printNotes bool) {
	tabby.AddLine(
		f.formatID(todo.ID, todo.IsPriority),
		f.formatCompleted(todo.Completed),
		f.formatInformation(todo),
		f.formatDue(todo.Due, todo.IsPriority, todo.Completed),
		f.formatSubject(todo.Subject, todo.IsPriority))
	if printNotes {
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

func (f *SimpleScreenPrinter) formatDue(due string, isPriority bool, completed bool) string {
	if due == "" {
		return fmt.Sprint("          ")
	}
	dueTime, _ := time.Parse("2006-01-02", due)

	if isPriority {
		return f.printPriorityDue(dueTime, completed)
	}
	return f.printDue(dueTime, completed)
}

func (f *SimpleScreenPrinter) formatInformation(todo *Todo) string {
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
	return fmt.Sprint(strings.Join(information, " "))
}

func (f *SimpleScreenPrinter) printDue(due time.Time, completed bool) string {
	if isToday(due) {
		return fmt.Sprint("today     ")
	} else if isTomorrow(due) {
		return fmt.Sprint("tomorrow  ")
	} else if isPastDue(due) && !completed {
		return fmt.Sprint(due.Format("Mon Jan 02"))
	}
	return fmt.Sprint(due.Format("Mon Jan 02"))
}

func (f *SimpleScreenPrinter) printPriorityDue(due time.Time, completed bool) string {
	if isToday(due) {
		return fmt.Sprint("today     ")
	} else if isTomorrow(due) {
		return fmt.Sprint("tomorrow  ")
	} else if isPastDue(due) && !completed {
		return fmt.Sprint(due.Format("Mon Jan 02"))
	}
	return fmt.Sprint(due.Format("Mon Jan 02"))
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
