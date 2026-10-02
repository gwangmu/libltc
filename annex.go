package libltc

import (
	"errors"
	"fmt"
	"slices"

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

	attachToObjs []IMainObj
	extraNoteOfObjs []IMainObj
}

//-- interface IWarningObj

func (this *Annex) Summary() (ret string) {
	if this.Title != "" {
		ret = "an annex '" + this.Title + "'"
	} else {
		ret = "a untitled annex"
	}

	extraStr := warning.BuildSummaryString(
		WSE{"ID", this.GetLocalID()},
		WSE{"format", this.format}, 
		WSE{"encoding", this.encoding},
	)
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}
	
	return
}

func (this Annex) IsUnknown() bool {
	return this.Title == "" && this.GetLocalID() == "" &&
		this.format == "" && this.encoding == ""
}

//-- interface IDiagnosable

func (this *Annex) DiagnoseLocal() (warns warning.Warnings) {
	if _, ok := coder.Decoders[this.encoding]; !ok {
		warns.Add("@0@ specifies unsupported encoding '%s'", this, this.encoding)
	}

	for akey, avals := range this.attrs {
		if akey == "AttachTo" || akey == "ExtraNoteOf" {
			newAvals := []string{}
			for _, aval := range avals {
				if aval == "" {
					warns.Add("a 'ContinuedFrom' attribute in @0@ is empty. removed.", this)
				} else {
					newAvals = append(newAvals, aval)
				}
			}
			this.attrs[akey] = newAvals
		}
	}

	return
}

