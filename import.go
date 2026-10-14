package libltc

import (
	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

type Import struct {
	MajorCommon

	Note Note

	link string
	startDate *Time
	endDate *Time
	offsetDate *Time
	categories *[]string		// nil: any categories
	
	loadedChart *Chart
	cache importCache
}

type importCache struct {
	events []*Event
	annexs []*Annex
	imports []*Import
	preOffsetEvents []*Event
	initialWarns warning.Warnings
}

func (this *Import) Equal(that IMajor) bool {
	return this.GetCommon() == that.GetCommon()
}

//-- interface IMajor

func (this *Import) resolveReferenceTo(c *Chart) {
	// No back- or cross-reference. So, nothing.
	return
}

func (this *Import) unresolveReferenceTo(c *Chart) {
	// No back- or cross-reference. So, nothing.
	return
}

//-- interface IWarningObj

func (this *Import) Summary() (ret string) {
	if this.IsUnresolved() {
		ret = "an unresolved import"
	} else {
		ret = "an import"
	}

	maybeOldID := ""
	if this.HasAttrKey("OldID") && len(this.GetAttrs("OldID")) > 0 {
		maybeOldID = this.GetAttrs("OldID")[0]
	}

	extraStr := ""
	if this.IsUnresolved() {
		extraStr = warning.BuildSummaryString(
			wse{"ID", this.getUnresRelQualID()},
			wse{"OldID", maybeOldID},
		)
	} else {
		extraStr = warning.BuildSummaryString(
			wse{"ID", this.GetIntrinsicID()},
			wse{"OldID", maybeOldID},
			wse{"Link", this.GetLink()},
			wse{"Started", this.GetStartDate()}, 
			wse{"Ended", this.GetEndDate()},
			wse{"Offset", this.GetOffsetDate()},
		)
	}
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}
	
	return
}

func (this Import) IsUnknown() bool {
	return this.GetIntrinsicID() == "" && this.GetStartDate().IsUnknown() &&
		this.GetEndDate().IsUnknown() && this.GetOffsetDate().IsUnknown() &&
		this.GetLink() == ""
}

//-- interface IDiagnosable

func (this *Import) DiagnoseLocal() (warns warning.Warnings) {
	if this.IsLoaded() {
		warns.Concat(this.cache.initialWarns)
	} else {
		warns.Add("@0@ cannot load the linked chart.", this)
	}

	if this.GetStartDate().IsInfiniteFuture() && this.GetEndDate().isFloating() {
		this.SetEndDate(InfiniteFuture)
	} else if this.GetEndDate().IsInfinitePast() && this.GetStartDate().isFloating() {
		this.SetStartDate(InfinitePast)
	}
       
	amb, inv := this.correctStartEndDates(this.GetStartDate())
	if amb {
        warns.Add("@0@ had an ambiguous date. corrected.", this)
    }
	if inv {
        warns.Add("@0@ had an inverted start/end date pair. corrected.", this)
    }

	return
}

func (this *Import) DiagnoseNonLocal() (warns warning.Warnings) {
	if this.GetChart() == nil {
		warns.Add("@0@ is not associated to any chart.", this)
		return
	}

	// TODO: move this to MajorCommon
	// Check invalid ID.
	if this.numID == mNID_Invalid {
		newNumID := this.chart.getNextmNumberID(MOK_Import)
		this.numID = newNumID
		warns.Add("@0@ had an invalid ID. auto-corrected to 'i%d'.", this, newNumID)
	}

	// TODO: move this to MajorCommon
	// Check duplicate ID.
	for _, iobj := range this.chart.GetImports() {
		if iobj != this && iobj.numID == this.numID {
			this.AddAttr("OldID", this.GetIntrinsicID())
			newNumID := this.chart.getNextmNumberID(MOK_Import)
			this.numID = newNumID
			warns.Add("@0@ had a duplicated ID. auto-corrected to 'i%d'.", this, newNumID)
			break
		}
	}

	return
}

//-- method (getters)

func (this *Import) GetLink() string {
	return this.link
}

