package libltc

import (
	"errors"
)

func (ltcChart *LTCChart) Export() (*LTCFile, []LTCWarning, error) {
	// TODO: print warnings during conversion (and what it did to alleviate it).
}
