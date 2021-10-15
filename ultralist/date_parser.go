package ultralist

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// DateParser is a thing that parses relative, arbitrary, and absolute dates, and returns a date in the format of yyyy-mm-dd
type DateParser struct{}

// ParseDate takes a date from a Filter and turns it into a string with the format of yyyy-mm-dd
func (dp *DateParser) ParseDate(dateString string, pivotDay time.Time) (date time.Time, err error) {
	switch dateString {
	case "none":
		return time.Time{}, nil
	case "yesterday", "yes":
		return bod(pivotDay).AddDate(0, 0, -1), nil
	case "today", "tod":
		return bod(pivotDay), nil
	case "tomorrow", "tom", "agenda":
		return bod(pivotDay).AddDate(0, 0, 1), nil
	case "monday", "mon":
		return dp.monday(pivotDay), nil
	case "tuesday", "tue":
		return dp.tuesday(pivotDay), nil
	case "wednesday", "wed":
		return dp.wednesday(pivotDay), nil
	case "thursday", "thu":
		return dp.thursday(pivotDay), nil
	case "friday", "fri":
		return dp.friday(pivotDay), nil
	case "saturday", "sat":
		return dp.saturday(pivotDay), nil
	case "sunday", "sun":
		return dp.sunday(pivotDay), nil
	case "lastweek":
		n := bod(pivotDay)
		return dp.getNearestMonday(n).AddDate(0, 0, -7), nil
	case "nextweek":
		n := bod(pivotDay)
		return dp.getNearestMonday(n).AddDate(0, 0, 7), nil
	case "thisweek":
		n := bod(pivotDay)
		return dp.getNearestMonday(n), nil
	case "thismonth":
		n := bod(pivotDay)
		return dp.getNearestFirstOfMonth(n), nil
	case "lastmonth":
		n := bod(pivotDay)
		return dp.getNearestFirstOfMonth(n).AddDate(0,-1,0), nil
	case "nextmonth":
		n := bod(pivotDay)
		return dp.getNearestFirstOfMonth(n).AddDate(0,1,0), nil
	}
	return dp.parseSpecificDate(dateString, pivotDay)
}

func (dp *DateParser) monday(day time.Time) time.Time {
	mon := dp.getNearestMonday(day)
	return dp.thisOrNextWeek(mon, day)
}

func (dp *DateParser) tuesday(day time.Time) time.Time {
	tue := dp.getNearestMonday(day).AddDate(0, 0, 1)
	return dp.thisOrNextWeek(tue, day)
}

func (dp *DateParser) wednesday(day time.Time) time.Time {
	wed := dp.getNearestMonday(day).AddDate(0, 0, 2)
	return dp.thisOrNextWeek(wed, day)
}

func (dp *DateParser) thursday(day time.Time) time.Time {
	thu := dp.getNearestMonday(day).AddDate(0, 0, 3)
	return dp.thisOrNextWeek(thu, day)
}

func (dp *DateParser) friday(day time.Time) time.Time {
	fri := dp.getNearestMonday(day).AddDate(0, 0, 4)
	return dp.thisOrNextWeek(fri, day)
}

func (dp *DateParser) saturday(day time.Time) time.Time {
	sat := dp.getNearestMonday(day).AddDate(0, 0, 5)
	return dp.thisOrNextWeek(sat, day)
}

func (dp *DateParser) sunday(day time.Time) time.Time {
	sun := dp.getNearestMonday(day).AddDate(0, 0, 6)
	return dp.thisOrNextWeek(sun, day)
}

func (dp *DateParser) thisOrNextWeek(day time.Time, pivotDay time.Time) time.Time {
	if day.Before(pivotDay) {
		return day.AddDate(0, 0, 7)
	}
	return day
}

func (dp *DateParser) parseSpecificDate(date string, pivot time.Time) (time.Time, error) {
	yearStr := strconv.Itoa(pivot.Year())
	dateWithYear := fmt.Sprintf("%s-%s", yearStr, date)

	if res, err := time.Parse("2006-Jan2", dateWithYear); err == nil {
		if res.After(pivot) {
			return res, nil
		}
		return res.AddDate(1, 0, 0), nil
	}

	if res, err := time.Parse("2006-2Jan", dateWithYear); err == nil {
		if res.After(pivot) {
			return res, nil
		}
		return res.AddDate(1, 0, 0), nil
	}

	errString := fmt.Sprintf("Could not parse the date you gave me: '%s'", date)

	return time.Time{}, errors.New(errString)
}

func (dp *DateParser) getNearestMonday(t time.Time) time.Time {
	for {
		if t.Weekday() != time.Monday {
			t = t.AddDate(0, 0, -1)
		} else {
			return t
		}
	}
}

func (dp *DateParser) getNearestFirstOfMonth(t time.Time) time.Time {
	for {
		if t.Day() != 1 {
			t = t.AddDate(0, 0, -1)
		} else {
			return t
		}
	}
}
