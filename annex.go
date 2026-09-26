package libltc

import (
	"errors"
	"fmt"

	"libltc/coder"
	"libltc/file"
)

// `Swap` kind of methods: only for the objects that can be "unresolved."

type Annex struct {
	Note Note
	Title string

	common MainObjCommon

	format string
	encoding string
	rawData []byte
	encodedData string

	attachedToObjs []MainObj
	extraNoteOfObjs []MainObj
}

//-- interface MainObj

func (this *Annex) Common() *MainObjCommon {
	return &this.common
}

func (this *Annex) IsUnresolved() bool {
	return this.common.IsUnresolved()
}

//-- interface WarningObj

func (this *Annex) Summary() string {
	// TODO
}

func (this *Annex) IsUnknown() bool {
	// TODO
}

//-- method (getters and setters)

func (this *Annex) GetFormat() string {
	return this.format
}

func (this *Annex) GetEncoding() string {
	return this.encoding
}

func (this *Annex) GetRawData() []byte {
	return this.rawData
}

func (this *Annex) GetEncodedData() string {
	return this.encodedData
}

func (this *Annex) GetAttachedToObjects() []MainObj {
	return this.attachedToObjs
}

func (this *Annex) GetExtraNoteOfObjects() []MainObj {
	return this.extraNoteOfObjs
}

func (this *Annex) SetRawData(format string, encoding string, data []byte) error {
	// Early-encode and fail fast.
	if encoder, ok := Encoders[encoding]; ok {
		if encoded, err := encoder(data); err == nil {
			this.format = format
			this.encoding = encoding
			this.rawData = data
			this.encodedData = encoded
			return nil
		} else {
			return errors.New("Encoding failed")
		}
	} else {
		return errors.New(fmt.Sprintf("Unrecognized encoding '%s'", encoding))
	}
}

func (this *Annex) SetFormat(format string) {
	this.format = format
}

func (this *Annex) SetEncoding(encoding string) error {
	if encoder, ok := Encoders[encoding]; ok {
		if encoded, err := encoder(this.rawData); err == nil {
			this.encoding = encoding
			this.encodedData = encoded
			return nil
		} else {
			return errors.New("Encoding failed")
		}
	} else {
		return errors.New(fmt.Sprintf("Unrecognized encoding '%s'", encoding))
	}
}

func (this *Annex) HasAttachedToObject(obj MainObj) bool {
	for _, elem := range this.attachedToObjs {
		if elem == obj {
			return true
		}
	}
	return false
}

func (this *Annex) AddAttachedToObject(obj MainObj) {
	if !this.HasAttachedToObject(obj) {
		this.attachedToObjs = append(this.attachedToObjs, obj)
	}
}

func (this *Annex) SwapAttachedToObject(oldobj MainObj, newobj MainObj) {
	for i, elem := range this.attachedToObjs {
		if elem == oldobj {
			this.attachedToObjs[i] = newobj
			return
		}
	}
}

func (this *Annex) RemoveAttachedToObject(obj MainObj) {
	for _, elem := range this.attachedToObjs {
		if elem == obj {
			this.attachedToObjs = append(this.attachedToObjs[:i], this.attachedToObjs[i+1:]...) 
			break
		}
	}
}

func (this *Annex) HasExtraNoteOfObject(obj MainObj) bool {
	for _, elem := range this.extraNoteOfObjs {
		if elem == obj {
			return true
		}
	}
	return false
}

func (this *Annex) AddExtraNoteOfObject(obj MainObj) {
	if !this.HasExtraNoteOfObject(obj) {
		this.extraNoteOfObjs = append(this.extraNoteOfObjs, obj)
	}
}

func (this *Annex) SwapExtraNoteOfObject(oldobj MainObj, newobj MainObj) {
	for i, elem := range this.extraNoteOfObjs {
		if elem == oldobj {
			this.extraNoteOfObjs[i] = newobj
			return
		}
	}
}

func (this *Annex) RemoveExtraNoteOfObject(obj MainObj) {
	for _, elem := range this.extraNoteOfObjs {
		if elem == obj {
			this.extraNoteOfObjs = append(this.extraNoteOfObjs[:i], this.extraNoteOfObjs[i+1:]...) 
			break
		}
	}
}

//-- method (relational)

func (this *Annex) MakeExtraNoteOf(obj MainObj) {
	// TODO
}

func (this *Annex) UnmakeExtraNoteOf(obj MainObj) {
	// TODO
}

func (this *Annex) MakeAttachTo(obj MainObj) {
	// TODO
}

func (this *Annex) UnmakeAttachTo(obj MainObj) {
	// TODO
}

//-- method (import and export)

func (this *Annex) Import(feobj *file.Annex) error {
	// TODO
}

func (this *Annex) Export() (*file.Annex, error) {
	// TODO
}

//-- method (creation and disposal)

func CreateEmptyAnnex() *Annex {
	return &Annex{
		Common: MainObj{
			kind: MOK_Annex,

			chart: nil,
			parent: nil,
			numID: NID_Invalid, 
			fullID: "",

			extraNoteAnnexs: []*Annex{},
			attachedAnnexs: []*Annex{},
			attrs: map[string][]string{},
		}

		Note: Note{},
		Title: "",

		format: "",
		encoding: "",
		rawData: []byte{},
		encodedData: "",

		attachedToObjs: []*MainObj{},
		extraNoteOfObjs: []*MainObj{},
	}
}
