package todolist

import (
	"regexp"
	"time"
)

type DateFilter struct {
	Todos    []*Todo
	Location *time.Location
}

func NewDateFilter(todos []*Todo) *DateFilter {
	return &DateFilter{Todos: todos, Location: time.Now().Location()}
}

func filterOnDue(todo *Todo) string {
	return todo.Due
}

func filterOnCompletedDate(todo *Todo) string {
	return todo.CompletedDateToDate()
}

func (f *DateFilter) FilterDate(input string) []*Todo {
	agendaRegex, _ := regexp.Compile(`agenda.*$`)
	if agendaRegex.MatchString(input) {
		return f.filterAgenda(bod(time.Now()))
	}

	overdueRegex, _ := regexp.Compile(`overdue.*$`)
	if overdueRegex.MatchString(input) {
		return f.filterOverdue(bod(time.Now()))
	}

	// filter due items
	r, _ := regexp.Compile(`due .*$`)
	match := r.FindString(input)
	switch {
	case match == "due tod" || match == "due today":
		return f.filterDueToday(bod(time.Now()))
	case match == "due tom" || match == "due tomorrow":
		return f.filterDueTomorrow(bod(time.Now()))
	case match == "due sun" || match == "due sunday":
		return f.filterDueDay(bod(time.Now()), time.Sunday)
	case match == "due mon" || match == "due monday":
		return f.filterDueDay(bod(time.Now()), time.Monday)
	case match == "due tue" || match == "due tuesday":
		return f.filterDueDay(bod(time.Now()), time.Tuesday)
	case match == "due wed" || match == "due wednesday":
		return f.filterDueDay(bod(time.Now()), time.Wednesday)
	case match == "due thu" || match == "due thursday":
		return f.filterDueDay(bod(time.Now()), time.Thursday)
	case match == "due fri" || match == "due friday":
		return f.filterDueDay(bod(time.Now()), time.Friday)
	case match == "due sat" || match == "due saturday":
		return f.filterDueDay(bod(time.Now()), time.Saturday)
	case match == "due this week":
		return f.filterThisWeek(bod(time.Now()))
	case match == "due next week":
		return f.filterNextWeek(bod(time.Now()))
	case match == "due last week":
		return f.filterLastWeek(bod(time.Now()))
	}

	// filter completed items
	r, _ = regexp.Compile(`completed .*$`)
	match = r.FindString(input)
	switch {
	case match == "completed tod" || match == "completed today":
		return f.filterCompletedToday(bod(time.Now()))
	case match == "completed this week":
		return f.filterCompletedThisWeek(bod(time.Now()))
	}

	return f.Todos
}

func (f *DateFilter) filterAgenda(pivot time.Time) []*Todo {
	var ret []*Todo

	for _, todo := range f.Todos {
		if todo.Due == "" || todo.Completed {
			continue
		}
		dueTime, _ := time.ParseInLocation("2006-01-02", todo.Due, f.Location)
		if dueTime.Before(pivot) || todo.Due == pivot.Format("2006-01-02") {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *DateFilter) filterToExactDate(pivot time.Time, filterOn func(*Todo) string) []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if filterOn(todo) == pivot.Format("2006-01-02") {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *DateFilter) filterDueToday(pivot time.Time) []*Todo {
	return f.filterToExactDate(pivot, filterOnDue)
}

func (f *DateFilter) filterDueTomorrow(pivot time.Time) []*Todo {
	pivot = pivot.AddDate(0, 0, 1)
	return f.filterToExactDate(pivot, filterOnDue)
}

func (f *DateFilter) filterCompletedToday(pivot time.Time) []*Todo {
	return f.filterToExactDate(pivot, filterOnCompletedDate)
}

func (f *DateFilter) filterDueDay(pivot time.Time, day time.Weekday) []*Todo {
	var weeklyOffset = 0
	if int(pivot.Weekday()) > int(day) {
		weeklyOffset = 7
	}
	pivot = f.FindSunday(pivot).AddDate(0, 0, int(day)+weeklyOffset)
	return f.filterToExactDate(pivot, filterOnDue)
}

func (f *DateFilter) filterBetweenDatesInclusive(begin, end time.Time, filterOn func(*Todo) string) []*Todo {
	var ret []*Todo

	for _, todo := range f.Todos {
		dueTime, _ := time.ParseInLocation("2006-01-02", filterOn(todo), f.Location)
		if (begin.Before(dueTime) || begin.Equal(dueTime)) && end.After(dueTime) {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *DateFilter) filterThisWeek(pivot time.Time) []*Todo {

	begin := bod(f.FindSunday(pivot))
	end := begin.AddDate(0, 0, 7)

	return f.filterBetweenDatesInclusive(begin, end, filterOnDue)
}

func (f *DateFilter) filterCompletedThisWeek(pivot time.Time) []*Todo {

	begin := bod(f.FindSunday(pivot))
	end := begin.AddDate(0, 0, 7)

	return f.filterBetweenDatesInclusive(begin, end, filterOnCompletedDate)
}

func (f *DateFilter) filterBetweenDatesExclusive(begin, end time.Time) []*Todo {
	var ret []*Todo

	for _, todo := range f.Todos {
		dueTime, _ := time.ParseInLocation("2006-01-02", todo.Due, f.Location)
		if begin.Before(dueTime) && end.After(dueTime) {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *DateFilter) filterNextWeek(pivot time.Time) []*Todo {

	begin := f.FindSunday(pivot).AddDate(0, 0, 7)
	end := begin.AddDate(0, 0, 7)

	return f.filterBetweenDatesExclusive(begin, end)
}

func (f *DateFilter) filterLastWeek(pivot time.Time) []*Todo {

	begin := f.FindSunday(pivot).AddDate(0, 0, -7)
	end := begin.AddDate(0, 0, 7)

	return f.filterBetweenDatesExclusive(begin, end)
}

func (f *DateFilter) filterOverdue(pivot time.Time) []*Todo {
	var ret []*Todo

	pivotDate := pivot.Format("2006-01-02")

	for _, todo := range f.Todos {
		dueTime, _ := time.ParseInLocation("2006-01-02", todo.Due, f.Location)
		if dueTime.Before(pivot) && pivotDate != todo.Due {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *DateFilter) FindSunday(pivot time.Time) time.Time {
	switch pivot.Weekday() {
	case time.Sunday:
		return pivot
	case time.Monday:
		return pivot.AddDate(0, 0, -1)
	case time.Tuesday:
		return pivot.AddDate(0, 0, -2)
	case time.Wednesday:
		return pivot.AddDate(0, 0, -3)
	case time.Thursday:
		return pivot.AddDate(0, 0, -4)
	case time.Friday:
		return pivot.AddDate(0, 0, -5)
	case time.Saturday:
		return pivot.AddDate(0, 0, -6)
	}
	return pivot
}
