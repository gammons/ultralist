package ultralist

// Filter holds the parsed filtering results from an input string
type Filter struct {
	Subject    string
	Archived   bool
	IsPriority bool
	Completed  bool

	Due       string
	DueBefore string
	DueAfter  string

	Contexts []string
	Projects []string
	Status   []string

	ExcludeContexts []string
	ExcludeProjects []string
	ExcludeStatus   []string

	CompletedAt []string

	HasCompleted   bool
	HasCompletedAt bool
	HasArchived    bool
	HasIsPriority  bool

	HasDueBefore bool
	HasDue       bool
	HasDueAfter  bool

	HasStatus        bool
	HasProjectFilter bool
	HasContextFilter bool

	HasRecur   bool
	Recur      string
	RecurUntil string
}

// LastStatus returns the last status from the filter
func (f *Filter) LastStatus() string {
	if len(f.Status) == 0 {
		return ""
	}
	return f.Status[len(f.Status)-1]
}
