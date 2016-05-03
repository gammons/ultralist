package todolist

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
