package libltc

import (
	"slices"

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

	if this.offsetDate != nil {
		t := *this.offsetDate
		if t.IsInfinite() || t.isFloating() {
			this.offsetDate = nil
			warns.Add("@0@ had a wrong offset date '%s'. changed to no offset.", this, t.String())
		}
	}
       
	amb, inv := this.correctStartEndDateMasks(this.GetStartDate())
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
		// Only check against the object on the same import level.
		if iobj != this && iobj.parent == this.parent && iobj.numID == this.numID {
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

func (this *Import) GetAllImportedEvents() []*Event {
	return this.cache.events
}

func (this *Import) GetAllImportedAnnexs() []*Annex {
	return this.cache.annexs
}

func (this *Import) GetAllImportedImports() []*Import {
	return this.cache.imports
}

func (this *Import) GetDirectImportedEvents() []*Event {
	ret := []*Event{}
	for _, obj := range this.cache.events {
		if obj.GetParent() == this {
			ret = append(ret, obj)
		}
	}
	return ret
}

func (this *Import) GetDirectImportedAnnexs() []*Annex {
	ret := []*Annex{}
	for _, obj := range this.cache.annexs {
		if obj.GetParent() == this {
			ret = append(ret, obj)
		}
	}
	return ret
}

func (this *Import) GetDirectImportedImports() []*Import {
	ret := []*Import{}
	for _, obj := range this.cache.imports {
		if obj.GetParent() == this {
			ret = append(ret, obj)
		}
	}
	return ret
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
	var moreWarns warning.Warnings

	if this.link == v {
		return
	}

	this.link = v

	// Unexport itself from the chart.
	moreWarns = this.unexportMaskedObjects()
	warns.Concat(moreWarns)

	// Reload `loadedChart`
	moreWarns = this.reloadLink()
	warns.Concat(moreWarns)

	// Offset all imported objects
	moreWarns = this.offsetLoadedObjects()
	warns.Concat(moreWarns)

	// Export itself to the chart.
	moreWarns = this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) SetStartDate(t Time) (warns warning.Warnings) {
	var moreWarns warning.Warnings

	if this.mask.startDate != nil && this.mask.startDate.Equal(t) {
		return
	}

	// Unexport itself from the chart.
	moreWarns = this.unexportMaskedObjects()
	warns.Concat(moreWarns)

	// Update `startDate` to 't' (fix order, etc).
	oldStartDate := this.GetStartDate()
	t.SetUsage(TUK_Start)
    this.mask.startDate = &t
	this.sanitizeStartEndDateMasks(true)
	this.correctStartEndDateMasks(oldStartDate)

	// Export itself to the chart.
	moreWarns = this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) UnsetStartDate() (warns warning.Warnings) {
	var moreWarns warning.Warnings

	if this.mask.startDate == nil {
		return
	}

	// Unexport itself from the chart.
	moreWarns = this.unexportMaskedObjects()
	warns.Concat(moreWarns)

	// Update `startDate` to floating.
	oldStartDate := this.GetStartDate()
    this.mask.startDate = nil 
	this.sanitizeStartEndDateMasks(true)
	this.correctStartEndDateMasks(oldStartDate)

	// Export itself to the chart.
	moreWarns = this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) SetEndDate(t Time) (warns warning.Warnings) {
	var moreWarns warning.Warnings

	if this.mask.endDate != nil && this.mask.endDate.Equal(t) {
		return
	}

	// Unexport itself from the chart.
	moreWarns = this.unexportMaskedObjects()
	warns.Concat(moreWarns)

	// Update `endDate` to 't' (fix order, etc).
	oldStartDate := this.GetStartDate()
	t.SetUsage(TUK_End)
    this.mask.endDate = &t
	this.sanitizeStartEndDateMasks(false)
    this.correctStartEndDateMasks(oldStartDate)

	// Export itself to the chart.
	moreWarns = this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) UnsetEndDate() (warns warning.Warnings) {
	var moreWarns warning.Warnings

	if this.mask.endDate == nil {
		return
	}

	// Unexport itself from the chart.
	moreWarns = this.unexportMaskedObjects()
	warns.Concat(moreWarns)

	// Update `endDate` to floating.
	oldStartDate := this.GetStartDate()
    this.mask.endDate = nil
	this.sanitizeStartEndDateMasks(false)
    this.correctStartEndDateMasks(oldStartDate)

	// Export itself to the chart.
	moreWarns = this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) SetOffsetDate(t Time) (warns warning.Warnings) {
	var moreWarns warning.Warnings

	if this.offsetDate != nil && this.offsetDate.Equal(t) {
		return
	}
		
	if t.IsInfinite() || t.isFloating() {
		warns.Add("cannot set a wrong offset date '%s'.", this, t.String())
		return
	}

	// Un-offset all imported objects.
	moreWarns = this.unoffsetLoadedObjects()
	warns.Concat(moreWarns)

	// Update 'offsetDate'.
	this.offsetDate = &t

	// Offset all imported objects.
	moreWarns = this.offsetLoadedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) UnsetOffsetDate() (warns warning.Warnings) {
	var moreWarns warning.Warnings

	if this.offsetDate == nil {
		return
	}

	// Un-offset all imported objects.
	moreWarns = this.unoffsetLoadedObjects()
	warns.Concat(moreWarns)

	// Unset 'offsetDate'.
	this.offsetDate = nil

	// Offset all imported objects.
	moreWarns = this.offsetLoadedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) SetImportCategories(v []string) (warns warning.Warnings) {
	var moreWarns warning.Warnings

	// Return if it's the same.
	if this.mask.importCates != nil && slices.Equal(*this.mask.importCates, v) {
		return
	}

	// Unexport itself from the chart.
	moreWarns = this.unexportMaskedObjects()
	warns.Concat(moreWarns)

	// Set `importCates`. 
    this.mask.importCates = &v 

	// Re-export.
	moreWarns = this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) AddImportCategory(cat string) (warns warning.Warnings) {
	// Return if it's already in.
	if this.mask.importCates != nil && slices.Contains(*this.mask.importCates, cat) {
		return
	}

	// Add 'cat' to `categories`. 
	if this.mask.importCates == nil {
		this.mask.importCates = &[]string{ cat }
	} else {
		*this.mask.importCates = append(*this.mask.importCates, cat)
	}

	// Return if the category is empty.
	if this.loadedChart != nil && len(this.loadedChart.GetEventsInCategory(cat)) == 0 {
		return
	}

	// Re-export.
	moreWarns := this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) RemoveImportCategory(cat string) (warns warning.Warnings) {
	// Remove 'cat' from `categories`. 
	if this.mask.importCates != nil {
		idx := slices.Index(*this.mask.importCates, cat)
		if idx != -1 {
			*this.mask.importCates = append((*this.mask.importCates)[:idx], (*this.mask.importCates)[idx+1:]...)
		} else {
			return
		}
	} else {
		return
	}

	// Return if the category is empty.
	if this.loadedChart != nil && len(this.loadedChart.GetEventsInCategory(cat)) == 0 {
		return
	}

	// Re-export.
	moreWarns := this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) UnsetImportCategories() (warns warning.Warnings) {
	var moreWarns warning.Warnings

	// Return if it's nil.
	if this.mask.importCates == nil {
		return
	}

	// Unexport itself from the chart.
	moreWarns = this.unexportMaskedObjects()
	warns.Concat(moreWarns)

	// Set `importCates`. 
    this.mask.importCates = nil 

	// Re-export.
	moreWarns = this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) SetExclCategories(v []string) (warns warning.Warnings) {
	var moreWarns warning.Warnings

	// Return if it's the same.
	if slices.Equal(this.mask.exclCates, v) {
		return
	}

	// Unexport itself from the chart.
	moreWarns = this.unexportMaskedObjects()
	warns.Concat(moreWarns)

	// Set `exclCates`. 
    this.mask.exclCates = v 

	// Re-export.
	moreWarns = this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) AddExclCategory(cat string) (warns warning.Warnings) {
	// Return if it's already in.
	if slices.Contains(this.mask.exclCates, cat) {
		return
	}

	// Add 'cat' to `exclCates`. 
	this.mask.exclCates = append(this.mask.exclCates, cat)

	// Return if the category is empty.
	if this.loadedChart != nil && len(this.loadedChart.GetEventsInCategory(cat)) == 0 {
		return
	}

	// Re-export.
	moreWarns := this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
}

func (this *Import) RemoveExclCategory(cat string) (warns warning.Warnings) {
	// Remove 'cat' from `categories`. 
	idx := slices.Index(this.mask.exclCates, cat)
	if idx != -1 {
		this.mask.exclCates = append(this.mask.exclCates[:idx], this.mask.exclCates[idx+1:]...)
	} else {
		return
	}

	// Return if the category is empty.
	if this.loadedChart != nil && len(this.loadedChart.GetEventsInCategory(cat)) == 0 {
		return
	}

	// Re-export.
	moreWarns := this.exportMaskedObjects()
	warns.Concat(moreWarns)

	return
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
	moreWarns = cimport.reloadLink()
	warns.Concat(moreWarns)
	moreWarns = cimport.offsetLoadedObjects()
	warns.Concat(moreWarns)
	moreWarns = cimport.exportMaskedObjects()
	warns.Concat(moreWarns)

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

func (this *Import) sanitizeStartEndDateMasks(fixedStart bool) {
	nts, nte, sanitized := sanitizeStartEndDatePair(this.GetStartDate(), this.GetEndDate(), fixedStart)

	if sanitized {
		if fixedStart {
			this.mask.endDate = &nte
		} else {
			this.mask.startDate = &nts
		}
	}
}

func (this *Import) correctStartEndDateMasks(oldStartDate Time) (ambiguous bool, inverted bool) {
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
	if this.loadedChart == nil || this.offsetDate == nil {
		return
	} else {
		toff := *this.offsetDate
		
		if toff.IsInfinite() || toff.isFloating() {
			warns.Add("!!!INTERNAL WARN!!! @0@ has a wrong offset date '%s'. cannot offset.", this, toff.String())
			return
		}

		for _, eobj := range this.loadedChart.GetEvents() {
			eobj.shiftStartEndDatesRaw(toff)
		}
		for _, aobj := range this.loadedChart.GetAnnexs() {
			aobj.shiftCreateDateRaw(toff)
		}
		// NOTE: nested imported objects must have been included in the objects
		// returned by the 'loadedChart', and the import objects themselves 
		// don't have to be offset.
		return
	}
}

func (this *Import) unoffsetLoadedObjects() (warns warning.Warnings) {
	if this.loadedChart == nil || this.offsetDate == nil {
		return
	} else {
		toff := this.offsetDate.Negate()
		
		if toff.IsInfinite() || toff.isFloating() {
			warns.Add("!!!INTERNAL WARN!!! @0@ has a wrong offset date '%s'. cannot offset.", this, toff.String())
			return
		}
		
		for _, eobj := range this.loadedChart.GetEvents() {
			eobj.shiftStartEndDatesRaw(toff)
		}
		for _, aobj := range this.loadedChart.GetAnnexs() {
			aobj.shiftCreateDateRaw(toff)
		}
		return
	}
}

func (this *Import) exportMaskedObjects() (warns warning.Warnings) {
	// FIXME: cache first. resolve only new caches.

	if this.loadedChart == nil {
		return
	}

	this.cache = importChartCache{}

	maskStartDate, maskEndDate := InfinitePast, InfiniteFuture
	if this.mask.startDate != nil {
		maskStartDate = *this.mask.startDate
	}
	if this.mask.endDate != nil {
		maskEndDate = *this.mask.endDate
	}

	// Cache events.
	for _, category := range this.loadedChart.GetEventCategories() {
		if slices.Contains(this.mask.exclCates, category) {
			continue
		}
		if this.mask.importCates != nil && !slices.Contains(*this.mask.importCates, category) {
			continue
		}
		nevs := this.loadedChart.GetEventsInCategoryBetween(category, maskStartDate, maskEndDate)
		this.cache.events = append(this.cache.events, nevs...)
	}

	// Cache annexs.
	this.cache.annexs = this.loadedChart.GetAnnexsBetween(maskStartDate, maskEndDate)

	// Cache imports. (include all)
	this.cache.imports = this.loadedChart.GetImports()

	if this.chart != nil {
		this.chart.addImportedObjects(this)
	}

	return
}

func (this *Import) unexportMaskedObjects() (warns warning.Warnings) {
	// FIXME: cache first. unresolve only old caches. 

	if this.chart != nil {
		this.chart.removeImportedObjects(this)
	}
	
	this.cache = importChartCache{}

	return
}

func (this *Import) reloadLink() (warns warning.Warnings) {
	chart, moreWarns := CreateChart(this.link)
	warns.Concat(moreWarns)

	for _, mobj := range chart.GetAllMajors(false) {
		// 'parent == nil' means it's on the top level. Make 'this' the new parent.
		if mobj.GetParent() == nil {
			mobj.setParent(this)
		}
	}

	this.loadedChart = chart

	return
}
