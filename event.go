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

	if this.eventEmbedLink != nil && this.embeddedEvent == nil {
		warns.Add("Cannot load the embedded event of @0@.", this)
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
	// TODO: use `localTitle` if it was defined.
	// TODO: otherwise, use the title of `embeddedEvent`.
	// TODO: otherwise, use the stringified name of `embeddedChart`.
	panic("Unimplemented")
}

func (this *Event) GetEmbeddedEvent() *Event {
	// TODO: return a embedded event (nil if `Embed` is invalid)
	// TODO: eagerly load `embeddedEvent` upon the file load because it may
	//       determine `{Start,End}Date`.
	panic("Unimplemented")
}

func (this *Event) GetEmbedEventLink() string {
	// TODO
	panic("Unimplemented")
}

func (this *Event) GetSubchart() *Chart {
	// TODO: return a subchart (nil if `Embed` is valid or `Subchart` is invalid).
	// TODO: lazy-load `embeddedChart` if it's nil.
	// TODO: on lazy-load, update `parent`s of the objects inside.
	// TODO: on lazy-load, invoke `resolveReferenceTo` for all subchart objs.
	panic("Unimplemented")
}

func (this *Event) GetEmbedChartLink() string {
	// TODO
	panic("Unimplemented")
}

func (this *Event) GetStartDate() Time {
	// TODO: use 'localStartDate` if it was defined.
	// TODO: otherwise, use the `StartDate` of the embedded event.
	// TODO: otherwise, return unknown.
	panic("Unimplemented")
}

func (this *Event) GetEndDate() Time {
	// TODO: use 'localEndDate` if it was defined.
	// TODO: otherwise, use the `EndDate` of the embedded event.
	// TODO: otherwise, return unknown.
	panic("Unimplemented")
}

func (this *Event) GetLocalStartDate() *Time {
	// TODO: return 'localStartDate` if it was defined.
	// TODO: otherwise, return nil.
	panic("Unimplemented")
}

func (this *Event) GetLocalEndDate() *Time {
	// TODO: return 'localEndDate` if it was defined.
	// TODO: otherwise, return nil.
	panic("Unimplemented")
}

func (this *Event) GetCategory() string{
	return this.category
}

func (this *Event) GetContinuedFromEvents() []*Event {
	// TODO
	panic("Unimplemented")
}

func (this *Event) GetContinuedToEvents() []*Event {
	// TODO
	panic("Unimplemented")
}

func (this *Event) HasContinuedFromEvent(eobj *Event) bool {
	// TODO
	panic("Unimplemented")
}

func (this *Event) HasContinuedToEvent(eobj *Event) bool {
	// TODO
	panic("Unimplemented")
}

//-- method (setters)

func (this *Event) SetEmbedEventLink(uri string) {
	// TODO: invalidate non-nill embed event.
	// TODO: eagerly load 'embeddedEvent'
	panic("Unimplemented")
}

func (this *Event) SetEmbedChartLink(uri string) {
	// TODO: invalidate non-nill subchart.
	panic("Unimplemented")
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
	// TODO: call 'GetChart().changeEventCategory()'.
	return nil
}

func (this *Event) addContinuedFromEvent(eobj *Event) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) resolveContinuedFromEvent(qualid string, newobj *Event) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) unresolveContinuedFromEvent(qualid string) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) removeContinuedFromEvent(eobj *Event) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) addContinuedToEvent(eobj *Event) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) removeContinuedToEvent(eobj *Event) {
	// TODO
	panic("Unimplemented")
}

//-- method (high-level operation)

func (this *Event) SetContinuedFrom(eobj *Event) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) UnsetContinuedFrom(eobj *Event) {
	// TODO
	panic("Unimplemented")
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
		efevent, moreWarns, err := file.LoadFromURI[file.Event](*event.eventEmbedLink)
		warns.Concat(moreWarns)
		if err != nil {
			warns.Add("@0@ cannot load an embedded event.", event)
			event.embeddedEvent = CreateEmptyEvent()
		} else {
			ecevent, moreWarns := createEventFromParsed(*efevent)
			warns.Concat(moreWarns)
			event.embeddedEvent = ecevent
		}
	}

	moreWarns = event.DiagnoseLocal()
	warns.Concat(moreWarns)

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
