package ultralist

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

var (
	blue            = color.New(color.FgBlue)
	blueBold        = color.New(color.FgBlue, color.Bold)
	cyan            = color.New(color.FgCyan)
	magenta         = color.New(color.FgMagenta)
	magentaBold     = color.New(color.FgMagenta, color.Bold)
	nocolor         = color.New()
	red             = color.New(color.FgRed)
	redBold         = color.New(color.FgRed, color.Bold)
	white           = color.New(color.FgWhite)
	whiteBold       = color.New(color.FgWhite, color.Bold)
	yellow          = color.New(color.FgYellow)
	yellowBold      = color.New(color.FgYellow, color.Bold)
	projectRegex, _ = regexp.Compile(`\+[\p{L}\d_]+`)
	contextRegex, _ = regexp.Compile(`\@[\p{L}\d_]+`)
)

// ScreenPrinter is the default struct of this file
type ScreenPrinter struct {
	Writer *tabwriter.Writer
}

// NewScreenPrinter creates a new screeen printer.
func NewScreenPrinter() *ScreenPrinter {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	formatter := &ScreenPrinter{Writer: w}
	return formatter
}

// Print prints the output of ultralist to the terminal screen.
func (f *ScreenPrinter) Print(groupedTodos *GroupedTodos, printNotes bool) {
	var keys []string
	for key := range groupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Fprintf(f.Writer, "\n %s\n", cyan.Sprint(key))
		for _, todo := range groupedTodos.Groups[key] {
			f.printTodo(todo)
			if printNotes {
				for nid, note := range todo.Notes {
					fmt.Fprintf(f.Writer, "   %s\t\t\t%s\t\n",
						cyan.Sprint(strconv.Itoa(nid)), white.Sprint(note))
				}
			}
		}
	}
	f.Writer.Flush()
}

func (f *ScreenPrinter) printTodo(todo *Todo) {
	fmt.Fprintf(f.Writer, " %s\t%s\t%s\t%s\t\n",
		f.formatID(todo.ID, todo.IsPriority),
		f.formatCompleted(todo.Completed),
		f.formatDue(todo.Due, todo.IsPriority, todo.Completed),
		f.formatSubject(todo.Subject, todo.IsPriority))
}

func (f *ScreenPrinter) formatID(ID int, isPriority bool) string {
	if isPriority {
		return yellowBold.Sprint(strconv.Itoa(ID))
	}
	return yellow.Sprint(strconv.Itoa(ID))
}

func (f *ScreenPrinter) formatCompleted(completed bool) string {
	if completed {
		return "[x]"
	}
	return "[ ]"
}

func (f *ScreenPrinter) formatDue(due string, isPriority bool, completed bool) string {
	if due == "" {
		return nocolor.Sprint("          ")
	}
	dueTime, _ := time.Parse("2006-01-02", due)

	if isPriority {
		return f.printPriorityDue(dueTime, completed)
	}
	return f.printDue(dueTime, completed)

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
