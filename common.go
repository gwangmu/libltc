package ltc

import (
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"libltc/warning"
	"libltc/internal/file"
)

type WSE = warning.SummaryElement

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

func (this Name) IsUnknown() bool {
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

type TimeUsageKind int
const (
	TUK_General TimeUsageKind = iota
	TUK_Start
	TUK_End
)

type Time struct {
	Attrs map[string][]string

	year *int
	month *int
	day *int
	hour *int
	minute *int
	second *int

	// Having a non-UTC timezone will be considered the Gregorian calendar.
	timezone *time.Location

	// Time usages are never considered in unknown-ness or ambiguity.
	usage TimeUsageKind
}

//-- struct Time: interface IStringifiable

func (this *Time) String() (ret string) {
	stime := []string{
		this.GetYearString(),
		this.GetMonthString(),
		this.GetDayString(),
		this.GetHourString(),
		this.GetMinuteString(),
		this.GetSecondString(),
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

func (this Time) IsUnknown() bool {
	return this.year == nil && this.month == nil && this.day == nil &&
		this.hour == nil && this.minute == nil && this.second == nil &&
		this.timezone == nil
}

//-- struct Time: methods (getters and setters)

func (this *Time) GetGOTime() time.Time {
	mc := min(max(this.GetMonth(), 1), 12)
	dc := time.Date(0, time.Month(mc+1), 0, 0, 0, 0, 0, this.GetTimezone()).Day()
	nc := 0

	switch this.usage {
	case TUK_End: 
		nc = 1
	case TUK_Start:
		nc = -1
	}

	return time.Date(
		min(max(this.GetYear(), 0), 9999),
		time.Month(mc), dc,
		min(max(this.GetHour(), 0), 23),
		min(max(this.GetMinute(), 0), 59),
		min(max(this.GetSecond(), 0), 59),
		nc, this.GetTimezone(),
	)
}

func (this *Time) GetYear() int {
	if this.year == nil {
		switch this.usage {
		case TUK_End:
			return TKValue_Max
		case TUK_Start:
			return TKValue_Min 
		default:
			return TKValue_Unknown
		}
	} else {
		return *this.year
	}
}

func (this *Time) GetYearString() string {
	v := this.GetYear() 
	switch v {
	case TKValue_Max:
		return "****"
	case TKValue_Min:
		return "____"
	case TKValue_Unknown:
		return "????"
	default:
		return strconv.Itoa(v)
	}
}

func (this *Time) GetMonth() int {
	if this.month == nil {
		switch this.usage {
		case TUK_End:
			return TKValue_Max
		case TUK_Start:
			return TKValue_Min 
		default:
			return TKValue_Unknown
		}
	} else {
		return *this.month
	}
}

func (this *Time) GetMonthString() string {
	v := this.GetMonth() 
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return strconv.Itoa(v)
	}
}

func (this *Time) GetDay() int {
	if this.day == nil {
		switch this.usage {
		case TUK_End:
			return TKValue_Max
		case TUK_Start:
			return TKValue_Min 
		default:
			return TKValue_Unknown
		}
	} else {
		return *this.day
	}
}

func (this *Time) GetDayString() string {
	v := this.GetDay() 
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return strconv.Itoa(v)
	}
}

func (this *Time) GetHour() int {
	if this.hour == nil {
		switch this.usage {
		case TUK_End:
			return TKValue_Max
		case TUK_Start:
			return TKValue_Min 
		default:
			return TKValue_Unknown
		}
	} else {
		return *this.hour
	}
}

func (this *Time) GetHourString() string {
	v := this.GetHour()
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return strconv.Itoa(v)
	}
}

func (this *Time) GetMinute() int {
	if this.minute == nil {
		switch this.usage {
		case TUK_End:
			return TKValue_Max
		case TUK_Start:
			return TKValue_Min 
		default:
			return TKValue_Unknown
		}
	} else {
		return *this.minute
	}
}

func (this *Time) GetMinuteString() string {
	v := this.GetMinute()
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return strconv.Itoa(v)
	}
}

