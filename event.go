package libltc

import (
	"regexp"
	"slices"

	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

type Event struct {
	MainObjCommon

	Note Note 

	localTitle *string
	localStartDate *Time
	localEndDate *Time

	category string

	chartEmbedLink *string
	embeddedChart IChart
	eventEmbedLink *string
	embeddedEvent *Event

	contFromEvents []*Event
	contToEvents []*Event
}

//-- interface IWarningObj

func (this *Event) Summary() (ret string) {
	if this.GetTitle() != "" {
		ret = "an event '" + this.GetTitle() + "'"
	} else {
		ret = "a untitled event"
	}

	extraStr := warning.BuildSummaryString(
		WSE{"ID", this.GetLocalID()},
		WSE{"category", this.GetCategory()},
		WSE{"started", this.GetStartDate()}, 
		WSE{"ended", this.GetEndDate()},
		WSE{"embedding event", this.GetEmbedEventLink()},
		WSE{"embedding chart", this.GetEmbedChartLink()},
	)
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}
	
	return
}

func (this Event) IsUnknown() bool {
	return this.GetTitle() == "" && this.GetLocalID() == "" &&
		this.GetCategory() == "" && this.GetStartDate().IsUnknown() &&
		this.GetEndDate().IsUnknown() && this.GetEmbedEventLink() == "" &&
		this.GetEmbedChartLink() == ""
}

//-- interface IDiagnosable

func (this *Event) DiagnoseLocal() (warns warning.Warnings) {
	if this.kind != MOK_Event || this.numID == NID_Invalid {
		this.kind = MOK_Event
		this.numID = 0
		warns.Add("@0@ has an invalid ID. auto-corrected to 'e0'.", this)
	}
	
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

    if this.GetStartDate().IsAfter(this.GetEndDate()) {
        warns.Add("@0@ has an inverted start/end date pair. correcting...", this)
        this.correctLocalStartEndDates()
    }

	return
}

