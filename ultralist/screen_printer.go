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
	blue    = color.New(color.FgBlue)
	cyan    = color.New(color.FgCyan)
	magenta = color.New(color.FgMagenta)
	nocolor = color.New()
	red     = color.New(color.FgRed)
	white   = color.New(color.FgWhite)
	yellow  = color.New(color.FgYellow)
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
	if todo.IsPriority {
		yellow.Add(color.Bold, color.Italic)
	}
	fmt.Fprintf(f.Writer, " %s\t%s\t%s\t%s\t\n",
		yellow.Sprint(strconv.Itoa(todo.ID)),
		f.formatCompleted(todo.Completed),
		f.formatDue(todo.Due, todo.IsPriority, todo.Completed),
		f.formatSubject(todo.Subject, todo.IsPriority))
}

func (f *ScreenPrinter) formatDue(due string, isPriority bool, completed bool) string {
	if isPriority {
		blue.Add(color.Bold, color.Italic)
		red.Add(color.Bold, color.Italic)
	}

	if due == "" {
		return nocolor.Sprint("          ")
	}
	dueTime, err := time.Parse("2006-01-02", due)

	if err != nil {
		fmt.Println(err)
		fmt.Println("This may due to the corruption of .todos.json file.")
		os.Exit(-1)
	}

	if isToday(dueTime) {
		return blue.Sprint("today     ")
	} else if isTomorrow(dueTime) {
		return blue.Sprint("tomorrow  ")
	} else if isPastDue(dueTime) && !completed {
		return red.Sprint(dueTime.Format("Mon Jan 02"))
	}
	return blue.Sprint(dueTime.Format("Mon Jan 02"))
}

func (f *ScreenPrinter) formatSubject(subject string, isPriority bool) string {
	if isPriority {
		red.Add(color.Bold, color.Italic)
		magenta.Add(color.Bold, color.Italic)
		white.Add(color.Bold, color.Italic)
	}

	splitted := strings.Split(subject, " ")
	projectRegex, _ := regexp.Compile(`\+[\p{L}\d_]+`)
	contextRegex, _ := regexp.Compile(`\@[\p{L}\d_]+`)

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

func (f *ScreenPrinter) formatCompleted(completed bool) string {
	if completed {
		return "[x]"
	}
	return "[ ]"
}
