package libltc

import (
	"errors"
	"time"
)

type Chart struct {
	Filepath string
	Version Version
	Note Note

	Setting Setting
	Subject Subject

	events []*Event		// Sorted by StartDate (unknown first), incl. imported events.
	annexs []*Annex
	imports []*Import
}

type Setting struct {
	DisplayLanguage string 
	CalendarSystem string
	NoteFormat string
}

type Subject struct {
	Name Name
	StartDate Time
	EndDate Time
	Sex string
}

//-- interface WarningObj

func (this *Chart) Summary() (ret string) {
	ret = "the LTC chart"
	if !this.Subject.IsUnknown() {
		ret += " for " + this.Subject.Summary()
	} else if this.Filepath != "" {
		ret += " from the file at \'" + this.Filepath "\'"
	}
	return
}

func (this *Chart) IsUnknown() bool {
	return this.Subject.IsUnknown() && this.Filepath == ""
}

//-- interface ltcTOMLPrintable

func (this *Chart) PrintTOML() string {
	// TODO
}

//-- method (main object association)

func (this *Chart) getNextNumberID(kind MainObjKind) (NumberID, error) {
	mainobjArr := []MainObj{}

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
	var maxNumID NumberID = 0
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

func (this *Chart) CreateEvent() *Event {
	// TODO
}

func (this *Chart) CreateAnnex() *Annex {
	// TODO
}

func (this *Chart) CreateImport() *Import {
	// TODO
}

func (this *Chart) RemoveEvent(eobj *Event) {
	// TODO: dispose of any possible links to other objs.
	// TODO: remove itself from the chart.
}

func (this *Chart) RemoveAnnex(eobj *Annex) {
	// TODO: dispose of any possible links to other objs.
	// TODO: remove itself from the chart.
}

func (this *Chart) RemoveImport(eobj *Import) {
	// TODO: dispose of any possible links to other objs.
	// TODO: remove itself from the chart.
}

func (this *Chart) 

//-- method (creation)

func CreateEmptyChart() *Chart {
	// TODO
}
