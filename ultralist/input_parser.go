package ultralist

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// InputParser parses text to extract a Filter struct
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
	filter := &Filter{
		HasStatus:        false,
		HasCompleted:     false,
		HasCompletedAt:   false,
		HasIsPriority:    false,
		HasProjectFilter: false,
		HasContextFilter: false,
		HasDueBefore:     false,
		HasDue:           false,
		HasDueAfter:      false,
		HasRecur:         false,
	}

	dateParser := &DateParser{}

	var subjectMatches []string

	cr, _ := regexp.Compile(`\@[\p{L}\d_-]+`)
	filter.Contexts = p.matchWords(input, cr)

	pr, _ := regexp.Compile(`\+[\p{L}\d_-]+`)
	filter.Projects = p.matchWords(input, pr)

	for _, word := range strings.Split(input, " ") {
		match := false
		r, _ := regexp.Compile(`archived:.*$`)
		if r.MatchString(word) {
			filter.HasArchived = true
			filter.Archived = p.parseBoolString(r.FindString(word)[9:])
			match = true
		}

		r, _ = regexp.Compile(`priority:.*$`)
		if r.MatchString(word) {
			filter.HasIsPriority = true
			filter.IsPriority = p.parseBoolString(r.FindString(word)[9:])
			match = true
		}

		r, _ = regexp.Compile(`completed:.*$`)
		if r.MatchString(word) {
			filter.HasCompleted = true
			filter.Completed = p.parseBoolString(r.FindString(word)[10:])
			match = true
		}

		r, _ = regexp.Compile(`duebefore:.*$`)
		if r.MatchString(word) {
			filter.HasDueBefore = true
			dueDate, err := dateParser.ParseDate(r.FindString(word)[10:], time.Now())
			if err != nil {
				return filter, err
			}

			if dueDate.IsZero() {
				filter.DueBefore = ""
			} else {
				filter.DueBefore = dueDate.Format(DATE_FORMAT)
			}
			match = true
		}

		r, _ = regexp.Compile(`due:.*$`)
		if r.MatchString(word) {
			filter.HasDue = true

			dueDate, err := dateParser.ParseDate(r.FindString(word)[4:], time.Now())
			if err != nil {
				return filter, err
			}

			if dueDate.IsZero() {
				filter.Due = ""
			} else {
				switch word {
					case "due:agenda":
						filter.HasDueBefore = true
						filter.HasDue = false
						filter.DueBefore = dueDate.Format(DATE_FORMAT)
					case "due:lastweek", "due:thisweek", "due:nextweek":
						filter.HasDue = false
						filter.HasDueBefore = true
						filter.DueBefore = dueDate.AddDate(0,0,7).Format(DATE_FORMAT)
						filter.HasDueAfter = true
						filter.DueAfter = dueDate.AddDate(0,0,-1).Format(DATE_FORMAT)
					case "due:lastmonth", "due:thismonth", "due:nextmonth":
						filter.HasDue = false
						filter.HasDueBefore = true
						filter.DueBefore = dueDate.AddDate(0,1,0).Format(DATE_FORMAT)
						filter.HasDueAfter = true
						filter.DueAfter = dueDate.AddDate(0,0,-1).Format(DATE_FORMAT)
					default:
						filter.Due = dueDate.Format(DATE_FORMAT)
				}
			}
			match = true
		}

		r, _ = regexp.Compile(`dueafter:.*$`)
		if r.MatchString(word) {
			filter.HasDueAfter = true
			dueDate, err := dateParser.ParseDate(r.FindString(word)[9:], time.Now())
			if err != nil {
				return filter, err
			}

			if dueDate.IsZero() {
				filter.DueAfter = ""
			} else {
				filter.DueAfter = dueDate.Format(DATE_FORMAT)
			}
			match = true
		}

		r, _ = regexp.Compile(`status:.*$`)
		if r.MatchString(word) {
			filter.HasStatus = true
			filter.Status, filter.ExcludeStatus = p.parseString(r.FindString(word)[7:])
			match = true
		}

		r, _ = regexp.Compile(`completedat:.*$`)
		if r.MatchString(word) {
			filter.HasCompletedAt = true
			filter.CompletedAt, _ = p.parseString(r.FindString(word)[12:])
			match = true
		}

		r, _ = regexp.Compile(`project:.*$`)
		if r.MatchString(word) {
			filter.HasProjectFilter = true
			filter.Projects, filter.ExcludeProjects = p.parseString(r.FindString(word)[8:])
			match = true
		}

		r, _ = regexp.Compile(`context:.*$`)
		if r.MatchString(word) {
			filter.HasContextFilter = true
			filter.Contexts, filter.ExcludeContexts = p.parseString(r.FindString(word)[8:])
			match = true
		}

		r, _ = regexp.Compile(`recur:.*$`)
		if r.MatchString(word) {
			match = true

			filter.HasRecur = true
			filter.Recur = r.FindString(word)[6:]

			if filter.Recur == "none" {
				filter.Recur = ""
			}

			r := &Recurrence{}
			if !r.ValidRecurrence(filter.Recur) {
				return filter, fmt.Errorf("I could not understand the recurrence you gave me: '%s'", filter.Recur)
			}
		}

		r, _ = regexp.Compile(`until:.*$`)
		if r.MatchString(word) {
			date, err := dateParser.ParseDate(r.FindString(word)[6:], time.Now())
			if err != nil {
				return filter, err
			}
			match = true

			filter.RecurUntil = date.Format(DATE_FORMAT)
		}

		if !match {
			subjectMatches = append(subjectMatches, word)
		}
	}

	filter.Subject = strings.Join(subjectMatches, " ")

	return filter, nil
}

func (p *InputParser) parseString(input string) ([]string, []string) {
	var include []string
	var exclude []string
	for _, str := range strings.Split(input, ",") {
		if strings.HasPrefix(str, "-") {
			exclude = append(exclude, str[1:])
		} else {
			include = append(include, str)
		}
	}
	return include, exclude
}

func (p *InputParser) parseBoolString(input string) bool {
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
