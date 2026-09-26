package libltc

import (
	"errors"
	"time"
)

type IChart interface {
	asChart() *Chart
	resolveReferenceTo(IMainObj)
}

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
	mainobjArr := []IMainObj{}

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

//-- method (getters)

func (this *Chart) GetEvents() []*Event {
	// TODO
}

func (this *Chart) GetAnnexs() []*Annex {
	// TODO
}

func (this *Chart) GetImports() []*Import {
	// TODO
}

//-- method (main object manipulation)

func (this *Chart) GetObject(qualId string) IMainObj {
	// TODO: return obj by qualified id.
}

func (this *Chart) GetEventsBetween(start Time, end Time, inclusive bool) []*Event {
	// TODO: inclusive = false: start < estart && eend < end
	// TODO: inclusive = true: start <= eend || estart <= end
}

func (this *Chart) GetEventsPerCategory() map[string][]*Event {
	// TODO
}

func (this *Chart) GetEventCategories() []string {
	// TODO
}

func (this *Chart) HasObject(obj IMainObj, onlyLocal bool) bool {
	// TODO
}

func (this *Chart) AddObject(obj IMainObj) {
	// TODO: add already-created obj. update chart and id.
	// TODO: renewobject implied.
}

func (this *Chart) RenewObject(obj IMainObj) {
	// TODO: fix relational fields between objs.
	// TODO: lib user should call this if `obj` was directly changed, NOT via chart.
	// TODO: for import objects, this will trigger re-import.
}

func (this *Chart) RemoveObject(obj IMainObj) {
	// TODO: dispose of any possible links to other objs.
	// TODO: remove itself from the chart.
}

//-- method (creation)

func CreateEmptyChart() *Chart {
	// TODO
}
