package todolist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFilterToday(t *testing.T) {
	assert := assert.New(t)

	var todos []*Todo
	todayTodo := &Todo{Id: 1, Subject: "one", Due: time.Now().Format("2006-01-02")}
	tomorrowTodo := &Todo{Id: 2, Subject: "two", Due: time.Now().AddDate(0, 0, 1).Format("2006-01-02")}
	todos = append(todos, todayTodo)
	todos = append(todos, tomorrowTodo)

	filter := NewDateFilter(todos)
	filtered := filter.filterDueToday(time.Now())

	assert.Equal(1, len(filtered))
	assert.Equal(1, filtered[0].Id)
}

func TestFilterTomorrow(t *testing.T) {
	assert := assert.New(t)

	var todos []*Todo
	todayTodo := &Todo{Id: 1, Subject: "one", Due: time.Now().Format("2006-01-02")}
	tomorrowTodo := &Todo{Id: 2, Subject: "two", Due: time.Now().AddDate(0, 0, 1).Format("2006-01-02")}
	todos = append(todos, todayTodo)
	todos = append(todos, tomorrowTodo)

	filter := NewDateFilter(todos)
	filtered := filter.filterDueTomorrow(time.Now())

	assert.Equal(1, len(filtered))
	assert.Equal(2, filtered[0].Id)
}

func TestFilterCompletedToday(t *testing.T) {
	assert := assert.New(t)

	var todos []*Todo
	todoNo1 := &Todo{Id: 1, Subject: "one", Due: time.Now().Format("2006-01-02")}
	todoNo2 := &Todo{Id: 2, Subject: "two", Due: time.Now().Format("2006-01-02")}

	todos = append(todos, todoNo1)
	todos = append(todos, todoNo2)

	filter := NewDateFilter(todos)
	filtered := filter.filterCompletedToday(time.Now())

	assert.Equal(0, len(filtered))

	todoNo1.Complete()
	filtered = filter.filterCompletedToday(time.Now())

	assert.Equal(1, len(filtered))
	assert.Equal(1, filtered[0].Id)

	todoNo1.Uncomplete()
	todoNo2.Complete()
	filtered = filter.filterCompletedToday(time.Now())

	assert.Equal(1, len(filtered))
	assert.Equal(2, filtered[0].Id)

}

func TestFilterThisWeek(t *testing.T) {
	assert := assert.New(t)

	var todos []*Todo
	lastWeekTodo := &Todo{Id: 1, Subject: "two", Due: time.Now().AddDate(0, 0, -7).Format("2006-01-02")}
	todayTodo := &Todo{Id: 2, Subject: "one", Due: time.Now().Format("2006-01-02")}
	nextWeekTodo := &Todo{Id: 3, Subject: "two", Due: time.Now().AddDate(0, 0, 8).Format("2006-01-02")}
	todos = append(todos, lastWeekTodo)
	todos = append(todos, todayTodo)
	todos = append(todos, nextWeekTodo)

	filter := NewDateFilter(todos)
	filtered := filter.filterThisWeek(time.Now())

	assert.Equal(1, len(filtered))
	assert.Equal(2, filtered[0].Id)
}

func TestFilterCompletedThisWeek(t *testing.T) {
	assert := assert.New(t)

	var todos []*Todo
	lastWeekTodo := &Todo{Id: 1, Subject: "two", Due: time.Now().AddDate(0, 0, -7).Format("2006-01-02")}
	todayTodo := &Todo{Id: 2, Subject: "one", Due: time.Now().Format("2006-01-02")}
	nextWeekTodo := &Todo{Id: 3, Subject: "two", Due: time.Now().AddDate(0, 0, 8).Format("2006-01-02")}
	todos = append(todos, lastWeekTodo)
	todos = append(todos, todayTodo)
	todos = append(todos, nextWeekTodo)

	filter := NewDateFilter(todos)
	filtered := filter.filterCompletedThisWeek(time.Now())

	assert.Equal(0, len(filtered))

	todayTodo.Complete()
	filtered = filter.filterCompletedThisWeek(time.Now())

	assert.Equal(1, len(filtered))
	assert.Equal(2, filtered[0].Id)

}

func TestFilterOverdue(t *testing.T) {
	assert := assert.New(t)

	var todos []*Todo
	lastWeekTodo := &Todo{Id: 1, Subject: "one", Due: time.Now().AddDate(0, 0, -7).Format("2006-01-02")}
	todayTodo := &Todo{Id: 2, Subject: "two", Due: bod(time.Now()).Format("2006-01-02")}
	tomorrowTodo := &Todo{Id: 3, Subject: "three", Due: time.Now().AddDate(0, 0, 1).Format("2006-01-02")}

	todos = append(todos, lastWeekTodo)
	todos = append(todos, todayTodo)
	todos = append(todos, tomorrowTodo)

	filter := NewDateFilter(todos)
	filtered := filter.filterOverdue(bod(time.Now()))

	assert.Equal(1, len(filtered))
	assert.Equal(1, filtered[0].Id)
}

func TestFilterDueDay(t *testing.T) {
	assert := assert.New(t)

	var todos []*Todo
	df := &DateFilter{}
	sunday := df.FindSunday(time.Now())

	monday := sunday.AddDate(0, 0, 1)
	tuesday := sunday.AddDate(0, 0, 2)
	wednesday := sunday.AddDate(0, 0, 3)

	tuesdayNext := tuesday.AddDate(0, 0, 7)

	mondayTodo := &Todo{Id: 1, Subject: "mondayTodo", Due: monday.Format("2006-01-02")}
	tuesdayTodo := &Todo{Id: 2, Subject: "tuesdayTodo", Due: tuesday.Format("2006-01-02")}

	nextTuesdayTodo := &Todo{Id: 3, Subject: "nextTuesdayTodo", Due: tuesdayNext.Format("2006-01-02")}

	todos = append(todos, mondayTodo)
	todos = append(todos, tuesdayTodo)
	todos = append(todos, nextTuesdayTodo)

	// set the users day to be Monday.
	filter := NewDateFilter(todos)
	filtered := filter.filterDueDay(monday, time.Tuesday)

	assert.Equal(1, len(filtered))
	assert.Equal(tuesdayTodo.Subject, filtered[0].Subject)

	// set the users day to be Wednesday.
	filtered = filter.filterDueDay(wednesday, time.Tuesday)

	assert.Equal(1, len(filtered))
	assert.Equal(nextTuesdayTodo.Subject, filtered[0].Subject)
}

func TestFilterAgenda(t *testing.T) {
	assert := assert.New(t)

	var todos []*Todo

	completedTodo := &Todo{Id: 1, Subject: "completed", Completed: true, Due: time.Now().Format("2006-01-02")}
	uncompletedTodo := &Todo{Id: 2, Subject: "uncompleted", Due: time.Now().Format("2006-01-02")}

	todos = append(todos, completedTodo)
	todos = append(todos, uncompletedTodo)

	filter := NewDateFilter(todos)

	filtered := filter.filterAgenda(time.Now())

	assert.Equal(1, len(filtered))
	assert.Equal(2, filtered[0].Id)
}
