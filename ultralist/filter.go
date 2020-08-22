package ultralist

// Filter holds the parsed filtering results from an input string
type Filter struct {
	Subject    string
	Archived   bool
	IsPriority bool
	Completed  bool

	HasCompleted   bool
	HasCompletedAt bool
	HasArchived    bool
	HasIsPriority  bool
	HasDue         bool
	HasStatus      bool

	Projects    []string
	Contexts    []string
	Due         []string
	CompletedAt []string
	Status      []string
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
