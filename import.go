package libltc

import (
	"errors"

	"github.com/gwangmu/libltc/internal/file"
)

type Import struct {
	common MainObjCommon

	Note Note

	link string
	startDate *Time
	endDate *Time
	offsetDate *Time
	categories *[]string		// nil: any categories
	
	loadedChart *Chart
	importedEvents []*Event
	importedAnnexs []*Annex
	importedImports []*Import

	preOffsetEvents []*Event
}

//-- interface MainObj

func (this *Import) Common() *MainObjCommon {
	return &this.common
}

func (this *Import) IsUnresolved() bool {
	return this.common.IsUnresolved()
}

func (this *Import) Clone(preserveChart bool, preserveID bool) IMainObj {
	newobj := *this
	if !preserveChart {
		newobj.common.chart = nil
	}
	if !preserveID {
		newobj.common.numID = NID_Invalid
		newobj.common.fullID = ""
	}
	return &newobj
}

//-- interface WarningObj

func (this *Import) Summary() string {
	// TODO
}

func (this *Import) IsUnknown() bool {
	// TODO
}

//-- method (getters)

func (this *Import) GetLink() string {
	// TODO
}

func (this *Import) GetStartDate() *Time {
	// TODO
}

func (this *Import) GetEndDate() *Time {
	// TODO
}

func (this *Import) GetOffsetDate() *Time {
	// TODO
}

func (this *Import) GetCategories() *[]string {
	// TODO
}

func (this *Import) GetImportedChart() *Chart {
	// TODO
}

func (this *Import) GetImportedEvents() []*Event {
	// TODO
}

func (this *Import) GetImportedAnnexs() []*Event {
	// TODO
}

func (this *Import) GetImportedImports() []*Import {
	// TODO
}

//-- method (setters)

func (this *Import) SetLink(v string) {
	// TODO: unresolveReferenceTo all imported objects.
	// TODO: reload `loadedChart`
	// TODO: if successful, clear `imported*`, replace `loadedChart`
	// TODO: clone main objects (ID preserved)
	// TODO: set parents of cloned objects to `this` 
	// TODO: mask period in cloned events
	// TODO: save it to `cachedEvents`
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
}

func (this *Import) SetStartDate(v *Time) {
	// TODO: unresolveReferenceTo all imported objects.
	// TODO: clone main objects (ID preserved)
	// TODO: set parents of cloned objects to `this` 
	// TODO: mask period in cloned events
	// TODO: save it to `cachedEvents`
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
}

func (this *Import) SetEndDate(v *Time) {
	// TODO: unresolveReferenceTo all existing objects.
	// TODO: clone main objects (ID preserved)
	// TODO: set parents of cloned objects to `this` 
	// TODO: mask period in cloned events
	// TODO: save it to `cachedEvents`
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
}

func (this *Import) SetOffsetDate(v *Time) {
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
}

func (this *Import) SetCategories(v *[]string) {
	// TODO: unresolveReferenceTo all imported objects.
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
}
