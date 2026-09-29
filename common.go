package libltc

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

//-- Common interfaces

type IStringifiable interface {
	ToString() string
}

//-- struct Name

type Name struct {
	First string
	Middle string
	Last string
	Attrs map[string][]string
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

//-- struct Name : method (creation)

func createNameFromParsed(o file.Name) Name {
	return Name{
		First: o.First,
		Middle: o.Middle,
		Last: o.Last,
		Attrs: convAttrsFileToChart(o.Attrs),
	}
}

//-- struct Time

type Time struct {
	Timezone *time.Location
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

	return time.Date(nyear, nmonth, nday, nhour, nminute, nsecond, 0, this.Timezone)
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
		this.Timezone = date.Location()
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

//-- struct Time: method (creation)

func createTimeFromParsed(optr *file.Time) (ret Time) {
	if optr == nil {
		return
	} 
	
	o := *optr

	if o.Year != nil {
		nyear := *o.Year
		ret.year = &nyear
	}

	if o.Month != nil {
		nmonth := *o.Month
		ret.month = &nmonth
	}
	
	if o.Day != nil {
		nday := *o.Day
		ret.day = &nday
	}
	
	if o.Hour != nil {
		nhour := *o.Hour
		ret.hour = &nhour
	}
	
	if o.Minute != nil {
		nminute := *o.Minute
		ret.minute = &nminute
	}
	
	if o.Second != nil {
		nsecond := *o.Second
		ret.second = &nsecond
	}

	if o.Timezone != nil {
		loc, tzerr := time.LoadLocation(*o.Timezone)

		if tzerr == nil {
			ret.Timezone = loc
		} else {
			utcloc, tzerr := time.LoadLocation("UTC")
			if tzerr == nil {
				ret.Timezone = utcloc
			} else {
				panic("UTC not loaded")
			}
		}
	}

	ret.Attrs = convAttrsFileToChart(o.Attrs)

	return
}

func CreateTimeFromTag(tagstr string) (ret Time) {
	re := regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})(?:\s+(\d{2}):(\d{2})(?::(\d{2}))?)?(?:\s*\((.*)\))?$`)
	matches := re.FindStringSubmatch(tagstr)

	if len(matches) < 8 {
		panic("Time tag matching unexpectedly too short")
	}

	nyear, yerr := strconv.Atoi(matches[1])
	nmonth, merr := strconv.Atoi(matches[2])
	nday, derr := strconv.Atoi(matches[3])

	if yerr == nil && merr == nil && derr == nil {
		ret.SetYear(nyear)
		ret.SetMonth(nmonth)
		ret.SetDay(nday)

		nhour, herr := strconv.Atoi(matches[4])
		nminute, mmerr := strconv.Atoi(matches[5])
		nsecond, serr := strconv.Atoi(matches[6])

		if herr == nil && mmerr == nil {
			ret.SetHour(nhour)
			ret.SetMinute(nminute)

			if serr == nil {
				ret.SetSecond(nsecond)
			} else {
				ret.SetSecond(0)
			}
		}

		loc, tzerr := time.LoadLocation(matches[7])

		if tzerr == nil {
			ret.Timezone = loc
		} else {
			utcloc, tzerr := time.LoadLocation("UTC")
			if tzerr == nil {
				ret.Timezone = utcloc
			} else {
				panic("UTC not loaded")
			}
		}
	}

	return
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
	Title string	// `Title` of normal notes: ignored
	Snippets []NoteSnippet
}

//-- struct Note: interface IStringifiable

func (this *Note) ToString() string {
	// TODO: title (if non-empty), first snippet -- just Text, others -- ToString()
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

//-- struct Note: interface ITOMLPrintable

func (this *Note) PrintTOML() string {
	// TODO
	panic("Unimplemented")
}

//-- struct Note: method (creation)

func CreateNoteFromString(notestr string) (Note, warning.Warnings) {
	note := Note{}
	snippet := NoteSnippet{}
	warns := warning.Warnings{}

	notelines := strings.Split(notestr, "\n")
	for i, line := range notelines {
		if i == 0 {
			re := regexp.MustCompile(`^\s*<!--\s*Title:(.*)-->\s*$`)
			if matches := re.FindStringSubmatch(line); matches != nil {
				note.Title = strings.TrimSpace(matches[1])
			}
		}

		re := regexp.MustCompile(`^\s*<!--\s*Edit:(.*)-->\s*$`)
		if matches := re.FindStringSubmatch(line); matches != nil {
			snippet.Text = strings.TrimRight(snippet.Text, "\r\n")
			note.Snippets = append(note.Snippets, snippet)

			snippet = NoteSnippet{}
			snippet.Time = CreateTimeFromTag(strings.TrimSpace(matches[1]))

			if snippet.Time.IsUnknown() {
				warns.Add("Unknown time tag in @0@.", &note)
			}
		}
	}
	note.Snippets = append(note.Snippets, snippet)

	return note, warns
}

//-- method (utils)

func convAttrsFileToChart(attrs []string) (cattrs map[string][]string) {
	for _, attr := range attrs {
		key, value, _ := strings.Cut(attr, ":")
		if _, ok := cattrs[key]; !ok {
			cattrs[key] = []string{}
		}
		cattrs[key] = append(cattrs[key], value)
	}
	return
}

func convAttrsChartToFile(attrs map[string][]string) (fattrs []string) {
	for akey, avals := range attrs {
		for _, aval := range avals {
			if aval != "" {
				fattrs = append(fattrs, akey + ":" + aval)
			} else {
				fattrs = append(fattrs, akey)
			}
		}
	}
	return
}