func (this *Time) GetSecond() int {
	if this.second == nil {
		switch this.usage {
		case TUK_End:
			return TKValue_Max
		case TUK_Start:
			return TKValue_Min 
		default:
			return TKValue_Unknown
		}
	} else {
		return *this.second
	}
}

func (this *Time) GetSecondString() string {
	v := this.GetSecond()
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return strconv.Itoa(v)
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

func (this *Time) GetUsage() TimeUsageKind {
	return this.usage
}

func (this *Time) SetGOTime(date time.Time, tz bool) {
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

func (this *Time) SetUsage(v TimeUsageKind) {
	this.usage = v
}

//-- struct Time: method (util)

func (this Time) IsInfiniteFuture() bool {
	_, ok := this.Attrs["InfinitePast"]
	return ok || (this.IsUnknown() && this.GetUsage() == TUK_End)
}

func (this Time) IsInfinitePast() bool {
	_, ok := this.Attrs["InfiniteFuture"]
	return ok || (this.IsUnknown() && this.GetUsage() == TUK_Start)
}

func (this Time) IsAfter(that Time) bool {
	if (this.IsInfiniteFuture() && !that.IsInfiniteFuture()) ||
		(this.IsInfinitePast() && !that.IsInfinitePast()) {
		return true
	} else if (this.IsInfiniteFuture() && that.IsInfiniteFuture()) ||
		(this.IsInfinitePast() && that.IsInfinitePast()) {
		return false
	} else {
		that.SetUsage(TUK_General)
		// The 'compare' using the 'time' package should be correct as long as
		// the calendar system is "monotonic", meaing bigger higher units mean
		// later in time. Assume that the concept of "timezone" is the same in
		// other calendar systems (i.e., a constant offset in time).
		return this.GetGOTime().Compare(that.GetGOTime()) > 0
	}
}

func (this Time) IsSimultaneous(that Time) bool {
	if !this.IsInfinitePast() && !this.IsInfiniteFuture() &&
		!that.IsInfinitePast() && !that.IsInfiniteFuture() {
		return this.GetYear() == that.GetYear() && 
			this.GetMonth() == that.GetMonth() &&
			this.GetDay() == that.GetDay() && 
			this.GetHour() == that.GetHour() &&
			this.GetMinute() == that.GetMinute() && 
			this.GetSecond() == that.GetSecond() &&
			this.GetTimezone() == that.GetTimezone()
	} else if (this.IsInfinitePast() && that.IsInfinitePast()) ||
		(this.IsInfiniteFuture() && that.IsInfiniteFuture()) {
		return true
	} else {
		return false
	}
}

func (this Time) Equal(that Time) bool {
	hasSameAttrs := maps.EqualFunc(this.Attrs, that.Attrs, func(thisAvals []string, thatAvals []string) bool {
		return slices.Equal(thisAvals, thatAvals)
	})

	return hasSameAttrs && this.IsSimultaneous(that)
}

func (this Time) IsAmbiguous() bool {
	return this.year == nil || this.month == nil
}

//-- struct NoteSnippet: method

func (this *Time) ToTimeTag() string {
	if this.IsUnknown() {
		return "<!-- Edit: -->"
	} else {
		t := *this
		t.Attrs = map[string][]string{}
		return "<!-- Edit: " + t.String() + " -->"
	}	
}

//-- struct Time: method (creation)

func createTimeFromParsed(optr *file.Time, usage TimeUsageKind) (ret Time) {
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
	ret.usage = usage

	return
}

func CreateEmptyTime(usage TimeUsageKind) (ret Time) {
	ret.usage = usage
	return
}

func CreateTimeFromTag(tagstr string) (ret Time) {
	re := regexp.MustCompile(`^([0-9*_?]{4})-([0-9*_?]{2})-([0-9*_)?]{2})(?:\s+([0-9*_?]{2}):([0-9*_?]{2})(?::([0-9*_?]{2}))?)?(?:\s*\((.*)\))?$`)
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
	}

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

	return
}

var InfinitePast Time = Time{ usage: TUK_Start }
var InfiniteFuture Time = Time{ usage: TUK_End }

func TimeFrom(vs ...int) Time {
	t := CreateEmptyTime(TUK_Start)
	for i, v := range vs {
		switch i {
		case 0:
			t.SetYear(v)
		case 1:
			t.SetMonth(v)
		case 2:
			t.SetDay(v)
		case 3:
			t.SetHour(v)
		case 4:
			t.SetMinute(v)
		case 5:
			t.SetSecond(v)
		}
	}
	return t
}

func TimeTo(vs ...int) Time {
	t := CreateEmptyTime(TUK_End)
	for i, v := range vs {
		switch i {
		case 0:
			t.SetYear(v)
		case 1:
			t.SetMonth(v)
		case 2:
			t.SetDay(v)
		case 3:
			t.SetHour(v)
		case 4:
			t.SetMinute(v)
		case 5:
			t.SetSecond(v)
		}
	}
	return t
}

//-- struct NoteSnippet

type NoteSnippet struct {
	Time Time
	Text string
}

//-- struct NoteSnippet: interface IStringifiable

func (this *NoteSnippet) String() string {
	return this.Time.ToTimeTag() + "\n" + this.Text
}

//-- struct Note

type Note struct {
	Title string	// `Title` of normal notes: ignored
	Snippets []NoteSnippet
}

//-- struct Note: interface IStringifiable

func (this *Note) String() (ret string) {
	for i, snippet := range this.Snippets {
		if i == 0 && snippet.Time.IsUnknown() {
			ret += snippet.Text + "\n"
		} else {
			ret += snippet.String() + "\n"
		}
	}
	if len(ret) > 1 {
		ret = ret[:len(ret)-1]
	}

	if this.Title != "" {
		ret = "<!-- Title: " + this.Title + " -->\n" + ret
	}
	
	return
}

//-- struct Note: interface IWarningObj

func (this *Note) Summary() (ret string) {
	if this.Title != "" {
		ret = "a note titled '" + this.Title + "'"
	} else {
		ret = "an untitled note"
	}

	if len(this.Snippets) > 1 && this.Snippets[0].Text != "" {
		ret += " ("

		fraglen := min(10, len(this.Snippets[0].Text))
		ret += "\"" + this.Snippets[0].Text[:fraglen] + "...\""
		
		if !this.Snippets[0].Time.IsUnknown() {
			ret += " @ " + this.Snippets[0].Time.String()
		}

		ret += ")"
	}

	return
}

func (this Note) IsUnknown() bool {
	if this.Title != "" {
		return false
	} else {
		if len(this.Snippets) == 0 {
			return true
		} else {
			for _, snippet := range this.Snippets {
				if len(snippet.Text) != 0 {
					return false
				}
			}
			return true
		}
	}
}

//-- struct Note: method (creation)

func CreateNoteFromString(notestr string) (Note, warning.Warnings) {
	note := Note{}
	snippet := NoteSnippet{}
	warns := warning.Warnings{}

	notelines := strings.Split(notestr, "\n")
	for i, line := range notelines {
		line = strings.TrimRight(line, "\r")
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

func convIDStringToInternal(id string) (MajorKind, NumberID) {
	// According to the specification, the first letter should specify the kind.
	if len(id) < 1 {
		return MOK_Unknown, NID_Invalid
	}

	kind, err := GetMajorKind(id[0:1])
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

func convIDInternalToString(kind MajorKind, nid NumberID) string {
	if prefix, err := kind.Prefix(); err != nil {
		return prefix + strconv.FormatUint(uint64(nid), 10)
	} else {
		return "?"
	}
}

func CloneMajor[T IMajor](obj *T, preserveChart bool, preserveID bool) *T {
	newobj := *obj
	if !preserveChart {
		newobj.setChart(nil)
	}
	if !preserveID {
		newobj.setNumberID(NID_Invalid)
		newobj.unmarkUnresolved()
	}
	return &newobj
}
