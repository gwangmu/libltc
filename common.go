package libltc

import (
	"errors"
	"time"
)


type LTCName struct {
	First string
	Middle string
	Last string
}

func (this LTCName) String() string {
	names = []string{}
	if this.First != "" {
		names = append(names, this.First)
	}
	if this.Middle != "" {
		names = append(names, this.Middle)
	}
	if this.Last != "" {
		names = append(names, this.Last)
	}

	return string.Join(names, " ")
}


type LTCTime struct {
	Timezone time.Location

	year *int
	month *int
	day *int
	hour *int
	minute *int
	second *int
}

func (this LTCTime) GetTime() time.Time {
	var nyear, nday, nhour, nminute, nsecond int
	var nmonth time.Month

	if this.year == nil {
		nyear = 0
	} else {
		nyear = *this.year
	}

	if this.month == nil {
		nmonth = time.January
	} else {
		nmonth = *this.month
	}

	if this.day == nil {
		nday = 1
	} else {
		nday = *this.day
	}

	if this.hour == nil {
		nhour = 0 
	} else {
		nhour = *this.hour
	}

	if this.minute == nil {
		nminute = 0
	} else {
		nminute = *this.minute
	}

	if this.second == nil {
		nsecond = 0
	} else {
		nsecond = *this.second
	}

	return time.Date(nyear, nmonth, nday, nhour, nminute, nsecond, this.Timezone)
}

func (this LTCTime) GetYear() (int, error) {
	if this.year == nil {
		return 0, errors.New("Year not specified")
	} else {
		return this.year, nil
	}
}

func (this LTCTime) GetMonth() (int, error) {
	if this.month == nil {
		return 0, errors.New("Month not specified")
	} else {
		return this.month, nil
	}
}

func (this LTCTime) GetDay() (int, error) {
	if this.day == nil {
		return 0, errors.New("Day not specified")
	} else {
		return this.day, nil
	}
}

func (this LTCTime) GetHour() (int, error) {
	if this.hour == nil {
		return 0, errors.New("Hour not specified")
	} else {
		return this.hour, nil
	}
}

func (this LTCTime) GetMinute() (int, error) {
	if this.minute == nil {
		return 0, errors.New("Minute not specified")
	} else {
		return this.minute, nil
	}
}

func (this LTCTime) GetSecond() (int, error) {
	if this.second == nil {
		return 0, errors.New("Second not specified")
	} else {
		return this.second, nil
	}
}

func (this LTCTime) SetTime(date time.Time, tz bool) {
	this.year = &data.Year()
	this.month = &int(data.Month())
	this.day = &data.Day()
	this.hour = &data.Hour()
	this.minute = &data.Minute()
	this.second = &data.Second()

	if tz {
		this.Timezone = *tz.Location()
	}
}

func (this LTCTime) UnsetTime() {
	this.year = nil
	this.month= nil
	this.day = nil
	this.hour = nil
	this.minute = nil
	this.second = nil
}

func (this LTCTime) SetYear(v int) {
	this.year = &v
}

func (this LTCTime) UnsetYear() {
	this.year = nil
}

func (this LTCTime) SetMonth(v int) {
	this.month = &v
}

func (this LTCTime) UnsetMonth() {
	this.month = nil
}

func (this LTCTime) SetDay(v int) {
	this.day = &v
}

func (this LTCTime) UnsetDay() {
	this.day = nil
}

func (this LTCTime) SetHour(v int) {
	this.hour = &v
}

func (this LTCTime) UnsetHour() {
	this.hour = nil
}

func (this LTCTime) SetMinute(v int) {
	this.minute = &v
}

func (this LTCTime) UnsetMinute() {
	this.minute = nil
}

func (this LTCTime) SetSecond(v int) {
	this.second = &v
}

func (this LTCTime) UnsetSecond() {
	this.second = nil
}


type LTCNoteSnippet struct {
	Time LTCTime
	Text string
}

type LTCNote struct {
	Title string
	Snippets []LTCNOteSnippet
}
