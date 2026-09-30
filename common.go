package libltc

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

//-- Common interfaces

type IStringifiable interface {
	String() string
}

//-- struct Name

type Name struct {
	First string
	Middle string
	Last string
	Attrs map[string][]string
}

//-- struct Name: interface IStringifiable

func (this *Name) String() string {
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
	strname := this.String()
	if len(this.Attrs) != 0 {
		addname := strings.Join(convAttrsChartToFile(this.Attrs), ", ")
		strname += " (" + addname + ")"
	}
	return strname
}

func (this *Name) IsUnknown() bool {
	return this.First == "" && this.Middle == "" && this.Last == "" &&
		len(this.Attrs) == 0
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
	Attrs map[string][]string

	year *int
	month *int
	day *int
	hour *int
	minute *int
	second *int
	timezone *time.Location
}

//-- struct Time: interface IStringifiable

func (this *Time) String() (ret string) {
	stime := []string{
		strconv.Itoa(this.GetYear()),
		strconv.Itoa(this.GetMonth()),
		strconv.Itoa(this.GetDay()),
		strconv.Itoa(this.GetHour()),
		strconv.Itoa(this.GetMinute()),
		strconv.Itoa(this.GetSecond()),
	}

	sextra := []string{}
	if this.timezone != nil {
		sextra = append(sextra, this.timezone.String())
	}
	if len(this.Attrs) != 0 {
		sextra = append(sextra, convAttrsChartToFile(this.Attrs)...)
	}

	ret += strings.Join(stime[0:3], "-")
	if this.hour != nil || this.minute != nil || this.second == nil {
		ret += " " + strings.Join(stime[3:6], ":")
	}
	if len(sextra) != 0 {
		ret += " (" + strings.Join(sextra, ", ")
	}

	return
}

//-- struct Time: interface WarningObj

func (this *Time) Summary() string {
	return this.String()
}

func (this *Time) IsUnknown() bool {
	return this.year == nil && this.month == nil && this.day == nil &&
		this.hour == nil && this.minute == nil && this.second == nil &&
		this.timezone == nil
}

//-- struct Time: methods (getters and setters)

func (this *Time) GetTime() time.Time {
	return time.Date(
		this.GetYear(),
		time.Month(this.GetMonth()),
		this.GetDay(),
		this.GetHour(),
		this.GetMinute(),
		this.GetSecond(),
		0, 
		this.GetTimezone(),
	)
}

func (this *Time) GetYear() int {
	if this.year == nil {
		return 1 
	} else {
		return *this.year
	}
}

func (this *Time) GetMonth() int {
	if this.month == nil {
		return 1
	} else {
		return *this.month
	}
}

func (this *Time) GetDay() int {
	if this.day == nil {
		return 1
	} else {
		return *this.day
	}
}

func (this *Time) GetHour() int {
	if this.hour == nil {
		return 0
	} else {
		return *this.hour
	}
}

func (this *Time) GetMinute() int {
	if this.minute == nil {
		return 0
	} else {
		return *this.minute
	}
}

func (this *Time) GetSecond() int {
	if this.second == nil {
		return 0
	} else {
		return *this.second
	}
}

func (this *Time) GetTimezone() *time.Location {
	if this.timezone == nil {
		utcloc, err := time.LoadLocation("UTC")
		if err == nil {
			return utcloc
		} else {
			panic("UTC not loaded")
		}
	} else {
		return this.timezone
	}
}

func (this *Time) SetTime(date time.Time, tz bool) {
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
		this.timezone = date.Location()
	}
}

func (this *Time) UnsetTime() {
	this.year = nil
	this.month= nil
	this.day = nil
	this.hour = nil
	this.minute = nil
	this.second = nil
}

func (this *Time) SetYear(v int) {
	this.year = &v
}

func (this *Time) UnsetYear() {
	this.year = nil
}

func (this *Time) SetMonth(v int) {
	this.month = &v
}

func (this *Time) UnsetMonth() {
	this.month = nil
}

func (this *Time) SetDay(v int) {
	this.day = &v
}

func (this *Time) UnsetDay() {
	this.day = nil
}

func (this *Time) SetHour(v int) {
	this.hour = &v
}

func (this *Time) UnsetHour() {
	this.hour = nil
}

func (this *Time) SetMinute(v int) {
	this.minute = &v
}

func (this *Time) UnsetMinute() {
	this.minute = nil
}

func (this *Time) SetSecond(v int) {
	this.second = &v
}

func (this *Time) UnsetSecond() {
	this.second = nil
}

func (this *Time) SetTimezone(v *time.Location) {
	this.timezone = v
}

//-- struct Time: method (util)

func (this *Time) IsAfter(that *Time) bool {
	// The 'compare' using the 'time' package should be correct as long as
	// the calendar system is "monotonic", meaing bigger higher units mean
	// later in time. Assume that the concept of "timezone" is the same in
	// other calendar systems (i.e., a constant offset in time).
	return this.GetTime().Compare(that.GetTime()) > 0
}

func (this *Time) IsAmbiguous() bool {
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
			ret.timezone = loc
		} else {
			utcloc, tzerr := time.LoadLocation("UTC")
			if tzerr == nil {
				ret.timezone = utcloc
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
			ret.timezone = loc
		} else {
			utcloc, tzerr := time.LoadLocation("UTC")
			if tzerr == nil {
				ret.timezone = utcloc
			} else {
				panic("UTC not loaded")
			}
		}
	}

	return
}

//-- struct NoteSnippet

type NoteSnippet struct {
	Time Time
	Text string
}

//-- struct NoteSnippet: interface IStringifiable

func (this *NoteSnippet) String() string {
	// TODO
	panic("Unimplemented")
}

//-- struct Note

type Note struct {
	Title string	// `Title` of normal notes: ignored
	Snippets []NoteSnippet
}

//-- struct Note: interface IStringifiable

func (this *Note) String() string {
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

func convIDStringToInternal(id string) (MainObjKind, NumberID) {
	// According to the specification, the first letter should specify the kind.
	if len(id) < 1 {
		return MOK_Unknown, NID_Invalid
	}

	kind, err := GetMainObjKind(id[0:1])
	if err != nil {
		return MOK_Unknown, NID_Invalid
	} else if len(id) < 2 {
		return kind, NID_Invalid
	}

	nid, err := strconv.ParseUint(id[1:], 10, 64)
	if err != nil {
		return kind, NID_Invalid
	}

	return kind, NumberID(nid)
}

func convIDInternalToString(kind MainObjKind, nid NumberID) string {
	if prefix, err := kind.Prefix(); err != nil {
		return prefix + strconv.FormatUint(uint64(nid), 10)
	} else {
		return "?"
	}
}

func CloneMainObject[T IMainObj](obj *T, preserveChart bool, preserveID bool) *T {
	newobj := *obj
	if !preserveChart {
		newobj.setChart(nil)
	}
	if !preserveID {
		newobj.setNumberID(NID_Invalid)
		newobj.setFullID("")
	}
	return &newobj
}
