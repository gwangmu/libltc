package libltc

import (
	"errors"

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
	// TODO
	panic("Unimplemented")
}

func (this *Event) IsUnknown() bool {
	// TODO
	panic("Unimplemented")
}

//-- interface IDiagnosable

func (this *Event) DiagnoseLocal() (warns warning.Warnings) {
	if this.kind != MOK_Event || this.numID == NID_Invalid {
		warns.Add("@0@ has an invalid ID.", this)
	}

	if this.eventEmbedLink != nil && this.embeddedEvent == nil {
		warns.Add("Cannot load the embedded event of @0@.", this)
	}

	return
}

func (this *Event) DiagnoseNonLocal() (warns warning.Warnings) {
	if this.GetChart() == nil {
		warns.Add("@0@ is not associated to any chart.", this)
		return
	}

	// TODO
	panic("Unimplemented")
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

func (this *Event) GetEmbedLink() string {
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

func (this *Event) GetSubchartLink() string {
	// TODO
	panic("Unimplemented")
}

func (this *Event) GetStartDate() *Time {
	// TODO: use 'localStartDate` if it was defined.
	// TODO: otherwise, use the `StartDate` of the embedded event.
	// TODO: otherwise, return unknown.
	panic("Unimplemented")
}

func (this *Event) GetEndDate() *Time {
	// TODO: use 'localEndDate` if it was defined.
	// TODO: otherwise, use the `EndDate` of the embedded event.
	// TODO: otherwise, return unknown.
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

func (this *Event) HasContinuedFromEvent() bool {
	// TODO
	panic("Unimplemented")
}

func (this *Event) HasContinuedToEvent() bool {
	// TODO
	panic("Unimplemented")
}

//-- method (setters)

func (this *Event) SetEmbedLink(uri string) {
	// TODO: invalidate non-nill embed event.
	// TODO: eagerly load 'embeddedEvent'
	panic("Unimplemented")
}

func (this *Event) SetSubchartLink(uri string) {
	// TODO: invalidate non-nill subchart.
	panic("Unimplemented")
}

func (this *Event) SetLocalStartDate(t *Time) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) SetLocalEndDate(t *Time) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) SetCategory(c string) error {
	// WARNING: directly calling this won't guarantee any consistency in
	// the associated chart, as in the chart may say 'the category is A'
	// but the event says 'no, the category is B'. Use the interface
	// provided by the associated chart.
	if this.chart != nil {
		return errors.New("Already associated to a chart.")
	} else {
		this.category = c
		return nil
	}
}

func (this *Event) addContinuedFromEvent(eobj *Event) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) swapContinuedFromEvent(oldobj *Event, newobj *Event) {
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
		// TODO: eagerly load embedded event. warn if it failed.
		panic("Unimplemented")
	}

	moreWarns = event.DiagnoseLocal()
	warns.Concat(moreWarns)

	return event, warns
}

func CreateEmptyEvent() *Event {
	return &Event{
		MainObjCommon: CreateEmptyMainObjCommon(MOK_Event),

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
