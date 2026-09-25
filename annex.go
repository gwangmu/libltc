package libltc

import (
	"errors"
	"libltc/coder"
)

type LTCAnnex struct {
	Note LTCNote
	Title string
	Format string
	Encoding string
	Data []byte 				// Decoded data (potentially binary)

	common ltcMainObjCommon
}

//-- TODO: interface LTCMainObj

func (this LTCAnnex) GetChart() *LTCChart {
	return this.chart
}

func (this LTCAnnex) GetLocalID() string {
	return "a" + strconv.Itoa(this.annexID)
}

func (this LTCAnnex) GetQualifiedID() (ret string) {
	if this.ref.parent != nil {
		ret = this.ref.parent.GetQualifiedID() + '/'
	}

	ret += this.GetLocalID()
	return
}

//-- TODO: method (creation)
