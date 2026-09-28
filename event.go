package libltc

import (
	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

type Event struct {
	common MainObjCommon

	Note Note 

	localTitle *string
	localStartDate *Time
	localEndDate *Time

	subchartEmbedLink *string
	loadedSubchart IChart
	eventEmbedLink *string
	loadedEvent *Event

	contFromEvents []*Event
	contToEvents []*Event
}

//-- interface IMainObj

func (this *Event) Common() *MainObjCommon {
	return &this.common
}

func (this *Event) IsUnresolved() bool {
	return this.common.IsUnresolved()
}

func (this *Event) Clone(preserveChart bool, preserveID bool) IMainObj {
	newobj := *this
	if !preserveChart {
		newobj.common.chart = nil
	}
	if !preserveID {
		newobj.common.numID = NID_Invalid
		newobj.common.fullID = ""
	}
	return &newobj
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

//-- method (getters)

func (this *Event) GetTitle() string {
	// TODO: use `localTitle` if it was defined.
	// TODO: otherwise, use the title of `loadedEvent`.
	// TODO: otherwise, use the stringified name of `loadedSubchart`.
	panic("Unimplemented")
}

func (this *Event) GetEmbeddedEvent() *Event {
	// TODO: return a embedded event (nil if `Embed` is invalid)
	// TODO: eagerly load `loadedEvent` upon the file load because it may
	//       determine `{Start,End}Date`.
	panic("Unimplemented")
}

func (this *Event) GetEmbedLink() string {
	// TODO
	panic("Unimplemented")
}

func (this *Event) GetSubchart() *Chart {
	// TODO: return a subchart (nil if `Embed` is valid or `Subchart` is invalid).
	// TODO: lazy-load `loadedSubchart` if it's nil.
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
	// TODO: eagerly load 'loadedEvent'
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

//-- method (import and export)

func (this *Event) Import(feobj *file.Event) (warning.Warnings, error) {
	// TODO
	panic("Unimplemented")
}

func (this *Event) Export() (*file.Event, warning.Warnings, error) {
	// TODO
	panic("Unimplemented")
}

//-- method (creation)

func CreateEmptyEvent() *Event {
	return &Event{
		common: MainObjCommon{
			kind: MOK_Event,

			chart: nil,
			parent: nil,
			numID: NID_Invalid,
			fullID: "", 

			extraNoteAnnexs: []*Annex{},
			attachedAnnexs: []*Annex{},
			attrs: map[string][]string{},
		},

		Note: Note{},

		localTitle: nil, 
		localStartDate: nil,
		localEndDate: nil,

		subchartEmbedLink: nil,
		loadedSubchart: nil,
		eventEmbedLink: nil,
		loadedEvent: nil,

		contFromEvents: []*Event{},
		contToEvents: []*Event{},
	} 
}
