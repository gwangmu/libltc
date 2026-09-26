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

	events []*LTCEvent		// Sorted by StartDate (unknown first), incl. imported events.
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

func (this *LTCChart) getNextNumberID(kind LTCMainObjKind) (LTCNumberID, error) {
	mainobjArr := []LTCMainObj{}

	switch kind {
	case MOK_Event:
		for _, o := range this.events {
			mainobjArr = append(mainobjArr, o.Common)
		}
	case MOK_Annex:
		for _, o := range this.annexs {
			mainobjArr = append(mainobjArr, o.Common)
		}
	case MOK_Import:
		for _, o := range this.imports {
			mainobjArr = append(mainobjArr, o.Common)
		}
	default:
		return NID_Invalid, errors.New("Unrecognized main object kind")
	}

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

//-- method (main object manipulation)

func (this *LTCChart) CreateEvent() *LTCEvent {
	// TODO
}

func (this *LTCChart) CreateAnnex() *LTCAnnex {
	// TODO
}

func (this *LTCChart) CreateImport() *LTCImport {
	// TODO
}

func (this *LTCChart) RemoveEvent(eobj *LTCEvent) {
	// TODO: dispose of any possible links to other objs.
	// TODO: remove itself from the chart.
}

func (this *LTCChart) RemoveAnnex(eobj *LTCAnnex) {
	// TODO: dispose of any possible links to other objs.
	// TODO: remove itself from the chart.
}

func (this *LTCChart) RemoveImport(eobj *LTCImport) {
	// TODO: dispose of any possible links to other objs.
	// TODO: remove itself from the chart.
}

func (this *LTCChart) 

//-- method (creation)

func CreateEmptyChart() *LTCChart {
	// TODO
}
