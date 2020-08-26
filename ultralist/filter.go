package ultralist

// Filter holds the parsed filtering results from an input string
type Filter struct {
	Subject    string
	Archived   bool
	IsPriority bool
	Completed  bool

	HasCompleted     bool
	HasCompletedAt   bool
	HasArchived      bool
	HasIsPriority    bool
	HasDue           bool
	HasStatus        bool
	HasProjectFilter bool
	HasContextFilter bool

	Contexts []string
	Due      []string
	Projects []string
	Status   []string

	ExcludeContexts []string
	ExcludeDue      []string
	ExcludeProjects []string
	ExcludeStatus   []string

	CompletedAt []string
}

// LastStatus returns the last status from the filter
func (f *Filter) LastStatus() string {
	if len(f.Status) == 0 {
		return ""
	}
	return f.Status[len(f.Status)-1]
}

// LastDue returns the last due from the filter
func (f *Filter) LastDue() string {
	if len(f.Due) == 0 {
		return ""
	}
	return f.Due[len(f.Due)-1]
}
