package file

import (
	"strconv"
	"strings"

	"github.com/gwangmu/libltc/warning"
)

type WSE = warning.SummaryElement

type IFileObject interface {
	File | Setting | Subject | Event | Annex | Import | Name | Time
}

type IDiagnosable interface {
	Diagnose() warning.Warnings
}

type IDiagnosablePtr[T any] interface {
	*T
	IDiagnosable
}

//-- struct File

type File struct {
	Filepath string 	`toml:"-"`
	Version string 		`toml:"FormatVersion"`
	Note string 		`toml:"Note,omitempty"`
	Setting Setting 	`toml:"Setting"`
	Subject Subject 	`toml:"Subject"`
	Event []Event		`toml:"Event,omitempty"`
	Annex []Annex		`toml:"Annex,omitempty"`
	Import []Import		`toml:"Import,omitempty"`
}

//-- struct File: interface IWarningObj

func (this *File) Summary() (ret string) {
	ret = "the LTC file"
	if this.Filepath != "" {
		ret += " at '" + this.Filepath + "'"
	} else if !this.Subject.IsUnknown() {
		ret += " for " + this.Subject.Summary()
	}
	return
}

func (this *File) IsUnknown() bool {
	return this.Filepath == "" && this.Subject.IsUnknown()
}

//-- struct File: interface IDiagnosable

func (this *File) Diagnose() warning.Warnings {
	warns := warning.Warnings{}

	warns.Concat(this.Setting.Diagnose())

	if this.Subject.IsUnknown() {
		warns.Add("@0@ has an unknown subject.", this)
	}
	warns.Concat(this.Subject.Diagnose())

	for _, event := range this.Event {
		warns.Concat(event.Diagnose())
	}

	for _, annex := range this.Annex {
		warns.Concat(annex.Diagnose())
	}

	for _, fimport := range this.Import {
		warns.Concat(fimport.Diagnose())
	}

	return warns
}

//-- struct Setting

type Setting struct {
	CalendarSystem string	`toml:"CalendarSystem"`
	NoteFormat string		`toml:"NoteFormat"`
	Categories []string		`toml:"Categories,omitempty"`
	Attrs []string 			`toml:"Attrs,omitempty"`
}

//-- struct Setting: interface IWarningObj

func (this *Setting) Summary() string {
	return "the LTC setting"
}

func (this *Setting) IsUnknown() bool {
	return false
}

//-- struct Setting: interface IDiagnosable

func (this *Setting) Diagnose() warning.Warnings {
	return warning.Warnings{}
}

//-- struct Subject

type Subject struct {
	Name Name			`toml:"Name"`
	StartDate *Time		`toml:"StartDate"`
	EndDate *Time		`toml:"EndDate,omitempty"`
	Sex string			`toml:"Sex,omitempty"`
}

//-- struct Subject: interface IWarningObj

