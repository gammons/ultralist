package ultralist

import (
	"fmt"
	"strings"
	"time"
)

type ViewPrinter struct{}

func (v *ViewPrinter) FormatID(todo *Todo) string {
	return fmt.Sprintf("[%s]%-3d[%s]", ColorYellow, todo.ID, ColorForeground)
}

func (v *ViewPrinter) FormatSubject(todo *Todo) string {
	splitted := strings.Split(todo.Subject, " ")
	coloredWords := []string{}

	for _, word := range splitted {
		if contextRegex.MatchString(word) {
			coloredWords = append(
				coloredWords,
				fmt.Sprintf("[%s]%s[%s]", ColorRed, word, ColorForeground),
			)
		} else if projectRegex.MatchString(word) {
			coloredWords = append(
				coloredWords,
				fmt.Sprintf("[%s]%s[%s]", ColorPurple, word, ColorForeground),
			)
		} else {
			coloredWords = append(
				coloredWords,
				fmt.Sprintf("[%s]%s[%s]", ColorForeground, word, ColorForeground),
			)
		}
	}

	return strings.Join(coloredWords, " ")
}

func (v *ViewPrinter) FormatCompleted(todo *Todo) string {
	if todo.Completed {
		return fmt.Sprintf("[%s][✔][%s]", ColorForeground, ColorForeground)
	}
	return fmt.Sprintf("[%s][ ][%s]", ColorForeground, ColorForeground)
}

func (v *ViewPrinter) FormatDue(todo *Todo) string {
	due, _ := time.Parse(todo.Due, DATE_FORMAT)
	if todo.DueToday() {
		return fmt.Sprintf("[%s]%-10s[%s]", ColorBlue, "today", ColorForeground)
	} else if todo.DueTomorrow() {
		return fmt.Sprintf("[%s]%-10s[%s]", ColorBlue, "tomorrow", ColorForeground)
	} else if todo.Due != "" && todo.PastDue() {
		return fmt.Sprintf("[%s]%-10s[%s]", ColorRed, due.Format("Mon Jan 02"), ColorForeground)
	}

	if todo.Due == "" {
		return ""
	} else {
		return fmt.Sprintf("[%s]%-10s[%s]", ColorBlue, due.Format("Mon Jan 02"), ColorForeground)
	}
}

func (v *ViewPrinter) FormatStatus(todo *Todo) string {
	return fmt.Sprintf("[%s]%-10s[%s]", ColorGreen, todo.Status, ColorForeground)
}
