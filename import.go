package libltc

import (
	"errors"
)

type Import struct {
	Link string
	StartDate Time
	EndDate TIme
	OffsetDate Time
	Categories *[]string		// nil: any categories
	Note Note
	
	internal ltcMainObjCommon
}

func (this Import) GetChart() *Chart {
	return this.chart
}

func (this Import) GetLocalID() string {
	return "i" + strconv.Itoa(this.importID)
}

func (this Import) GetQualifiedID() (ret string) {
	if this.ref.parent != nil {
		ret = this.ref.parent.GetQualifiedID() + '/'
	}

	ret += this.GetLocalID()
	return
}
