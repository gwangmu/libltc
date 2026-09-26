package libltc

import (
	"errors"
	"fmt"

	"libltc/coder"
	"libltc/file"
)

// `Swap` kind of methods: only for the objects that can be "unresolved."

type LTCAnnex struct {
	Note LTCNote
	Title string

	common LTCMainObjCommon

	format string
	encoding string
	rawData []byte
	encodedData string

	attachedToObjs []LTCMainObj
	extraNoteOfObjs []LTCMainObj
}

//-- interface LTCMainObj

func (this *LTCAnnex) Common() *LTCMainObjCommon {
	return &this.common
}

func (this *LTCAnnex) IsUnresolved() bool {
	return this.common.IsUnresolved()
}

//-- interface LTCWarningObj

func (this *LTCAnnex) Summary() string {
	// TODO
}

func (this *LTCAnnex) IsUnknown() bool {
	// TODO
}

//-- method (getters and setters)

func (this *LTCAnnex) GetFormat() string {
	return this.format
}

func (this *LTCAnnex) GetEncoding() string {
	return this.encoding
}

func (this *LTCAnnex) GetRawData() []byte {
	return this.rawData
}

func (this *LTCAnnex) GetEncodedData() string {
	return this.encodedData
}

func (this *LTCAnnex) GetAttachedToObjects() []LTCMainObj {
	return this.attachedToObjs
}

func (this *LTCAnnex) GetExtraNoteOfObjects() []LTCMainObj {
	return this.extraNoteOfObjs
}

func (this *LTCAnnex) SetRawData(format string, encoding string, data []byte) error {
	// Early-encode and fail fast.
	if encoder, ok := LTCEncoders[encoding]; ok {
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

func (this *LTCAnnex) SetFormat(format string) {
	this.format = format
}

func (this *LTCAnnex) SetEncoding(encoding string) error {
	if encoder, ok := LTCEncoders[encoding]; ok {
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

func (this *LTCAnnex) HasAttachedToObject(obj LTCMainObj) bool {
	for _, elem := range this.attachedToObjs {
		if elem == obj {
			return true
		}
	}
	return false
}

func (this *LTCAnnex) AddAttachedToObject(obj LTCMainObj) {
	if !this.HasAttachedToObject(obj) {
		this.attachedToObjs = append(this.attachedToObjs, obj)
	}
}

func (this *LTCAnnex) SwapAttachedToObject(oldobj LTCMainObj, newobj LTCMainObj) {
	for i, elem := range this.attachedToObjs {
		if elem == oldobj {
			this.attachedToObjs[i] = newobj
			return
		}
	}
}

func (this *LTCAnnex) RemoveAttachedToObject(obj LTCMainObj) {
	for _, elem := range this.attachedToObjs {
		if elem == obj {
			this.attachedToObjs = append(this.attachedToObjs[:i], this.attachedToObjs[i+1:]...) 
			break
		}
	}
}

func (this *LTCAnnex) HasExtraNoteOfObject(obj LTCMainObj) bool {
	for _, elem := range this.extraNoteOfObjs {
		if elem == obj {
			return true
		}
	}
	return false
}

func (this *LTCAnnex) AddExtraNoteOfObject(obj LTCMainObj) {
	if !this.HasExtraNoteOfObject(obj) {
		this.extraNoteOfObjs = append(this.extraNoteOfObjs, obj)
	}
}

func (this *LTCAnnex) SwapExtraNoteOfObject(oldobj LTCMainObj, newobj LTCMainObj) {
	for i, elem := range this.extraNoteOfObjs {
		if elem == oldobj {
			this.extraNoteOfObjs[i] = newobj
			return
		}
	}
}

func (this *LTCAnnex) RemoveExtraNoteOfObject(obj LTCMainObj) {
	for _, elem := range this.extraNoteOfObjs {
		if elem == obj {
			this.extraNoteOfObjs = append(this.extraNoteOfObjs[:i], this.extraNoteOfObjs[i+1:]...) 
			break
		}
	}
}

//-- method (import and export)

func (this *LTCAnnex) Import(feobj *file.LTCAnnex, chart *LTCChart) error {
	// TODO
}

func (this *LTCAnnex) Export() (*file.LTCAnnex, error) {
	// TODO
}

//-- method (creation and disposal)

func CreateEmptyAnnex() *LTCAnnex {
	return &LTCAnnex{
		Common: LTCMainObj{
			kind: MOK_Annex,

			chart: nil,
			parent: nil,
			numID: NID_Invalid, 
			fullID: "",

			extraNoteAnnexs: []*LTCAnnex{},
			attachedAnnexs: []*LTCAnnex{},
			attrs: map[string][]string{},
		}

		Note: LTCNote{},
		Title: "",

		format: "",
		encoding: "",
		rawData: []byte{},
		encodedData: "",

		attachedToObjs: []*LTCMainObj{},
		extraNoteOfObjs: []*LTCMainObj{},
	}
}
