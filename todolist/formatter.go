package todolist

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

type Formatter struct {
	GroupedTodos *GroupedTodos
	Writer       *tabwriter.Writer
}

func NewFormatter(todos *GroupedTodos) *Formatter {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	formatter := &Formatter{GroupedTodos: todos, Writer: w}
	return formatter
}

func (f *Formatter) Print() {
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
	f.Writer.Flush()
}

func (f *Formatter) printTodo(todo *Todo) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Fprintf(f.Writer, " %s\t%s\t%s\t%s\t\n",
		yellow(strconv.Itoa(todo.Id)),
		f.formatCompleted(todo.Completed),
		f.formatDue(todo.Due),
		f.formatSubject(todo.Subject))
}

func (f *Formatter) formatDue(due string) string {
	blue := color.New(color.FgBlue).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if due == "" {
		return blue(" ")
	}
	dueTime, err := time.Parse("2006-01-02", due)

	if err != nil {
		panic(err)
	}

	if isToday(dueTime) {
		return blue("today")
	} else if isTomorrow(dueTime) {
		return blue("tomorrow")
	} else if isPastDue(dueTime) {
		return red(dueTime.Format("Mon Jan 2"))
	} else {
		return blue(dueTime.Format("Mon Jan 2"))
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

func (f *Formatter) formatSubject(subject string) string {
	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	splitted := strings.Split(subject, " ")
	projectRegex, _ := regexp.Compile(`\+\w+`)
	contextRegex, _ := regexp.Compile(`\@\w+`)

	coloredWords := []string{}

	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, magenta(word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, red(word))
		} else {
			coloredWords = append(coloredWords, word)
		}
	}
	return strings.Join(coloredWords, " ")

}

func (f *Formatter) formatCompleted(completed bool) string {
	if completed {
		return "[x]"
	} else {
		return "[ ]"
	}
}
