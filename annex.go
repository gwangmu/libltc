package libltc

import (
	"errors"
	"fmt"

	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/coder"
	"github.com/gwangmu/libltc/internal/file"
)

type Annex struct {
	MainObjCommon

	Note Note
	Title string

	format string
	encoding string
	rawData []byte
	encodedData string
	isDecoded bool

	attachedToObjs []IMainObj
	extraNoteOfObjs []IMainObj
}

//-- interface IWarningObj

func (this *Annex) Summary() string {
	// TODO
	panic("Unimplemented")
}

func (this *Annex) IsUnknown() bool {
	// TODO
	panic("Unimplemented")
}

//-- interface IDiagnosable

func (this *Annex) DiagnoseLocal() (warns warning.Warnings) {
	if _, ok := coder.Decoders[this.encoding]; !ok {
		warns.Add("@0@ specifies unsupported encoding '%s'", this, this.encoding)
	}

	return
}

func (this *Annex) DiagnoseNonLocal() (warns warning.Warnings) {
	if this.GetChart() == nil {
		warns.Add("@0@ is not associated to any chart.", this)
		return
	}

	// TODO
	panic("Unimplemented")
}

//-- method (getters)

func (this *Annex) GetFormat() string {
	return this.format
}

func (this *Annex) GetEncoding() string {
	return this.encoding
}

func (this *Annex) GetRawData() ([]byte, error) {
	if !this.isDecoded {
		if decoder, ok := coder.Decoders[this.encoding]; !ok {
			return []byte{}, errors.New("No decoder for 'Encoding'.")
		} else {
			if rawData, err := decoder(this.encodedData); err != nil {
				return []byte{}, errors.New("Failed to decode data.")
			} else {
				this.rawData = rawData
				this.isDecoded = true
			}
		}
	}
	return this.rawData, nil
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

func (this *Annex) SetEncoding(encoding string) {
	this.encoding = encoding
}

func (this *Annex) SetRawData(data []byte) error {
	// Early-encoding and fast-fail.
	encoded, err := tryEncode(this.encoding, data)
	if err == nil {
		this.rawData = data
		this.encodedData = encoded
		this.isDecoded = true
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

// `Swap` kind of methods: only for the objects that can be "unresolved."

func (this *Annex) swapAttachedToObject(oldobj IMainObj, newobj IMainObj) {
	for i, elem := range this.attachedToObjs {
		if elem == oldobj {
			this.attachedToObjs[i] = newobj
			return
		}
	}
}

func (this *Annex) removeAttachedToObject(obj IMainObj) {
	for i, elem := range this.attachedToObjs {
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
	for i, elem := range this.extraNoteOfObjs {
		if elem == obj {
			this.extraNoteOfObjs = append(this.extraNoteOfObjs[:i], this.extraNoteOfObjs[i+1:]...) 
			break
		}
	}
}

//-- method (high-level operation)

func (this *Annex) SetAnnex(format string, encoding string, data []byte) error {
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
	panic("Unimplemented")
}

func (this *Annex) UnsetExtraNoteOf(obj IMainObj) {
	// TODO
	panic("Unimplemented")
}

func (this *Annex) SetAttachTo(obj IMainObj) {
	// TODO
	panic("Unimplemented")
}

func (this *Annex) UnsetAttachTo(obj IMainObj) {
	// TODO
	panic("Unimplemented")
}

//-- method (creation)

func createAnnexFromParsed(o file.Annex) (*Annex, warning.Warnings) {
	annex := CreateEmptyAnnex()
	warns := warning.Warnings{}

	annex.kind, annex.numID = convIDStringToInternal(o.ID)
	annex.attrs = convAttrsFileToChart(o.Attrs)

	note, moreWarns := CreateNoteFromString(o.Note)
	warns.Concat(moreWarns)
	annex.Note = note

	annex.Title = o.Title
	annex.format = o.Format
	annex.encoding = o.Encoding

	annex.encodedData = o.Data
	annex.isDecoded = false		// NOTE: lazy decode.
	
	moreWarns = annex.DiagnoseLocal()
	warns.Concat(moreWarns)

	return annex, warns
}

func CreateEmptyAnnex() *Annex {
	return &Annex{
		MainObjCommon: CreateEmptyMainObjCommon(MOK_Annex),

		Note: Note{},
		Title: "",

		format: "",
		encoding: "",
		rawData: []byte{},
		encodedData: "",
		isDecoded: true,

		attachedToObjs: []IMainObj{},
		extraNoteOfObjs: []IMainObj{},
	}
}

//-- method (private)

func tryEncode(encoding string, data []byte) (string, error) {
	if encoder, ok := coder.Encoders[encoding]; ok {
		if encoded, err := encoder(data); err == nil {
			return encoded, nil
		} else {
			return "", errors.New("Encoding failed")
		}
	} else {
		return "", errors.New(fmt.Sprintf("Unrecognized encoding '%s'", encoding))
	}
}
