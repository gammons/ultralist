package ultralist

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

type PrettyScreenPrinter struct{}

func (p *PrettyScreenPrinter) Print(groupedTodos *GroupedTodos, printNotes bool, showStatus bool) {
	var keys []string
	for key := range groupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)

	t.AppendHeader(table.Row{"ID", "Comp", "Due", "Status", "Subject"})

	for idx, key := range keys {
		t.AppendRow(table.Row{key})

		for _, todo := range groupedTodos.Groups[key] {
			t.AppendRow(
				table.Row{
					p.formatID(todo.ID, todo.IsPriority),
					p.formatCompleted(todo.Completed),
					p.formatDue(todo.Due, todo.IsPriority, todo.Completed),
					p.formatStatus(todo.Status, todo.IsPriority),
					p.formatSubject(todo.Subject, todo.IsPriority),
				},
			)
		}

		if idx < len(keys) {
			t.AppendSeparator()
		}
	}
	fmt.Println(t.Render() + "\n")
}

func (f *PrettyScreenPrinter) formatID(ID int, isPriority bool) string {
	if isPriority {
		return yellowBold.Sprint(strconv.Itoa(ID))
	}
	return yellow.Sprint(strconv.Itoa(ID))
}

func (f *PrettyScreenPrinter) formatCompleted(completed bool) string {
	if completed {
		return white.Sprint("[✔]")
	}
	return white.Sprint("[ ]")
}

func (f *PrettyScreenPrinter) formatDue(due string, isPriority bool, completed bool) string {
	if due == "" {
		return white.Sprint("          ")
	}
	dueTime, _ := time.Parse(DATE_FORMAT, due)

	if isPriority {
		return f.printPriorityDue(dueTime, completed)
	}
	return f.printDue(dueTime, completed)
}

func (f *PrettyScreenPrinter) formatStatus(status string, isPriority bool) string {
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

func (f *PrettyScreenPrinter) formatInformation(todo *Todo) string {
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

func (f *PrettyScreenPrinter) printDue(due time.Time, completed bool) string {
	if isToday(due) {
		return blue.Sprint("today     ")
	} else if isTomorrow(due) {
		return blue.Sprint("tomorrow  ")
	} else if isPastDue(due) && !completed {
		return red.Sprint(due.Format("Mon Jan 02"))
	}
	return blue.Sprint(due.Format("Mon Jan 02"))
}

func (f *PrettyScreenPrinter) printPriorityDue(due time.Time, completed bool) string {
	if isToday(due) {
		return blueBold.Sprint("today     ")
	} else if isTomorrow(due) {
		return blueBold.Sprint("tomorrow  ")
	} else if isPastDue(due) && !completed {
		return redBold.Sprint(due.Format("Mon Jan 02"))
	}
	return blueBold.Sprint(due.Format("Mon Jan 02"))
}

func (f *PrettyScreenPrinter) formatSubject(subject string, isPriority bool) string {
	splitted := strings.Split(subject, " ")

	if isPriority {
		return f.printPrioritySubject(splitted)
	}
	return f.printSubject(splitted)
}

func (f *PrettyScreenPrinter) printPrioritySubject(splitted []string) string {
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

func (f *PrettyScreenPrinter) printSubject(splitted []string) string {
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
