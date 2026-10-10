package ltc

import (
	"regexp"

	"libltc/warning"
	"libltc/internal/file"
)

type Event struct {
	MajorCommon

	Note Note 

	localTitle *string
	localStartDate *Time
	localEndDate *Time

	category string

	chartEmbedLink *string
	embeddedChart IChart
	eventEmbedLink *string
	embeddedEvent *Event

    ContinuedFrom	// Specifying 'ContinuedFrom'
    ContinuedTo		// Pointed by 'ContinuedFrom'
}

type ChartLoader = func (*file.File) (*Chart, warning.Warnings)

func (this *Event) Equal(that IMajor) bool {
	return this.GetCommon() == that.GetCommon()
}

//-- interface IMajor

func (this *Event) resolveReferenceTo(c IChart) {
	for _, category := range c.GetEventCategories() {
		for _, eobj := range c.GetEventsInCategory(category) {
			this.ContinuedFrom.Link(eobj)			
		}
	}
}

func (this *Event) unresolveReferenceTo(c IChart) {
	for _, category := range c.GetEventCategories() {
		for _, eobj := range c.GetEventsInCategory(category) {
			this.ContinuedFrom.Unlink(eobj)			
		}
	}
}

func (this *Event) initialize() {
	*this = Event{
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
    this.enclosing = this
	this.ContinuedFrom.initialize(this)
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
			WSE{"ID", this.getUnresRelQualID()},
			WSE{"OldID", maybeOldID},
		)
	} else {
		extraStr = warning.BuildSummaryString(
			WSE{"ID", this.GetIntrinsicID()},
			WSE{"OldID", maybeOldID},
			WSE{"Category", this.GetCategory()},
			WSE{"Started", this.GetStartDate()}, 
			WSE{"Ended", this.GetEndDate()},
			WSE{"Embedded-event", this.GetEmbedEventLink()},
			WSE{"Embedded-chart", this.GetEmbedChartLink()},
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
       
	amb, inv := this.correctLocalStartEndDates()
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
	if this.numID == NID_Invalid {
		newNumID := this.chart.getNextNumberID(MOK_Event)
		this.numID = newNumID
		warns.Add("@0@ had an invalid ID. auto-corrected to 'e%d'.", this, newNumID)
	}

	// Check duplicate ID.
	for _, category := range this.chart.GetEventCategories() {
		for _, eobj := range this.chart.GetEventsInCategory(category) {
			if eobj != this && eobj.numID == this.numID {
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

				if eobjFrom != nil && !eobjFrom.ContinuedTo.Has(this) {
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
		if sc, _ := this.GetSubchart(nil); sc != nil {
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

func (this *Event) GetSubchart(cloader ChartLoader) (*Chart, warning.Warnings) {
	warns := warning.Warnings{}

	// Lazy-load `embeddedChart` if it's nil (but shouldn't).
	if this.IsChartEmbedding() && this.embeddedChart == nil {
		moreWarns := this.loadEmbeddedChart(cloader)
		warns.Concat(moreWarns)
	}

	return this.embeddedChart.asChart(), warns
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

func (this *Event) GetStartDate() Time {
	if this.localStartDate != nil {
		return *this.localStartDate
	} else if this.embeddedEvent != nil {
		return this.embeddedEvent.GetStartDate()
	} else {
		return CreateEmptyTime(TUK_Start) 
	}
}

func (this *Event) GetEndDate() Time {
	if this.localEndDate != nil {
		return *this.localEndDate
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

func (this *Event) SetLocalStartDate(t *Time) {
    this.localStartDate = t
    this.correctLocalStartEndDates()
}

func (this *Event) SetLocalEndDate(t *Time) {
    this.localEndDate = t
    this.correctLocalStartEndDates()
}

// Convenience wrapper to 'SetLocalStartDate'.
func (this *Event) SetStartDate(t Time) {
	this.SetLocalStartDate(&t)
}

// Convenience wrapper to 'SetLocalEndDate'.
func (this *Event) SetEndDate(t Time) {
	this.SetLocalEndDate(&t)
}

func (this *Event) SetCategory(c string) error {
	oldc := this.category
	this.category = c
	if this.chart != nil {
		this.chart.changeEventCategory(this, oldc)
	}
	return nil
}

//-- method (high-level operation)

// Public wrapper of 'addContinuedFromEvent'.
func (this *Event) SetContinuedFrom(eobj *Event) {
    this.ContinuedFrom.Link(eobj)
}

// Public wrapper of 'removeContinuedFromEvent'.
func (this *Event) UnsetContinuedFrom(eobj *Event) {
	this.ContinuedFrom.Unreserve(eobj.GetAbsQualifiedID())
}

// Public wrapper of 'removeContinuedFromEvent'.
func (this *Event) UnsetContinuedFromByID(absQualID string) {
	this.ContinuedFrom.Unreserve(absQualID)
}

//-- method (creation)

func createEventFromParsed(o file.Event) (*Event, warning.Warnings) {
	event := CreateEmptyEvent()
	warns := warning.Warnings{}

    kind, numID := convIDStringToInternal(o.ID)
    event.numID = numID
	event.attrs = convAttrsFileToChart(o.Attrs)

	if kind != MOK_Event || numID == NID_Invalid {
		event.numID = NID_Invalid
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
    eobj := Event{}
    eobj.initialize()
    return &eobj
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

func (this *Event) loadEmbeddedChart(cloader ChartLoader) (warns warning.Warnings) {
	if cloader == nil {
		warns.Add("!!!INTERNAL WARN!!! chart loader not provided")
		return
	}

	if this.IsChartEmbedding() && this.embeddedChart != nil {
		efile, moreWarns, err := file.LoadFromURI[file.File](*this.chartEmbedLink)
		warns.Concat(moreWarns)
		if err != nil {
			warns.Add("@0@ cannot load an embedded chart.", this)
			embeddedChart, moreWarns := cloader(nil)
			warns.Concat(moreWarns)
			this.embeddedChart = embeddedChart
		} else {
			echart, moreWarns := cloader(efile)
			warns.Concat(moreWarns)
			this.embeddedChart = echart

			// Update 'parent's of all subchart objects.
			for _, obj := range echart.GetAllMajors(false) {
				obj.SetParent(this)
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

func (this *Event) correctLocalStartEndDates() (ambiguous bool, inverted bool) {
    oldStartDate := this.GetStartDate()

	// Correct ambiguity.
	if !this.GetStartDate().IsInfinite() && this.GetStartDate().IsAmbiguous() && !this.GetEndDate().IsAmbiguous() {
		newStartDate := this.GetStartDate()
		if this.GetStartDate().year == nil && this.GetStartDate().month == nil {
			// Can still be inverted. To be treated below.
			newStartDate.SetYear(this.GetEndDate().GetYear())
			newStartDate.SetMonth(this.GetEndDate().GetMonth())
		} else if this.GetStartDate().year == nil {
			if this.GetEndDate().GetMonth() >= this.GetStartDate().GetMonth() {
				newStartDate.SetYear(this.GetEndDate().GetYear())
			} else {
				newStartDate.SetYear(this.GetEndDate().GetYear() - 1)
			}
		} else /* if this.GetStartDate().month == nil */ {
			newStartDate.SetMonth(this.GetEndDate().GetMonth())
		}

		this.localStartDate = &newStartDate
		ambiguous = true
	} else if !this.GetEndDate().IsInfinite() && this.GetEndDate().IsAmbiguous() && !this.GetStartDate().IsAmbiguous() {
		newEndDate := this.GetEndDate()
		if this.GetEndDate().year == nil && this.GetEndDate().month == nil {
			// Can still be inverted. To be treated below.
			newEndDate.SetYear(this.GetStartDate().GetYear())
			newEndDate.SetMonth(this.GetStartDate().GetMonth())
		} else if this.GetEndDate().year == nil {
			if this.GetEndDate().GetMonth() >= this.GetStartDate().GetMonth() {
				newEndDate.SetYear(this.GetStartDate().GetYear())
			} else {
				newEndDate.SetYear(this.GetStartDate().GetYear() + 1)
			}
		} else /* if this.GetEndDate().month == nil */ {
			newEndDate.SetMonth(this.GetStartDate().GetMonth())
		}
		this.localEndDate = &newEndDate
		ambiguous = true
	}

	// Correct order.
	if this.GetStartDate().IsAfter(this.GetEndDate()) {
		hlsd := this.HasLocalStartDate()
		hled := this.HasLocalEndDate()
		
		if hlsd && hled {
			orgStartDate := this.GetLocalStartDate()
			this.localStartDate = this.GetLocalEndDate()
			this.localEndDate = orgStartDate
			this.GetLocalStartDate().SetUsage(TUK_Start)
			this.GetLocalEndDate().SetUsage(TUK_End)
		} else if hlsd && !hled {
			this.localEndDate = this.GetLocalStartDate()
			this.localStartDate = nil
			this.GetLocalEndDate().SetUsage(TUK_End)
		} else if !hlsd && hled {
			this.localStartDate = this.GetLocalEndDate()
			this.localEndDate = nil
			this.GetLocalStartDate().SetUsage(TUK_Start)
		} else { //if !hlsd && !hled 
            // This case shouldn't happen (unless the embedded event is broken),
            // but if it does, create local dates to "overshadow" the problem.
            newStartDate := this.GetEndDate()
            newEndDate := this.GetStartDate()
            newStartDate.SetUsage(TUK_Start)
            newEndDate.SetUsage(TUK_End)
            this.localStartDate = &newStartDate
            this.localEndDate = &newEndDate
		}

		inverted = true
	}

    // Notify chart if 'StartDate' was changed.
    if !oldStartDate.Equal(this.GetStartDate()) && this.chart != nil {
        this.chart.changeEventStartDate(this)
    }

	return
}
