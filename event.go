package libltc

import (
	"regexp"

	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

type Event struct {
	MajorCommon

	Note Note 

	localTitle *string
	localStartDate *Time
	localEndDate *Time

	timeOffed bool
	preoffLocalStartDate *Time
	preoffLocalEndDate *Time

	category string

	chartEmbedLink *string
	embeddedChart *Chart
	eventEmbedLink *string
	embeddedEvent *Event

    ContinuedFrom	// Specifying 'ContinuedFrom'
    ContinuedTo		// Pointed by 'ContinuedFrom'
}

func (this *Event) Equal(that IMajor) bool {
	return this.GetCommon() == that.GetCommon()
}

//-- interface IMajor

func (this *Event) resolveReferenceTo(c *Chart) {
	for _, category := range c.GetEventCategories() {
		for _, eobj := range c.GetEventsInCategory(category) {
			this.ContinuedFrom.Link(eobj, false)	
		}
	}
}

func (this *Event) unresolveReferenceTo(c *Chart) {
	for _, category := range c.GetEventCategories() {
		for _, eobj := range c.GetEventsInCategory(category) {
			this.ContinuedFrom.Unlink(eobj)			
		}
	}
}

//-- interface IWarningObj

func (this *Event) Summary() (ret string) {
	if this.IsUnresolved() {
		ret = "an unresolved event"
	} else if this.GetTitle() != "" {
		ret = "the event '" + this.GetTitle() + "'"
	} else {
		ret = "a untitled event"
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
			wse{"Category", this.GetCategory()},
			wse{"Started", this.GetStartDate()}, 
			wse{"Ended", this.GetEndDate()},
			wse{"Embedded-event", this.GetEmbedEventLink()},
			wse{"Embedded-chart", this.GetEmbedChartLink()},
		)
	}
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}
	
	return
}

func (this Event) IsUnknown() bool {
	return this.GetTitle() == "" && this.GetIntrinsicID() == "" &&
		this.GetCategory() == "" && this.GetStartDate().IsUnknown() &&
		this.GetEndDate().IsUnknown() && this.GetEmbedEventLink() == "" &&
		this.GetEmbedChartLink() == ""
}

//-- interface IDiagnosable

func (this *Event) DiagnoseLocal() (warns warning.Warnings) {
	if this.IsEventEmbedding() && this.IsChartEmbedding() {
		warns.Add("@0@ attempted to embed both an event and a chart. Favoring event...", this)
		this.chartEmbedLink = nil
		this.embeddedChart = nil
	}

	for akey, avals := range this.attrs {
		if akey == "ContinuedFrom" {
			newAvals := []string{}
			for _, aval := range avals {
				if aval == "" {
					warns.Add("a 'ContinuedFrom' attribute in @0@ is empty. removed.", this)
				} else {
					reEmbedChart := regexp.MustCompile(`e[0-9]+/`)
					if reEmbedChart.MatchString(aval) {
						warns.Add("@0@ attempts to continue from a subchart event (%s), which is not recommended.", this, aval)
					}
					newAvals = append(newAvals, aval)
				}
			}
			this.attrs[akey] = newAvals
		}
	}

	if this.GetStartDate().IsInfiniteFuture() && this.GetEndDate().IsFloating() {
		this.SetEndDate(InfiniteFuture)
	} else if this.GetEndDate().IsInfinitePast() && this.GetStartDate().IsFloating() {
		this.SetStartDate(InfinitePast)
	}
       
	amb, inv := this.correctLocalStartEndDates(this.GetStartDate())
	if amb {
        warns.Add("@0@ had an ambiguous date. corrected.", this)
    }
	if inv {
        warns.Add("@0@ had an inverted start/end date pair. corrected.", this)
    }

	return
}

