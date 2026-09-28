package libltc

import (
	"github.com/gwangmu/libltc/warning"
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

//-- interface IMainObj

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

//-- interface IWarningObj

func (this *Import) Summary() string {
	// TODO
	panic("Unimplemented")
}

func (this *Import) IsUnknown() bool {
	// TODO
	panic("Unimplemented")
}

//-- interface IDiagnosable

func (this *Import) DiagnoseLocal() (warns warning.Warnings) {
	// TODO
	panic("Unimplemented")
}

func (this *Import) DiagnoseNonLocal() (warns warning.Warnings) {
	if this.Common().GetChart() == nil {
		warns.Add("@0@ is not associated to any chart.", this)
		return
	}

	// TODO
	panic("Unimplemented")
}

//-- method (getters)

func (this *Import) GetLink() string {
	// TODO
	panic("Unimplemented")
}

func (this *Import) GetStartDate() *Time {
	// TODO
	panic("Unimplemented")
}

func (this *Import) GetEndDate() *Time {
	// TODO
	panic("Unimplemented")
}

func (this *Import) GetOffsetDate() *Time {
	// TODO
	panic("Unimplemented")
}

func (this *Import) GetCategories() *[]string {
	// TODO
	panic("Unimplemented")
}

func (this *Import) GetImportedChart() *Chart {
	// TODO
	panic("Unimplemented")
}

func (this *Import) GetImportedEvents() []*Event {
	// TODO
	panic("Unimplemented")
}

func (this *Import) GetImportedAnnexs() []*Event {
	// TODO
	panic("Unimplemented")
}

func (this *Import) GetImportedImports() []*Import {
	// TODO
	panic("Unimplemented")
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
	panic("Unimplemented")
}

func (this *Import) SetStartDate(v *Time) {
	// TODO: unresolveReferenceTo all imported objects.
	// TODO: clone main objects (ID preserved)
	// TODO: set parents of cloned objects to `this` 
	// TODO: mask period in cloned events
	// TODO: save it to `cachedEvents`
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
	panic("Unimplemented")
}

func (this *Import) SetEndDate(v *Time) {
	// TODO: unresolveReferenceTo all existing objects.
	// TODO: clone main objects (ID preserved)
	// TODO: set parents of cloned objects to `this` 
	// TODO: mask period in cloned events
	// TODO: save it to `cachedEvents`
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
	panic("Unimplemented")
}

func (this *Import) SetOffsetDate(v *Time) {
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	panic("Unimplemented")
}

func (this *Import) SetCategories(v *[]string) {
	// TODO: unresolveReferenceTo all imported objects.
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
	panic("Unimplemented")
}

//-- method (import and export)

func (this *Import) Import(feobj *file.Import) (warning.Warnings, error) {
	// TODO
	panic("Unimplemented")
}

func (this *Import) Export() (*file.Import, warning.Warnings, error) {
	// TODO
	panic("Unimplemented")
}

//-- method (creation)

func createImportFromParsed(o file.Import) (*Import, warning.Warnings) {
	// TODO
	panic("Unimplemented")
}

func CreateEmptyImport() *Import {
	return &Import{
		common: CreateEmptyMainObjCommon(MOK_Import),

		Note: Note{},

		link: "",
		startDate: nil,
		endDate: nil,
		offsetDate: nil,
		categories: nil,

		loadedChart: nil,
		importedEvents: []*Event{},
		importedAnnexs: []*Annex{},
		importedImports: []*Import{},

		preOffsetEvents: []*Event{},
	}
}
