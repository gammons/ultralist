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
type TviewPrinter struct {
	Writer         *io.Writer
	UnicodeSupport bool
}

// NewScreenPrinter creates a new screeen printer.
func NewTviewPrinter(unicodeSupport bool) *TviewPrinter {
	w := new(io.Writer)
	formatter := &TviewPrinter{Writer: w, UnicodeSupport: unicodeSupport}
	return formatter
}

// Print prints the output of ultralist to the terminal screen.
func (f *TviewPrinter) Print(groupedTodos *GroupedTodos, printNotes bool) {
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
			f.printTodo(tabby, todo, printNotes)
		}
		tabby.AddLine()
	}
	tabby.Print()
}

func (f *TviewPrinter) printTodo(tabby *tabby.Tabby, todo *Todo, printNotes bool) {
	tabby.AddLine(
		f.FormatID(todo.ID, todo.IsPriority),
		f.FormatCompleted(todo.Completed),
		f.FormatDue(todo.Due, todo.IsPriority, todo.Completed),
		f.FormatSubject(todo.Subject, todo.IsPriority))
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

func (f *TviewPrinter) FormatID(ID int, isPriority bool) string {
	if isPriority {
		return fmt.Sprintf("[yellow]%s[white]", strconv.Itoa(ID))
	}
	return fmt.Sprintf("[yellow]%s[white]", strconv.Itoa(ID))
}

func (f *TviewPrinter) FormatCompleted(completed bool) string {
	if completed {
		if f.UnicodeSupport {
			return fmt.Sprint("[white][✔]")
		} else {
			return fmt.Sprint("[white][x]")
		}
	}
	return fmt.Sprint("[white][ ]")
}

func (f *TviewPrinter) FormatDue(due string, isPriority bool, completed bool) string {
	if due == "" {
		return fmt.Sprint("          ")
	}
	dueTime, _ := time.Parse("2006-01-02", due)

	if isPriority {
		return f.printPriorityDue(dueTime, completed)
	}
	return f.printDue(dueTime, completed)

}

func (f *TviewPrinter) printDue(due time.Time, completed bool) string {
	if isToday(due) {
		return fmt.Sprint("[blue]today[white]     ")
	} else if isTomorrow(due) {
		return fmt.Sprint("[blue]tomorrow[white]  ")
	} else if isPastDue(due) && !completed {
		return fmt.Sprint(due.Format("[red]Mon Jan 02[white]"))
	}
	return fmt.Sprint(due.Format("[blue]Mon Jan 02[white]"))
}

func (f *TviewPrinter) printPriorityDue(due time.Time, completed bool) string {
	if isToday(due) {
		return fmt.Sprint("[blue]today[white]     ")
	} else if isTomorrow(due) {
		return fmt.Sprint("[blue]tomorrow[white]  ")
	} else if isPastDue(due) && !completed {
		return fmt.Sprint(due.Format("[red]Mon Jan 02[white]"))
	}
	return blueBold.Sprint(due.Format("[red]Mon Jan 02[white]"))
}

func (f *TviewPrinter) FormatSubject(subject string, isPriority bool) string {
	splitted := strings.Split(subject, " ")

	if isPriority {
		return f.printPrioritySubject(splitted)
	}
	return f.printSubject(splitted)
}

func (f *TviewPrinter) printPrioritySubject(splitted []string) string {
	coloredWords := []string{}
	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprintf("[purple]%s[white]", word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprintf("[red]%s[white]", word))
		} else {
			coloredWords = append(coloredWords, fmt.Sprint(word))
		}
	}
	return strings.Join(coloredWords, " ")
}

func (f *TviewPrinter) printSubject(splitted []string) string {
	coloredWords := []string{}
	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprintf("[purple]%s[white]", word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprintf("[red]%s[white]", word))
		} else {
			coloredWords = append(coloredWords, fmt.Sprintf(word))
		}
	}
	fmt.Println("coloredWords = ", strings.Join(coloredWords, " "))
	return strings.Join(coloredWords, " ")
}