func (this *Event) DiagnoseNonLocal() (warns warning.Warnings) {
	if this.chart == nil {
		warns.Add("@0@ is not associated to any chart.", this)
		return
	}

	// Check invalid ID.
	if this.numID == mNID_Invalid {
		newNumID := this.chart.getNextNumberID(MOK_Event)
		this.numID = newNumID
		warns.Add("@0@ had an invalid ID. auto-corrected to 'e%d'.", this, newNumID)
	}

	// Check duplicate ID.
	for _, category := range this.chart.GetEventCategories() {
		for _, eobj := range this.chart.GetEventsInCategory(category) {
			// Only check against the object on the same import level.
			if eobj != this && eobj.parent == this.parent && eobj.numID == this.numID {
				this.AddAttr("OldID", this.GetIntrinsicID())
				newNumID := this.chart.getNextNumberID(MOK_Event)
				this.numID = newNumID
				warns.Add("@0@ had a duplicated ID. auto-corrected to 'e%d'.", this, newNumID)
				break
			}
		}
	}

	for akey, avals := range this.attrs {
		if akey == "ContinuedFrom" {
			for _, aval := range avals {
				var eobjFrom *Event
				for _, eobjIn := range this.ContinuedFrom.Get() {
					if eobjIn.GetRelQualifiedID(this.GetChart()) == aval {
						eobjFrom = eobjIn
						break
					}
				}

				if eobjFrom == nil || eobjFrom.IsUnresolved() {
					warns.Add("@0@ has a dangling 'ContinuedFrom' to '%s'.", this, aval)
				} 

				if eobjFrom != nil && !eobjFrom.IsUnresolved() && !eobjFrom.ContinuedTo.Has(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'ContinuedFrom' to '%s', but it doesn't reference back. corrected.", this, aval)
                    eobjFrom.ContinuedTo.add(this)
				}
			}
		}
	}

	return
}

//-- method (getters)

func (this *Event) GetTitle() string {
	// Use `localTitle` if it was defined.
	// Otherwise, use the title of `embeddedEvent` or `embeddedChart`.
	if this.localTitle != nil {
		return *this.localTitle
	} else if this.IsEventEmbedding() {
		return this.embeddedEvent.GetTitle()
	} else if this.IsChartEmbedding() {
		if sc, _ := this.GetSubchart(); sc != nil {
			return sc.GetStringifiedSubjectName()
		} else {
			return "(unloaded subchart name)"
		}
	} else {
		return ""
	}
}

func (this *Event) GetEmbeddedEvent() *Event {
	return this.embeddedEvent
}

func (this *Event) GetEmbedEventLink() string {
	if this.IsEventEmbedding() {
		return *this.eventEmbedLink
	} else {
		return ""
	}
}

func (this *Event) IsEventEmbedding() bool {
	return this.eventEmbedLink != nil
}

func (this *Event) IsSubchartLoaded() bool {
	return this.embeddedChart != nil
}

func (this *Event) GetSubchart() (*Chart, warning.Warnings) {
	warns := warning.Warnings{}

	// Lazy-load `embeddedChart` if it's nil (but shouldn't).
	if this.IsChartEmbedding() && this.embeddedChart == nil {
		moreWarns := this.loadEmbeddedChart()
		warns.Concat(moreWarns)
	}

	return this.embeddedChart, warns
}

func (this *Event) GetEmbedChartLink() string {
	if this.IsChartEmbedding() {
		return *this.chartEmbedLink
	} else {
		return ""
	}
}

func (this *Event) IsChartEmbedding() bool {
	return this.chartEmbedLink != nil
}

// NOTE: 'Get{Start,End}Date' should be either both floating or both non-floating.
func (this *Event) GetStartDate() Time {
	if this.localStartDate != nil {
		return this.localStartDate.Clone()
	} else if this.embeddedEvent != nil {
		return this.embeddedEvent.GetStartDate()
	} else {
		return CreateEmptyTime(TUK_Start) 
	}
}

// NOTE: 'Get{Start,End}Date' should be either both floating or both non-floating.
func (this *Event) GetEndDate() Time {
	if this.localEndDate != nil {
		return this.localEndDate.Clone()
	} else if this.embeddedEvent != nil {
		return this.embeddedEvent.GetEndDate()
	} else {
		return CreateEmptyTime(TUK_End) 
	}
}

func (this *Event) GetLocalStartDate() *Time {
	return this.localStartDate
}

func (this *Event) GetLocalEndDate() *Time {
	return this.localEndDate
}

func (this *Event) HasLocalStartDate() bool {
	return this.localStartDate != nil
}

func (this *Event) HasLocalEndDate() bool {
	return this.localEndDate != nil
}

func (this *Event) GetCategory() string{
	return this.category
}

//-- method (setters)

func (this *Event) SetTitle(t string) {
	this.localTitle = &t
}

func (this *Event) UnsetLocalTitle() {
	this.localTitle = nil
}

func (this *Event) SetEmbedEventLink(uri string) (warns warning.Warnings) {
	this.eventEmbedLink = &uri

	// Eagerly load embedded event. Warn if it failed.
	moreWarns := this.loadEmbeddedEvent()
	warns.Concat(moreWarns)

	return
}

func (this *Event) SetEmbedChartLink(uri string) {
	if this.embeddedChart != nil {
		this.unloadEmbeddedChart()
	}

	// Lazy-load subchart.
	this.chartEmbedLink = &uri
	this.embeddedChart = nil
}

