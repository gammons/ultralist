package todolist

import "time"

func AddIfNotThere(arr []string, items []string) []string {
	for _, item := range items {
		there := false
		for _, arrItem := range arr {
			if item == arrItem {
				there = true
			}
		}
		if !there {
			arr = append(arr, item)
		}
	}
	return arr
}

func AddTodoIfNotThere(arr []*Todo, item *Todo) []*Todo {
	there := false
	for _, arrItem := range arr {
		if item.Id == arrItem.Id {
			there = true
		}
	}
	if !there {
		arr = append(arr, item)
	}
	return arr
}

func bod(t time.Time) time.Time {
	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func timestamp(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	return time.Date(year, month, day, hour, min, sec, 0, t.Location())
}

func getNearestMonday(t time.Time) time.Time {
	for {
		if t.Weekday() != time.Monday {
			t = t.AddDate(0, 0, -1)
		} else {
			return t
		}
	}
}

func pluralize(count int, singular, plural string) string {
	if count > 1 {
		return plural
	}
	return singular
}

func isToday(t time.Time) bool {
	nowYear, nowMonth, nowDay := time.Now().Date()
	timeYear, timeMonth, timeDay := t.Date()
	return nowYear == timeYear &&
		nowMonth == timeMonth &&
		nowDay == timeDay
}

func isTomorrow(t time.Time) bool {
	nowYear, nowMonth, nowDay := time.Now().AddDate(0, 0, 1).Date()
	timeYear, timeMonth, timeDay := t.Date()
	return nowYear == timeYear &&
		nowMonth == timeMonth &&
		nowDay == timeDay
}

func isPastDue(t time.Time) bool {
	return time.Now().After(t)
}
