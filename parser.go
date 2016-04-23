package todolist

import (
	"fmt"
	"regexp"
	"strings"
)

type Parser struct{}

func (p *Parser) Parse(input string) *Todo {
	todo := NewTodo()
	todo.Subject = p.Subject(input)
	todo.Projects = p.Projects(input)
	return todo
}

func (p *Parser) Subject(input string) string {
	if strings.Contains(input, " due") {
		index := strings.LastIndex(input, " due")
		return input[0:index]
	} else {
		return input
	}
}

func (p *Parser) Projects(input string) []string {
	r, err := regexp.Compile(`\+\w+`)
	if err != nil {
		fmt.Printf("There was a problem with the regexp", err)
	}
	results := r.FindAllString(input, -1)
	ret := []string{}

	for _, val := range results {
		fmt.Println("Result is ", val)
		ret = append(ret, val[1:len(val)])
	}
	return ret
}