func (this *Event) UnsetLocalStartDate() {
	oldStartDate := this.GetStartDate()
	this.localStartDate = nil
	this.sanitizeLocalStartEndDates(true)
    this.correctLocalStartEndDates(oldStartDate)
}

func (this *Event) UnsetLocalEndDate() {
	oldStartDate := this.GetStartDate()
	this.localEndDate = nil
	this.sanitizeLocalStartEndDates(false)
    this.correctLocalStartEndDates(oldStartDate)
}

// Combination of 'UnsetLocalStartDate' and 'UnsetLocalEndDate'.
func (this *Event) UnsetLocalDate() {
	this.UnsetLocalStartDate()
	this.UnsetLocalEndDate()
}

// Convenience wrapper to 'SetLocalStartDate'.
func (this *Event) SetStartDate(t Time) {
	oldStartDate := this.GetStartDate()
	t.SetUsage(TUK_Start)
    this.localStartDate = &t
	this.sanitizeLocalStartEndDates(true)
	this.correctLocalStartEndDates(oldStartDate)
}

// Convenience wrapper to 'SetLocalEndDate'.
func (this *Event) SetEndDate(t Time) {
	oldStartDate := this.GetStartDate()
	t.SetUsage(TUK_End)
    this.localEndDate = &t
	this.sanitizeLocalStartEndDates(false)
    this.correctLocalStartEndDates(oldStartDate)
}

// Combination of 'SetStartDate' and 'SetEndDate'.
func (this *Event) SetDate(t Time) {
	this.SetStartDate(t)
	this.SetEndDate(t)
}

func (this *Event) SetCategory(c string) error {
	oldc := this.category
	this.category = c
	if this.chart != nil {
		this.chart.changeEventCategory(this, oldc)
	}
	return nil
}

//-- method (creation)

func createEventFromParsed(o file.Event) (*Event, warning.Warnings) {
	event := CreateEmptyEvent()
	warns := warning.Warnings{}

    kind, numID := convIDStringToInternal(o.ID)
    event.numID = numID
	event.attrs = convAttrsFileToChart(o.Attrs)

	if kind != MOK_Event || numID == mNID_Invalid {
		event.numID = mNID_Invalid
		event.AddAttr("OldID", o.ID)
	}

	note, moreWarns := CreateNoteFromString(o.Note)
	warns.Concat(moreWarns)
	event.Note = note

	if o.Title != nil {
		localTitle := *o.Title
		event.localTitle = &localTitle
	}

	if o.StartDate != nil {
		startDate := createTimeFromParsed(o.StartDate, TUK_Start)
		event.localStartDate = &startDate
	}

	if o.EndDate != nil {
		endDate := createTimeFromParsed(o.EndDate, TUK_End)
		event.localEndDate = &endDate
	}

	if event.GetStartDate().IsFloating() && !event.GetEndDate().IsFloating() {
		infPast := event.GetStartDate()
		if infPast.Attrs == nil {
			infPast.Attrs = map[string][]string{ "InfinitePast": []string{""} }
		} else if _, ok := infPast.Attrs["InfinitePast"]; !ok {
			infPast.Attrs["InfinitePast"] = []string{""}
		}
		event.SetStartDate(infPast)
	} else if !event.GetStartDate().IsFloating() && event.GetEndDate().IsFloating() {
		infFuture := event.GetEndDate()
		if infFuture.Attrs == nil {
			infFuture.Attrs = map[string][]string{ "InfiniteFuture": []string{""} }
		} else if _, ok := infFuture.Attrs["InfiniteFuture"]; !ok {
			infFuture.Attrs["InfiniteFuture"] = []string{""}
		}
		event.SetEndDate(infFuture)
	}

	event.category = o.Category

	if o.EmbedChart != nil {
		embedChart := *o.EmbedChart
		event.chartEmbedLink = &embedChart
	}

	if o.EmbedEvent != nil {
		embedEvent := *o.EmbedEvent
		event.eventEmbedLink = &embedEvent

		// Eagerly load embedded event. Warn if it failed.
		moreWarns := event.loadEmbeddedEvent()
		warns.Concat(moreWarns)
	}

	// Diagnose and partially auto-correct.
	moreWarns = event.DiagnoseLocal()
	warns.Concat(moreWarns)

	// Create unresolved references.
	for akey, avals := range event.attrs {
		if akey == "ContinuedFrom" {
			for _, aval := range avals {
				event.ContinuedFrom.Reserve(aval)
			}
		}
	}

	return event, warns
}

func CreateEmptyEvent() *Event {
	eobj := &Event{
		MajorCommon: *CreateEmptyMajorCommon(),

		Note: Note{},

		localTitle: nil, 
		localStartDate: nil,
		localEndDate: nil,

		chartEmbedLink: nil,
		embeddedChart: nil,
		eventEmbedLink: nil,
		embeddedEvent: nil,
	} 
    eobj.enclosing = eobj
	eobj.ContinuedFrom.initialize(eobj)
    return eobj
}

