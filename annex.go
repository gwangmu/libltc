package libltc

import (
	"errors"
	"fmt"

	"github.com/gwangmu/libltc/internal/coder"
	"github.com/gwangmu/libltc/internal/file"
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

	attachedToObjs []IMainObj
	extraNoteOfObjs []IMainObj
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

//-- method (getters)

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

func (this *Annex) GetAttachedToObjects() []IMainObj {
	return this.attachedToObjs
}

func (this *Annex) GetExtraNoteOfObjects() []IMainObj {
	return this.extraNoteOfObjs
}

func (this *Annex) HasAttachedToObject(obj IMainObj) bool {
	for _, elem := range this.attachedToObjs {
		if elem == obj {
			return true
		}
	}
	return false
}

func (this *Annex) HasExtraNoteOfObject(obj IMainObj) bool {
	for _, elem := range this.extraNoteOfObjs {
		if elem == obj {
			return true
		}
	}
	return false
}

//-- method (setters)

func (this *Annex) SetFormat(format string) {
	this.format = format
}

func (this *Annex) SetEncoding(encoding string) error {
	encoded, err := tryEncode(encoding, this.rawData)
	if err == nil {
		this.encoding = encoding
		this.encodedData = encoded
		return nil
	} else {
		return err
	}
}

func (this *Annex) SetRawData(data []byte) error {
	encoded, err := tryEncode(this.encoding, data)
	if err == nil {
		this.rawData = data
		this.encodedData = encoded
		return nil
	} else {
		return err
	}
}

func (this *Annex) addAttachedToObject(obj IMainObj) {
	if !this.HasAttachedToObject(obj) {
		this.attachedToObjs = append(this.attachedToObjs, obj)
	}
}

func (this *Annex) swapAttachedToObject(oldobj IMainObj, newobj IMainObj) {
	for i, elem := range this.attachedToObjs {
		if elem == oldobj {
			this.attachedToObjs[i] = newobj
			return
		}
	}
}

func (this *Annex) removeAttachedToObject(obj IMainObj) {
	for _, elem := range this.attachedToObjs {
		if elem == obj {
			this.attachedToObjs = append(this.attachedToObjs[:i], this.attachedToObjs[i+1:]...) 
			break
		}
	}
}

func (this *Annex) addExtraNoteOfObject(obj IMainObj) {
	if !this.HasExtraNoteOfObject(obj) {
		this.extraNoteOfObjs = append(this.extraNoteOfObjs, obj)
	}
}

func (this *Annex) swapExtraNoteOfObject(oldobj IMainObj, newobj IMainObj) {
	for i, elem := range this.extraNoteOfObjs {
		if elem == oldobj {
			this.extraNoteOfObjs[i] = newobj
			return
		}
	}
}

func (this *Annex) removeExtraNoteOfObject(obj IMainObj) {
	for _, elem := range this.extraNoteOfObjs {
		if elem == obj {
			this.extraNoteOfObjs = append(this.extraNoteOfObjs[:i], this.extraNoteOfObjs[i+1:]...) 
			break
		}
	}
}

//-- method (high-level operation)

func (this *Annex) SetRawData(format string, encoding string, data []byte) error {
	// Early-encode and fail fast.
	encoded, err := tryEncode(encoding, data)
	if err == nil {
		this.format = format
		this.encoding = encoding
		this.rawData = data
		this.encodedData = encoded
		return nil
	} else {
		return err
	}
}

func (this *Annex) SetExtraNoteOf(obj IMainObj) {
	// TODO
}

func (this *Annex) UnsetExtraNoteOf(obj IMainObj) {
	// TODO
}

func (this *Annex) SetAttachTo(obj IMainObj) {
	// TODO
}

func (this *Annex) UnsetAttachTo(obj IMainObj) {
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
		Common: MainObjCommon{
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

		attachedToObjs: []IMainObj{},
		extraNoteOfObjs: []IMainObj{},
	}
}

//-- method (private)

func tryEncode(encoding string, data []byte) (string, error) {
	if encoder, ok := Encoders[encoding]; ok {
		if encoded, err := encoder(data); err == nil {
			return encoded, nil
		} else {
			return "", errors.New("Encoding failed")
		}
	} else {
		return "", errors.New(fmt.Sprintf("Unrecognized encoding '%s'", encoding))
	}
}
