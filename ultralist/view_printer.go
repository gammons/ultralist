package ultralist

import (
	"fmt"
	"strings"
	"time"
)

type ViewPrinter struct{}

func (v *ViewPrinter) FormatID(todo *Todo) string {
	return fmt.Sprintf("[#F4BF75]%-3d[#d0d0d0]", todo.ID)
}

func (v *ViewPrinter) FormatSubject(todo *Todo) string {
	splitted := strings.Split(todo.Subject, " ")
	coloredWords := []string{}

	for _, word := range splitted {
		if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprintf("[#AC4142]%s[#d0d0d0]", word))
		} else if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, fmt.Sprintf("[#AA759F]%s[#d0d0d0]", word))
		} else {
			coloredWords = append(coloredWords, fmt.Sprintf("[#d0d0d0]%s[#d0d0d0]", word))
		}
	}

	return strings.Join(coloredWords, " ")
}

func (v *ViewPrinter) FormatCompleted(todo *Todo) string {
	if todo.Completed {
		return "[#d0d0d0][✔][#d0d0d0]"
	}
	return "[#d0d0d0][ ][#d0d0d0]"
}

func (v *ViewPrinter) FormatDue(todo *Todo) string {
	due, _ := time.Parse(todo.Due, DATE_FORMAT)
	if todo.DueToday() {
		return fmt.Sprintf("[#6A9FB5]%-10s[#d0d0d0]", "today")
	} else if todo.DueTomorrow() {
		return fmt.Sprintf("[#6A9FB5]%-10s[#d0d0d0]", "tomorrow")
	} else if todo.Due != "" && todo.PastDue() {
		return fmt.Sprintf("[#AC4142]%-10s[#d0d0d0]", due.Format("Mon Jan 02"))
	}

	if todo.Due == "" {
		return fmt.Sprintf("[#6A9FB5]%-10s[#d0d0d0]", "")
	} else {
		return fmt.Sprintf("[#6A9FB5]%-10s[#d0d0d0]", due.Format("Mon Jan 02"))
	}
}

func (v *ViewPrinter) FormatStatus(todo *Todo) string {
	return fmt.Sprintf("[#90A959]%-10s[#d0d0d0]", todo.Status)
}
