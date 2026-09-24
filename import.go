package libltc

import (
	"errors"
)

type LTCImport struct {
	Link string
	StartDate LTCTime
	EndDate LTCTIme
	OffsetDate LTCTime
	Categories *[]string		// nil: any categories
	Note LTCNote
	
	internal ltcMainObjCommon
}

func (this LTCImport) GetChart() *LTCChart {
	return this.chart
}

func (this LTCImport) GetLocalID() string {
	return "i" + strconv.Itoa(this.importID)
}

func (this LTCImport) GetQualifiedID() (ret string) {
	if this.ref.parent != nil {
		ret = this.ref.parent.GetQualifiedID() + '/'
	}

	ret += this.GetLocalID()
	return
}
