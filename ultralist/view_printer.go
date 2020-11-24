package ultralist

import (
	"fmt"
	"strings"
	"time"
)

type ViewPrinter struct{}

func (v *ViewPrinter) FormatID(todo *Todo) string {
	return fmt.Sprintf("[%s::%s]%-3d[%s::-]", ColorYellow, v.bold(todo), todo.ID, ColorForeground)
}

func (v *ViewPrinter) FormatSubject(todo *Todo) string {
	splitted := strings.Split(todo.Subject, " ")
	coloredWords := []string{}

	for _, word := range splitted {
		if contextRegex.MatchString(word) {
			coloredWords = append(
				coloredWords,
				fmt.Sprintf("[%s::%s]%s[%s::-]", ColorRed, v.bold(todo), word, ColorForeground),
			)
		} else if projectRegex.MatchString(word) {
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

func (v *ViewPrinter) FormatCompleted(todo *Todo) string {
	if todo.Completed {
		return fmt.Sprintf("[%s::%s][✔][%s::-]", ColorForeground, v.bold(todo), ColorForeground)
	}
	return fmt.Sprintf("[%s::%s][ ][%s::-]", ColorForeground, v.bold(todo), ColorForeground)
}

func (v *ViewPrinter) FormatDue(todo *Todo) string {
	due, _ := time.Parse(todo.Due, DATE_FORMAT)
	if todo.DueToday() {
		return fmt.Sprintf("[%s::%s]%-10s[%s::-]", ColorBlue, v.bold(todo), "today", ColorForeground)
	} else if todo.DueTomorrow() {
		return fmt.Sprintf("[%s::%s]%-10s[%s::-]", ColorBlue, v.bold(todo), "tomorrow", ColorForeground)
	} else if todo.Due != "" && todo.PastDue() {
		return fmt.Sprintf("[%s::%s]%-10s[%s::-]", ColorRed, v.bold(todo), due.Format("Mon Jan 02"), ColorForeground)
	}

	if todo.Due == "" {
		return fmt.Sprintf("%-10s", "")
	} else {
		return fmt.Sprintf("[%s::%s]%-10s[%s::-]", ColorBlue, v.bold(todo), due.Format("Mon Jan 02"), ColorForeground)
	}
}

func (v *ViewPrinter) FormatStatus(todo *Todo) string {
	return fmt.Sprintf("[%s::%s]%-10s[%s::-]", ColorGreen, v.bold(todo), todo.Status, ColorForeground)
}

func (v *ViewPrinter) bold(todo *Todo) string {
	if todo.IsPriority {
		return "b"
	}
	return "-"
}
