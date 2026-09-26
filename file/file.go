package file

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type wii = struct {string; interface{}}

type LTCTOMLPrintable interface {
	PrintTOML() string
}

//-- struct LTCFile

type LTCFile struct {
	Filepath string 		`toml:"-"`
	Version string 			`toml:"FormatVersion"`
	Note string 			`toml:"Note,omitempty"`
	Setting LTCSetting 		`toml:"Setting"`
	Subject LTCSubject 		`toml:"Subject"`
	Event []LTCEvent		`toml:"Event,omitempty"`
	Annex []LTCAnnex		`toml:"Annex,omitempty"`
	Import []LTCImport		`toml:"Import,omitempty"`
}

//-- struct LTCFile: interface LTCTOMLPrintable

func (this *LTCFile) PrintTOML() string {
	// TODO
}

//-- struct LTCFile: interface LTCWarningObj

func (this *LTCFile) Summary() (ret string) {
	ret = "the LTC file"
	if this.Filepath != "" {
		ret += " at \'" + this.Filepath + "\'"
	} else if !this.Subject.IsUnknown() {
		ret += " for " + this.Subject.Summary()
	}
	return
}

func (this *LTCFile) IsUnknown() bool {
	return this.Filepath == "" && this.Subject.IsUnknown()
}

//-- struct LTCSetting

type LTCSetting struct {
	CalendarSystem string	`toml:"CalendarSystem"`
	NoteFormat string		`toml:"NoteFormat"`
	Categories []string		`toml:"Categories,omitempty"`
}

//-- struct LTCSetting: interface LTCTOMLPrintable

func (this *LTCSetting) PrintTOML() string {
	// TODO
}

//-- struct LTCSetting: interface LTCWarningObj

func (this *LTCSetting) Summary() string {
	return "the LTC setting"
}

func (this *LTCSetting) IsUnknown() bool {
	return false
}

//-- struct LTCSubject

type LTCSubject struct {
	Name LTCName			`toml:"Name"`
	StartDate LTCTime		`toml:"StartDate"`
	EndDate *LTCTime		`toml:"StartDate,omitempty"`
	Sex string				`toml:"Sex,omitempty"`
}

//-- struct LTCSubject: interface LTCTOMLPrintable

func (this *LTCSubject) PrintTOML() string {
	// TODO
}

//-- struct LTCSubject: interface LTCWarningObj

func (this *LTCSubject) Summary() (ret string) {
	ret = this.Name.Summary()

	extraStr := getWarningInfoString(
		wii{"", this.Sex}, 
		wii{"started", this.StartDate},
		wii{"ended", this.EndDate},
	)
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}
	return
}

func (this *LTCSubject) IsUnknown() bool {
	return this.Name.IsUnknown() && this.StartDate.IsUnknown() && 
		this.EndDate.IsUnknown() && this.Sex == ""
}

//-- struct LTCEvent

type LTCEvent struct {
	ID string				`toml:"ID"`
	Title string			`toml:"Title"`
	Category string			`toml:"Category,omitempty"`
	Subchart *string 		`toml:"Subchart,omitempty"`
	Embed *string 			`toml:"Embed,omitempty"`
	StartDate *LTCTime		`toml:"StartDate,omitempty"`
	EndDate *LTCTime		`toml:"EndDate,omitempty"`
	Note *string			`toml:"Note,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

//-- struct LTCEvent: interface LTCTOMLPrintable

func (this *LTCEvent) PrintTOML() string {
	// TODO
}

//-- struct LTCEvent: interface LTCWarningObj

func (this *LTCEvent) Summary() (ret string) {
	if this.Title != "" {
		ret = "the event \'" + this.Title + "\'"
	} else {
		ret = "the untitled event"
	}

	extraStr := ""
	if this.Title == "" {
		extraStr = getWarningInfoString(
			wii{"ID", this.ID},
			wii{"category", this.Category}, 
			wii{"embedding chart", this.Subchart},
			wii{"started", this.StartDate}, 
			wii{"ended", this.EndDate},
		)
	} else {
		extraStr = getWarningInfoString(
			wii{"ID", this.ID},
			wii{"category", this.Category},
			wii{"embedding chart", this.Subchart},
		)
	}
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}

	return
}

func (this *LTCSubject) IsUnknown() bool {
	return this.ID == "" && this.Title == "" && 
		this.Category == "" && this.Subchart == "" &&
		(this.StartDate == nil || this.StartDate.IsUnknown()) &&
		(this.EndDate == nil ||| this.EndDate.IsUnknown())
}

//-- struct LTCAnnex

type LTCAnnex struct {
	ID string 				`toml:"ID"`
	Format string 			`toml:"Format,omitempty"`
	Encoding string 		`toml:"Encoding,omitempty"`
	Data string 			`toml:"Data"`
	Note string				`toml:"Note,omitempty"`
	Attrs []string 			`toml:"Attrs,omitempty"`
}

//-- struct LTCAnnex: interface LTCTOMLPrintable

func (this *LTCAnnex) PrintTOML() string {
	// TODO
}

//-- struct LTCAnnex: interface LTCWarningObj

func (this *LTCAnnex) Summary() (ret string) {
	ret = "an annex"

	extraStr := getWarningInfoString(
		wii{"ID", this.ID},
		wii{"format", this.Format}, 
	)
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	} else {
		ret = "an unknown annex"
	}

	return
}

func (this *LTCAnnex) IsUnknown() bool {
	return this.ID == "" && this.Format == ""
}

//-- struct LTCImport

type LTCImport struct {
	ID string				`toml:"ID"`
	Title string 			`toml:"Title,omitempty"`
	Link string 			`toml:"Link"`
	StartDate *LTCTime		`toml:"StartDate,omitempty"`
	EndDate *LTCTime		`toml:"EndDate,omitempty"`
	OffsetDate *LTCTime		`toml:"OffsetDate,omitempty"`
	Categories []string		`toml:"Categories,omitempty"`
	Note string				`toml:"Note,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

//-- struct LTCImport: interface LTCTOMLPrintable

func (this *LTCImport) PrintTOML() string {
	// TODO
}

//-- struct LTCImport: interface LTCWarningObj

func (this *LTCImport) Summary() (ret string) {
	if this.Title != "" {
		ret = "the import '" + this.Title + "'"
	} else {
		ret = "an import"
	}

	extraStr := getWarningInfoString(
		wii{"ID", this.ID},
		wii{"from", this.StartDate},
		wii{"to", this.EndDate},
		wii{"with offset", this.OffsetDate}
	)
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	} else if this.Title == "" {
		ret = "an unknown import"
	}

	return
}

