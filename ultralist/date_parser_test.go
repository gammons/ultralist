package ultralist

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNone(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2010-01-01")

	res, err := dp.ParseDate("none", date)

	if err != nil {
		fmt.Println("err is ", err)
		t.Fail()
	}

	assert.Equal(true, res.IsZero())
}

func TestToday(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2010-01-01")

	res, err := dp.ParseDate("today", date)
	formattedRes := res.Format(DATE_FORMAT)

	if err != nil {
		fmt.Println("err is ", err)
		t.Fail()
	}

	expectedTime, _ := time.Parse(DATE_FORMAT, "2010-01-01")
	expected := expectedTime.Format(DATE_FORMAT)
	assert.Equal(expected, formattedRes)
}

func TestTomorrow(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2010-01-01")

	res, err := dp.ParseDate("tomorrow", date)
	formattedRes := res.Format(DATE_FORMAT)

	if err != nil {
		fmt.Println("err is ", err)
		t.Fail()
	}

	expectedTime, _ := time.Parse(DATE_FORMAT, "2010-01-02")
	expected := expectedTime.Format(DATE_FORMAT)
	assert.Equal(expected, formattedRes)
}

// If today is Sunday the 16th, then saying "due:mon" should return Monday the 17th
func TestMondayThisWeek(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2020-08-16")

	res, err := dp.ParseDate("mon", date)
	formattedRes := res.Format(DATE_FORMAT)

	if err != nil {
		fmt.Println("err is ", err)
		t.Fail()
	}

	expectedTime, _ := time.Parse(DATE_FORMAT, "2020-08-17")
	expected := expectedTime.Format(DATE_FORMAT)
	assert.Equal(expected, formattedRes)
}

// If today is Thursday the 20th, then saying "due:mon" should return Monday the 24th
func TestMondayNextWeek(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2020-08-20")

	res, err := dp.ParseDate("mon", date)
	formattedRes := res.Format(DATE_FORMAT)

	if err != nil {
		fmt.Println("err is ", err)
		t.Fail()
	}

	expectedTime, _ := time.Parse(DATE_FORMAT, "2020-08-24")
	expected := expectedTime.Format(DATE_FORMAT)
	assert.Equal(expected, formattedRes)
}

// if specific date is after the pivot date, then keep the same year
func TestSpecificUSDateAfterPivot(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2019-08-11")

	res, err := dp.ParseDate("Aug12", date)

	if err != nil {
		fmt.Println("err is ", err)
		t.Fail()
	}

	formattedRes := res.Format(DATE_FORMAT)
	expectedTime, _ := time.Parse(DATE_FORMAT, "2019-08-12")
	expected := expectedTime.Format(DATE_FORMAT)
	assert.Equal(expected, formattedRes)
}

// if specific date is before the pivot date, then make it next year
func TestSpecificUSDateBeforePivot(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2019-08-20")

	res, err := dp.ParseDate("Aug12", date)

	if err != nil {
		fmt.Println("err is ", err)
		t.Fail()
	}

	formattedRes := res.Format(DATE_FORMAT)
	expectedTime, _ := time.Parse(DATE_FORMAT, "2020-08-12")
	expected := expectedTime.Format(DATE_FORMAT)
	assert.Equal(expected, formattedRes)
}

// if specific date is after the pivot date, then keep the same year
func TestSpecificEUDateAfterPivot(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2019-08-11")

	res, err := dp.ParseDate("12Aug", date)

	if err != nil {
		fmt.Println("err is ", err)
		t.Fail()
	}

	formattedRes := res.Format(DATE_FORMAT)
	expectedTime, _ := time.Parse(DATE_FORMAT, "2019-08-12")
	expected := expectedTime.Format(DATE_FORMAT)
	assert.Equal(expected, formattedRes)
}

// if specific date is before the pivot date, then make it next year
func TestSpecificEUDateBeforePivot(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2019-08-20")

	res, err := dp.ParseDate("12Aug", date)

	if err != nil {
		fmt.Println("err is ", err)
		t.Fail()
	}

	formattedRes := res.Format(DATE_FORMAT)
	expectedTime, _ := time.Parse(DATE_FORMAT, "2020-08-12")
	expected := expectedTime.Format(DATE_FORMAT)
	assert.Equal(expected, formattedRes)
}

func TestParsingError(t *testing.T) {
	assert := assert.New(t)
	dp := &DateParser{}
	date, _ := time.Parse(DATE_FORMAT, "2019-08-20")

	_, err := dp.ParseDate("asdf", date)

	if err == nil {
		fmt.Println("expected an error")
		t.Fail()
	}

	assert.Equal(err.Error(), "Could not parse the date you gave me: 'asdf'")
}
