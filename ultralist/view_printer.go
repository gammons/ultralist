package ultralist

import (
	"fmt"
	"strings"
)

type ViewPrinter struct{}

func (v *ViewPrinter) FormatID(todo *Todo) string {
	return fmt.Sprintf("[#F4BF75]%v[#d0d0d0]", todo.ID)
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
