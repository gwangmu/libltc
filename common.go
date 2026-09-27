package libltc

import (
	"errors"
	"strings"
	"time"
)

//-- Common interfaces

type Stringifiable interface {
	ToString() string
}

//-- struct Name

type Name struct {
	First string
	Middle string
	Last string
}

//-- struct Name: interface Stringifiable

func (this *Name) ToString() string {
	names := []string{}
	if this.First != "" {
		names = append(names, this.First)
	}
	if this.Middle != "" {
		names = append(names, this.Middle)
	}
	if this.Last != "" {
		names = append(names, this.Last)
	}

	return strings.Join(names, " ")
}

//-- struct Name: interface WarningObj

func (this *Name) Summary() string {
	// TODO
	panic("Unimplemented")
}

func (this *Name) IsUnknown() bool {
	// TODO
	panic("Unimplemented")
}

//-- struct Name: interface TOMLPrintable

func (this *Name) PrintTOML() string {
	// TODO
	panic("Unimplemented")
}

//-- struct Time

type Time struct {
	Timezone time.Location
	Attrs map[string][]string

	year *int
	month *int
	day *int
	hour *int
	minute *int
	second *int
}

//-- struct Time: methods (getters and setters)

func (this Time) GetTime() time.Time {
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
		nmonth = time.Month(*this.month)
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

	return time.Date(nyear, nmonth, nday, nhour, nminute, nsecond, 0, &this.Timezone)
}

func (this Time) GetYear() (int, error) {
	if this.year == nil {
		return 0, errors.New("Year not specified")
	} else {
		return *this.year, nil
	}
}

func (this Time) GetMonth() (int, error) {
	if this.month == nil {
		return 0, errors.New("Month not specified")
	} else {
		return *this.month, nil
	}
}

func (this Time) GetDay() (int, error) {
	if this.day == nil {
		return 0, errors.New("Day not specified")
	} else {
		return *this.day, nil
	}
}

func (this Time) GetHour() (int, error) {
	if this.hour == nil {
		return 0, errors.New("Hour not specified")
	} else {
		return *this.hour, nil
	}
}

func (this Time) GetMinute() (int, error) {
	if this.minute == nil {
		return 0, errors.New("Minute not specified")
	} else {
		return *this.minute, nil
	}
}

func (this Time) GetSecond() (int, error) {
	if this.second == nil {
		return 0, errors.New("Second not specified")
	} else {
		return *this.second, nil
	}
}

func (this Time) SetTime(date time.Time, tz bool) {
	nyear := date.Year()
	nmonth := int(date.Month())
	nday := date.Day()
	nhour := date.Hour()
	nminute := date.Minute()
	nsecond := date.Second()

	this.year = &nyear
	this.month = &nmonth
	this.day = &nday
	this.hour = &nhour
	this.minute = &nminute
	this.second = &nsecond

	if tz {
		this.Timezone = *date.Location()
	}
}

func (this Time) UnsetTime() {
	this.year = nil
	this.month= nil
	this.day = nil
	this.hour = nil
	this.minute = nil
	this.second = nil
}

func (this Time) SetYear(v int) {
	this.year = &v
}

func (this Time) UnsetYear() {
	this.year = nil
}

func (this Time) SetMonth(v int) {
	this.month = &v
}

func (this Time) UnsetMonth() {
	this.month = nil
}

func (this Time) SetDay(v int) {
	this.day = &v
}

func (this Time) UnsetDay() {
	this.day = nil
}

func (this Time) SetHour(v int) {
	this.hour = &v
}

func (this Time) UnsetHour() {
	this.hour = nil
}

func (this Time) SetMinute(v int) {
	this.minute = &v
}

func (this Time) UnsetMinute() {
	this.minute = nil
}

func (this Time) SetSecond(v int) {
	this.second = &v
}

func (this Time) UnsetSecond() {
	this.second = nil
}

//-- struct Time: interface Stringifiable

func (this Time) ToString() string {
	// TODO
	panic("Unimplemented")
}

//-- struct Time: interface WarningObj

func (this Time) Summary() string {
	// TODO
	panic("Unimplemented")
}

func (this Time) IsUnknown() bool {
	// TODO
	panic("Unimplemented")
}

//-- struct Time: interface TOMLPrintable

func (this *Time) PrintTOML() string {
	// TODO
	panic("Unimplemented")
}

//-- struct NoteSnippet

type NoteSnippet struct {
	Time Time
	Text string
}

//-- struct NoteSnippet: interface Stringifiable

func (this *NoteSnippet) ToString() string {
	// TODO
	panic("Unimplemented")
}

//-- struct Note

type Note struct {
	title string	// `Title` of normal notes: ignored
	snippets []*NoteSnippet
}

//-- struct Note: methods (getters and setters)

func (this *Note) GetTitle() string {
	return this.title
}

func (this *Note) SetTitle(t string) {
	this.title = t
}

func (this *Note) GetSnippets() []*NoteSnippet {
	if len(this.snippets) == 0 {
		return []*NoteSnippet{ &NoteSnippet{} }
	} else {
		return this.snippets
	}
}

func (this *Note) HasSnippet(s *NoteSnippet) bool {
	for _, elem := range this.snippets {
		if elem == s {
			return true
		}
	}
	return false
}

func (this *Note) AddSnippet(s *NoteSnippet) {
	this.snippets = append(this.snippets, s)
}

func (this *Note) RemoveSnippet(s *NoteSnippet) {
	for i, elem := range this.snippets {
		if elem == s {
			this.snippets = append(this.snippets[:i], this.snippets[i+1:]...)
			return
		}
	}
}

//-- struct Note: interface IStringifiable

func (this *Note) ToString() string {
	// TODO: first snippet -- just Text, others -- ToString()
	panic("Unimplemented")
}

//-- struct Note: interface IWarningObj

func (this *Note) Summary() string {
	// TODO
	panic("Unimplemented")
}

func (this *Note) IsUnknown() bool {
	// TODO
	panic("Unimplemented")
}

//-- struct Note: interface ITOMLPrinter

func (this *Note) PrintTOML() string {
	// TODO
	panic("Unimplemented")
}
