package todolist

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type App struct {
	TodoStore Store
	TodoList  *TodoList
}

func NewApp() *App {
	app := &App{TodoList: &TodoList{}, TodoStore: NewFileStore()}
	return app
}

func (a *App) InitializeRepo() {
	a.TodoStore.Initialize()
}

func (a *App) AddTodo(input string) {
	a.Load()
	parser := &Parser{}
	todo := parser.ParseNewTodo(input)
	if todo == nil {
		fmt.Println("I need more information. Try something like 'todo a chat with @bob due tom'")
		return
	}

	id := a.TodoList.NextId()
	a.TodoList.Add(todo)
	a.Save()
	fmt.Printf("Todo %d added.\n", id)
}

func (a *App) DeleteTodo(input string) {
	a.Load()
	ids := a.getIds(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Delete(ids...)
	a.Save()
	fmt.Printf("%s deleted.\n", pluralize(len(ids), "Todo", "Todos"))
}

func (a *App) CompleteTodo(input string) {
	a.Load()
	ids := a.getIds(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Complete(ids...)
	a.Save()
	fmt.Println("Todo completed.")
}

func (a *App) UncompleteTodo(input string) {
	a.Load()
	ids := a.getIds(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Uncomplete(ids...)
	a.Save()
	fmt.Println("Todo uncompleted.")
}

func (a *App) ArchiveTodo(input string) {
	a.Load()
	ids := a.getIds(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Archive(ids...)
	a.Save()
	fmt.Println("Todo archived.")
}

func (a *App) UnarchiveTodo(input string) {
	a.Load()
	ids := a.getIds(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Unarchive(ids...)
	a.Save()
	fmt.Println("Todo unarchived.")
}

func (a *App) EditTodo(input string) {
	a.Load()
	id := a.getId(input)
	if id == -1 {
		return
	}
	todo := a.TodoList.FindById(id)
	if todo == nil {
		fmt.Println("No such id.")
		return
	}
	parser := &Parser{}

	if parser.ParseEditTodo(todo, input) {
		a.Save()
		fmt.Println("Todo updated.")
	}
}

func (a *App) ExpandTodo(input string) {
	a.Load()
	id := a.getId(input)
	parser := &Parser{}
	if id == -1 {
		return
	}

	commonProject := parser.ExpandProject(input)
	todos := strings.LastIndex(input, ":")
	if commonProject == "" || len(input) <= todos+1 || todos == -1 {
		fmt.Println("I'm expecting a format like \"todolist ex <project>: <todo1>, <todo2>, ...\"")
		return
	}

	newTodos := strings.Split(input[todos+1:], ",")

	for _, todo := range newTodos {
		args := []string{"add ", commonProject, " ", todo}
		a.AddTodo(strings.Join(args, ""))
	}

	a.TodoList.Delete(id)
	a.Save()
	fmt.Println("Todo expanded.")
}

func (a *App) HandleNotes(input string) {
	a.Load()
	id := a.getId(input)
	if id == -1 {
		return
	}
	todo := a.TodoList.FindById(id)
	if todo == nil {
		fmt.Println("No such id.")
		return
	}
	parser := &Parser{}

	if parser.ParseAddNote(todo, input) {
		fmt.Println("Note added.")
	} else if parser.ParseDeleteNote(todo, input) {
		fmt.Println("Note deleted.")
	} else if parser.ParseEditNote(todo, input) {
		fmt.Println("Note edited.")
	} else if parser.ParseShowNote(todo, input) {
		groups := map[string][]*Todo{}
		groups[""] = append(groups[""], todo)
		formatter := NewFormatter(&GroupedTodos{Groups: groups})
		formatter.Print(true)
		return
	}
	a.Save()
}

func (a *App) ArchiveCompleted() {
	a.Load()
	for _, todo := range a.TodoList.Todos() {
		if todo.Completed {
			todo.Archive()
		}
	}
	a.Save()
	fmt.Println("All completed todos have been archived.")
}

func (a *App) ListTodos(input string) {
	a.Load()
	filtered := NewFilter(a.TodoList.Todos()).Filter(input)
	grouped := a.getGroups(input, filtered)

	formatter := NewFormatter(grouped)
	re, _ := regexp.Compile(`^ln`)
	formatter.Print(re.MatchString(input))
}

func (a *App) PrioritizeTodo(input string) {
	a.Load()
	ids := a.getIds(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Prioritize(ids...)
	a.Save()
	fmt.Println("Todo prioritized.")
}

func (a *App) UnprioritizeTodo(input string) {
	a.Load()
	ids := a.getIds(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Unprioritize(ids...)
	a.Save()
	fmt.Println("Todo un-prioritized.")
}

func (a *App) getId(input string) int {
	re, _ := regexp.Compile("\\d+")
	if re.MatchString(input) {
		id, _ := strconv.Atoi(re.FindString(input))
		return id
	}

	fmt.Println("Invalid id.")
	return -1
}

func (a *App) getIds(input string) (ids []int) {

	idGroups := strings.Split(input, ",")
	for _, idGroup := range idGroups {
		if rangedIds, err := a.parseRangedIds(idGroup); len(rangedIds) > 0 || err != nil {
			if err != nil {
				fmt.Printf("Invalid id group: %s.\n", input)
				continue
			}
			ids = append(ids, rangedIds...)
		} else if id := a.getId(idGroup); id != -1 {
			ids = append(ids, id)
		} else {
			fmt.Printf("Invalid id: %s.\n", idGroup)
		}
	}
	return ids
}

func (a *App) parseRangedIds(input string) (ids []int, err error) {
	rangeNumberRE, _ := regexp.Compile("(\\d+)-(\\d+)")
	if matches := rangeNumberRE.FindStringSubmatch(input); len(matches) > 0 {
		lowerID, _ := strconv.Atoi(matches[1])
		upperID, _ := strconv.Atoi(matches[2])
		if lowerID >= upperID {
			return ids, fmt.Errorf("Invalid id group: %s.\n", input)
		}
		for id := lowerID; id <= upperID; id++ {
			ids = append(ids, id)
		}
	}
	return ids, err
}

func (a *App) getGroups(input string, todos []*Todo) *GroupedTodos {
	grouper := &Grouper{}
	contextRegex, _ := regexp.Compile(`by c.*$`)
	projectRegex, _ := regexp.Compile(`by p.*$`)

	var grouped *GroupedTodos

	if contextRegex.MatchString(input) {
		grouped = grouper.GroupByContext(todos)
	} else if projectRegex.MatchString(input) {
		grouped = grouper.GroupByProject(todos)
	} else {
		grouped = grouper.GroupByNothing(todos)
	}
	return grouped
}

func (a *App) GarbageCollect() {
	a.Load()
	a.TodoList.GarbageCollect()
	a.Save()
	fmt.Println("Garbage collection complete.")
}

func (a *App) Load() error {
	todos, err := a.TodoStore.Load()
	if err != nil {
		return err
	}
	a.TodoList.Load(todos)
	return nil
}

func (a *App) Save() {
	a.TodoStore.Save(a.TodoList.Data)
}
