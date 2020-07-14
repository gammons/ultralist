package ultralist

import (
	"regexp"

	"github.com/fatih/color"
)

var (
	blue        = color.New(0, color.FgBlue)
	blueBold    = color.New(color.Bold, color.FgBlue)
	cyan        = color.New(0, color.FgCyan)
	cyanBold    = color.New(color.Bold, color.FgCyan)
	magenta     = color.New(0, color.FgMagenta)
	magentaBold = color.New(color.Bold, color.FgMagenta)
	red         = color.New(0, color.FgRed)
	redBold     = color.New(color.Bold, color.FgRed)
	white       = color.New(0, color.FgWhite)
	whiteBold   = color.New(color.Bold, color.FgWhite)
	yellow      = color.New(0, color.FgYellow)
	yellowBold  = color.New(color.Bold, color.FgYellow)
)

var (
	projectRegex, _ = regexp.Compile(`\+[\p{L}\d_]+`)
	contextRegex, _ = regexp.Compile(`\@[\p{L}\d_]+`)
)

// Printer is an interface for printing grouped todos.
type Printer interface {
	Print(*GroupedTodos, bool)
}
