package todolist

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/now"
)

type Parser struct{}

func (p *Parser) ParseNewTodo(input string) *Todo {
	todo := NewTodo()
	todo.Subject = p.Subject(input)
	todo.Projects = p.Projects(input)
	todo.Contexts = p.Contexts(input)
	if p.hasDue(input) {
		todo.Due = p.Due(input, time.Now())
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

func (p *Parser) Due(input string, day time.Time) string {
	r, _ := regexp.Compile(`due .*$`)

	res := r.FindString(input)
	res = res[4:len(res)]
	switch {
	case res == "none":
		return ""
	case res == "today":
		return now.BeginningOfDay().Format("2006-01-02")
	case res == "tomorrow" || res == "tom":
		return now.BeginningOfDay().AddDate(0, 0, 1).Format("2006-01-02")
	case res == "monday" || res == "mon":
		return p.monday(day)
	case res == "tuesday" || res == "tue":
		return p.tuesday(day)
	case res == "wednesday" || res == "wed":
		return p.wednesday(day)
	case res == "thursday" || res == "thu":
		return p.thursday(day)
	case res == "friday" || res == "fri":
		return p.friday(day)
	case res == "saturday" || res == "sat":
		return p.saturday(day)
	case res == "sunday" || res == "sun":
		return p.sunday(day)
	case res == "next week":
		n := now.BeginningOfDay()
		return now.New(n).Monday().AddDate(0, 0, 7).Format("2006-01-02")
	}
	return time.Now().Format("2006-01-02")
}

func (p *Parser) monday(day time.Time) string {
	mon := now.New(day).Monday()
	return p.thisOrNextWeek(mon, day)
}

func (p *Parser) tuesday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 1)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) wednesday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 2)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) thursday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 3)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) friday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 4)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) saturday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 5)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) sunday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 6)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) thisOrNextWeek(day time.Time, pivotDay time.Time) string {
	if day.Before(pivotDay) {
		return day.AddDate(0, 0, 7).Format("2006-01-02")
	} else {
		return day.Format("2006-01-02")
	}
}

func (p *Parser) matchWords(input string, r *regexp.Regexp) []string {
	results := r.FindAllString(input, -1)
	ret := []string{}

	for _, val := range results {
		ret = append(ret, val[1:len(val)])
	}
	return ret
}
