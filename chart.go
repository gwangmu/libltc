package libltc

import (
	"errors"
	"time"
)

// LTCChart and LTC[^F]* structs are internal representation of the LTC file.

type LTCVersion int 
const (
	Ver_26_09_1 LTCVersion = iota
)
var ltcVersionToStr = map[LTCVersion]string{
	Ver_26_09_1: "26.09.1",
}

const Ver_Unknown LTCVersion = -1
const ltcVersionUnknownStr := "?"

func (ver LTCVersion) String() (string, error) {
	if verstr, ok := ltcVersionToStr[ver]; ok {
		return verstr
	}
	return ltcVersionUnknownStr, errors.New("Unknown LTCVersion")
}

func (ver LTCVersion) Summary() string {
	if verstr, ok := ltcVersionToStr[ver]; ok {
		return verstr
	}
	return ltcVersionUnknownStr
}

func (ver LTCVersion) IsUnknown() bool {
	_, exists := ltcVersionToStr[ver]
	return !exists
}

func GetLTCVersion(reqverstr string) (outver LTCVersion, e error) {
	for ver, verstr := range ltcVersionToStr {
		if verstr == reqverstr {
			outver = ver
			e = nil
			return
		}
	}
	return Ver_Unknown, errors.New("Unknown version")
}

type LTCTimeKind int
const (
	TK_Year LTCTimeKind = iota
	TK_Month LTCTimeKind
	TK_Day LTCTimeKind
	TK_Hour LTCTimeKind
	TK_Minute LTCTimeKind
	TK_Second LTCTimeKind
)

type LTCChart struct {
	Filepath string
	Version LTCVersion
	Setting LTCSetting
	Entity LTCEntity
	Events []*LTCEvent	// Sorted by StartDate (unknown last)
	Annexs []*LTCAnnex
	Attrs map[string]string
}

func (this LTCChart) Summary() (ret string) {
	ret = "the LTC chart"
	if !this.Entity.IsUnknown() {
		ret += " for " + this.Entity.Summary()
	} else if this.Filepath != "" {
		ret += " from the file at \'" + this.Filepath "\'"
	}
	return
}

func (this LTCChart) IsUnknown() bool {
	return this.Entity.IsUnknown() && this.Filepath == ""
}

type LTCSetting struct {
	DisplayLanguage string 
	CalendarSystem string
	NoteFormat string
	Attrs map[string]string

	chart *LTCChart
}

func (this LTCSetting) GetChart() *LTCChart {
	return this.chart
}

type LTCEntity struct {
	Name LTCName
	Birthday LTCTime
	Sex string
	Attrs map[string]string

	chart *LTCChart
}

func (this LTCEntity) GetChart() *LTCChart {
	return this.chart
}

type LTCEvent struct {
	Title string
	StartDate LTCTime
	EndDate LTCTime
	Note string
	Attrs map[string]string

	chart *LTCChart
	eventID uint64 
}

func (this LTCEvent) GetChart() *LTCChart {
	return this.chart
}

func (this LTCEvent) GetID() string {
	return "e" + strconv.Itoa(this.eventID)
}

type LTCAnnex struct {
	Format string
	Data []byte 
	Attrs map[string]string

	chart *LTCChart
	annexID uint64 
}

func (this LTCAnnex) GetChart() *LTCChart {
	return this.chart
}

func (this LTCAnnex) GetID() string {
	return "a" + strconv.Itoa(this.annexID)
}

type LTCName struct {
	First string
	Middle string
	Last string
	Attrs map[string]string
}

type LTCTime struct {
	Attrs map[string]string

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

	return time.Date(nyear, nmonth, nday, nhour, nminute, nsecond, time.UTC)
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

func (this LTCTime) SetTime(date time.Time) {
	this.year = &data.Year()
	this.month = &int(data.Month())
	this.day = &data.Day()
	this.hour = &data.Hour()
	this.minute = &data.Minute()
	this.second = &data.Second()
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


func (this LTCChart) getNextEventID() (uint64, error) {
	// Next event ID = Last event ID number + 1
	// FIXME: error out if 'maxEventID' == UINT_MAX. Unlikely, but still.
	var maxEventID uint64 = 0
	for _, events := range this.Events {
		for _, event := range events {
			maxEventID = max(maxEventID, event.eventID)
		}
	}
	if (maxEventID == math.MaxUint64) {
		return 0, errors.New("Cannot get the next event ID")
	} else {
		return maxEventID + 1, nil
	}
}

func (this LTCChart) getNextAnnexID() (uint64, error) {
	// Next annex ID = Last annex ID number + 1
	// FIXME: error out if 'maxAnnexID' == UINT_MAX. Unlikely, but still.
	var maxAnnexID uint64 = 0
	for _, annex := range this.Annexs {
		maxAnnexID = max(maxAnnexID, annex.annexID)
	}
	if (maxAnnexID == math.MaxUint64) {
		return 0, errors.New("Cannot get the next annex ID")
	} else {
		return maxAnnexID + 1, nil
	}
}

func (this LTCChart) CreateEvent() (*LTCEvent, error) {
	newEventID, err := this.getNextEventID()
	if err != nil {
		return nil, err
	}

	return &LTCEvent{
		Title: "",
		StartDate: LTCTime{},
		EndDate: LTCTime{},
		Note: "",
		Attrs: map[string]string{},

		chart: &this,
		eventID: newEventID,
	}
}

func (this LTCChart) CreateAnnex() (*LTCAnnex, error) {
	newAnnexID, err := this.getNextAnnexID()
	if err != nil {
		return nil, err
	}

	return &LTCAnnex{
		Format: "text",
		Data: []byte{},
		Attrs: map[string]string{},

		chart: &this,
		annexID: newAnnexID,
	}
}
