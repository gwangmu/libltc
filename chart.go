package libltc

import (
	"errors"
	"time"
)


type LTCChart struct {
	Filepath string
	Version LTCVersion
	Note LTCNote

	Setting LTCSetting
	Subject LTCSubject

	events []*LTCEvent	// Sorted by StartDate (unknown first), incl. imported events.
	annexs []*LTCAnnex
	imports []*LTCImport
}

type LTCSetting struct {
	DisplayLanguage string 
	CalendarSystem string
	NoteFormat string
}

type LTCSubject struct {
	Name LTCName
	StartDate LTCTime
	EndDate LTCTime
	Sex string
}


func (this *LTCChart) getNextEventID() (uint64, error) {
	// Next event ID = Last event ID number + 1
	// FIX: error out if 'maxEventID' == UINT_MAX. Unlikely, but still.
	var maxEventID uint64 = 0
	for _, events := range this.Events {
		for _, event := range events {
			maxEventID = max(maxEventID, event.eventID)
		}
	}
	if (maxEventID == math.MaxUint64) {
		return 0, errors.New("Cannot get the next event ID")
	} else {
		return maxEventID + 1, nil
	}
}

func (this *LTCChart) getNextAnnexID() (uint64, error) {
	// Next annex ID = Last annex ID number + 1
	// FIX: error out if 'maxAnnexID' == UINT_MAX. Unlikely, but still.
	var maxAnnexID uint64 = 0
	for _, annex := range this.Annexs {
		maxAnnexID = max(maxAnnexID, annex.annexID)
	}
	if (maxAnnexID == math.MaxUint64) {
		return 0, errors.New("Cannot get the next annex ID")
	} else {
		return maxAnnexID + 1, nil
	}
}

func (this *LTCChart) Summary() (ret string) {
	ret = "the LTC chart"
	if !this.Entity.IsUnknown() {
		ret += " for " + this.Entity.Summary()
	} else if this.Filepath != "" {
		ret += " from the file at \'" + this.Filepath "\'"
	}
	return
}

func (this *LTCChart) IsUnknown() bool {
	return this.Entity.IsUnknown() && this.Filepath == ""
}

func (this *LTCChart) CreateEvent() (*LTCEvent, error) {
	newEventID, err := this.getNextEventID()
	if err != nil {
		return nil, err
	}

	return &LTCEvent{
		Title: "",
		StartDate: LTCTime{},
		EndDate: LTCTime{},
		Note: "",
		Attrs: map[string]string{},

		chart: &this,
		eventID: newEventID,
	}
}

func (this *LTCChart) CreateAnnex() (*LTCAnnex, error) {
	newAnnexID, err := this.getNextAnnexID()
	if err != nil {
		return nil, err
	}

	return &LTCAnnex{
		Format: "text",
		Data: []byte{},
		Attrs: map[string]string{},

		chart: &this,
		annexID: newAnnexID,
	}
}
