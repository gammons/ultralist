package ultralist

import (
	"errors"
	"regexp"
	"strings"
)

// Parser parses text to extract a Filter struct
type InputParser struct{}

/*

# status of not now
status:-now

# status of now OR next
status:now,next

# status not now
status:-now
status:next

# due today OR tomorrow
due:tod,tom
due:today due:tom

# filter todos before a certain date
due:<aug15

completed:true

priority:false

project:one,-two

*/

// Parse parses raw input and returns a Filter object
func (p *InputParser) Parse(input string) (*Filter, error) {
	if input == "" {
		return &Filter{}, errors.New("Could not parse input")
	}

	filter := &Filter{
		HasStatus:      false,
		HasCompleted:   false,
		HasCompletedAt: false,
		HasIsPriority:  false,
		HasDue:         false,
	}
	var subjectMatches []string

	cr, _ := regexp.Compile(`\@[\p{L}\d_-]+`)
	filter.Contexts = p.matchWords(input, cr)

	pr, _ := regexp.Compile(`\+[\p{L}\d_-]+`)
	filter.Projects = p.matchWords(input, pr)

	for _, word := range strings.Split(input, " ") {
		match := false
		r1, _ := regexp.Compile(`archived:.*$`)
		if r1.MatchString(word) {
			filter.HasArchived = true
			filter.Archived = p.parseBool(r1.FindString(word)[9:])
			match = true
		}

		r2, _ := regexp.Compile(`priority:.*$`)
		if r2.MatchString(word) {
			filter.HasIsPriority = true
			filter.IsPriority = p.parseBool(r2.FindString(word)[9:])
			match = true
		}

		r3, _ := regexp.Compile(`completed:.*$`)
		if r3.MatchString(word) {
			filter.HasCompleted = true
			filter.Completed = p.parseBool(r3.FindString(word)[10:])
			match = true
		}

		r4, _ := regexp.Compile(`due:.*$`)
		if r4.MatchString(word) {
			filter.HasDue = true
			filter.Due = p.parseString(r4.FindString(word)[4:])
			match = true
		}

		r5, _ := regexp.Compile(`status:.*$`)
		if r5.MatchString(word) {
			filter.HasStatus = true
			filter.Status = p.parseString(r5.FindString(word)[7:])
			match = true
		}

		r6, _ := regexp.Compile(`completedat:.*$`)
		if r6.MatchString(word) {
			filter.HasCompletedAt = true
			filter.Status = p.parseString(r6.FindString(word)[7:])
			match = true
		}

		if !match {
			subjectMatches = append(subjectMatches, word)
		}
	}

	filter.Subject = strings.Join(subjectMatches, " ")

	return filter, nil
}

func (p *InputParser) parseString(input string) []string {
	return strings.Split(input, ",")
}

func (p *InputParser) parseBool(input string) bool {
	return input == "true"
}

func (p *InputParser) matchWords(input string, r *regexp.Regexp) []string {
	results := r.FindAllString(input, -1)
	ret := []string{}

	for _, val := range results {
		ret = append(ret, val[1:])
	}
	return ret
}