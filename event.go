package libltc

import (
	"errors"
	"libltc/file"
)

type Event struct {
	Note Note 

	common MainObjCommon

	localTitle *string
	localStartDate *Time
	localEndDate *Time

	subchartLink *string
	loadedSubchart *Chart
	embedLink *string
	loadedEmbed *Event

	contFromEvents []*Event
	contToEvents []*Event
}

//-- interface MainObj

func (this *Event) Common() *MainObjCommon {
	return &this.common
}

func (this *Event) IsUnresolved() bool {
	return this.common.IsUnresolved()
}

//-- interface WarningObj

func (this *Event) Summary() string {
	// TODO
}

func (this *Event) IsUnknown() bool {
	// TODO
}

//-- method (getters and setters)

func (this *Event) GetTitle() string {
	// TODO: use `localTitle` if it was defined.
	// TODO: otherwise, use the title of `loadedEmbed`.
	// TODO: otherwise, use the stringified name of `loadedSubchart`.
}

func (this *Event) GetEmbeddedEvent() *Event {
	// TODO: return a embedded event (nil if `Embed` is invalid)
	// TODO: eagerly load `loadedEmbed` upon the file load because it may
	//       determine `{Start,End}Date`.
}

func (this *Event) GetEmbedLink() string {
	// TODO
}

func (this *Event) SetEmbedLink() {
	// TODO: invalidate non-nill embed event.
	// TODO: eagerly load 'loadedEmbed'
}

func (this *Event) GetSubchart() *Chart {
	// TODO: return a subchart (nil if `Embed` is valid or `Subchart` is invalid).
	// TODO: lazy-load `loadedSubchart` if it's nil.
	// TODO: on lazy-load, update `parent`s of the objects inside.
}

func (this *Event) GetSubchartLink() string {
	// TODO
}

func (this *Event) SetSubchartLink() {
	// TODO: invalidate non-nill subchart.
}

func (this *Event) GetStartDate() Time {
	// TODO: use 'localStartDate` if it was defined.
	// TODO: otherwise, use the `StartDate` of the embedded event.
	// TODO: otherwise, return unknown.
}

func (this *Event) GetContinuedFromEvents() []*Event {
	// TODO
}

func (this *Event) GetContinuedToEvents() []*Event {
	// TODO
}

func (this *Event) SetLocalStartDate(t Time) {
	// TODO
}

func (this *Event) UnsetLocalStartDate() {
	// TODO
}

func (this *Event) GetEndDate() Time {
	// TODO: use 'localEndDate` if it was defined.
	// TODO: otherwise, use the `EndDate` of the embedded event.
	// TODO: otherwise, return unknown.
}

func (this *Event) SetLocalEndDate(t Time) {
	// TODO
}

func (this *Event) UnsetLocalEndDate() {
	// TODO
}

func (this *Event) HasContinuedFromEvent() bool {
	// TODO
}

func (this *Event) AddContinuedFromEvent(eobj *Event) {
	// TODO
}

func (this *Event) SwapContinuedFromEvent(oldobj *Event, newobj *Event) {
	// TODO
}

func (this *Event) RemoveContinuedFromEvent(eobj *Event) {
	// TODO
}

func (this *Event) HasContinuedToEvent() bool {
	// TODO
}

func (this *Event) AddContinuedToEvent(eobj *Event) {
	// TODO
}

func (this *Event) RemoveContinuedToEvent(eobj *Event) {
	// TODO
}

//-- method (relational)

func (this *Event) MakeContinuedFrom(eobj *Event) {
	// TODO
}

func (this *Event) UnmakeContinuedFrom(eobj *Event) {
	// TODO
}

//-- method (import and export)

func (this *Event) Import(feobj *file.Event) error {
	// TODO
}

func (this *Event) Export() (*file.Event, error) {
	// TODO
}

//-- method (creation and disposal)

func CreateEmptyEvent() *Event {
	return &Event{
		Common: MainObjCommon{
			kind: MOK_Event,

			chart: nil,
			parent: nil,
			numID: NID_Invalid,
			fullID: "", 

			extraNoteAnnexs: []*Annex{},
			attachedAnnexs: []*Annex{},
			attrs: map[string][]string{},
		}

		Note: Note{},

		localTitle: "",
		localStartDate: nil,
		localEndDate: nil,

		subchartLink: nil,
		loadedSubchart: nil,
		embedLink: nil,
		loadedEmbed: nil,

		contFromEvents: []*Event{},
		contToEvents: []*Event{},
	} 
}
