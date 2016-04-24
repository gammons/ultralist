package main

import "fmt"
import "github.com/gammons/todolist/todolist"

func main() {
	store := todolist.NewFileStore()
	store.Load()
	for _, item := range store.Data {
		fmt.Println(item)
	}
}
