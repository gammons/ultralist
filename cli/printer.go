package cli

import (
	"regexp"

	"github.com/ultralist/ultralist/ultralist"
)

var (
	projectRegex, _ = regexp.Compile(ultralist.ProjectRegexp)
	contextRegex, _ = regexp.Compile(ultralist.ContextRegexp)
)

// Printer is an interface for printing grouped todos.
type Printer interface {
	Print(*ultralist.GroupedTodos)
}
