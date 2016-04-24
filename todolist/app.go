package todolist

type App struct {
	TodoStore Store
}

func NewApp() *App {
	app := &App{TodoStore: NewFileStore()}
	app.TodoStore.Load()
	return app
}

func (a *App) ListTodos(input string) {
	formatter := NewFormatter(a.TodoStore.Todos())
	formatter.Print()
}
