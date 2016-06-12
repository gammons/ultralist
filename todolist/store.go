package todolist

type Store interface {
	Initialize()
	Load()
	Save()
}
