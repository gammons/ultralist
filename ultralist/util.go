package ultralist

import (
	"fmt"
	"time"

	"github.com/twinj/uuid"
)

// AddIfNotThere is appending an item to an array if the item is not already present.
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

// AddTodoIfNotThere is appending an todo item to an todo array if the item is not already present.
func AddTodoIfNotThere(arr []*Todo, item *Todo) []*Todo {
	there := false
	for _, arrItem := range arr {
		if item.ID == arrItem.ID {
			there = true
		}
	}
	if !there {
		arr = append(arr, item)
	}
	return arr
}

func newUUID() string {
	return fmt.Sprintf("%s", uuid.NewV4())
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

func pluralize(count int, singular, plural string) string {
	if count > 1 {
		return plural
	}
	return singular
}
