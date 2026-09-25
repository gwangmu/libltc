package libltc

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type wii = struct {string; interface{}}

// LTCFile and LTCF* structs represent fields in the raw LTC file.

type LTCFile struct {
	Filepath string 		`toml:"-"`
	Version string 			`toml:"FormatVersion"`
	Note string 			`toml:"Note",omitempty"`
	Setting LTCFSetting 	`toml:"Setting"`
	Subject LTCFSubject 	`toml:"Subject"`
	Event []LTCFEvent		`toml:"Event,omitempty"`
	Annex []LTCFAnnex		`toml:"Annex,omitempty"`
	Import []LTCFImport		`toml:"Import,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

func (this LTCFile) Summary() (ret string) {
	ret = "the LTC file"
	if this.Filepath != "" {
		ret += " at \'" + this.Filepath + "\'"
	} else if !this.Subject.IsUnknown() {
		ret += " for " + this.Subject.Summary()
	}
	return
}

func (this LTCFile) IsUnknown() bool {
	return this.Filepath == "" && this.Subject.IsUnknown()
}

type LTCFSetting struct {
	DisplayLanguage string	`toml:"DisplayLanguage"`
	CalendarSystem string	`toml:"CalendarSystem"`
	NoteFormat string		`toml:"NoteFormat"`
	Categories []string		`toml:"Categories,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

func (this LTCFSetting) Summary() string {
	return "the LTC setting"
}

func (this LTCFSetting) IsUnknown() bool {
	return false
}

type LTCFSubject struct {
	Name LTCFName			`toml:"Name"`
	StartDate LTCFTime		`toml:"StartDate"`
	EndDate *LTCFTime		`toml:"StartDate,omitempty"`
	Sex string				`toml:"Sex,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

func (this LTCFSubject) Summary() (ret string) {
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

func (this LTCFSubject) IsUnknown() bool {
	return this.Name.IsUnknown() && this.StartDate.IsUnknown() && 
		this.EndDate.IsUnknown() && this.Sex == ""
}

type LTCFEvent struct {
	ID string				`toml:"ID"`
	Title string			`toml:"Title"`
	Category string			`toml:"Category,omitempty"`
	Subchart string 		`toml:"Subchart,omitempty"`
	Embed string 			`toml:"Embed,omitempty"`
	StartDate *LTCFTime		`toml:"StartDate,omitempty"`
	EndDate *LTCFTime		`toml:"EndDate,omitempty"`
	Note string				`toml:"Note,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

func (this LTCFEvent) Summary() (ret string) {
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

func (this LTCFSubject) IsUnknown() bool {
	return this.ID == "" && this.Title == "" && 
		this.Category == "" && this.Subchart == "" &&
		(this.StartDate == nil || this.StartDate.IsUnknown()) &&
		(this.EndDate == nil ||| this.EndDate.IsUnknown())
}

type LTCFAnnex struct {
	ID string 				`toml:"ID"`
	Format string 			`toml:"Format,omitempty"`
	Encoding string 		`toml:"Encoding,omitempty"`
	Data string 			`toml:"Data"`
	Note string				`toml:"Note,omitempty"`
	Attrs []string 			`toml:"Attrs,omitempty"`
}

func (this LTCFAnnex) Summary() (ret string) {
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

func (this LTCFAnnex) IsUnknown() bool {
	return this.ID == "" && this.Format == ""
}

type LTCFImport struct {
	ID string				`toml:"ID"`
	Title string 			`toml:"Title,omitempty"`
	Link string 			`toml:"Link"`
	StartDate *LTCFTime		`toml:"StartDate,omitempty"`
	EndDate *LTCFTime		`toml:"EndDate,omitempty"`
	OffsetDate *LTCFTime	`toml:"OffsetDate,omitempty"`
	Categories []string		`toml:"Categories,omitempty"`
	Note string				`toml:"Note,omitempty"`
	Attrs []string			`toml:"Attrs,omitempty"`
}

func (this LTCFImport) Summary() (ret string) {
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

func (this LTCFImport) IsUnknown() bool {
	return this.Title == "" && this.Link == "" && this.StartDate == nil &&
			this.EndDate == nil && this.OffsetDate == nil
}

type LTCFName struct {
	First string			`toml:"First"`
	Middle string			`toml:"Middle,omitempty"`
	Last string				`toml:"Last,omitempty"`
}

func (this LTCFName) Summary() string {
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

func (this LTCFName) IsUnknown() bool {
	return this.First == "" && this.Middle == "" && this.Last == ""
}

type LTCFTime struct {
	Year *int				`toml:"Year,omitempty"`
	Month *int				`toml:"Month,omitempty"`
	Day *int  				`toml:"Day,omitempty"`
	Hour *int				`toml:"Hour,omitempty"`
	Minute *int				`toml:"Minute,omitempty"`
	Second *int				`toml:"Second,omitempty"`
	Timezone *string		`toml:"Timezone,omitempty"`
}

func (this LTCFTime) Summary() (ret string) {
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

func (this LTCFTime) IsUnknown() bool {
	return this.Year == nil && this.Month == nil && this.Day == nil &&
		this.Hour == nil && this.Minute == nil && this.Second == nil &&
		this.Timezone == nil
}

func getDefaultLTCFile() LTCFile {
	return LTCFile{
		Version: "?",
		Setting: LTCFSetting{
			DisplayLanguage: "English",
			CalendarSystem: "Gregorian",
			NoteFormat: "Markdown",
			Categories: [],
			Attrs: [],
		},
		Subject: LTCFSubject{
			Name: LTCFName{
				First: "Unknown",
				Middle: "",
				Last: "",
				Attrs: [],
			},
			StartDate: LTCFTime{},
			EndDate: nil,
			Sex: "",
			Attrs: [],
		},
		Event: [],
		Annex: [],
		Attrs: [],
	}
}

func loadLTCFile(filepath string) (*LTCFile, []LTCWarning, error) {
	ltcFile := getDefaultLTCFile()
	if _, err := toml.DecodeFile(filepath, &ltcFile); err != nil {
		return nil, err
	}

	ltcFile.Filepath = filepath

	// TODO: Auto-fix some issues and report via LTCWarning.

	return &ltcFile, nil
}

func (ltcFile *LTCFile) save(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	if err := toml.NewEncoder(file).Encode(ltcFile); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
