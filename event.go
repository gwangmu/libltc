package libltc

import (
	"errors"
)

type LTCEvent struct {
	Common LTCMainObjCommon
	Note LTCNote 

	localTitle *string
	localStartDate *LTCTime
	localEndDate *LTCTime

	subchartLink *string
	loadedSubchart *LTCChart
	embedLink *string
	loadedEmbed *LTCEvent
}

//-- interface LTCWarningObj

func (this *LTCEvent) Summary() string {
	// TODO
}

func (this *LTCEvent) IsUnknown() bool {
	// TODO
}

//-- interface LTCTomlPrintable

func (this *LTCEvent) PrintTOML() string {
	// TODO
}

//-- method (getters and setters)

func (this *LTCEvent) GetTitle() string {
	// TODO: use `localTitle` if it was defined.
	// TODO: otherwise, use the title of `loadedEmbed`.
	// TODO: otherwise, use the stringified name of `loadedSubchart`.
}

func (this *LTCEvent) GetEmbeddedEvent() *LTCEvent {
	// TODO: return a embedded event (nil if `Embed` is invalid)
	// TODO: eagerly load `loadedEmbed` upon the file load because it may
	//       determine `{Start,End}Date`.
}

func (this *LTCEvent) GetEmbedLink() string {
	// TODO
}

func (this *LTCEvent) SetEmbedLink() {
	// TODO: invalidate non-nill embed event.
	// TODO: eagerly load 'loadedEmbed'
}

func (this *LTCEvent) GetSubchart() *LTCChart {
	// TODO: return a subchart (nil if `Embed` is valid or `Subchart` is invalid).
	// TODO: lazy-load `loadedSubchart` if it's nil.
	// TODO: on lazy-load, update `parent`s of the objects inside.
}

func (this *LTCEvent) GetSubchartLink() string {
	// TODO
}

func (this *LTCEvent) SetSubchartLink() {
	// TODO: invalidate non-nill subchart.
}

func (this *LTCEvent) GetStartDate() LTCTime {
	// TODO: use 'localStartDate` if it was defined.
	// TODO: otherwise, use the `StartDate` of the embedded event.
	// TODO: otherwise, return unknown.
}

func (this *LTCEvent) SetLocalStartDate(t LTCTime) {
	// TODO
}

func (this *LTCEvent) UnsetLocalStartDate() {
	// TODO
}

func (this *LTCEvent) GetEndDate() LTCTime {
	// TODO: use 'localEndDate` if it was defined.
	// TODO: otherwise, use the `EndDate` of the embedded event.
	// TODO: otherwise, return unknown.
}

func (this *LTCEvent) SetLocalEndDate(t LTCTime) {
	// TODO
}

func (this *LTCEvent) UnsetLocalEndDate() {
	// TODO
}

//-- method (creation)

func CreateEmptyEvent() *LTCEvent {
	return &LTCEvent{
		Common: LTCMainObjCommon{
			kind: MOK_Event,

			chart: nil,
			parent: nil,
			numID: NID_Invalid,

			extraNoteAnnexs: []*LTCAnnex{},
			attachedAnnexs: []*LTCAnnex{},
			attrs: map[string][]string{},
		}

		Note: LTCNote{},

		localTitle: "",
		localStartDate: nil,
		localEndDate: nil,

		subchartLink: nil,
		loadedSubchart: nil,
		embedLink: nil,
		loadedEmbed: nil,
	}
}
