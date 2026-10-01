package libltc

import (
	"regexp"

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

func (this *Event) Summary() string {
	// TODO: unimplemented
	return "an event"
}

func (this *Event) IsUnknown() bool {
	// TODO: unimplemented
	return false
}

//-- interface IDiagnosable

func (this *Event) DiagnoseLocal() (warns warning.Warnings) {
	if this.kind != MOK_Event || this.numID == NID_Invalid {
		this.kind = MOK_Event
		this.numID = 0
		warns.Add("@0@ has an invalid ID. auto-corrected to 'e0'.", this)
	}
	
	if this.embeddedChart != nil && this.embeddedEvent != nil {
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
					if eobjIn.GetQualifiedIDFrom(this.chart.GetEmbeddingEvent()) == aval {
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
	if this.eventEmbedLink != nil {
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
	if this.chartEmbedLink != nil && this.embeddedChart == nil {
		moreWarns := this.loadEmbeddedChart()
		warns.Concat(moreWarns)
	}

	return this.embeddedChart.asChart(), warns
}

func (this *Event) GetEmbedChartLink() string {
	if this.chartEmbedLink != nil {
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
		return Time{}
	}
}

func (this *Event) GetEndDate() Time {
	if this.localEndDate != nil {
		return *this.localEndDate
	} else if this.embeddedEvent != nil {
		return this.embeddedEvent.GetEndDate()
	} else {
		return Time{}
	}
}

func (this *Event) GetLocalStartDate() *Time {
	return this.localStartDate
}

func (this *Event) GetLocalEndDate() *Time {
	return this.localEndDate
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
	for _, oobj := range this.contFromEvents {
		if oobj == eobj {
			return true
		}
	}
	return false
}

func (this *Event) HasContinuedToEvent(eobj *Event) bool {
	for _, oobj := range this.contToEvents {
		if oobj == eobj {
			return true
		}
	}
	return false
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
	// TODO: call 'GetChart().changeEventStartDate()' if successful.
	// TODO: "successful": non-embedding event, unconditional.
	// 			embedding event, (before) local unused, (after) local used.
	//			embedding event, (before) local used, (after) local unused.
	//			embedding event, (before) local used, (after) local used.
	panic("Unimplemented")
}

func (this *Event) SetLocalEndDate(t *Time) {
	// TODO: swap if StartDate is after `t`.
	// TODO: if so, call `chart.changeEventStartDate()` according to above std.
	panic("Unimplemented")
}

func (this *Event) SetCategory(c string) error {
	this.category = c
	if this.chart != nil {
		this.chart.changeEventCategory(this)
	}
	return nil
}

func (this *Event) addContinuedFromEvent(eobj *Event) {
	// TODO: update contToEvent of eobj (if not unresolved)
	panic("Unimplemented")
}

func (this *Event) resolveContinuedFromEvent(qualid string, newobj *Event) {
	// TODO: update contToEvent of newobj 
	panic("Unimplemented")
}

func (this *Event) unresolveContinuedFromEvent(qualid string) {
	// TODO: update contToEvent of qualid 
	panic("Unimplemented")
}

func (this *Event) removeContinuedFromEvent(eobj *Event) {
	// TODO: update contToEvent of eobj
	panic("Unimplemented")
}

func (this *Event) addContinuedToEvent(eobj *Event) {
	found := false
	for _, elem := range this.contToEvents {
		if elem == eobj {
			found = true
			break
		}
	}

	if !found {
		this.contToEvents = append(this.contToEvents, eobj)
	}
}

func (this *Event) removeContinuedToEvent(eobj *Event) {
	idx := -1
	for i, elem := range this.contToEvents {
		if elem == eobj {
			idx = i
			break
		}
	}

	if idx != -1 {
		this.contToEvents = append(this.contToEvents[:idx], this.contToEvents[idx+1:]...)
	}
}

//-- method (high-level operation)

func (this *Event) SetContinuedFrom(eobj *Event) {
	this.addContinuedFromEvent(eobj)
}

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
		startDate := createTimeFromParsed(o.StartDate)
		event.localStartDate = &startDate
	}

	if o.EndDate != nil {
		endDate := createTimeFromParsed(o.EndDate)
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

func CreateUnresolvedEvent(fullid string) *Event {
	ret := CreateEmptyEvent()
	ret.fullID = fullid
	return ret
}

//-- method (util)

func (this *Event) loadEmbeddedEvent() (warns warning.Warnings) {
	if this.eventEmbedLink != nil {
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

func (this *Event) loadEmbeddedChart() (warns warning.Warnings) {
	if this.chartEmbedLink != nil {
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
}
