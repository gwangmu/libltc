package ltc

import (
	"errors"
	"fmt"

	"libltc/warning"
	"libltc/internal/coder"
	"libltc/internal/file"
)

type Annex struct {
	MajorCommon

	Note Note
	Title string

	format string
	encoding string
	rawData []byte
	isDecoded bool

	encodedData string
	indirLink *string

	AttachTo EPSource[*MajorCommon]		// Specifying 'AttachTo'
	ExtraNoteOf EPSource[*MajorCommon]	// Specifying 'ExtraNoteOf'
}

//-- interface IMajor

func (this *Annex) resolveReferenceTo(c IChart) {
	for _, cobj := range c.GetAllMajors(true) {
		CreateLink[AttachTo](this, cobj.GetCommon())
		CreateLink[ExtraNoteOf](this, cobj.GetCommon())
	}
}

func (this *Annex) unresolveReferenceTo(c IChart) {
	for _, cobj := range c.GetAllMajors(true) {
		BreakLink[AttachTo](this, cobj.GetCommon())
		BreakLink[ExtraNoteOf](this, cobj.GetCommon())
	}
}

func (this *Annex) initialize() {
	*this = Annex{
		MajorCommon: *CreateEmptyMajorCommon(),

		Note: Note{},
		Title: "",

		format: "",
		encoding: "",
		rawData: []byte{},
		encodedData: "",
		isDecoded: true,
	}
	this.enclosing = this
}

//-- interface IWarningObj

func (this *Annex) Summary() (ret string) {
	if this.Title != "" {
		ret = "an annex '" + this.Title + "'"
	} else {
		ret = "a untitled annex"
	}

	extraStr := warning.BuildSummaryString(
		WSE{"ID", this.GetIntrinsicID()},
		WSE{"format", this.format}, 
		WSE{"encoding", this.encoding},
	)
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}
	
	return
}

func (this Annex) IsUnknown() bool {
	return this.Title == "" && this.GetIntrinsicID() == "" &&
		this.format == "" && this.encoding == ""
}

//-- interface IDiagnosable