//-- method (util)

func (this *Event) loadEmbeddedEvent() (warns warning.Warnings) {
	if this.IsEventEmbedding() && this.embeddedEvent == nil {
		efevent, moreWarns, err := file.LoadFromURI[file.Event](*this.eventEmbedLink)
		if err != nil {
			warns.Add("@0@ cannot load an embedded event.", this)
			this.embeddedEvent = CreateEmptyEvent()
		} else {
			warns.Concat(moreWarns)
			ecevent, moreWarns := createEventFromParsed(*efevent)
			warns.Concat(moreWarns)
			this.embeddedEvent = ecevent
		}
	}
	return
}

func (this *Event) isEmbeddedEventLoaded() bool {
	return this.embeddedEvent != nil
}

func (this *Event) unloadEmbeddedEvent() {
	this.embeddedEvent = nil
}

func (this *Event) loadEmbeddedChart() (warns warning.Warnings) {
	if this.IsChartEmbedding() && this.embeddedChart == nil {
		efile, moreWarns, err := file.LoadFromURI[file.File](*this.chartEmbedLink)
		warns.Concat(moreWarns)
		if err != nil {
			warns.Add("@0@ cannot load an embedded chart.", this)
			embeddedChart, moreWarns := createChartFromParsed(nil)
			warns.Concat(moreWarns)
			this.embeddedChart = embeddedChart
		} else {
			echart, moreWarns := createChartFromParsed(efile)
			warns.Concat(moreWarns)
			this.embeddedChart = echart

			// Update 'parent's of all subchart objects.
			for _, obj := range echart.GetAllMajors(false) {
				obj.setParent(this)
			}

			// Resolve possible unresolved references.
			for _, obj := range echart.GetAllMajors(true) {
				this.chart.resolveReferenceTo(obj)
			}
		}
	}
	return
}

func (this *Event) unloadEmbeddedChart() {
	if this.embeddedChart != nil {
		// unresolve references.
		for _, obj := range this.embeddedChart.GetAllMajors(true) {
			this.chart.unresolveReferenceTo(obj)
		}
	}
	this.embeddedChart = nil
}

func (this *Event) isEmbeddedChartLoaded() bool {
	return this.embeddedChart != nil
}

func (this *Event) shiftStartEndDatesRaw(t Time) {
	oldStartDate := this.GetStartDate()

	if this.timeOffed {
		this.unshiftStartEndDatesRaw()
	} else {
		this.preoffLocalStartDate = this.localStartDate
		this.preoffLocalEndDate = this.localEndDate
		this.timeOffed = true
	}

	nts := this.GetStartDate().Add(t)
	this.localStartDate = nil
	if !this.GetStartDate().Equal(nts) {
		this.localStartDate = &nts
	}
	
	nte := this.GetEndDate().Add(t)
	this.localEndDate = nil
	if !this.GetEndDate().Equal(nte) {
		this.localEndDate = &nte
	}

    // Notify chart if 'StartDate' was changed.
    if !oldStartDate.Equal(this.GetStartDate()) && this.chart != nil {
        this.chart.changeEventStartDate(this)
    }
}

func (this *Event) unshiftStartEndDatesRaw() {
	if this.timeOffed {
		this.localStartDate = this.preoffLocalStartDate
		this.localEndDate = this.preoffLocalEndDate
		this.timeOffed = false
	}
}

// Ensure that both start and end are either floating or non-floating at the same time.
func (this *Event) sanitizeLocalStartEndDates(fixedStart bool) {
	nts, nte, sanitized := sanitizeStartEndDatePair(this.GetStartDate(), this.GetEndDate(), fixedStart)

	if sanitized {
		if fixedStart {
			this.localEndDate = &nte
		} else {
			this.localStartDate = &nts
		}
	}
}

func (this *Event) correctLocalStartEndDates(oldStartDate Time) (ambiguous bool, inverted bool) {
	nts, nte, amb, inv := correctStartEndDatePair(this.GetStartDate(), this.GetEndDate())

	// Reflect corrected results.
	if !this.GetStartDate().Equal(nts) {
		this.localStartDate = &nts
	}

	if !this.GetEndDate().Equal(nte) {
		this.localEndDate = &nte
	}

	ambiguous = amb
	inverted = inv

    // Notify chart if 'StartDate' was changed.
    if !oldStartDate.Equal(this.GetStartDate()) && this.chart != nil {
        this.chart.changeEventStartDate(this)
    }

	return
}
