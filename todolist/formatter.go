package todolist

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

type Formatter struct {
	Todos  []Todo
	Writer *tabwriter.Writer
}

func NewFormatter(todos []Todo) *Formatter {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	formatter := &Formatter{Todos: todos, Writer: w}
	return formatter
}

func (f *Formatter) Print() {
	for _, todo := range f.Todos {
		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Fprintf(f.Writer, "     \t%s\t%s\t%s\t\n",
			yellow(strconv.Itoa(todo.Id)),
			f.formatCompleted(todo.Completed),
			f.formatSubject(todo.Subject))
	}
	f.Writer.Flush()
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
