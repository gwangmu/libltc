package file

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type wii = struct {string; interface{}}

type ITOMLPrinter interface {
	PrintTOML() string
}

//-- struct File

type File struct {
	Filepath string 		`toml:"-"`
	Version string 			`toml:"FormatVersion"`
	Note string 			`toml:"Note,omitempty"`
	Setting Setting 		`toml:"Setting"`
	Subject Subject 		`toml:"Subject"`
	Event []Event		`toml:"Event,omitempty"`
	Annex []Annex		`toml:"Annex,omitempty"`
	Import []Import		`toml:"Import,omitempty"`
}

//-- struct File: interface ITOMLPrinter

func (this *File) PrintTOML() string {
	// TODO
}

//-- struct File: interface IWarningObj

func (this *File) Summary() (ret string) {
	ret = "the LTC file"
	if this.Filepath != "" {
		ret += " at \'" + this.Filepath + "\'"
	} else if !this.Subject.IsUnknown() {
		ret += " for " + this.Subject.Summary()
	}
	return
}

func (this *File) IsUnknown() bool {
	return this.Filepath == "" && this.Subject.IsUnknown()
}

//-- struct Setting

type Setting struct {
	CalendarSystem string	`toml:"CalendarSystem"`
	NoteFormat string		`toml:"NoteFormat"`
	Categories []string		`toml:"Categories,omitempty"`
}

//-- struct Setting: interface ITOMLPrinter

func (this *Setting) PrintTOML() string {
	// TODO
}

//-- struct Setting: interface IWarningObj

func (this *Setting) Summary() string {
	return "the LTC setting"
}

func (this *Setting) IsUnknown() bool {
	return false
}

//-- struct Subject

type Subject struct {
	Name Name			`toml:"Name"`
	StartDate Time		`toml:"StartDate"`
	EndDate *Time		`toml:"StartDate,omitempty"`
	Sex string				`toml:"Sex,omitempty"`
}

//-- struct Subject: interface ITOMLPrinter

func (this *Subject) PrintTOML() string {
	// TODO
}

//-- struct Subject: interface IWarningObj

func (this *Subject) Summary() (ret string) {
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

func (this *Subject) IsUnknown() bool {
	return this.Name.IsUnknown() && this.StartDate.IsUnknown() && 
		this.EndDate.IsUnknown() && this.Sex == ""
}

//-- struct Event

type Event struct {
	ID string				`toml:"ID"`
	Title string			`toml:"Title"`
	Category string			`toml:"Category,omitempty"`
	Subchart *string 		`toml:"Subchart,omitempty"`
	Embed *string 			`toml:"Embed,omitempty"`
	StartDate *Time		`toml:"StartDate,omitempty"`
	EndDate *Time		`toml:"EndDate,omitempty"`
	Note *string			`toml:"Note,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

//-- struct Event: interface ITOMLPrinter

func (this *Event) PrintTOML() string {
	// TODO
}

//-- struct Event: interface IWarningObj

func (this *Event) Summary() (ret string) {
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

func (this *Subject) IsUnknown() bool {
	return this.ID == "" && this.Title == "" && 
		this.Category == "" && this.Subchart == "" &&
		(this.StartDate == nil || this.StartDate.IsUnknown()) &&
		(this.EndDate == nil ||| this.EndDate.IsUnknown())
}

//-- struct Annex

type Annex struct {
	ID string 				`toml:"ID"`
	Format string 			`toml:"Format,omitempty"`
	Encoding string 		`toml:"Encoding,omitempty"`
	Data string 			`toml:"Data"`
	Note string				`toml:"Note,omitempty"`
	Attrs []string 			`toml:"Attrs,omitempty"`
}

//-- struct Annex: interface ITOMLPrinter

func (this *Annex) PrintTOML() string {
	// TODO
}

//-- struct Annex: interface IWarningObj

func (this *Annex) Summary() (ret string) {
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

func (this *Annex) IsUnknown() bool {
	return this.ID == "" && this.Format == ""
}

//-- struct Import

type Import struct {
	ID string				`toml:"ID"`
	Title string 			`toml:"Title,omitempty"`
	Link string 			`toml:"Link"`
	StartDate *Time		`toml:"StartDate,omitempty"`
	EndDate *Time		`toml:"EndDate,omitempty"`
	OffsetDate *Time		`toml:"OffsetDate,omitempty"`
	Categories []string		`toml:"Categories,omitempty"`
	Note string				`toml:"Note,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

//-- struct Import: interface ITOMLPrinter

func (this *Import) PrintTOML() string {
	// TODO
}

//-- struct Import: interface IWarningObj

func (this *Import) Summary() (ret string) {
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

func (this *Import) IsUnknown() bool {
	return this.Title == "" && this.Link == "" && this.StartDate == nil &&
			this.EndDate == nil && this.OffsetDate == nil
}

//-- struct Name

type Name struct {
	First string			`toml:"First"`
	Middle string			`toml:"Middle,omitempty"`
	Last string				`toml:"Last,omitempty"`
}

//-- struct Name: interface ITOMLPrinter

func (this *Name) PrintTOML() string {
	// TODO
}

//-- struct Name: interface IWarningObj

func (this *Name) Summary() string {
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

func (this *Name) IsUnknown() bool {
	return this.First == "" && this.Middle == "" && this.Last == ""
}

//-- struct Time

type Time struct {
	Year *int				`toml:"Year,omitempty"`
	Month *int				`toml:"Month,omitempty"`
	Day *int  				`toml:"Day,omitempty"`
	Hour *int				`toml:"Hour,omitempty"`
	Minute *int				`toml:"Minute,omitempty"`
	Second *int				`toml:"Second,omitempty"`
	Timezone *string		`toml:"Timezone,omitempty"`
	Attrs []string 			`toml:"Attrs,omitempty"`
}

//-- struct Time: interface ITOMLPrinter

func (this *Time) PrintTOML() string {
	// TODO
}

//-- struct Time: interface IWarningObj

func (this *Time) Summary() (ret string) {
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

func (this *Time) IsUnknown() bool {
	return this.Year == nil && this.Month == nil && this.Day == nil &&
		this.Hour == nil && this.Minute == nil && this.Second == nil &&
		this.Timezone == nil
}

//-- Package methods

func LoadFile(filepath string) (*File, []Warning, error) {
	ltcFile := getDefaultFile()
	if _, err := toml.DecodeFile(filepath, &ltcFile); err != nil {
		return nil, err
	}

	ltcFile.Filepath = filepath

	// TODO: Auto-fix some issues and report via Warning.

	return &ltcFile, nil
}
