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
	offsetDate *Time			// Unfortunately, using it assumes Gregorian calendars.
	categories *[]string		// nil: any categories
	exclCategories []string
	
	loadedChart *Chart
	cache importCache
}

type importCache struct {
	events []*Event
	annexs []*Annex
	imports []*Import
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

	// Check invalid ID.
	if this.numID == mNID_Invalid {
		newNumID := this.chart.getNextNumberID(MOK_Import)
		this.numID = newNumID
		warns.Add("@0@ had an invalid ID. auto-corrected to 'i%d'.", this, newNumID)
	}

	// Check duplicate ID.
	for _, iobj := range this.chart.GetImports() {
		if iobj != this && iobj.numID == this.numID {
			this.AddAttr("OldID", this.GetIntrinsicID())
			newNumID := this.chart.getNextNumberID(MOK_Import)
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


func (this *Import) GetImportCategories() *[]string {
	return this.categories
}

func (this *Import) GetExclCategories() []string {
	return this.exclCategories
}

func (this *Import) GetImportedChart() *Chart {
	return this.loadedChart
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
	// TODO: dissociate itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: (re)load `loadedChart`
	// TODO:  - (re)create `loadedChart`
	// TODO:  - set parents of imported objects to `this` 
	// TODO: offset all imported objects
	// TODO: associate itself to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) SetStartDate(v Time) {
	// TODO: dissociate itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: update `startDate` to 'v' (fix order, etc).
	// TODO: associate itself to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) UnsetStartDate() {
	// TODO: dissociate itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: update `startDate` to InfinitePast.
	// TODO: associate itself to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) SetEndDate(v Time) {
	// TODO: dissociate itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: update `endDate` to 'v' (fix order, etc).
	// TODO: associate itself to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) UnsetEndDate() {
	// TODO: dissociate itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: update `endDate` to InfiniteFuture.
	// TODO: associate itself to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) SetOffsetDate(v Time) {
	// TODO: return if `offsetDate` is the same.
	// TODO: un-offset all imported objects.
	// TODO: update 'offsetDate'.
	// TODO: offset all imported objects
	panic("Unimplemented")
}

func (this *Import) UnsetOffsetDate() {
	// TODO: return if `offsetDate` is nil
	// TODO: un-offset all imported objects.
	// TODO: update 'offsetDate'.
	// TODO: offset all imported objects
	panic("Unimplemented")
}

func (this *Import) SetImportCategories(v []string) {
	// TODO: return if 'v' equals to `categories`.
	// TODO: dissociate itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: set `categories`. 
	// TODO: associate itself to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) AddImportCategory(cat string) {
	// TODO: return if 'cat' already exists in `categories`.
	// TODO: add 'cat' to `categories`. 
	// TODO: return if no object in 'cat'.
	// TODO: associate 'cat' objects to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) RemoveImportCategory(cat string) {
	// TODO: return if 'cat' doesn't exist in `categories`.
	// TODO: remove 'cat' from `categories`. 
	// TODO: return if no object in 'cat'.
	// TODO: dissociate 'cat' objects from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	panic("Unimplemented")
}

func (this *Import) UnsetImportCategories() {
	// TODO: return if `categories` is nil.
	// TODO: dissociate itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: nil-ify `categories`. 
	// TODO: associate itself to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) SetExclCategories(v []string) {
	// TODO: return if 'v' equals to `exclCategories`.
	// TODO: dissociate itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: set `exclCategories`. 
	// TODO: associate itself to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) AddExclCategory(cat string) {
	// TODO: return if 'cat' already exists in `exclCategories`.
	// TODO: add 'cat' to `exclCategories`. 
	// TODO: return if no object in 'cat'.
	// TODO: associate 'cat' objects to the chart.
	// TODO:  - cache objects with period/category masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) RemoveExclCategory(cat string) {
	// TODO: return if 'cat' doesn't exist in `exclCategories`.
	// TODO: remove 'cat' from `exclCategories`. 
	// TODO: return if no object in 'cat'.
	// TODO: dissociate 'cat' objects from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
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

	// TODO: populate 'exclCategories'.
	// TODO: eagerly load imported objects. warn if it failed.
	// TODO: exclude 'exclCategories'.
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

func (this *Import) loadLink() (warns warning.Warnings) {
	if this.loadedChart == nil {
		// TODO: read above
		panic("Unimplemented")
	}
	return
}
