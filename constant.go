package libltc

import (
	"errors"
	"math"
)

type LTCVersion int 
const (
	Ver_26_09_1 LTCVersion = iota
)
var ltcVersionToStr = map[LTCVersion]string{
	Ver_26_09_1: "26.09.1",
}

const Ver_Unknown LTCVersion = -1
const ltcVersionUnknownStr := "?"

func (ver LTCVersion) String() (string, error) {
	if verstr, ok := ltcVersionToStr[ver]; ok {
		return verstr
	}
	return ltcVersionUnknownStr, errors.New("Unknown LTCVersion")
}

func (ver LTCVersion) Summary() string {
	if verstr, ok := ltcVersionToStr[ver]; ok {
		return verstr
	}
	return ltcVersionUnknownStr
}

func (ver LTCVersion) IsUnknown() bool {
	_, exists := ltcVersionToStr[ver]
	return !exists
}

func GetLTCVersion(reqverstr string) (outver LTCVersion, e error) {
	for ver, verstr := range ltcVersionToStr {
		if verstr == reqverstr {
			outver = ver
			e = nil
			return
		}
	}
	return Ver_Unknown, errors.New("Unknown version")
}

type LTCTimeKind int
const (
	TK_Year LTCTimeKind = iota
	TK_Month LTCTimeKind
	TK_Day LTCTimeKind
	TK_Hour LTCTimeKind
	TK_Minute LTCTimeKind
	TK_Second LTCTimeKind
)

type LTCNumberID uint64
const NID_Max = math.MaxUint64 - 1
const NID_Invalid = math.MaxUint64