func (this *Annex) DiagnoseLocal() (warns warning.Warnings) {
	if _, ok := coder.Decoders[this.encoding]; !ok {
		warns.Add("@0@ specifies unsupported encoding '%s'", this, this.encoding)
	}

	if len(this.encodedData) != 0 && this.IsIndirect() {
		warns.Add("@0@ attempted to specify both data itself and an indirect link. Favoring data...", this)
		this.indirLink = nil
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

	// Check invalid ID.
	if this.numID == NID_Invalid {
		newNumID := this.chart.getNextNumberID(MOK_Annex)
		this.numID = newNumID
		warns.Add("@0@ had an invalid ID. auto-corrected to 'e%d'.", newNumID)
	}

	// Check duplicate ID.
	for _, aobj := range this.chart.GetAnnexs() {
		if aobj.numID == this.numID {
			this.AddAttr("OldID", this.GetIntrinsicID())
			newNumID := this.chart.getNextNumberID(MOK_Annex)
			this.numID = newNumID
			warns.Add("@0@ had a duplicated ID. auto-corrected to 'a%d'.", newNumID)
			break
		}
	}

	for akey, avals := range this.attrs {
		if akey == "AttachTo" {
			for _, aval := range avals {
				var objTo *MajorCommon
				for _, aobjIn := range this.AttachTo.Get() {
					if aobjIn.GetRelQualifiedID(this.GetChart()) == aval {
						objTo = aobjIn
						break
					}
				}

				if objTo == nil || objTo.IsUnresolved() {
					warns.Add("@0@ has a dangling 'AttachTo' to '%s'.", this, aval)
				} 

				if objTo != nil && !objTo.Attached.Has(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'AttachTo' to '%s', but it doesn't reference back. corrected.", this, aval)
					objTo.Attached.add(this)
				}
			}
		} else if akey == "ExtraNoteOf" {
			for _, aval := range avals {
				var objOf *MajorCommon
				for _, aobjIn := range this.ExtraNoteOf.Get() {
					if aobjIn.GetRelQualifiedID(this.GetChart()) == aval {
						objOf = aobjIn
						break
					}
				}

				if objOf == nil || objOf.IsUnresolved() {
					warns.Add("@0@ has a dangling 'ExtraNoteOf' to '%s'.", this, aval)
				} 

				if objOf != nil && !objOf.ExtraNote.Has(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'ExtraNoteOf' to '%s', but it doesn't reference back. corrected.", this, aval)
					objOf.ExtraNote.add(this)
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
		if this.IsIndirect() {
			if loadedData, err := file.ReadFileFromURI(*this.indirLink); err != nil {
				return []byte{}, errors.New("Failed to load indirect link")
			} else {
				this.rawData = loadedData
				this.isDecoded = true
			}
		} else {
			if decoder, ok := coder.Decoders[this.encoding]; !ok {
				return []byte{}, errors.New("No decoder for encoding '" + this.encoding + "'")
			} else {
				if rawData, err := decoder(this.encodedData); err != nil {
					return []byte{}, errors.New("Failed to decode data")
				} else {
					this.rawData = rawData
					this.isDecoded = true
				}
			}
		}
	}
	return this.rawData, nil
}

func (this *Annex) GetEncodedData() (string, error) {
	if this.IsIndirect() {
		return "", errors.New("Cannot return encoded data for indirect annex")
	} else {
		return this.encodedData, nil
	}
}

func (this *Annex) IsIndirect() bool {
	return this.indirLink != nil
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
func (this *Annex) SetExtraNoteOf(obj *MajorCommon) {
	CreateLink[ExtraNoteOf](this, obj)
}

// Public wrapper of 'removeExtraNoteOfObject'.
func (this *Annex) UnsetExtraNoteOf(obj *MajorCommon) {
	UnreserveLink[ExtraNoteOf](this, obj.GetAbsQualifiedID())
}

// Public wrapper of 'removeExtraNoteOfObject'.
func (this *Annex) UnsetExtraNoteOfByID(absQualID string) {
	UnreserveLink[ExtraNoteOf](this, absQualID)
}

// Public wrapper of 'addAttachToObject'.
func (this *Annex) SetAttachTo(obj *MajorCommon) {
	CreateLink[AttachTo](this, obj)
}

// Public wrapper of 'removeAttachToObject'.
func (this *Annex) UnsetAttachTo(obj *MajorCommon) {
	UnreserveLink[AttachTo](this, obj.GetAbsQualifiedID())
}

// Public wrapper of 'removeAttachToObject'.
func (this *Annex) UnsetAttachToByID(absQualID string) {
	UnreserveLink[AttachTo](this, absQualID)
}

//-- method (creation)

func createAnnexFromParsed(o file.Annex) (*Annex, warning.Warnings) {
	annex := CreateEmptyAnnex()
	warns := warning.Warnings{}

	kind, numID := convIDStringToInternal(o.ID)
	annex.numID = numID
	annex.attrs = convAttrsFileToChart(o.Attrs)

	if kind != MOK_Annex || numID == NID_Invalid {
		annex.numID = NID_Invalid
		annex.AddAttr("OldID", o.ID)
	}

	note, moreWarns := CreateNoteFromString(o.Note)
	warns.Concat(moreWarns)
	annex.Note = note

	annex.Title = o.Title
	annex.format = o.Format
	annex.encoding = o.Encoding

	annex.encodedData = o.Data
	annex.isDecoded = false		// NOTE: lazy decode.

	if o.Indirect != nil {
		indirLink := *o.Indirect
		annex.indirLink = &indirLink
	}

	if annex.format == "" {
		annex.format = "txt"
	}

	if annex.encoding == "" {
		annex.encoding = "none"
	}
	
	// Diagnose and partially auto-correct.
	moreWarns = annex.DiagnoseLocal()
	warns.Concat(moreWarns)
	
	// Create unresolved references.
	for akey, avals := range annex.attrs {
		if akey == "AttachTo" {
			for _, aval := range avals {
				ReserveLink[AttachTo](annex, aval)
			}
		} else if akey == "ExtraNoteOf" {
			for _, aval := range avals {
				ReserveLink[ExtraNoteOf](annex, aval)
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