func (this *Import) GetStartDate() Time {
	if this.startDate != nil {
		return *this.startDate
	} else {
		return CreateEmptyTime(TUK_Start) 
	}
}

func (this *Import) GetEndDate() Time {
	if this.endDate != nil {
		return *this.endDate
	} else {
		return CreateEmptyTime(TUK_End) 
	}
}

func (this *Import) GetOffsetDate() Time {
	if this.offsetDate != nil {
		return *this.offsetDate
	} else {
		return CreateEmptyTime(TUK_General) 
	}
}

func (this *Import) GetImportedChart() *Chart {
	return this.loadedChart
}

func (this *Import) GetCategories() *[]string {
	return this.categories
}

func (this *Import) GetAllImportedEvents() []*Event {
	return this.cache.events
}

func (this *Import) GetAllImportedAnnexs() []*Annex {
	return this.cache.annexs
}

func (this *Import) GetAllImportedImports() []*Import {
	return this.cache.imports
}

func (this *Import) GetAllImportedObjects() []IMajor {
	ret := []IMajor{}
	for _, obj := range this.cache.events {
		ret = append(ret, obj)
	}
	for _, obj := range this.cache.annexs {
		ret = append(ret, obj)
	}
	for _, obj := range this.cache.imports {
		ret = append(ret, obj)
	}
	return ret
}

func (this *Import) IsLoaded() bool {
	return this.loadedChart != nil
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

func (this *Import) SetStartDate(v Time) {
	// TODO: unresolveReferenceTo all imported objects.
	// TODO: clone main objects (ID preserved)
	// TODO: set parents of cloned objects to `this` 
	// TODO: mask period in cloned events
	// TODO: save it to `cachedEvents`
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
	panic("Unimplemented")
}

func (this *Import) SetEndDate(v Time) {
	// TODO: unresolveReferenceTo all existing objects.
	// TODO: clone main objects (ID preserved)
	// TODO: set parents of cloned objects to `this` 
	// TODO: mask period in cloned events
	// TODO: save it to `cachedEvents`
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
	panic("Unimplemented")
}

func (this *Import) SetOffsetDate(v Time) {
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	panic("Unimplemented")
}

func (this *Import) SetCategories(v []string) {
	// TODO: unresolveReferenceTo all imported objects.
	// TODO: clone `cachedEvents`, mask categories, apply offsetDate
	// TODO: resolveReferenceTo all imported objects.
	panic("Unimplemented")
}

func (this *Import) UnsetStartDate() {
	// TODO: read above
	panic("Unimplemented")
}

func (this *Import) UnsetEndDate() {
	// TODO: read above
	panic("Unimplemented")
}

func (this *Import) UnsetOffsetDate() {
	// TODO: read above
	panic("Unimplemented")
}

func (this *Import) UnsetCategories() {
	// TODO: read above
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
		tmpcats := *o.Categories
		cimport.categories = &tmpcats
	}

	// TODO: eagerly load imported objects. warn if it failed.
	// TODO: exclude 'ExcludeCategory' categories.
	panic("Unimplemented")

	moreWarns = cimport.DiagnoseLocal()
	warns.Concat(moreWarns)

	return cimport, warns
}

func CreateEmptyImport() *Import {
	cimport := &Import{
		MajorCommon: *CreateEmptyMajorCommon(),

		Note: Note{},

		link: "",
		startDate: nil,
		endDate: nil,
		offsetDate: nil,
		categories: nil,

		loadedChart: nil,
		cache: importCache{},
	}
	cimport.enclosing = cimport
	return cimport
}

//-- method (util)

func (this *Import) correctStartEndDates(oldStartDate Time) (ambiguous bool, inverted bool) {
	nts, nte, amb, inv := correctStartEndDatePair(this.GetStartDate(), this.GetEndDate())

	// Reflect corrected results.
	if !this.GetStartDate().Equal(nts) {
		this.startDate = &nts
	}

	if !this.GetEndDate().Equal(nte) {
		this.endDate = &nte
	}

	ambiguous = amb
	inverted = inv

	return
}
