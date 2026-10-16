package libltc

import (
	"errors"
	"fmt"

	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/coder"
	"github.com/gwangmu/libltc/internal/file"
)

type Annex struct {
	MajorCommon

	Note Note
	Title string

	createDate Time

	format string
	encoding string
	rawData []byte
	isDecoded bool

	encodedData string
	indirLink *string

	AttachTo	// Specifying 'AttachTo'
	ExtraNoteOf	// Specifying 'ExtraNoteOf'
}

func (this *Annex) Equal(that IMajor) bool {
	return this.GetCommon() == that.GetCommon()
}

//-- interface IMajor

func (this *Annex) resolveReferenceTo(c *Chart) {
	for _, cobj := range c.GetAllMajors(true) {
		this.AttachTo.Link(cobj, false)
		this.ExtraNoteOf.Link(cobj, false)
	}
}

func (this *Annex) unresolveReferenceTo(c *Chart) {
	for _, cobj := range c.GetAllMajors(true) {
		this.AttachTo.Unlink(cobj)
		this.ExtraNoteOf.Unlink(cobj)
	}
}

//-- interface IWarningObj

func (this *Annex) Summary() (ret string) {
	if this.Title != "" {
		ret = "the annex '" + this.Title + "'"
	} else {
		ret = "a untitled annex"
	}

	maybeOldID := ""
	if this.HasAttrKey("OldID") && len(this.GetAttrs("OldID")) > 0 {
		maybeOldID = this.GetAttrs("OldID")[0]
	}

	extraStr := warning.BuildSummaryString(
		wse{"ID", this.GetIntrinsicID()},
		wse{"OldID", maybeOldID},
		wse{"Created", this.createDate}, 
		wse{"Format", this.format}, 
		wse{"Encoding", this.encoding},
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
	if this.numID == mNID_Invalid {
		newNumID := this.chart.getNextNumberID(MOK_Annex)
		this.numID = newNumID
		warns.Add("@0@ had an invalid ID. auto-corrected to 'a%d'.", this, newNumID)
	}

	// Check duplicate ID.
	for _, aobj := range this.chart.GetAnnexs() {
		// Only check against the object on the same import level.
		if aobj != this && aobj.parent == this.parent && aobj.numID == this.numID {
			this.AddAttr("OldID", this.GetIntrinsicID())
			newNumID := this.chart.getNextNumberID(MOK_Annex)
			this.numID = newNumID
			warns.Add("@0@ had a duplicated ID. auto-corrected to 'a%d'.", this, newNumID)
			break
		}
	}

	for akey, avals := range this.attrs {
		if akey == "AttachTo" {
			for _, aval := range avals {
				var objTo IMajor
				for _, aobjIn := range this.AttachTo.Get() {
					if aobjIn.GetRelQualifiedID(this.GetChart()) == aval {
						objTo = aobjIn
						break
					}
				}

				if objTo == nil || objTo.IsUnresolved() {
					warns.Add("@0@ has a dangling 'AttachTo' to '%s'.", this, aval)
				} 

				if objTo != nil && !objTo.IsUnresolved() && !objTo.GetCommon().Attached.Has(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'AttachTo' to '%s', but it doesn't reference back. corrected.", this, aval)
					objTo.GetCommon().Attached.add(this)
				}
			}
		} else if akey == "ExtraNoteOf" {
			for _, aval := range avals {
				var objOf IMajor
				for _, aobjIn := range this.ExtraNoteOf.Get() {
					if aobjIn.GetRelQualifiedID(this.GetChart()) == aval {
						objOf = aobjIn
						break
					}
				}

				if objOf == nil || objOf.IsUnresolved() {
					warns.Add("@0@ has a dangling 'ExtraNoteOf' to '%s'.", this, aval)
				} 

				if objOf != nil && !objOf.IsUnresolved() && !objOf.GetCommon().ExtraNote.Has(this) {
					warns.Add("!!!INTERNAL WARN!!! @0@ has 'ExtraNoteOf' to '%s', but it doesn't reference back. corrected.", this, aval)
					objOf.GetCommon().ExtraNote.add(this)
				}
			}
		}
	}

	return
}

//-- method (getters)

func (this *Annex) GetCreateDate() Time {
	return this.createDate
}

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

func (this *Annex) SetCreateDate(t Time) {
	oldCreateDate := this.createDate
	t.SetUsage(TUK_General)
	this.createDate = t
	if !oldCreateDate.Equal(t) && this.chart != nil {
		this.chart.changeAnnexCreateDate(this)
	}
}

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

//-- method (creation)

func createAnnexFromParsed(o file.Annex) (*Annex, warning.Warnings) {
	annex := CreateEmptyAnnex()
	warns := warning.Warnings{}

	kind, numID := convIDStringToInternal(o.ID)
	annex.numID = numID
	annex.attrs = convAttrsFileToChart(o.Attrs)

	if kind != MOK_Annex || numID == mNID_Invalid {
		annex.numID = mNID_Invalid
		annex.AddAttr("OldID", o.ID)
	}

	note, moreWarns := CreateNoteFromString(o.Note)
	warns.Concat(moreWarns)
	annex.Note = note

	annex.Title = o.Title
	annex.format = o.Format

	if o.CreateDate != nil {
		annex.createDate = createTimeFromParsed(o.CreateDate, TUK_General)
	}

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
				annex.AttachTo.Reserve(aval)
			}
		} else if akey == "ExtraNoteOf" {
			for _, aval := range avals {
				annex.ExtraNoteOf.Reserve(aval)
			}
		}
	}

	return annex, warns
}

func CreateEmptyAnnex() *Annex {
	aobj := &Annex{
		MajorCommon: *CreateEmptyMajorCommon(),

		Note: Note{},
		Title: "",

		createDate: CreateEmptyTime(TUK_General),

		format: "",
		encoding: "",
		rawData: []byte{},
		encodedData: "",
		isDecoded: true,
	}
	aobj.enclosing = aobj
	aobj.AttachTo.initialize(aobj)
	aobj.ExtraNoteOf.initialize(aobj)
	return aobj
}

//-- method (util)

func (this *Annex) shiftCreateDateRaw(t Time) {
	oldCreateDate := this.GetCreateDate()
	this.createDate = this.createDate.Add(t)
	if !oldCreateDate.Equal(t) && this.chart != nil {
		this.chart.changeAnnexCreateDate(this)
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
