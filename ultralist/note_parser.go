package ultralist

import (
	"fmt"
	"regexp"
	"strconv"
)

// NoteParser is a legacy class that parses various
type NoteParser struct{}

// ParseAddNote is adding a note to an todo.
func (p *NoteParser) ParseAddNote(todo *Todo, input string) bool {
	r, _ := regexp.Compile(`^an\s+\d+\s+(.*)`)
	matches := r.FindStringSubmatch(input)
	if len(matches) != 2 {
		return false
	}

	todo.Notes = append(todo.Notes, matches[1])
	return true
}

// ParseDeleteNote is deleting a note from an todo.
func (p *NoteParser) ParseDeleteNote(todo *Todo, input string) bool {
	r, _ := regexp.Compile(`^dn\s+\d+\s+(\d+)`)
	matches := r.FindStringSubmatch(input)
	if len(matches) != 2 {
		return false
	}

	rmid, err := p.getNoteID(matches[1])
	if err != nil {
		return false
	}

	for id := range todo.Notes {
		if id == rmid {
			todo.Notes = append(todo.Notes[:rmid], todo.Notes[rmid+1:]...)
			return true
		}
	}
	return false
}

// ParseEditNote is editing a note from an todo.
func (p *NoteParser) ParseEditNote(todo *Todo, input string) bool {
	r, _ := regexp.Compile(`^en\s+\d+\s+(\d+)\s+(.*)`)
	matches := r.FindStringSubmatch(input)
	if len(matches) != 3 {
		return false
	}

	edid, err := p.getNoteID(matches[1])
	if err != nil {
		return false
	}

	for id := range todo.Notes {
		if id == edid {
			todo.Notes[id] = matches[2]
			return true
		}
	}
	return false
}

// ParseShowNote is defining if notes should be shown or not.
func (p *NoteParser) ParseShowNote(todo *Todo, input string) bool {
	r, _ := regexp.Compile(`^n\s+\d+`)
	matches := r.FindStringSubmatch(input)
	if len(matches) != 1 {
		return false
	}
	return true
}

func (p *NoteParser) getNoteID(input string) (int, error) {
	ret, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("wrong note id")
		return -1, err
	}
	return ret, nil
}
