package libltc

import (
	"errors"
	"math"
)

//-- type Version

type Version int 
const (
	VER_26_09_1 Version = iota
)
const VER_Unknown Version = -1

//-- type Version: interface WarningObj

func (ver Version) Summary() string {
	switch ver {
	case VER_26_09_1:
		return "26.09.1"
	default:
		return "?"
	}
}

func (ver Version) IsUnknown() bool {
	return ver.Summary() == "?"
}

//-- type Version: method (creation)

func GetVersion(reqverstr string) Version {
	switch reqverstr {
	case "26.09.1":
		return VER_26_09_1
	default:
		return VER_Unknown
	}
}

//-- type TimeKind

type TimeKind int
const (
	TK_Year TimeKind = iota
	TK_Month TimeKind
	TK_Day TimeKind
	TK_Hour TimeKind
	TK_Minute TimeKind
	TK_Second TimeKind
)

//-- type NumberID

type NumberID uint64
const NID_Max = math.MaxUint64 - 1
const NID_Invalid = math.MaxUint64

//-- type MainObjKind

type MainObjKind int
const (
	MOK_Event MainObjKind = iota
	MOK_Annex MainObjKind
	MOK_Import MainObjKind
	MOK_Unknown MainObjKind
)

//-- type MainObjKind: method (stringify)

func (mok MainObjKind) Prefix() (string, error) {
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
