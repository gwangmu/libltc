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

	AttachTo Endpoint[MainObjCommon, *MainObjCommon]
	ExtraNoteOf Endpoint[MainObjCommon, *MainObjCommon]
}

//-- interface IMainObj

func (this *Annex) initialize() {
	*this = Annex{
		MainObjCommon: *CreateEmptyMainObjCommon(),

		Note: Note{},
		Title: "",

		format: "",
		encoding: "",
		rawData: []byte{},
		encodedData: "",
		isDecoded: true,
	}
	this.kind = MOK_Annex
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
				var objTo *MainObjCommon
				for _, aobjIn := range this.AttachTo.Get() {
					if aobjIn.GetQualifiedID() == this.convIDRelToAbs(aval) {
						objTo = aobjIn
						break
					}
				}

				if objTo == nil || objTo.IsUnresolved() {
					warns.Add("@0@ has a dangling 'AttachTo' to '%s'.", this, aval)
				} 

				if objTo != nil && !objTo.Attached.Has(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'AttachTo' to '%s', but it doesn't reference back. corrected.", this, aval)
					CreateLink[AttachTo](this, objTo)
				}
			}
		} else if akey == "ExtraNoteOf" {
			for _, aval := range avals {
				var objOf *MainObjCommon
				for _, aobjIn := range this.ExtraNoteOf.Get() {
					if aobjIn.GetQualifiedID() == this.convIDRelToAbs(aval) {
						objOf = aobjIn
						break
					}
				}

				if objOf == nil || objOf.IsUnresolved() {
					warns.Add("@0@ has a dangling 'ExtraNoteOf' to '%s'.", this, aval)
				} 

				if objOf != nil && !objOf.ExtraNote.Has(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'ExtraNoteOf' to '%s', but it doesn't reference back. corrected.", this, aval)
					CreateLink[ExtraNoteOf](this, objOf)
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

// `(un)resolve*` kind of methods: only for the objects that can be "unresolved."
// For the lists having `(un)resolve*`, its `remove*` function accepts a
// qualified ID instead of an object pointer because the underlying object
// in the list could not have been resolved yet. For the other lists, `remove*`
// functions accept object pointers.

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

// Public wrapper of 'addExtraNoteOfObject'.
func (this *Annex) SetExtraNoteOf(obj *MainObjCommon) {
	CreateLink[ExtraNoteOf](this, obj)
}

// Public wrapper of 'removeExtraNoteOfObject'.
func (this *Annex) UnsetExtraNoteOf(obj *MainObjCommon) {
	RemoveLink[ExtraNoteOf](this, obj.GetQualifiedID())
}

// Public wrapper of 'removeExtraNoteOfObject'.
func (this *Annex) UnsetExtraNoteOfByID(absQualID string) {
	RemoveLink[ExtraNoteOf](this, absQualID)
}

// Public wrapper of 'addAttachToObject'.
func (this *Annex) SetAttachTo(obj *MainObjCommon) {
	CreateLink[AttachTo](this, obj)
}

// Public wrapper of 'removeAttachToObject'.
func (this *Annex) UnsetAttachTo(obj *MainObjCommon) {
	RemoveLink[AttachTo](this, obj.GetQualifiedID())
}

// Public wrapper of 'removeAttachToObject'.
func (this *Annex) UnsetAttachToByID(absQualID string) {
	RemoveLink[AttachTo](this, absQualID)
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
				uobj := CreateEmptyMainObjCommon()
				uobj.markUnresolved(aval)
				annex.AttachTo.add(uobj)
			}
		} else if akey == "ExtraNoteOf" {
			for _, aval := range avals {
				uobj := CreateEmptyMainObjCommon()
				uobj.markUnresolved(aval)
				annex.ExtraNoteOf.add(uobj)
			}
		}
	}

	return annex, warns
}

func CreateEmptyAnnex() *Annex {
	aobj := Annex{}
	aobj.initialize()
	return &aobj
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
