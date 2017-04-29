package todolist

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

type Formatter interface {
	Print()
}

type ColorFormatter struct {
	GroupedTodos *GroupedTodos
	Writer       io.Writer
}

type WriterFormatter struct {
	ColorFormatter
}

func NewWriterFormatter(w io.Writer, todos *GroupedTodos) Formatter {
	colorFormatter := ColorFormatter{GroupedTodos: todos, Writer: w}
	formatter := WriterFormatter{colorFormatter}
	return Formatter(formatter)
}

func (f WriterFormatter) Print() {
	f.print()
}

type TabFormatter struct {
	ColorFormatter
	Writer *tabwriter.Writer
}

func NewFormatter(todos *GroupedTodos) Formatter {
	return NewTabFormatter(todos)
}

func NewTabFormatter(todos *GroupedTodos) Formatter {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	colorFormatter := ColorFormatter{GroupedTodos: todos, Writer: w}
	formatter := TabFormatter{colorFormatter, w}
	return Formatter(formatter)
}

func (f TabFormatter) Print() {
	f.print()
	f.Writer.Flush()
}

func (f ColorFormatter) print() {
	cyan := color.New(color.FgCyan).SprintFunc()

	var keys []string
	for key := range f.GroupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Fprintf(f.Writer, "\n %s\n", cyan(key))
		for _, todo := range f.GroupedTodos.Groups[key] {
			f.printTodo(todo)
		}
	}
}

func (f ColorFormatter) printTodo(todo *Todo) {
	yellow := color.New(color.FgYellow)
	if todo.IsPriority {
		yellow.Add(color.Bold, color.Italic)
	}
	fmt.Fprintf(f.Writer, " %s\t%s\t%s\t%s\t\n",
		yellow.SprintFunc()(strconv.Itoa(todo.Id)),
		f.formatCompleted(todo.Completed),
		f.formatDue(todo.Due, todo.IsPriority),
		f.formatSubject(todo.Subject, todo.IsPriority))
}

func (f ColorFormatter) formatDue(due string, isPriority bool) string {
	blue := color.New(color.FgBlue)
	red := color.New(color.FgRed)

	if isPriority {
		blue.Add(color.Bold, color.Italic)
		red.Add(color.Bold, color.Italic)
	}

	if due == "" {
		return blue.SprintFunc()(" ")
	}
	dueTime, err := time.Parse("2006-01-02", due)

	if err != nil {
		fmt.Println(err)
		fmt.Println("This may due to the corruption of .todos.json file.")
		os.Exit(-1)
	}

	if isToday(dueTime) {
		return blue.SprintFunc()("today")
	} else if isTomorrow(dueTime) {
		return blue.SprintFunc()("tomorrow")
	} else if isPastDue(dueTime) {
		return red.SprintFunc()(dueTime.Format("Mon Jan 2"))
	} else {
		return blue.SprintFunc()(dueTime.Format("Mon Jan 2"))
	}
}

func isToday(t time.Time) bool {
	nowYear, nowMonth, nowDay := time.Now().Date()
	timeYear, timeMonth, timeDay := t.Date()
	return nowYear == timeYear &&
		nowMonth == timeMonth &&
		nowDay == timeDay
}

func isTomorrow(t time.Time) bool {
	nowYear, nowMonth, nowDay := time.Now().AddDate(0, 0, 1).Date()
	timeYear, timeMonth, timeDay := t.Date()
	return nowYear == timeYear &&
		nowMonth == timeMonth &&
		nowDay == timeDay
}

func isPastDue(t time.Time) bool {
	return time.Now().After(t)
}

func (f ColorFormatter) formatSubject(subject string, isPriority bool) string {

	red := color.New(color.FgRed)
	magenta := color.New(color.FgMagenta)
	white := color.New(color.FgWhite)

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
			coloredWords = append(coloredWords, magenta.SprintFunc()(word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, red.SprintFunc()(word))
		} else {
			coloredWords = append(coloredWords, white.SprintFunc()(word))
		}
	}
	return strings.Join(coloredWords, " ")

}

func (f ColorFormatter) formatCompleted(completed bool) string {
	if completed {
		return "[x]"
	} else {
		return "[ ]"
	}
}
