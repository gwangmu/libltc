package libltc

import (
	"errors"
)

type LTCEvent struct {
	Note LTCNote 

	common LTCMainObjCommon

	localTitle *string
	localStartDate *LTCTime
	localEndDate *LTCTime

	subchartLink *string
	loadedSubchart *LTCChart
	embedLink *string
	loadedEmbed *LTCEvent
}

//-- interface LTCMainObj

func (this *LTCEvent) GetChart() *LTCChart {
	return this.common.GetChart()
}

func (this *LTCEvent) GetLocalID() string {
	return this.common.GetLocalID("e")
}

func (this *LTCEvent) GetQualifiedID() string {
	return this.common.GetQualifiedID("e")
}

func (this *LTCEvent) GetExtraNotes() []*LTCNote {
	return this.common.GetExtraNotes()
}

func (this *LTCEvent) AddExtraNoteAnnex(aobj *LTCAnnex) {
	if !aobj.HasAttr("ExtraNoteOf", this.GetQualifiedID()) {
		aobj.AddAttr("ExtraNoteOf", this.GetQualifiedID())
	}
	this.common.AddExtraNoteAnnex(aobj)
}

func (this *LTCEvent) RemoveExtraNoteAnnex(aobj *LTCAnnex) {
	this.common.RemoveExtraNoteAnnex(aobj)
	aobj.RemoveAttr("ExtraNoteOf", this.GetQualifiedID())
}

func (this *LTCEvent) GetAttachedAnnexs() []*LTCAnnex {
	return this.common.GetAttachedAnnexs()
}

func (this *LTCEvent) AddAttachedAnnex(aobj *LTCAnnex) {
	if !aobj.HasAttr("AttachTo", this.GetQualifiedID()) {
		aobj.AddAttr("AttachTo", this.GetQualifiedID())
	}
	this.common.AddAttachedAnnex(aobj)
}

func (this *LTCEvent) RemoveAttachedAnnex(aobj *LTCAnnex) {
	this.common.RemoveAttachedAnnex(aobj)
	aobj.RemoveAttr("AttachTo", this.GetQualifiedID())
}

func (this *LTCEvent) GetAttrs(key string) []string {
	return this.common.GetAttrs(key)
}

func (this *LTCEvent) HasAttr(key string, value string) bool {
	return this.common.HasAttr(key, value)
}

func (this *LTCEvent) AddAttr(key string, value string) {
	this.common.AddAttr(key, value)
}

func (this *LTCEvent) RemoveAttr(key string, value string) {
	this.common.RemoveAttr(key, value)
}

func (this *LTCEvent) getNumberID() LTCNumberID {
	return this.common.getNumberID()
}

//-- interface LTCWarningObj

func (this *LTCEvent) Summary() string {
	// TODO
}

func (this *LTCEvent) IsUnknown() bool {
	// TODO
}

//-- interface LTCTomlPrintable

func (this *LTCChart) PrintTOML() string {
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

func CreateEvent() *LTCEvent {
	return &LTCEvent{
		Note: LTCNote{},

		common: ltcMainObjCommon{
			chart: nil,
			parent: nil,
			numID: NID_Invalid,

			extraNoteAnnexs: []*LTCAnnex{},
			attachedAnnexs: []*LTCAnnex{},
			attrs: map[string][]string{},
		}

		localTitle: "",
		localStartDate: nil,
		localEndDate: nil,

		subchartLink: nil,
		loadedSubchart: nil,
		embedLink: nil,
		loadedEmbed: nil,
	}
}

func LoadEvent(uri string) (*LTCEvent, error) {
	// TODO
}