func (this *Subject) Summary() (ret string) {
	ret = this.Name.Summary()

	extraStr := warning.BuildSummaryString(
		WSE{"", this.Sex}, 
		WSE{"started", this.StartDate},
		WSE{"ended", this.EndDate},
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

//-- struct Subject: interface IDiagnosable

func (this *Subject) Diagnose() warning.Warnings {
	return warning.Warnings{}
}

//-- struct Event

type Event struct {
	ID string			`toml:"ID"`
	Title *string		`toml:"Title"`
	Category string		`toml:"Category,omitempty"`
	EmbedChart *string 	`toml:"EmbedChart,omitempty"`
	EmbedEvent *string 	`toml:"EmbedEvent,omitempty"`
	StartDate *Time		`toml:"StartDate,omitempty"`
	EndDate *Time		`toml:"EndDate,omitempty"`
	Note string			`toml:"Note,omitempty"`
	Attrs []string		`toml:"Attrs,omitempty"`
}

//-- struct Event: interface IWarningObj

func (this *Event) Summary() (ret string) {
	if this.Title != nil {
		ret = "the event '" + *this.Title + "'"
	} else {
		ret = "the untitled event"
	}

	extraStr := ""
	if this.Title == nil {
		extraStr = warning.BuildSummaryString(
			WSE{"ID", this.ID},
			WSE{"category", this.Category}, 
			WSE{"embedding chart", this.EmbedChart},
			WSE{"started", this.StartDate}, 
			WSE{"ended", this.EndDate},
		)
	} else {
		extraStr = warning.BuildSummaryString(
			WSE{"ID", this.ID},
			WSE{"category", this.Category},
			WSE{"embedding chart", this.EmbedChart},
		)
	}
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}

	return
}

func (this *Event) IsUnknown() bool {
	return this.ID == "" && this.Title == nil && 
		this.Category == "" && this.EmbedChart == nil &&
		(this.StartDate == nil || this.StartDate.IsUnknown()) &&
		(this.EndDate == nil || this.EndDate.IsUnknown())
}

//-- struct Event: interface IDiagnosable

func (this *Event) Diagnose() warning.Warnings {
	warns := warning.Warnings{}

	if this.EmbedChart == nil && this.EmbedEvent == nil {
		if this.Title == nil {
			warns.Add("@0@ has no title, but it's not an embedding event.", this)
		}
	}

	if this.StartDate != nil && this.EndDate != nil {
		if this.StartDate.IsAfter(this.EndDate) {
			warns.Add("@0@ has the start date (@1@) later than the end date (@2@). swapping dates...", this, this.StartDate, this.EndDate)
			tmpDatePtr := this.StartDate
			this.StartDate = this.EndDate
			this.EndDate = tmpDatePtr
		}
	}

	return warns
}

//-- struct Annex

type Annex struct {
	ID string 				`toml:"ID"`
	Title string 			`toml:"Title,omitempty"`
	Format string 			`toml:"Format,omitempty"`
	Encoding string 		`toml:"Encoding,omitempty"`
	Data string 			`toml:"Data"`
	Note string				`toml:"Note,omitempty"`
	Attrs []string 			`toml:"Attrs,omitempty"`
}

//-- struct Annex: interface IWarningObj

func (this *Annex) Summary() (ret string) {
	ret = "an annex"

	extraStr := warning.BuildSummaryString(
		WSE{"ID", this.ID},
		WSE{"format", this.Format}, 
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

//-- struct Annex: interface IDiagnosable

func (this *Annex) Diagnose() warning.Warnings {
	// TODO: Auto-fix some issues and report via Warning.
	panic("Unimplemented")
}

//-- struct Import

type Import struct {
	ID string				`toml:"ID"`
	Link string 			`toml:"Link"`
	StartDate *Time			`toml:"StartDate,omitempty"`
	EndDate *Time			`toml:"EndDate,omitempty"`
	OffsetDate *Time		`toml:"OffsetDate,omitempty"`
	Categories *[]string	`toml:"Categories,omitempty"`
	Note string				`toml:"Note,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

//-- struct Import: interface IWarningObj

func (this *Import) Summary() (ret string) {
	ret = "an import"

	extraStr := warning.BuildSummaryString(
		WSE{"ID", this.ID},
		WSE{"from", this.StartDate},
		WSE{"to", this.EndDate},
		WSE{"with offset", this.OffsetDate},
	)
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}

	return
}

func (this *Import) IsUnknown() bool {
	return this.Link == "" && this.StartDate == nil &&
			this.EndDate == nil && this.OffsetDate == nil
}

//-- struct Import: interface IDiagnosable

func (this *Import) Diagnose() warning.Warnings {
	// TODO: Auto-fix some issues and report via Warning.
	panic("Unimplemented")
}

//-- struct Name

type Name struct {
	First string			`toml:"First"`
	Middle string			`toml:"Middle,omitempty"`
	Last string				`toml:"Last,omitempty"`
	Attrs []string 			`toml:"Attrs,omitempty"`
}

//-- struct Name: interface IWarningObj

func (this *Name) Summary() string {
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

	if len(this.Attrs) != 0 {
		addname := strings.Join(this.Attrs, ", ")
		names = append(names, "(" + addname + ")")
	}

	if (len(names) > 0) {
		return strings.Join(names, " ")
	} else {
		return "an unknown name"
	}
}

func (this *Name) IsUnknown() bool {
	return this.First == "" && this.Middle == "" && this.Last == "" &&
		len(this.Attrs) == 0
}

//-- struct Name: interface IDiagnosable

func (this *Name) Diagnose() warning.Warnings {
	// TODO: Auto-fix some issues and report via Warning.
	panic("Unimplemented")
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

//-- struct Time: interface IWarningObj

func (this *Time) Summary() (ret string) {
	if this.IsUnknown() {
		return "an unknown time"
	} else {
		ret = "the time"
	}

	stime := []string{}
	if this.Year != nil {
		stime = append(stime, strconv.Itoa(*this.Year))
	} else {
		stime = append(stime, "????")
	}
	if this.Month != nil {
		stime = append(stime, strconv.Itoa(*this.Month))
	} else {
		stime = append(stime, "??")
	}
	if this.Day != nil {
		stime = append(stime, strconv.Itoa(*this.Day))
	} else {
		stime = append(stime, "??")
	}
	if this.Hour != nil {
		stime = append(stime, strconv.Itoa(*this.Hour))
	} else {
		stime = append(stime, "??")
	} 
	if this.Minute != nil {
		stime = append(stime, strconv.Itoa(*this.Minute))
	} else {
		stime = append(stime, "??")
	}
	if this.Second != nil {
		stime = append(stime, strconv.Itoa(*this.Second))
	}

	sextra := []string{}
	if this.Timezone != nil {
		sextra = append(sextra, *this.Timezone)
	}
	if len(this.Attrs) != 0 {
		sextra = append(sextra, this.Attrs...)
	}

	ret += " " + strings.Join(stime[0:3], "-")
	if this.Hour != nil || this.Minute != nil || this.Second == nil {
		ret += " " + strings.Join(stime[3:6], ":")
	}
	if len(sextra) != 0 {
		ret += " (" + strings.Join(sextra, ", ")
	}

	return
}

func (this *Time) IsUnknown() bool {
	return this.Year == nil && this.Month == nil && this.Day == nil &&
		this.Hour == nil && this.Minute == nil && this.Second == nil &&
		this.Timezone == nil
}

//-- struct Time: interface IDiagnosable

func (this *Time) Diagnose() warning.Warnings {
	// TODO: Auto-fix some issues and report via Warning.
	panic("Unimplemented")
}

//-- struct Time: method

func (this *Time) IsAfter(t *Time) bool {
	// TODO
	panic("Unimplemented")
}
