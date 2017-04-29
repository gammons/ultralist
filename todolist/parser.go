package todolist

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	input string
}

func (p *Parser) ParseNewTodo(input string) *Todo {
	r, _ := regexp.Compile(`^(add|a)(\\ |) `)
	input = r.ReplaceAllString(input, "")
	if input == "" {
		return nil
	}

	todo := NewTodo()
	todo.Subject = p.Subject(input)
	todo.Projects = p.Projects(input)
	todo.Contexts = p.Contexts(input)
	if p.hasDue(input) {
		todo.Due = p.Due(input, time.Now())
	}
	return todo
}

func (p Parser) Parse() (string, int, string) {
	input := p.input
	r := regexp.MustCompile(`(\w+) (\d+) (.*)`)
	matches := r.FindStringSubmatch(input)
	id, err := strconv.Atoi(matches[2])
	if err != nil {
		fmt.Println("Invalid id.")
		id = -1
	}

	return matches[1], id, matches[3]
}

func (p *Parser) Subject(input string) string {
	if strings.Contains(input, " due") {
		index := strings.LastIndex(input, " due")
		return input[0:index]
	} else {
		return input
	}
}

func (p *Parser) ExpandProject(input string) string {
	r, _ := regexp.Compile(`(ex|expand) +\d+ +\+[\p{L}\d_-]+:`)
	pattern := r.FindString(input)
	if len(pattern) == 0 {
		return ""
	}

	newProject := pattern[0 : len(pattern)-1]
	project := strings.Split(newProject, " ")
	return project[len(project)-1]
}

func (p *Parser) Projects(input string) []string {
	r, _ := regexp.Compile(`\+[\p{L}\d_-]+`)
	return p.matchWords(input, r)
}

func (p *Parser) Contexts(input string) []string {
	r, err := regexp.Compile(`\@[\p{L}\d_]+`)
	if err != nil {
		fmt.Println("regex error", err)
	}
	return p.matchWords(input, r)
}

func (p *Parser) hasDue(input string) bool {
	r1, _ := regexp.Compile(`due \w+$`)
	r2, _ := regexp.Compile(`due \w+ \d+$`)
	return (r1.MatchString(input) || r2.MatchString(input))
}

func (p *Parser) Due(input string, day time.Time) string {
	r, _ := regexp.Compile(`due .*$`)

	res := r.FindString(input)
	res = res[4:]
	switch res {
	case "none":
		return ""
	case "today", "tod":
		return bod(time.Now()).Format("2006-01-02")
	case "tomorrow", "tom":
		return bod(time.Now()).AddDate(0, 0, 1).Format("2006-01-02")
	case "monday", "mon":
		return p.monday(day)
	case "tuesday", "tue":
		return p.tuesday(day)
	case "wednesday", "wed":
		return p.wednesday(day)
	case "thursday", "thu":
		return p.thursday(day)
	case "friday", "fri":
		return p.friday(day)
	case "saturday", "sat":
		return p.saturday(day)
	case "sunday", "sun":
		return p.sunday(day)
	case "last week":
		n := bod(time.Now())
		return getNearestMonday(n).AddDate(0, 0, -7).Format("2006-01-02")
	case "next week":
		n := bod(time.Now())
		return getNearestMonday(n).AddDate(0, 0, 7).Format("2006-01-02")
	}
	return p.parseArbitraryDate(res, time.Now())
}

func (p *Parser) parseArbitraryDate(_date string, pivot time.Time) string {
	d1 := p.parseArbitraryDateWithYear(_date, pivot.Year())

	var diff1 time.Duration
	if d1.After(time.Now()) {
		diff1 = d1.Sub(pivot)
	} else {
		diff1 = pivot.Sub(d1)
	}
	d2 := p.parseArbitraryDateWithYear(_date, pivot.Year()+1)
	if d2.Sub(pivot) > diff1 {
		return d1.Format("2006-01-02")
	} else {
		return d2.Format("2006-01-02")
	}
}

func (p *Parser) parseArbitraryDateWithYear(_date string, year int) time.Time {
	res := strings.Join([]string{_date, strconv.Itoa(year)}, " ")
	if date, err := time.Parse("Jan 2 2006", res); err == nil {
		return date
	}

	if date, err := time.Parse("2 Jan 2006", res); err == nil {
		return date
	}
	fmt.Printf("Could not parse the date you gave me: %s\n", _date)
	fmt.Println("I'm expecting a date like \"Dec 22\" or \"22 Dec\".")
	fmt.Println("See http://todolist.site/#adding for more info.")
	os.Exit(-1)
	return time.Now()
}

func (p *Parser) monday(day time.Time) string {
	mon := getNearestMonday(day)
	return p.thisOrNextWeek(mon, day)
}

func (p *Parser) tuesday(day time.Time) string {
	tue := getNearestMonday(day).AddDate(0, 0, 1)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) wednesday(day time.Time) string {
	tue := getNearestMonday(day).AddDate(0, 0, 2)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) thursday(day time.Time) string {
	tue := getNearestMonday(day).AddDate(0, 0, 3)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) friday(day time.Time) string {
	tue := getNearestMonday(day).AddDate(0, 0, 4)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) saturday(day time.Time) string {
	tue := getNearestMonday(day).AddDate(0, 0, 5)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) sunday(day time.Time) string {
	tue := getNearestMonday(day).AddDate(0, 0, 6)
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
		ret = append(ret, val[1:])
	}
	return ret
}
