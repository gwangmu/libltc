package libltc

import (
	"fmt"
	"maps"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

type wse = warning.SummaryElement

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

func (this Name) Summary() string {
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

func (this Time) String() (ret string) {
	if this.IsUnknown() {
		if this.IsInfiniteFuture() {
			return "(inf)"
		} else if this.IsInfinitePast() {
			return "(-inf)"
		}
		return "(unknown)"
	}

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
	if this.hour != nil || this.minute != nil || this.second != nil {
		ret += " " + strings.Join(stime[3:6], ":")
	}
	if len(sextra) != 0 {
		ret += " (" + strings.Join(sextra, ", ") + ")"
	}

	return
}

//-- struct Time: interface WarningObj

func (this Time) Summary() string {
	return this.String()
}

func (this Time) IsUnknown() bool {
	return this.year == nil && this.month == nil && this.day == nil &&
		this.hour == nil && this.minute == nil && this.second == nil &&
		this.timezone == nil
}

//-- struct Time: methods (getters and setters)

func (this Time) GetGOTime() time.Time {
	yc := this.GetYear()

	mc := this.GetMonth() 
	switch mc {
	case TKValue_Max:
		// FIXME: currently using the Gragorian calendar max.
		mc = 12
	case TKValue_Min:
		mc = 1
	}

	dc := this.GetDay()
	switch dc {
	case TKValue_Max:
		// FIXME: currently using the Gragorian calendar max.
		dc = time.Date(0, time.Month(mc+1), 0, 0, 0, 0, 0, this.GetTimezone()).Day()
	case TKValue_Min:
		dc = 1
	}

	hc := this.GetHour()
	switch hc {
	case TKValue_Max:
		hc = 23
	case TKValue_Min:
		hc = 0
	}

	mmc := this.GetMinute()
	switch mmc {
	case TKValue_Max:
		mmc = 59
	case TKValue_Min:
		mmc = 0
	}

	sc := this.GetSecond()
	switch sc {
	case TKValue_Max:
		sc = 59
	case TKValue_Min:
		sc = 0
	}

	nc := 0
	switch this.usage {
	case TUK_End: 
		nc = 1
	case TUK_Start:
		nc = -1
	}

	return time.Date(yc, time.Month(mc), dc, hc, mmc, sc, nc, this.GetTimezone())
}

func (this Time) GetYear() int {
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

func (this Time) GetYearString() string {
	v := this.GetYear() 
	switch v {
	case TKValue_Max:
		return "****"
	case TKValue_Min:
		return "____"
	case TKValue_Unknown:
		return "????"
	default:
		return fmt.Sprintf("%04d", v)
	}
}

func (this Time) GetMonth() int {
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

func (this Time) GetMonthString() string {
	v := this.GetMonth() 
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return fmt.Sprintf("%02d", v)
	}
}

func (this Time) GetDay() int {
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

func (this Time) GetDayString() string {
	v := this.GetDay() 
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return fmt.Sprintf("%02d", v)
	}
}

func (this Time) GetHour() int {
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

func (this Time) GetHourString() string {
	v := this.GetHour()
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return fmt.Sprintf("%02d", v)
	}
}

func (this Time) GetMinute() int {
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

func (this Time) GetMinuteString() string {
	v := this.GetMinute()
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return fmt.Sprintf("%02d", v)
	}
}

func (this Time) GetSecond() int {
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

func (this Time) GetSecondString() string {
	v := this.GetSecond()
	switch v {
	case TKValue_Max:
		return "**"
	case TKValue_Min:
		return "__"
	case TKValue_Unknown:
		return "??"
	default:
		return fmt.Sprintf("%02d", v)
	}
}

func (this Time) GetTimezone() *time.Location {
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

func (this Time) GetUsage() TimeUsageKind {
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
	_, ok := this.Attrs["InfiniteFuture"]
	return ok && this.IsUnknown()
}

func (this Time) IsInfinitePast() bool {
	_, ok := this.Attrs["InfinitePast"]
	return ok && this.IsUnknown()
}

func (this Time) IsInfinite() bool {
	return this.IsInfiniteFuture() || this.IsInfinitePast()
}

func (this Time) isFloating() bool {
	return this.IsUnknown() && !this.IsInfinite()
}

// Return value is undefined if either of the times is unknown (non-inf) or ambiguous.
// Infinite futures are after themselves. Same for infinite pasts.
func (this Time) IsAfter(that Time) bool {
	if this.IsInfiniteFuture() || that.IsInfinitePast() {
		return true
	} else if this.IsInfinitePast() || that.IsInfiniteFuture() {
		return false 
	} else {
		if this.usage == that.usage {
			// Give 'this' a little more edge.
			that.SetUsage(TUK_General)
		}
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

// NOTE: infinite times are both unknown and ambiguous.
func (this Time) IsAmbiguous() bool {
	return this.year == nil || this.month == nil
}

// Limitation 1: this method assumes Gregorian calendars.
// Supporting non-Gregorian calendars is future work.
// Limitation 2: the timezone of RHS is ignored. (i.e., UTC assumed)
// Limitation 3: Adding infinity or adding to infinity result in infinity.
// Adding infinite future to infinite past is undefined.
// Hyperreal-style infinity sum is maybe useful, but it's also future work.
// The returning time's usage follows that of LHS (i.e., 'this').
func (this Time) Add(that Time) Time {
	if this.isFloating() || that.isFloating() {
		ret := UnknownTime
		ret.SetUsage(this.GetUsage())
		return ret
	} else if this.IsInfinitePast() || that.IsInfinitePast() {
		ret := InfinitePast
		ret.SetUsage(this.GetUsage())
		return ret
	} else if this.IsInfiniteFuture() || that.IsInfiniteFuture() {
		ret := InfiniteFuture
		ret.SetUsage(this.GetUsage())
		return ret
	} else {
		that.SetUsage(TUK_General)
		that.SetTimezone(nil)
		
		// Convert them to GO times.
		got1 := this.GetGOTime()
		got2 := that.GetGOTime()

		// Add years, months, and days.
		// Those are too big, so they may cause an overflow.
		got1.AddDate(got2.Year(), int(got2.Month()), got2.Day())

		// Add smaller bits using duration.
		got2 = time.Date(got1.Year(), got1.Month(), got1.Day(),
			got2.Hour(), got2.Minute(), got2.Second(), 0, got1.Location())
		dur := got2.Sub(got1)
		got1 = got1.Add(dur)

		// Convert back to LTC time.
		ret := this
		ret.SetGOTime(got1, true)
		return ret
	}	
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

	if !slices.Contains(o.Attrs, "InfiniteFuture") &&
		!slices.Contains(o.Attrs, "InfinitePast") {
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
	}

	ret.Attrs = convAttrsFileToChart(o.Attrs)
	ret.usage = usage

	return
}

func CreateEmptyTime(usage TimeUsageKind) (ret Time) {
	ret.usage = usage
	return
}

func CreateTimeFromTag(tagstr string) Time {
	ret := CreateEmptyTime(TUK_General)
	re := regexp.MustCompile(`^([0-9*_?]{4})(?:-([0-9*_?]{2}))?(?:-([0-9*_)?]{2}))?(?:\s+([0-9*_?]{2}):([0-9*_?]{2})(?::([0-9*_?]{2}))?)?(?:\s*\((.*)\))?`)
	matches := re.FindStringSubmatch(tagstr)

	if matches == nil {
		matches = make([]string, 8)
	}

	if len(matches) < 8 {
		panic("Time tag matching unexpectedly too short")
	}

	if nyear, yerr := strconv.Atoi(matches[1]); yerr == nil {
		ret.SetYear(nyear)
		if nmonth, merr := strconv.Atoi(matches[2]); merr == nil {
			ret.SetMonth(nmonth)
			if nday, derr := strconv.Atoi(matches[3]); derr == nil {
				ret.SetDay(nday)
				if nhour, herr := strconv.Atoi(matches[4]); herr == nil {
					ret.SetHour(nhour)
					if nminute, mmerr := strconv.Atoi(matches[5]); mmerr == nil {
						ret.SetMinute(nminute)
						if nsecond, serr := strconv.Atoi(matches[6]); serr == nil {
							ret.SetSecond(nsecond)
						} else {
							ret.SetSecond(0)
						}
					}
				}

				if matches[7] != "" {
					// Use the timezone if it's available.
					if loc, err := time.LoadLocation(matches[7]); err == nil {
						ret.timezone = loc
					}	
				}
			}
		}
	}

	return ret
}

var InfinitePast Time = Time{ Attrs: map[string][]string{ "InfinitePast": []string{""} } }
var InfiniteFuture Time = Time{ Attrs: map[string][]string{ "InfiniteFuture": []string{""} } }
var UnknownTime Time = Time{} 

func TimeAt(vs ...int) Time {
	t := CreateEmptyTime(TUK_General)
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

func (this Note) Summary() (ret string) {
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
				continue
			}
		}

		re := regexp.MustCompile(`^\s*<!--\s*Edit:(.*)-->\s*$`)
		if matches := re.FindStringSubmatch(line); matches != nil {
			// Append new snippet. As a special case, if it's before the first 
			// snippet, just update the time.
			if len(note.Snippets) != 0 || len(snippet.Text) != 0 {
				snippet.Text = strings.TrimRight(snippet.Text, "\r\n")
				note.Snippets = append(note.Snippets, snippet)
				snippet = NoteSnippet{}
			}

			snippet.Time = CreateTimeFromTag(strings.TrimSpace(matches[1]))

			if i != 0 && snippet.Time.IsUnknown() {
				warns.Add("Unknown time tag in @0@.", &note)
			}
			continue
		}

		snippet.Text += line + "\n"
	}
	snippet.Text = strings.TrimRight(snippet.Text, "\r\n")
	note.Snippets = append(note.Snippets, snippet)

	return note, warns
}

//-- method (utils)

func convAttrsFileToChart(attrs []string) (cattrs map[string][]string) {
	cattrs = map[string][]string{}
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

func convIDStringToInternal(id string) (MajorKind, mNumberID) {
	// According to the specification, the first letter should specify the kind.
	if len(id) < 1 {
		return MOK_Unknown, mNID_Invalid
	}

	kind, err := GetMajorKind(id[0:1])
	if err != nil {
		return MOK_Unknown, mNID_Invalid
	} else if len(id) < 2 {
		return kind, mNID_Invalid
	}

	nid, err := strconv.ParseUint(id[1:], 10, 64)
	if err != nil {
		return kind, mNID_Invalid
	}

	return kind, mNumberID(nid)
}

func convIDInternalToString(kind MajorKind, nid mNumberID) string {
	if prefix, err := kind.Prefix(); err == nil {
		return prefix + strconv.FormatUint(uint64(nid), 10)
	} else {
		return "?"
	}
}

func sanitizeStartEndDatePair(ts Time, te Time, fixedStart bool) (nts Time, nte Time, sanitized bool) {
	nts = ts
	nte = te
	sanitized = false

	if fixedStart {
		if ts.isFloating() {
			nte = UnknownTime
			nte.SetUsage(TUK_End)
			sanitized = true
		} else if te.isFloating() {
			nte = InfiniteFuture
			nte.SetUsage(TUK_End)
			sanitized = true
		}
	} else {
		if te.isFloating() {
			nts = UnknownTime
			nts.SetUsage(TUK_Start)
			sanitized = true
		} else if ts.isFloating() {
			nts = InfinitePast
			nts.SetUsage(TUK_Start)
			sanitized = true
		}
	}

	return
}

// Only outcome: none, ambiguous, inverted, or ambigous AND THEN inverted.
func correctStartEndDatePair(ts Time, te Time) (nts Time, nte Time, ambiguous bool, inverted bool) {
	nts = ts
	nte = te

	// Correct ambiguity.
	if !ts.isFloating() && !ts.IsInfinite() && ts.IsAmbiguous() && !te.IsAmbiguous() {
		newStartDate := ts
		if ts.year == nil && ts.month == nil {
			// Can still be inverted. To be treated below.
			newStartDate.SetYear(te.GetYear())
			newStartDate.SetMonth(te.GetMonth())
		} else if ts.year == nil {
			if te.GetMonth() >= ts.GetMonth() {
				newStartDate.SetYear(te.GetYear())
			} else {
				newStartDate.SetYear(te.GetYear() - 1)
			}
		} else /* if ts.month == nil */ {
			newStartDate.SetMonth(te.GetMonth())
		}

		nts = newStartDate
		ambiguous = true
	} else if !te.isFloating() && !te.IsInfinite() && te.IsAmbiguous() && !ts.IsAmbiguous() {
		newEndDate := te
		if te.year == nil && te.month == nil {
			// Can still be inverted. To be treated below.
			newEndDate.SetYear(ts.GetYear())
			newEndDate.SetMonth(ts.GetMonth())
		} else if te.year == nil {
			if te.GetMonth() >= ts.GetMonth() {
				newEndDate.SetYear(ts.GetYear())
			} else {
				newEndDate.SetYear(ts.GetYear() + 1)
			}
		} else /* if te.month == nil */ {
			newEndDate.SetMonth(ts.GetMonth())
		}
		nte = newEndDate
		ambiguous = true
	}

	// Correct order.
	if !nts.isFloating() && !nte.isFloating() {
		if nts.IsAfter(nte) {
			tmpt := nts
			tmpsu := nts.GetUsage()
			tmpeu := nte.GetUsage()
			nts = nte
			nte = tmpt
			nts.SetUsage(tmpsu)
			nte.SetUsage(tmpeu)
			inverted = true
		}
	}

	return
}

//-- vvv Constants vvv

//-- type mTimeKind

type mTimeKind int
const (
	mTK_Year mTimeKind = iota
	mTK_Month
	mTK_Day
	mTK_Hour
	mTK_Minute
	mTK_Second
)

const TKValue_Max int = math.MaxInt
const TKValue_Min int = math.MinInt 
const TKValue_Unknown int = -1 
