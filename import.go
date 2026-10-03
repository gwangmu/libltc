package libltc

import (
	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

type Import struct {
	MainObjCommon

	Note Note

	link string
	startDate *Time
	endDate *Time
	offsetDate *Time
	categories []string		// nil: any categories
	
	loadedChart *Chart
	importedEvents []*Event
	importedAnnexs []*Annex
	importedImports []*Import

	preOffsetEvents []*Event
}

//-- interface IWarningObj

func (this *Import) Summary() string {
	// TODO
	panic("Unimplemented")
}

func (this Import) IsUnknown() bool {
	// TODO
	panic("Unimplemented")
}

//-- interface IDiagnosable

func (this *Import) DiagnoseLocal() (warns warning.Warnings) {
	// TODO: recursive to 'loadedChart'
	panic("Unimplemented")
}

func (this *Import) DiagnoseNonLocal() (warns warning.Warnings) {
	if this.GetChart() == nil {
		warns.Add("@0@ is not associated to any chart.", this)
		return
	}

	// TODO: check duplicate ID between 'this' and other chart objs.
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

func (this *Import) GetCategories() []string {
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

func (this *Import) SetCategories(v []string) {
	// TODO: unresolveReferenceTo all imported objects.
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
	panic("Unimplemented")
}

//-- method (creation)

func createImportFromParsed(o file.Import) (*Import, warning.Warnings) {
	cimport := CreateEmptyImport()
	warns := warning.Warnings{}

	kind, numID := convIDStringToInternal(o.ID)
	cimport.numID = numID
	cimport.attrs = convAttrsFileToChart(o.Attrs)

	if kind != MOK_Import {
		warns.Add("@0@ has a wrong kind prefix. Fixing...")
	}

	note, moreWarns := CreateNoteFromString(o.Note)
	warns.Concat(moreWarns)
	cimport.Note = note

	cimport.link = o.Link
	// TODO: eagerly load imported objects. warn if it failed.
	panic("Unimplemented")

	if o.StartDate != nil {
		startDate := createTimeFromParsed(o.StartDate, TUK_Start)
		cimport.startDate = &startDate
	}

	if o.EndDate != nil {
		endDate := createTimeFromParsed(o.EndDate, TUK_End)
		cimport.endDate = &endDate
	}

	if o.OffsetDate != nil {
		endDate := createTimeFromParsed(o.OffsetDate, TUK_General)
		cimport.offsetDate = &endDate
	}

	if o.Categories != nil {
		cimport.categories = append([]string{}, o.Categories...)
	}

	moreWarns = cimport.DiagnoseLocal()
	warns.Concat(moreWarns)

	return cimport, warns
}

func CreateEmptyImport() *Import {
	cimport := &Import{
		MainObjCommon: *CreateEmptyMainObjCommon(),

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
	cimport.enclosing = cimport
	return cimport
}
