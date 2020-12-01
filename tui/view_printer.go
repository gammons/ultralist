package tui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ultralist/ultralist/ultralist"
)

type ViewPrinter struct{}

func (v *ViewPrinter) FormatID(todo *ultralist.Todo) string {
	return fmt.Sprintf("[%s::%s]%-3d[%s::-]", ColorYellow, v.bold(todo), todo.ID, ColorForeground)
}

func (v *ViewPrinter) FormatSubject(todo *ultralist.Todo) string {
	if todo.Archived {
		return fmt.Sprintf("[%s::%s]%s[%s::-]", ColorGray, v.bold(todo), todo.Subject, ColorForeground)
	}

	splitted := strings.Split(todo.Subject, " ")
	coloredWords := []string{}

	for _, word := range splitted {
		if match, _ := regexp.MatchString(ultralist.ContextRegexp, word); match == true {
			coloredWords = append(
				coloredWords,
				fmt.Sprintf("[%s::%s]%s[%s::-]", ColorRed, v.bold(todo), word, ColorForeground),
			)
		} else if match, _ := regexp.MatchString(ultralist.ProjectRegexp, word); match == true {
			coloredWords = append(
				coloredWords,
				fmt.Sprintf("[%s::%s]%s[%s::-]", ColorPurple, v.bold(todo), word, ColorForeground),
			)
		} else {
			coloredWords = append(
				coloredWords,
				fmt.Sprintf("[%s::%s]%s[%s::-]", ColorForeground, v.bold(todo), word, ColorForeground),
			)
		}
	}

	return strings.Join(coloredWords, " ")
}

func (v *ViewPrinter) FormatCompleted(todo *ultralist.Todo) string {
	if todo.Completed {
		return fmt.Sprintf("[%s::%s][✔][%s::-]", ColorForeground, v.bold(todo), ColorForeground)
	}
	return fmt.Sprintf("[%s::%s][ ][%s::-]", ColorForeground, v.bold(todo), ColorForeground)
}

func (v *ViewPrinter) FormatDue(todo *ultralist.Todo) string {
	due, _ := time.Parse(todo.Due, ultralist.DateFormat)

	if todo.DueToday() {
		return fmt.Sprintf(
			"[%s::%s]%-10s[%s::-]",
			v.color(todo, ColorBlue),
			v.bold(todo),
			"today",
			ColorForeground,
		)
	} else if todo.DueTomorrow() {
		return fmt.Sprintf(
			"[%s::%s]%-10s[%s::-]",
			v.color(todo, ColorBlue),
			v.bold(todo),
			"tomorrow",
			ColorForeground,
		)
	} else if todo.Due != "" && todo.PastDue() {
		return fmt.Sprintf(
			"[%s::%s]%-10s[%s::-]",
			v.color(todo, ColorRed),
			v.bold(todo),
			due.Format("Mon Jan 02"),
			ColorForeground,
		)
	}

	if todo.Due == "" {
		return fmt.Sprintf("%-10s", "")
	} else {
		return fmt.Sprintf(
			"[%s::%s]%-10s[%s::-]",
			v.color(todo, ColorBlue),
			v.bold(todo),
			due.Format("Mon Jan 02"),
			ColorForeground,
		)
	}
}

func (v *ViewPrinter) FormatStatus(todo *ultralist.Todo) string {
	return fmt.Sprintf(
		"[%s::%s]%-10s[%s::-]",
		v.color(todo, ColorGreen),
		v.bold(todo),
		todo.Status,
		ColorForeground,
	)
}

func (v *ViewPrinter) bold(todo *ultralist.Todo) string {
	if todo.IsPriority {
		return "b"
	}
	return "-"
}

func (v *ViewPrinter) color(todo *ultralist.Todo, defaultColor string) string {
	if todo.Archived {
		return ColorGray
	}
	return defaultColor
}