func (this *LTCImport) IsUnknown() bool {
	return this.Title == "" && this.Link == "" && this.StartDate == nil &&
			this.EndDate == nil && this.OffsetDate == nil
}

//-- struct LTCName

type LTCName struct {
	First string			`toml:"First"`
	Middle string			`toml:"Middle,omitempty"`
	Last string				`toml:"Last,omitempty"`
}

//-- struct LTCName: interface LTCTOMLPrintable

func (this *LTCName) PrintTOML() string {
	// TODO
}

//-- struct LTCName: interface LTCWarningObj

func (this *LTCName) Summary() string {
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

	if (len(names) > 0) {
		return string.Join(names, " ")
	} else {
		return "an unknown name"
	}
}

func (this *LTCName) IsUnknown() bool {
	return this.First == "" && this.Middle == "" && this.Last == ""
}

//-- struct LTCTime

type LTCTime struct {
	Year *int				`toml:"Year,omitempty"`
	Month *int				`toml:"Month,omitempty"`
	Day *int  				`toml:"Day,omitempty"`
	Hour *int				`toml:"Hour,omitempty"`
	Minute *int				`toml:"Minute,omitempty"`
	Second *int				`toml:"Second,omitempty"`
	Timezone *string		`toml:"Timezone,omitempty"`
	Attrs []string 			`toml:"Attrs,omitempty"`
}

//-- struct LTCTime: interface LTCTOMLPrintable

func (this *LTCTime) PrintTOML() string {
	// TODO
}

//-- struct LTCTime: interface LTCWarningObj

func (this *LTCTime) Summary() (ret string) {
	if this.IsUnknown() {
		return "an unknown time"
	}

	var syear, smonth, sday string
	var stime []string

	if this.Year != nil {
		syear = strconv.Itoa(this.Year)
	} else {
		syear = "????"
	}
	if this.Month != nil {
		smonth = strconv.Itoa(this.Month)
	} else {
		smonth = "??"
	}
	if this.Day != nil {
		sday = strconv.Itoa(this.Day)
	} else {
		sday = "??"
	}

	if this.Hour == nil && this.Minute == nil && this.Second == nil {
		return syear + "/" + smonth + "/" + sday
	}

	if this.Hour != nil {
		stime = append(stime, strconv.Itoa(this.Hour))
	} else if this.Minute != nil || this.Second != nil {
		shour = append(stime, "??")
	} 
	if this.Minute != nil {
		stime = append(stime, strconv.Itoa(this.Minute))
	} else if this.Second != nil {
		stime = append(stime, "??")
	}
	if this.Second != nil {
		stime = append(stime, strconv.Itoa(this.Second))
	}

	return syear + "/" + smonth + "/" + sday + " " + strings.Join(stime, ":")
}

func (this *LTCTime) IsUnknown() bool {
	return this.Year == nil && this.Month == nil && this.Day == nil &&
		this.Hour == nil && this.Minute == nil && this.Second == nil &&
		this.Timezone == nil
}

//-- Package methods

func LoadLTCFile(filepath string) (*LTCFile, []LTCWarning, error) {
	ltcFile := getDefaultLTCFile()
	if _, err := toml.DecodeFile(filepath, &ltcFile); err != nil {
		return nil, err
	}

	ltcFile.Filepath = filepath

	// TODO: Auto-fix some issues and report via LTCWarning.

	return &ltcFile, nil
}
