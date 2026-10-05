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
	TK_Month
	TK_Day
	TK_Hour
	TK_Minute
	TK_Second
)

const YearMax int = math.MaxInt
const MonthMax int = 12			// Gregorian max
const DayMax int = 31			// Gregorian max
const HourMax int = 23
const MinuteMax int = 59 
const SecondMax int = 59 

const YearMin int = 1
const MonthMin int = 1
const DayMin int = 1
const HourMin int = 0
const MinuteMin int = 0
const SecondMin int = 0

//-- type NumberID

type NumberID uint64
const NID_Max = math.MaxUint64 - 1
const NID_Invalid = math.MaxUint64

//-- type MainObjKind

type MainObjKind int
const (
	MOK_Event MainObjKind = iota
	MOK_Annex
	MOK_Import
	MOK_Unknown
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

//-- type MainObjKind: method (creation)

func GetMainObjKind(prefix string) (MainObjKind, error) {
	switch prefix {
	case "e":
		return MOK_Event, nil
	case "a":
		return MOK_Annex, nil
	case "i":
		return MOK_Import, nil
	default:
		return MOK_Unknown, errors.New("Bogus main object prefix")
	}
}

//-- type IntervalKind 

// Let's say there is an "interval" (start ~ end) with respect to another
// "interval" (estart ~ eend).
//
//   |-------|==========|--------|
//   ^       ^          ^        ^
// start   estart      eend     end
//
// Normal: include equals at the boundary. (start <= estart && eend < end)
// Extended: include stradding ones. (start <= eend || estart < end)

type IntervalKind int
const (
	IK_Normal IntervalKind = iota
	IK_Extended
)
