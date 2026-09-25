package libltc

import (
	"errors"
	"time"
)

type LTCChart struct {
	Filepath string
	Version LTCVersion
	Note LTCNote

	Setting LTCSetting
	Subject LTCSubject

	events []*LTCEvent	// Sorted by StartDate (unknown first), incl. imported events.
	annexs []*LTCAnnex
	imports []*LTCImport
}

type LTCSetting struct {
	DisplayLanguage string 
	CalendarSystem string
	NoteFormat string
}

type LTCSubject struct {
	Name LTCName
	StartDate LTCTime
	EndDate LTCTime
	Sex string
}

//-- interface LTCWarningObj

func (this *LTCChart) Summary() (ret string) {
	ret = "the LTC chart"
	if !this.Subject.IsUnknown() {
		ret += " for " + this.Subject.Summary()
	} else if this.Filepath != "" {
		ret += " from the file at \'" + this.Filepath "\'"
	}
	return
}

func (this *LTCChart) IsUnknown() bool {
	return this.Subject.IsUnknown() && this.Filepath == ""
}

//-- interface ltcTOMLPrintable

func (this *LTCChart) PrintTOML() string {
	// TODO
}

//-- method (main object association)

func (this *LTCChart) LinkEventObject(o *LTCEvent) error {
	// TODO: grant new ID, update chart
	// TODO: if embed, do embed
}

func (this *LTCChart) LinkAnnexObject(o *LTCAnnex) error {
	// TODO: grant new ID, update chart
	// TODO: do extra note or attach
}

func (this *LTCChart) LinkImportObject(o *LTCImport) error {
	// TODO: grant new ID, update chart
	// TODO: do import
}

func (this *LTCChart) RenewEventObject(o *LTCEvent) error {
	// TODO: similar to add, but no ID/Chart update - just contents.
}

func (this *LTCChart) RenewAnnexObject(o *LTCAnnex) error {
	// TODO: similar to add, but no ID/Chart update - just contents.
}

func (this *LTCChart) RenewImportObject(o *LTCImport) error {
	// TODO: similar to add, but no ID/Chart update - just contents.
}

func (this *LTCChart) UnlinkEventObject(o *LTCEvent) error {
	// TODO: don't remove obj
}

func (this *LTCChart) UnlinkAnnexObject(o *LTCAnnex) error {
	// TODO: don't remove obj
}

func (this *LTCChart) UnlinkImportObject(o *LTCImport) error {
	// TODO: don't remove obj
}

func (this *LTCChart) getNextNumberID(mainobjArr *[]ltcMainObjCommon) (LTCNumberID, error) {
	// Next next ID = Max existing numID + 1
	// FIX: error out if 'maxNumID' == UINT_MAX - 1. Unlikely, but still.
	var maxNumID LTCNumberID = 0
	for _, mainobj := range mainobjArr {
		maxNumID = max(maxNumID, event.getNumberID())
	}
	if (maxNumID == NID_Max) {
		return 0, errors.New("Cannot get the next number ID")
	} else {
		return maxNumID + 1, nil
	}
}

//-- method (creation)

func CreateChart() *LTCChart {
	// TODO
}

func LoadChart(uri string) (*LTCChart, error) {
	// TODO
}