func (this *Annex) DiagnoseNonLocal() (warns warning.Warnings) {
	if this.GetChart() == nil {
		warns.Add("@0@ is not associated to any chart.", this)
		return
	}

	if this.kind != MOK_Annex || this.numID == NID_Invalid {
		panic("call DiagnoseLocal() first.")
	}

	for _, aobj := range this.chart.GetAnnexs() {
		if aobj.numID == this.numID {
			newNumID, err := this.chart.getNextNumberID(MOK_Annex)
			if err != nil {
				panic("Cannot get new number ID.")
			}
			this.numID = newNumID
			warns.Add("@0@ has a duplicated ID. auto-corrected to 'a%d'.", newNumID)
			break
		}
	}

	for akey, avals := range this.attrs {
		if akey == "AttachTo" {
			for _, aval := range avals {
				var objTo IMainObj
				for _, aobjIn := range this.attachToObjs {
					if aobjIn.GetQualifiedID() == this.convIDRelToAbs(aval) {
						objTo = aobjIn
						break
					}
				}

				if objTo == nil || objTo.IsUnresolved() {
					warns.Add("@0@ has a dangling 'AttachTo' to '%s'.", this, aval)
				} 

				if objTo != nil && !objTo.HasAttachedAnnex(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'AttachTo' to '%s', but it doesn't reference back. corrected.", this, aval)
					objTo.addAttachedAnnex(this)
				}
			}
		} else if akey == "ExtraNoteOf" {
			for _, aval := range avals {
				var objOf IMainObj
				for _, aobjIn := range this.extraNoteOfObjs {
					if aobjIn.GetQualifiedID() == this.convIDRelToAbs(aval) {
						objOf = aobjIn
						break
					}
				}

				if objOf == nil || objOf.IsUnresolved() {
					warns.Add("@0@ has a dangling 'ExtraNoteOf' to '%s'.", this, aval)
				} 

				if objOf != nil && !objOf.HasExtraNoteAnnex(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'ExtraNoteOf' to '%s', but it doesn't reference back. corrected.", this, aval)
					objOf.addExtraNoteAnnex(this)
				}
			}
		}
	}

	return
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

func (this *Annex) GetAttachToObjects() []IMainObj {
	return this.attachToObjs
}

func (this *Annex) GetExtraNoteOfObjects() []IMainObj {
	return this.extraNoteOfObjs
}

func (this *Annex) HasAttachToObject(obj IMainObj) bool {
	return slices.Contains(this.attachToObjs, obj)
}

func (this *Annex) HasExtraNoteOfObject(obj IMainObj) bool {
	return slices.Contains(this.extraNoteOfObjs, obj)
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

func (this *Annex) addAttachToObject(obj IMainObj) {
	if !this.HasAttachToObject(obj) {
		this.attachToObjs = append(this.attachToObjs, obj)
	}
	if !obj.IsUnresolved() {
		obj.addAttachedAnnex(this)
	}
}

// `(un)resolve*` kind of methods: only for the objects that can be "unresolved."

func (this *Annex) resolveAttachToObject(absQualID string, newobj IMainObj) {
	for i, elem := range this.attachToObjs {
		if elem.IsUnresolved() && elem.GetQualifiedID() == absQualID {
			this.attachToObjs[i] = newobj
			newobj.addAttachedAnnex(this)
			return
		}
	}
}

func (this *Annex) unresolveAttachToObject(absQualID string) {
	for i, elem := range this.attachToObjs {
		if !elem.IsUnresolved() && elem.GetQualifiedID() == absQualID {
			this.attachToObjs[i].removeAttachedAnnex(this)
			this.attachToObjs[i] = CreateUnresolvedMainObjCommon(MOK_Unknown, this.convIDAbsToRel(absQualID))
			return
		}
	}
}

func (this *Annex) removeAttachToObject(obj IMainObj) {
	for i, elem := range this.attachToObjs {
		if elem == obj {
			this.attachToObjs = append(this.attachToObjs[:i], this.attachToObjs[i+1:]...) 
			break
		}
	}
	if !obj.IsUnresolved() {
		obj.removeAttachedAnnex(this)
	}
}

func (this *Annex) addExtraNoteOfObject(obj IMainObj) {
	if !this.HasExtraNoteOfObject(obj) {
		this.extraNoteOfObjs = append(this.extraNoteOfObjs, obj)
	}
	if !obj.IsUnresolved() {
		obj.addExtraNoteAnnex(this)
	}
}

func (this *Annex) resolveExtraNoteOfObject(absQualID string, newobj IMainObj) {
	for i, elem := range this.extraNoteOfObjs {
		if elem.IsUnresolved() && elem.GetQualifiedID() == absQualID {
			this.extraNoteOfObjs[i] = newobj
			newobj.addExtraNoteAnnex(this)
			return
		}
	}
}

func (this *Annex) unresolveExtraNoteOfObject(absQualID string) {
	for i, elem := range this.extraNoteOfObjs {
		if !elem.IsUnresolved() && elem.GetQualifiedID() == absQualID {
			this.extraNoteOfObjs[i].removeExtraNoteAnnex(this)
			this.extraNoteOfObjs[i] = CreateUnresolvedMainObjCommon(MOK_Unknown, this.convIDAbsToRel(absQualID))
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
	if !obj.IsUnresolved() {
		obj.removeExtraNoteAnnex(this)
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
	this.addExtraNoteOfObject(obj)
}

func (this *Annex) UnsetExtraNoteOf(obj IMainObj) {
	this.removeExtraNoteOfObject(obj)
}

func (this *Annex) SetAttachTo(obj IMainObj) {
	this.addAttachToObject(obj)
}

func (this *Annex) UnsetAttachTo(obj IMainObj) {
	this.removeAttachToObject(obj)
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
	
	// Diagnose and partially auto-correct.
	moreWarns = annex.DiagnoseLocal()
	warns.Concat(moreWarns)
	
	// Create unresolved references.
	for akey, avals := range annex.attrs {
		if akey == "AttachTo" {
			for _, aval := range avals {
				uobj := CreateUnresolvedMainObjCommon(MOK_Unknown, aval)
				annex.addAttachToObject(uobj)
			}
		} else if akey == "ExtraNoteOf" {
			for _, aval := range avals {
				uobj := CreateUnresolvedMainObjCommon(MOK_Unknown, aval)
				annex.addExtraNoteOfObject(uobj)
			}
		}
	}

	return annex, warns
}

func CreateEmptyAnnex() *Annex {
	return &Annex{
		MainObjCommon: *CreateEmptyMainObjCommon(MOK_Annex),

		Note: Note{},
		Title: "",

		format: "",
		encoding: "",
		rawData: []byte{},
		encodedData: "",
		isDecoded: true,

		attachToObjs: []IMainObj{},
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
