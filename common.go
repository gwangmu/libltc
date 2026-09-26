package libltc

import (
	"errors"
	"time"
)

//-- Common interfaces

type LTCStringifiable interface {
	ToString() string
}

//-- struct LTCName

type LTCName struct {
	First string
	Middle string
	Last string
}

//-- struct LTCName: interface LTCStringifiable

func (this *LTCName) ToString() string {
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

//-- struct LTCName: interface LTCWarningObj

func (this *LTCName) Summary() string {
	// TODO
}

func (this *LTCName) IsUnknown() bool {
	// TODO
}

//-- struct LTCName: interface LTCTOMLPrintable

func (this *LTCName) PrintTOML() string {
	// TODO
}

//-- struct LTCTime

type LTCTime struct {
	Timezone time.Location
	Attrs map[string][]string

	year *int
	month *int
	day *int
	hour *int
	minute *int
	second *int
}

//-- struct LTCTime: methods (getters and setters)

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

//-- struct LTCTime: interface LTCStringifiable

func (this LTCTime) ToString() string {
	// TODO
}

//-- struct LTCTime: interface LTCWarningObj

func (this LTCTime) Summary() string {
	// TODO
}

func (this LTCTime) IsUnknown() bool {
	// TODO
}

//-- struct LTCTime: interface LTCTOMLPrintable

func (this *LTCTime) PrintTOML() string {
	// TODO
}

//-- struct LTCNoteSnippet

type LTCNoteSnippet struct {
	Time LTCTime
	Text string
}

//-- struct LTCNoteSnippet: interface LTCStringifiable

func (this *LTCNoteSnippet) ToString() string {
	// TODO
}

//-- struct LTCNote

type LTCNote struct {
	title string	// `Title` of normal notes: ignored
	snippets []*LTCNoteSnippet
}

//-- struct LTCNote: methods (getters and setters)

func (this *LTCNote) GetTitle() string {
	return this.title
}

func (this *LTCNote) SetTitle(t string) {
	this.title = t
}

func (this *LTCNote) GetSnippets() []*LTCNoteSnippet {
	if len(this.snippets) == 0 {
		return []*LTCNoteSnippet{ &LTCNoteSnippet{} }
	} else
		return this.snippets
	}
}

func (this *LTCNote) HasSnippet(s *LTCNoteSnippet) bool {
	for i, elem := range this.snippets {
		if elem == s {
			return true
		}
	}
	return false
}

func (this *LTCNote) AddSnippet(s *LTCNoteSnippet) {
	this.snippets = append(this.snippets, s)
}

func (this *LTCNote) RemoveSnippet(s *LTCNoteSnippet) {
	for i, elem := range this.snippets {
		if elem == s {
			this.snippets = append(this.snippets[:i], this.snippets[i+1:]...)
			return
		}
	}
}

//-- struct LTCNote: interface LTCStringifiable

func (this *LTCNote) ToString() string {
	// TODO: first snippet -- just Text, others -- ToString()
}

//-- struct LTCNote: interface LTCWarningObj

func (this *LTCNote) Summary() string {
	// TODO
}

func (this *LTCNote) IsUnknown() bool {
	// TODO
}

//-- struct LTCNote: interface LTCTOMLPrintable

func (this *LTCNote) PrintTOML() string {
	// TODO
}
