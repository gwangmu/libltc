package libltc

import (
	"errors"
	"math"
)

//-- type LTCVersion

type LTCVersion int 
const (
	VER_26_09_1 LTCVersion = iota
)
const VER_Unknown LTCVersion = -1

//-- type LTCVersion: interface LTCWarningObj

func (ver LTCVersion) Summary() string {
	switch ver {
	case VER_26_09_1:
		return "26.09.1"
	default:
		return "?"
	}
}

func (ver LTCVersion) IsUnknown() bool {
	return ver.Summary() == "?"
}

//-- type LTCVersion: method (creation)

func GetLTCVersion(reqverstr string) LTCVersion {
	switch reqverstr {
	case "26.09.1":
		return VER_26_09_1
	default:
		return VER_Unknown
	}
}

//-- type LTCTimeKind

type LTCTimeKind int
const (
	TK_Year LTCTimeKind = iota
	TK_Month LTCTimeKind
	TK_Day LTCTimeKind
	TK_Hour LTCTimeKind
	TK_Minute LTCTimeKind
	TK_Second LTCTimeKind
)

//-- type LTCNumberID

type LTCNumberID uint64
const NID_Max = math.MaxUint64 - 1
const NID_Invalid = math.MaxUint64

//-- type LTCMainObjKind

type LTCMainObjKind int
const (
	MOK_Event LTCMainObjKind = iota
	MOK_Annex LTCMainObjKind
	MOK_Import LTCMainObjKind
)

//-- type LTCMainObjKind: method (stringify)

func (mok LTCMainObjKind) Prefix() (string, error) {
	switch mok {
	case MOK_Event:
		return "e", nil
	case MOK_Annex:
		return "a", nil
	case MOK_Import:
		return "i", nil
	default:
		return "?", errors.New("Bogus main object kind")
	}
}
