package libltc

import (
	"errors"
)

type LTCAnnex struct {
	Format string
	Data []byte 				// Decoded data (potentially binary)
	Note LTCNote

	internal ltcMainObjCommon
}

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

