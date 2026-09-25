package libltc

import (
	"errors"
	"fmt"
	"libltc/coder"
)

type LTCAnnex struct {
	Common ltcMainObjCommon

	Note LTCNote
	Title string

	format string
	encoding string
	rawData []byte
	encodedData string
}

//-- interface LTCWarningObj

func (this *LTCAnnex) Summary() string {
	// TODO
}

func (this *LTCAnnex) IsUnknown() bool {
	// TODO
}

//-- interface LTCTomlPrintable

func (this *LTCAnnex) PrintTOML() string {
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

//-- method (creation)

func CreateEmptyAnnex() *LTCAnnex {
	return &LTCAnnex{
		Common: LTCMainObjCommon{
			kind: MOK_Annex,

			chart: nil,
			parent: nil,
			numID: NID_Invalid,

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
	}
}

func CreateAnnexFromData(format string, encoding string, data []byte) (*LTCAnnex, error) {
	aobj := CreateEmptyAnnex()
	if err := aobj.SetData(format, encoding, data); err != nil {
		return nil, err
	} else {
		return aobj, nil
	}
}
