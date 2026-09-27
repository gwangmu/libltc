package libltc

import (
	"errors"

	"github.com/gwangmu/libltc/warning"
)

type IChart interface {
	asChart() *Chart
	resolveReferenceTo(IMainObj)
	unresolveReferenceTo(IMainObj)
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

//-- interface IChart

func (this *Chart) asChart() *Chart {
	return this
}

func (this *Chart) resolveReferenceTo(obj IMainObj) {
	// TODO: resolve dangling references to `obj` in existing objs.
	// TODO: assume GetQualifiedID() of `obj` is valid.
	panic("Unimplemented")
}

func (this *Chart) unresolveReferenceTo(obj IMainObj) {
	// TODO: change references to `obj` dangling ref.
	// TODO: assume GetQualifiedID() of `obj` is valid.
	panic("Unimplemented")
}

//-- interface IWarningObj

func (this *Chart) Summary() (ret string) {
	ret = "the LTC chart"
	if !this.Subject.IsUnknown() {
		ret += " for " + this.Subject.Summary()
	} else if this.Filepath != "" {
		ret += " from the file at '" + this.Filepath + "'"
	}
	return
}

func (this *Chart) IsUnknown() bool {
	return this.Subject.IsUnknown() && this.Filepath == ""
}

//-- method (main object association)

func (this *Chart) getNextNumberID(kind MainObjKind) (NumberID, error) {
	mainobjArr := []*MainObjCommon{}

	switch kind {
	case MOK_Event:
		for _, o := range this.events {
			mainobjArr = append(mainobjArr, o.Common())
		}
	case MOK_Annex:
		for _, o := range this.annexs {
			mainobjArr = append(mainobjArr, o.Common())
		}
	case MOK_Import:
		for _, o := range this.imports {
			mainobjArr = append(mainobjArr, o.Common())
		}
	default:
		return NID_Invalid, errors.New("Unrecognized main object kind")
	}

	// Next next ID = Max existing numID + 1
	// FIX: error out if 'maxNumID' == UINT_MAX - 1. Unlikely, but still.
	var maxNumID NumberID = 0
	for _, mainobj := range mainobjArr {
		maxNumID = max(maxNumID, mainobj.getNumberID())
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
	panic("Unimplemented")
}

func (this *Chart) GetAnnexs() []*Annex {
	// TODO
	panic("Unimplemented")
}

func (this *Chart) GetImports() []*Import {
	// TODO
	panic("Unimplemented")
}

//-- method (main object manipulation)

func (this *Chart) GetObject(qualId string) IMainObj {
	// TODO: return obj by qualified id.
	panic("Unimplemented")
}

func (this *Chart) GetEventsBetween(start Time, end Time, inclusive bool) []*Event {
	// TODO: inclusive = false: start < estart && eend < end
	// TODO: inclusive = true: start <= eend || estart <= end
	panic("Unimplemented")
}

func (this *Chart) GetEventsPerCategory() map[string][]*Event {
	// TODO
	panic("Unimplemented")
}

func (this *Chart) GetEventCategories() []string {
	// TODO
	panic("Unimplemented")
}

func (this *Chart) HasObject(obj IMainObj, onlyLocal bool) bool {
	// TODO
	panic("Unimplemented")
}

func (this *Chart) AddObject(obj IMainObj) {
	// TODO: add already-created obj. update chart and id.
	// TODO: for import objects, add `imported*` to the chart, too.
	// TODO: for import objects, resorveReferenceTo all imported objs.
	panic("Unimplemented")
}

func (this *Chart) RemoveObject(obj IMainObj) {
	// TODO: dispose of any possible links to other objs.
	// TODO: for import objects, unresorveReferenceTo all imported objs.
	// TODO: for import objects, remove `imported*` from the chart, too.
	// TODO: remove itself from the chart.
	panic("Unimplemented")
}

//-- method (creation)

func CreateEmptyChart() *Chart {
	// TODO
	panic("Unimplemented")
}

//-- struct Setting 

type WSE = warning.SummaryElement

type Setting struct {
	CalendarSystem string
	NoteFormat string
	Attrs map[string][]string
}

//-- struct Setting: interface IWarningObj

func (this *Setting) Summary() string {
	return warning.BuildSummaryString(
		WSE{"calendar", this.CalendarSystem},
		WSE{"note", this.NoteFormat},
	)
}

func (this *Setting) IsUnknown() bool {
	return false
}

//-- struct Subject

type Subject struct {
	Name Name
	StartDate Time
	EndDate Time
	Sex string
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
