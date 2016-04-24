package todolist

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/now"
)

type Parser struct{}

func (p *Parser) Parse(input string) *Todo {
	todo := NewTodo()
	todo.Subject = p.Subject(input)
	todo.Projects = p.Projects(input)
	todo.Contexts = p.Contexts(input)
	if p.hasDue(input) {
		todo.FormattedDue = p.Due(input)
	}
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
	r, _ := regexp.Compile(`\+\w+`)
	return p.matchWords(input, r)
}

func (p *Parser) Contexts(input string) []string {
	r, err := regexp.Compile(`\@\w+`)
	if err != nil {
		fmt.Println("regex error", err)
	}
	return p.matchWords(input, r)
}

func (p *Parser) hasDue(input string) bool {
	r, _ := regexp.Compile(`due \w+$`)
	return r.MatchString(input)
}

func (p *Parser) Due(input string) time.Time {
	r, _ := regexp.Compile(`due .*$`)

	res := r.FindString(input)
	res = res[4:len(res)]
	switch {
	case res == "today":
		return now.BeginningOfDay()
	case res == "tomorrow" || res == "tom":
		return now.BeginningOfDay().AddDate(0, 0, 1)
	case res == "monday" || res == "mon":
		n := now.BeginningOfDay()
		return now.New(n).Monday().AddDate(0, 0, 7)
	case res == "tuesday" || res == "tue":
		n := now.BeginningOfDay()
		return now.New(n).Monday().AddDate(0, 0, 1)
	case res == "wednesday" || res == "wed":
		n := now.BeginningOfDay()
		return now.New(n).Monday().AddDate(0, 0, 2)
	case res == "next week":
		n := now.BeginningOfDay()
		return now.New(n).Monday().AddDate(0, 0, 7)
	}
	//return now.Parse(input)
	return time.Now()
}

func (p *Parser) matchWords(input string, r *regexp.Regexp) []string {
	results := r.FindAllString(input, -1)
	ret := []string{}

	for _, val := range results {
		ret = append(ret, val[1:len(val)])
	}
	return ret
}
