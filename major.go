package ltc

import (
	"errors"
	"math"
	"strings"
)

type IMajor interface {
	Equal(IMajor) bool

	GetKind() MajorKind
	GetCommon() *MajorCommon
	GetEnclosingObject() IMajor

	GetChart() IChart
	setChart(c IChart)

	getNumberID() NumberID 
	setNumberID(nid NumberID)

	getUnresRelQualID() string
	markUnresolved(unresRelQualID string, refcer IMajor)
	unmarkUnresolved()
	IsUnresolved() bool

	GetIntrinsicID() string 
	GetRelQualifiedID(baseChart IChart) string
	GetAbsQualifiedID() string

	GetParent() IMajor
	SetParent(p IMajor)

	IsImported() bool	// Imported by an import object?
	IsEmbedded() bool	// Embedded in an event object? (incl. subchart objs)
	IsIncluded() bool	// Imported or embedded?
	IsLocallyImported() bool	// Imported within a subchart? 

	GetAttrs(key string) []string
	HasAttrKey(key string) bool
	AddAttr(key string, value string)
	RemoveAttr(key string, value string)

	convIDRelToAbs(relid string) string
	convIDAbsToRel(absid string) string

	resolveReferenceTo(c IChart)
	unresolveReferenceTo(c IChart)

	initialize()
}

type IMajorPtr[T any] interface {
	*T
	IMajor
}

type MajorCommon struct {
	enclosing IMajor			// Enclosing object of this common struct

	chart IChart				// Associated chart
	parent IMajor				// Included by...
	numID NumberID 				// Numeric part of ID

	unresRelQualID string		// **UNRESOLVED** local qualified ID
	unresRefcer IMajor 			// Referencer of **UNRESOLVED** obj (=this) 

	attrs map[string][]string	// Attributes

	ExtraNote	// Pointed by 'ExtraNoteOf' 
	Attached 	// Pointed by 'AttachTo'
}

func (this *MajorCommon) Equal(that IMajor) bool {
	return this == that.GetCommon()
}

//-- struct MajorCommon: method

func (this *MajorCommon) GetKind() MajorKind {
	switch this.enclosing.(type) {
	case *Event:
		return MOK_Event
	case *Annex:
		return MOK_Annex
	case *Import:
		return MOK_Import
	default:
		return MOK_Unknown
	}
}

func (this *MajorCommon) GetCommon() *MajorCommon {
	return this
}

func (this *MajorCommon) GetEnclosingObject() IMajor {
	return this.enclosing
}

// Relationship Taxonomy:
//  - Chart and main objects: "association"
//  - Main object and main object: "parental"
//  - Subchart and subchart event object: "embedding"

func (this *MajorCommon) GetChart() IChart {
	return this.chart
}

func (this *MajorCommon) setChart(c IChart) {
	this.chart = c
}

func (this *MajorCommon) getNumberID() NumberID {
	return this.numID
}

func (this *MajorCommon) setNumberID(nid NumberID) {
	this.numID = nid
}

func (this *MajorCommon) getUnresRelQualID() string {
	return this.unresRelQualID
}

func (this *MajorCommon) markUnresolved(unresRelQualID string, refcer IMajor) {
	this.unresRelQualID = unresRelQualID
	this.unresRefcer = refcer 
}

func (this *MajorCommon) unmarkUnresolved() {
	this.unresRelQualID = ""
}

func (this MajorCommon) IsUnresolved() bool {
	return this.unresRelQualID != ""
}

func (this *MajorCommon) GetIntrinsicID() string {
	if this.IsUnresolved() {
		unresRelQualID := this.unresRelQualID
		lastidx := strings.LastIndex(unresRelQualID, "/")

		if lastidx == -1 {
			return unresRelQualID
		} else if len(unresRelQualID) < lastidx + 1 {
			return "?"
		} else {
			return unresRelQualID[lastidx+1:]
		}
	} else {
		if this.numID == NID_Invalid {
			return "?"
		} else {
			return convIDInternalToString(this.GetKind(), this.numID)
		}
	}
}

func (this *MajorCommon) GetRelQualifiedID(baseChart IChart) (ret string) {
	if this.IsUnresolved() {
		ret = this.GetAbsQualifiedID()
		if this.unresRefcer != nil {
			if c := this.unresRefcer.GetChart(); c != nil {
				if eceobj := c.GetEmbeddingEvent(); eceobj != nil {
					ret = strings.TrimPrefix(eceobj.GetAbsQualifiedID(), ret)
				}
			}
		}
		return
	} else {
		if this.parent != nil && (this.GetChart() != baseChart || 
			this.parent.GetChart() == baseChart) {
			ret = this.parent.GetRelQualifiedID(baseChart) + "/"
		}
		ret += this.GetIntrinsicID()
		return
	}
}

