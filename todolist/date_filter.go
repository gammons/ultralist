package todolist

import (
	"regexp"
	"time"

	"github.com/jinzhu/now"
)

type DateFilter struct {
	Todos    []*Todo
	Location *time.Location
}

func NewDateFilter(todos []*Todo) *DateFilter {
	return &DateFilter{Todos: todos, Location: time.Now().Location()}
}

func (f *DateFilter) FilterDate(input string) []*Todo {
	agendaRegex, _ := regexp.Compile(`agenda .*$`)
	if agendaRegex.MatchString(input) {
		return f.filterAgenda(now.BeginningOfDay())
	}

	r, _ := regexp.Compile(`due .*$`)
	match := r.FindString(input)
	switch {
	case match == "due tod" || match == "due today":
		return f.filterToday(now.BeginningOfDay())
	case match == "due tom" || match == "due tomorrow":
		return f.filterTomorrow(now.BeginningOfDay())
	case match == "due sun" || match == "due sunday":
		return f.filterDay(now.BeginningOfDay(), time.Sunday)
	case match == "due mon" || match == "due monday":
		return f.filterDay(now.BeginningOfDay(), time.Monday)
	case match == "due tue" || match == "due tuesday":
		return f.filterDay(now.BeginningOfDay(), time.Tuesday)
	case match == "due wed" || match == "due wednesday":
		return f.filterDay(now.BeginningOfDay(), time.Wednesday)
	case match == "due thu" || match == "due thursday":
		return f.filterDay(now.BeginningOfDay(), time.Thursday)
	case match == "due fri" || match == "due friday":
		return f.filterDay(now.BeginningOfDay(), time.Friday)
	case match == "due sat" || match == "due saturday":
		return f.filterDay(now.BeginningOfDay(), time.Saturday)
	case match == "due this week":
		return f.filterThisWeek(now.BeginningOfDay())
	case match == "due next week":
		return f.filterNextWeek(now.BeginningOfDay())
	case match == "overdue":
		return f.filterOverdue(now.BeginningOfDay())
	}
	return f.Todos
}

func (f *DateFilter) filterAgenda(pivot time.Time) []*Todo {
	var ret []*Todo

	for _, todo := range f.Todos {
		dueTime, _ := time.ParseInLocation("2006-01-02", todo.Due, f.Location)
		if dueTime.Before(pivot) || todo.Due == pivot.Format("2006-01-02") {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *DateFilter) filterToday(pivot time.Time) []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if todo.Due == pivot.Format("2006-01-02") {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *DateFilter) filterDay(pivot time.Time, day time.Weekday) []*Todo {
	var ret []*Todo
	filtered := f.filterThisWeek(pivot)
	for _, todo := range filtered {
		dueTime, _ := time.ParseInLocation("2006-01-02", todo.Due, f.Location)
		if dueTime.Weekday() == day {
			ret = append(ret, todo)
		}

	}
	return ret
}

func (f *DateFilter) filterTomorrow(pivot time.Time) []*Todo {
	var ret []*Todo
	pivot = pivot.AddDate(0, 0, 1)
	for _, todo := range f.Todos {
		if todo.Due == pivot.Format("2006-01-02") {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *DateFilter) filterThisWeek(pivot time.Time) []*Todo {
	var ret []*Todo

	begin := now.New(f.findSunday(pivot)).BeginningOfDay()
	end := begin.AddDate(0, 0, 7)

	for _, todo := range f.Todos {
		dueTime, _ := time.ParseInLocation("2006-01-02", todo.Due, f.Location)
		if (begin.Before(dueTime) || begin.Equal(dueTime)) && end.After(dueTime) {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *DateFilter) filterNextWeek(pivot time.Time) []*Todo {
	var ret []*Todo

	begin := f.findSunday(pivot).AddDate(0, 0, 7)
	end := begin.AddDate(0, 0, 7)

	for _, todo := range f.Todos {
		dueTime, _ := time.ParseInLocation("2006-01-02", todo.Due, f.Location)
		if begin.Before(dueTime) && end.After(dueTime) {
			ret = append(ret, todo)
		}
	}
	return ret
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

func (f *DateFilter) findSunday(pivot time.Time) time.Time {
	switch now.New(pivot).Weekday() {
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
