package cli

import (
	"regexp"

	"github.com/ultralist/ultralist/ultralist"
)

var (
	projectRegex, _ = regexp.Compile(`\+[\p{L}\d_]+`)
	contextRegex, _ = regexp.Compile(`\@[\p{L}\d_]+`)
)

// Printer is an interface for printing grouped todos.
type Printer interface {
	Print(*ultralist.GroupedTodos)
}