func (this *MajorCommon) GetAbsQualifiedID() (ret string) {
	if this.IsUnresolved() {
		unresRelQualID := this.unresRelQualID
		prefixid := ""
		if this.unresRefcer != nil {
			if c := this.unresRefcer.GetChart(); c != nil {
				if eceobj := c.GetEmbeddingEvent(); eceobj != nil {
					prefixid = eceobj.GetAbsQualifiedID() + "/"
				}
			}
		}
		return prefixid + unresRelQualID
	} else {
		if this.parent != nil {
			ret = this.parent.GetAbsQualifiedID() + "/"
		}
		ret += this.GetIntrinsicID()
		return
	}
}

func (this *MajorCommon) GetParent() IMajor {
	return this.parent
}

func (this *MajorCommon) SetParent(p IMajor) {
	this.parent = p
}

func (this *MajorCommon) IsImported() bool {
	if this.parent != nil {
		if this.parent.GetKind() == MOK_Import {
			return true
		} else {
			return this.parent.IsImported()
		}
	} else {
		return false
	}
}

func (this *MajorCommon) IsEmbedded() bool {
	if this.parent != nil {
		if this.parent.GetKind() == MOK_Event {
			return true
		} else {
			return this.parent.IsEmbedded()
		}
	} else {
		return false
	}
}

func (this *MajorCommon) IsIncluded() bool {
	return this.IsImported() || this.IsEmbedded()
}

func (this *MajorCommon) IsLocallyImported() bool {
	return this.parent != nil && this.parent.GetKind() == MOK_Import 
}

func (this *MajorCommon) GetAttrs(key string) []string {
	if attrlist, ok := this.attrs[key]; ok {
		// `attrlist` CANNOT be an empty list. (at least [""])
		return attrlist
	} else {
		return []string{}
	}
}

func (this *MajorCommon) HasAttrKey(key string) bool {
	_, ok := this.attrs[key]
	return ok 
}
	
// `{Add,Remove}Attr` are passive methods; they don't update other
// objects that may be referenced by the added/removed attribute.

func (this *MajorCommon) AddAttr(key string, value string) {
	if _, ok := this.attrs[key]; !ok {
		this.attrs[key] = []string{}
	}
	this.attrs[key] = append(this.attrs[key], value)
}

func (this *MajorCommon) RemoveAttr(key string, value string) {
	if attrlist, ok := this.attrs[key]; ok {
		for i, elem := range attrlist {
			if elem == value {
				newattrlist := append(attrlist[:i], attrlist[i+1:]...)
				if len(newattrlist) != 0 {
					this.attrs[key] = newattrlist
				} else {
					delete(this.attrs, key)
				}
				return
			}
		}
	}
}

func (this *MajorCommon) convIDRelToAbs(relid string) string {
	if this.parent != nil {
		return this.parent.GetAbsQualifiedID() + "/" + relid
	} else {
		return relid
	}
}

func (this *MajorCommon) convIDAbsToRel(absid string) string {
	if this.parent != nil {
		return strings.TrimPrefix(absid, this.parent.GetAbsQualifiedID() + "/")
	} else {
		return absid
	}
}

func (this *MajorCommon) resolveReferenceTo(c IChart) {
	// No EPSource. Nothing to do.
	return
}

func (this *MajorCommon) unresolveReferenceTo(c IChart) {
	// No EPSource. Nothing to do.
	return
}

func (this *MajorCommon) initialize() {
	*this = MajorCommon{
		enclosing: this,

		chart: nil,
		parent: nil,
		numID: NID_Invalid,
		unresRelQualID: "",

		attrs: map[string][]string{},
	}
}

// Unresolved = non-empty `unresRelQualID`
// `unresRelQualID` should be a "qualified" ID, meaning if an object was included,
// the `unresRelQualID` here should prepend the qualified ID of the chart-embedding event object.

func CreateEmptyMajorCommon() *MajorCommon {
	obj := MajorCommon{}
	obj.initialize()
	return &obj
}

//-- vvv Constants vvv

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

const TKValue_Max int = math.MaxInt
const TKValue_Min int = math.MinInt 
const TKValue_Unknown int = -1 

//-- type NumberID

type NumberID uint64
const NID_Max = math.MaxUint64 - 1
const NID_Invalid = math.MaxUint64

//-- type MajorKind

type MajorKind int
const (
	MOK_Event MajorKind = iota
	MOK_Annex
	MOK_Import
	MOK_Unknown
)

//-- type MajorKind: method (stringify)

func (mok MajorKind) Prefix() (string, error) {
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

//-- type MajorKind: method (creation)

func GetMajorKind(prefix string) (MajorKind, error) {
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

/*
//-- type IntervalKind 

// Let's say there is an "interval" (start ~ end) with respect to another
// "interval" (estart ~ eend).
//
//   |-------|==========|--------|
//   ^       ^          ^        ^
// start   estart      eend     end
//
// Inclusive: include equals at the boundary. (start <= estart && eend <= end)
// Extended: include stradding ones. (start <= eend || estart <= end)

type IntervalKind int
const (
	IK_Inclusive IntervalKind = iota
	IK_Extended
)
*/