func (this *Event) DiagnoseNonLocal() (warns warning.Warnings) {
	if this.chart == nil {
		warns.Add("@0@ is not associated to any chart.", this)
		return
	}

	if this.kind != MOK_Event || this.numID == NID_Invalid {
		panic("call DiagnoseLocal() first.")
	}

	for _, category := range this.chart.GetEventCategories() {
		for _, eobj := range this.chart.GetEventsInCategory(category) {
			if eobj.numID == this.numID {
				newNumID, err := this.chart.getNextNumberID(MOK_Event)
				if err != nil {
					panic("Cannot get new number ID.")
				}
				this.numID = newNumID
				warns.Add("@0@ has a duplicated ID. auto-corrected to 'e%d'.", newNumID)
				break
			}
		}
	}

	for akey, avals := range this.attrs {
		if akey == "ContinuedFrom" {
			for _, aval := range avals {
				var eobjFrom *Event
				for _, eobjIn := range this.contFromEvents {
					if eobjIn.GetQualifiedID() == this.convIDRelToAbs(aval) {
						eobjFrom = eobjIn
						break
					}
				}

				if eobjFrom == nil || eobjFrom.IsUnresolved() {
					warns.Add("@0@ has a dangling 'ContinuedFrom' to '%s'.", this, aval)
				} 

				if eobjFrom != nil && !eobjFrom.HasContinuedToEvent(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'ContinuedFrom' to '%s', but it doesn't reference back. corrected.", this, aval)
					eobjFrom.addContinuedToEvent(this)
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
	} else if this.embeddedEvent != nil {
		return this.embeddedEvent.GetTitle()
	} else if this.embeddedChart != nil {
		return this.embeddedChart.GetStringifiedSubjectName()
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

func (this *Event) GetSubchart() (*Chart, warning.Warnings) {
	warns := warning.Warnings{}

	// Lazy-load `embeddedChart` if it's nil (but shouldn't).
	if this.IsChartEmbedding() && this.embeddedChart == nil {
		moreWarns := this.loadEmbeddedChart()
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

func (this *Event) GetContinuedFromEvents() []*Event {
	return this.contFromEvents
}

func (this *Event) GetContinuedToEvents() []*Event {
	return this.contToEvents
}

func (this *Event) HasContinuedFromEvent(eobj *Event) bool {
	return slices.Contains(this.contFromEvents, eobj)
}

func (this *Event) HasContinuedToEvent(eobj *Event) bool {
    return slices.Contains(this.contToEvents, eobj)
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
	this.category = c
	if this.chart != nil {
		this.chart.changeEventCategory(this)
	}
	return nil
}

func (this *Event) addContinuedFromEvent(eobj *Event) {
	// TODO: update contToEvent of qualid (if not unresolved)
	panic("Unimplemented")
}

func (this *Event) resolveContinuedFromEvent(absQualID string, newobj *Event) {
	// TODO: update contToEvent of newobj 
	panic("Unimplemented")
}

func (this *Event) unresolveContinuedFromEvent(absQualID string) {
	// TODO: update contToEvent of absQualID 
	panic("Unimplemented")
}

func (this *Event) removeContinuedFromEvent(eobj *Event) {
	// TODO: update contToEvent of eobj
	panic("Unimplemented")
}

func (this *Event) addContinuedToEvent(eobj *Event) {
	if !this.HasContinuedToEvent(eobj) {
		this.contToEvents = append(this.contToEvents, eobj)
	}
}

func (this *Event) removeContinuedToEvent(eobj *Event) {
	if idx := slices.Index(this.contToEvents, eobj); idx != -1 {
		this.contToEvents = append(this.contToEvents[:idx], this.contToEvents[idx+1:]...)
	}
}

//-- method (high-level operation)

// Public wrapper of 'addContinuedFromEvent'.
func (this *Event) SetContinuedFrom(eobj *Event) {
	this.addContinuedFromEvent(eobj)
}

// Public wrapper of 'removeContinuedFromEvent'.
func (this *Event) UnsetContinuedFrom(eobj *Event) {
	this.removeContinuedFromEvent(eobj)
}

//-- method (creation)

func createEventFromParsed(o file.Event) (*Event, warning.Warnings) {
	event := CreateEmptyEvent()
	warns := warning.Warnings{}

	event.kind, event.numID = convIDStringToInternal(o.ID)
	event.attrs = convAttrsFileToChart(o.Attrs)

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
				ueobj := CreateUnresolvedEvent(aval)
				event.addContinuedFromEvent(ueobj)
			}
		}
	}

	return event, warns
}

func CreateEmptyEvent() *Event {
	return &Event{
		MainObjCommon: *CreateEmptyMainObjCommon(MOK_Event),

		Note: Note{},

		localTitle: nil, 
		localStartDate: nil,
		localEndDate: nil,

		chartEmbedLink: nil,
		embeddedChart: nil,
		eventEmbedLink: nil,
		embeddedEvent: nil,

		contFromEvents: []*Event{},
		contToEvents: []*Event{},
	} 
}

func CreateUnresolvedEvent(relQualID string) *Event {
	ret := CreateEmptyEvent()
	ret.setUnresRelQualID(relQualID)
	return ret
}

//-- method (util)

func (this *Event) loadEmbeddedEvent() (warns warning.Warnings) {
	if this.IsEventEmbedding() && this.embeddedEvent != nil {
		efevent, moreWarns, err := file.LoadFromURI[file.Event](*this.eventEmbedLink)
		warns.Concat(moreWarns)
		if err != nil {
			warns.Add("@0@ cannot load an embedded event.", this)
			this.embeddedEvent = CreateEmptyEvent()
		} else {
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
	if this.IsChartEmbedding() && this.embeddedChart != nil {
		efile, moreWarns, err := file.LoadFromURI[file.File](*this.chartEmbedLink)
		warns.Concat(moreWarns)
		if err != nil {
			warns.Add("@0@ cannot load an embedded chart.", this)
			this.embeddedChart = CreateEmptyChart()
		} else {
			echart, moreWarns := createChartFromParsed(efile)
			warns.Concat(moreWarns)
			this.embeddedChart = echart

			// Update 'parent's of all subchart objects.
			for _, obj := range this.chart.GetAllMainObjects() {
				obj.SetParent(this)
			}

			// Resolve possible unresolved references.
			for _, obj := range this.embeddedChart.GetAllMainObjects() {
				this.chart.resolveReferenceTo(obj)
			}
		}
	}
	return
}

func (this *Event) unloadEmbeddedChart() {
	if this.embeddedChart != nil {
		// unresolve references.
		for _, obj := range this.embeddedChart.GetAllMainObjects() {
			this.chart.unresolveReferenceTo(obj)
		}
	}
	this.embeddedChart = nil
}

func (this *Event) isEmbeddedChartLoaded() bool {
	return this.embeddedChart != nil
}

func (this *Event) correctLocalStartEndDates() {
    oldStartDate := this.GetStartDate()

	// Correct ambiguity.
	// Assume 12-month years. (Sorry non-Gregorian calendars)
	if this.GetStartDate().IsAmbiguous() && !this.GetEndDate().IsAmbiguous() {
		newStartDate := this.GetEndDate()
		newStartDate.SetUsage(TUK_Start)
		newStartDate.SetMonth(newStartDate.GetMonth() - 1)
		if newStartDate.GetMonth() <= 0 {
			newStartDate.SetYear(newStartDate.GetYear() - 1)
			newStartDate.SetMonth(12)
		}
		this.SetLocalStartDate(&newStartDate)
	} else if this.GetEndDate().IsAmbiguous() && !this.GetStartDate().IsAmbiguous() {
		newEndDate := this.GetStartDate()
		newEndDate.SetUsage(TUK_End)
		newEndDate.SetMonth(newEndDate.GetMonth() + 1)
		if newEndDate.GetMonth() >= 12 {
			newEndDate.SetYear(newEndDate.GetYear() + 1)
			newEndDate.SetMonth(1)
		}
		this.SetLocalEndDate(&newEndDate)
	}

	// Correct order.
	if this.GetStartDate().IsAfter(this.GetEndDate()) {
		hlsd := this.HasLocalStartDate()
		hled := this.HasLocalEndDate()
		
		if hlsd && hled {
			orgStartDate := this.GetLocalStartDate()
			this.SetLocalStartDate(this.GetLocalEndDate())
			this.SetLocalEndDate(orgStartDate)
			this.GetLocalStartDate().SetUsage(TUK_Start)
			this.GetLocalEndDate().SetUsage(TUK_End)
		} else if hlsd && !hled {
			this.SetLocalEndDate(this.GetLocalStartDate())
			this.SetLocalStartDate(nil)
			this.GetLocalEndDate().SetUsage(TUK_End)
		} else if !hlsd && hled {
			this.SetLocalStartDate(this.GetLocalEndDate())
			this.SetLocalEndDate(nil)
			this.GetLocalStartDate().SetUsage(TUK_Start)
		} else { //if !hlsd && !hled 
            // This case shouldn't happen (unless the embedded event is broken),
            // but if it does, create local dates to "overshadow" the problem.
            newStartDate := this.GetEndDate()
            newEndDate := this.GetStartDate()
            newStartDate.SetUsage(TUK_Start)
            newEndDate.SetUsage(TUK_End)
            this.SetLocalStartDate(&newStartDate)
            this.SetLocalEndDate(&newEndDate)
		}
	}

    // Notify chart if 'StartDate' was changed.
    if !oldStartDate.Equal(this.GetStartDate()) && this.chart != nil {
        this.chart.changeEventStartDate(this)
    }
}
