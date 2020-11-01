package ultralist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeeklyCompleteWeekBefore(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-10-17")
	// Wed, Oct 28
	todo := &Todo{Recur: "weekly", Due: "2020-10-28"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-04", nextTodo.Due)
}

func TestWeeklyCompleteDayBefore(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-10-27")
	// Wed, Oct 28
	todo := &Todo{Recur: "weekly", Due: "2020-10-28"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-04", nextTodo.Due)
}

func TestWeeklyCompleteDayOf(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-10-28")
	todo := &Todo{Recur: "weekly", Due: "2020-10-28"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-04", nextTodo.Due)
}

func TestWeeklyCompleteDayAfter(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-10-29")
	todo := &Todo{Recur: "weekly", Due: "2020-10-28"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-04", nextTodo.Due)
}

func TestWeeklyCompleteWeekAfter(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-11-10")
	todo := &Todo{Recur: "weekly", Due: "2020-10-28"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-11", nextTodo.Due)
}

func TestMonthlyCompletedBefore(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-08-28")
	todo := &Todo{Recur: "monthly", Due: "2020-10-28"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-28", nextTodo.Due)
}

func TestMonthlyCompletedDayOf(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-10-28")
	todo := &Todo{Recur: "monthly", Due: "2020-10-28"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-28", nextTodo.Due)
}

func TestMonthlyCompletedAfter(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-11-29")
	todo := &Todo{Recur: "monthly", Due: "2020-10-28"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-12-28", nextTodo.Due)
}

func TestWeekDaysOnMonday(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-11-02")
	todo := &Todo{Recur: "weekdays", Due: "2020-11-02"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-03", nextTodo.Due)
}

func TestWeekDaysOnFriday(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-11-06")
	todo := &Todo{Recur: "weekdays", Due: "2020-11-06"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-09", nextTodo.Due)
}

func TestDailySameDay(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-11-06")
	todo := &Todo{Recur: "daily", Due: "2020-11-06"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-07", nextTodo.Due)
}

func TestDailyOverdue(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-11-09")
	todo := &Todo{Recur: "daily", Due: "2020-11-06"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-10", nextTodo.Due)
}

func TestDailyEarly(t *testing.T) {
	assert := assert.New(t)

	r := &Recurrence{}
	pivot, _ := time.Parse(DATE_FORMAT, "2020-11-04")
	todo := &Todo{Recur: "daily", Due: "2020-11-06"}

	nextTodo := r.NextRecurringTodo(todo, pivot)

	assert.Equal("2020-11-07", nextTodo.Due)
}
