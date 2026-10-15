package libltc

import (
	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

type Import struct {
	MajorCommon

	Note Note

	link string
	loadedChart *Chart
	offsetDate *Time			// Unfortunately, using it assumes Gregorian calendars.
	cache importChartCache
	mask importChartMask
}

type importChartCache struct {
	events []*Event
	annexs []*Annex
	imports []*Import
}

type importChartMask struct {
	startDate *Time
	endDate *Time
	importCates *[]string		// nil: any categories
	exclCates []string
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
	if this.mask.startDate != nil {
		return *this.mask.startDate
	} else {
		return CreateEmptyTime(TUK_Start) 
	}
}

func (this *Import) GetEndDate() Time {
	if this.mask.endDate != nil {
		return *this.mask.endDate
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
	return this.mask.importCates
}

func (this *Import) GetExclCategories() []string {
	return this.mask.exclCates
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

func (this *Import) SetLink(v string) (warns warning.Warnings) {
	var err error

	// Unexport itself from the chart.
	moreWarns := this.unexportMaskedObjects()
	warns.Concat(moreWarns)

	// Reload `loadedChart`
	moreWarns, err = this.reloadLink()
	warns.Concat(moreWarns)
	if err != nil {
		return
	}

	// Offset all imported objects
	moreWarns = this.offsetLoadedObjects()
	warns.Concat(moreWarns)

	// Export itself to the chart.
	moreWarns = this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) SetStartDate(v Time) (warns warning.Warnings) {
	// TODO: unexport itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: update `startDate` to 'v' (fix order, etc).
	// TODO: export itself to the chart.
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) UnsetStartDate() (warns warning.Warnings) {
	// TODO: unexport itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: update `startDate` to InfinitePast.
	// TODO: export itself to the chart.
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) SetEndDate(v Time) (warns warning.Warnings) {
	// TODO: unexport itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: update `endDate` to 'v' (fix order, etc).
	// TODO: export itself to the chart.
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) UnsetEndDate() (warns warning.Warnings) {
	// TODO: unexport itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: update `endDate` to InfiniteFuture.
	// TODO: export itself to the chart.
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) SetOffsetDate(v Time) (warns warning.Warnings) {
	// TODO: return if `offsetDate` is the same.
	// TODO: un-offset all imported objects.
	// TODO: update 'offsetDate'.
	// TODO: offset all imported objects
	panic("Unimplemented")
}

func (this *Import) UnsetOffsetDate() (warns warning.Warnings) {
	// TODO: return if `offsetDate` is nil
	// TODO: un-offset all imported objects.
	// TODO: update 'offsetDate'.
	// TODO: offset all imported objects
	panic("Unimplemented")
}

func (this *Import) SetImportCategories(v []string) (warns warning.Warnings) {
	// TODO: return if 'v' equals to `categories`.
	// TODO: unexport itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: set `categories`. 
	// TODO: export itself to the chart.
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) AddImportCategory(cat string) (warns warning.Warnings) {
	// TODO: return if 'cat' already exists in `categories`.
	// TODO: add 'cat' to `categories`. 
	// TODO: return if no object in 'cat'.
	// TODO: export 'cat' objects to the chart.
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) RemoveImportCategory(cat string) (warns warning.Warnings) {
	// TODO: return if 'cat' doesn't exist in `categories`.
	// TODO: remove 'cat' from `categories`. 
	// TODO: return if no object in 'cat'.
	// TODO: unexport 'cat' objects from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	panic("Unimplemented")
}

func (this *Import) UnsetImportCategories() (warns warning.Warnings) {
	// TODO: return if `categories` is nil.
	// TODO: unexport itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: nil-ify `categories`. 
	// TODO: export itself to the chart.
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) SetExclCategories(v []string) (warns warning.Warnings) {
	// TODO: return if 'v' equals to `exclCategories`.
	// TODO: unexport itself from the chart.
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	// TODO: set `exclCategories`. 
	// TODO: export itself to the chart.
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) AddExclCategory(cat string) (warns warning.Warnings) {
	// TODO: return if 'cat' already exists in `exclCategories`.
	// TODO: add 'cat' to `exclCategories`. 
	// TODO: return if no object in 'cat'.
	// TODO: export 'cat' objects to the chart.
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) RemoveExclCategory(cat string) (warns warning.Warnings) {
	// TODO: return if 'cat' doesn't exist in `exclCategories`.
	// TODO: remove 'cat' from `exclCategories`. 
	// TODO: return if no object in 'cat'.
	// TODO: unexport 'cat' objects from the chart.
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
		cimport.mask.startDate = &startDate
	}

	if o.EndDate != nil {
		endDate := createTimeFromParsed(o.EndDate, TUK_End)
		cimport.mask.endDate = &endDate
	}

	if o.OffsetDate != nil {
		endDate := createTimeFromParsed(o.OffsetDate, TUK_General)
		cimport.offsetDate = &endDate
	}

	if o.Categories != nil {
		tmpcats := *o.Categories
		cimport.mask.importCates = &tmpcats
	}

	if exclcats, ok := cimport.attrs["ExcludeCategory"]; ok {
		for _, exclcat := range exclcats {
			cimport.mask.exclCates = append(cimport.mask.exclCates, exclcat)
		}
	}

	moreWarns = cimport.DiagnoseLocal()
	warns.Concat(moreWarns)

	// Eagerly load linked chart.
	moreWarnsR, errR := cimport.reloadLink();
	warns.Concat(moreWarnsR)
	if errR != nil {
		moreWarnsO := cimport.offsetLoadedObjects()
		warns.Concat(moreWarnsO)

		moreWarnsE := cimport.exportMaskedObjects()
		warns.Concat(moreWarnsE)
	}

	return cimport, warns
}

func CreateEmptyImport() *Import {
	cimport := &Import{
		MajorCommon: *CreateEmptyMajorCommon(),

		Note: Note{},

		link: "",
		loadedChart: nil,
		offsetDate: nil,
		cache: importChartCache{},
		mask: importChartMask{},
	}
	cimport.enclosing = cimport
	return cimport
}

//-- method (util)

func (this *Import) correctStartEndDates(oldStartDate Time) (ambiguous bool, inverted bool) {
	nts, nte, amb, inv := correctStartEndDatePair(this.GetStartDate(), this.GetEndDate())

	// Reflect corrected results.
	if !this.GetStartDate().Equal(nts) {
		this.mask.startDate = &nts
	}

	if !this.GetEndDate().Equal(nte) {
		this.mask.endDate = &nte
	}

	ambiguous = amb
	inverted = inv

	return
}

func (this *Import) offsetLoadedObjects() (warns warning.Warnings) {
	// TODO: things to do:
	// TODO:  - offset all objects in the chart.
	panic("Unimplemented")
}

func (this *Import) unoffsetLoadedObjects() (warns warning.Warnings) {
	// TODO: things to do:
	// TODO:  - un-offset all objects in the chart.
	panic("Unimplemented")
}

func (this *Import) exportMaskedObjects() (warns warning.Warnings) {
	// TODO: cache first. resolve only new caches.
	// TODO: resolve: only if associated to a chart.
	// TODO: things to do:
	// TODO:  - cache objects with period/category/exclCategory masks
	// TODO:    use the offset period mask
	// TODO:  - resolveReferenceTo objects.
	panic("Unimplemented")
}

func (this *Import) unexportMaskedObjects() (warns warning.Warnings) {
	// TODO: cache first. unresolve only old caches. 
	// TODO: unresolve: only if associated to a chart.
	// TODO: things to do:
	// TODO:  - unresolveReferenceTo objects.
	// TODO:  - shot down caches.
	panic("Unimplemented")
}

func (this *Import) reloadLink() (warns warning.Warnings, err error) {
	// TODO: things to do:
	// TODO:  - (re)create `loadedChart`
	// TODO:  - set parents of imported objects to `this` 
	panic("Unimplemented")
	return
}
