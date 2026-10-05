package ltc

import (
	"slices"
	"strings"
)

type IMainObj interface {
	GetKind() MainObjKind
	GetCommon() *MainObjCommon
	GetEnclosingObject() IMainObj

	GetChart() *Chart
	setChart(c *Chart)

	getNumberID() NumberID 
	setNumberID(nid NumberID)

	getUnresRelQualID() string
	markUnresolved(unresRelQualID string, parent IMainObj)
	unmarkUnresolved()
	IsUnresolved() bool

	GetIntrinsicID() string 
	GetRelQualifiedID(baseChart *Chart) string
	GetAbsQualifiedID() string

	GetParent() IMainObj
	SetParent(p IMainObj)

	IsImported() bool	// Imported by an import object?
	IsEmbedded() bool	// Embedded in an event object? (incl. subchart objs)
	IsIncluded() bool	// Imported or embedded?
	IsLocallyImported() bool	// Imported within a subchart? 

	GetAttrs(key string) []string
	HasAttr(key string, value string) bool
	AddAttr(key string, value string)
	RemoveAttr(key string, value string)

	convIDRelToAbs(relid string) string
	convIDAbsToRel(absid string) string

	resolveReferenceTo(c IChart)
	unresolveReferenceTo(c IChart)

	initialize()
}

type IMainObjPtr[T any] interface {
	*T
	IMainObj
}

type MainObjCommon struct {
	enclosing IMainObj			// Enclosing object of this common struct

	chart IChart				// Associated chart
	parent IMainObj				// Included by...
	numID NumberID 				// Numeric part of ID
	unresRelQualID string		// **UNRESOLVED** local qualified ID

	attrs map[string][]string	// Attributes

	ExtraNote EPSink[*Annex]	// Pointed by 'ExtraNoteOf' 
	Attached EPSink[*Annex]		// Pointed by 'AttachTo'
}

//-- struct MainObjCommon: method

func (this *MainObjCommon) GetKind() MainObjKind {
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

func (this *MainObjCommon) GetCommon() *MainObjCommon {
	return this
}

func (this *MainObjCommon) GetEnclosingObject() IMainObj {
	return this.enclosing
}

// Relationship Taxonomy:
//  - Chart and main objects: "association"
//  - Main object and main object: "parental"
//  - Subchart and subchart event object: "embedding"

func (this *MainObjCommon) GetChart() *Chart {
	return this.chart.asChart()
}

func (this *MainObjCommon) setChart(c *Chart) {
	this.chart = c
}

func (this *MainObjCommon) getNumberID() NumberID {
	return this.numID
}

func (this *MainObjCommon) setNumberID(nid NumberID) {
	this.numID = nid
}

func (this *MainObjCommon) getUnresRelQualID() string {
	return this.unresRelQualID
}

func (this *MainObjCommon) markUnresolved(unresRelQualID string, parent IMainObj) {
	this.unresRelQualID = unresRelQualID
	this.parent = parent
}

func (this *MainObjCommon) unmarkUnresolved() {
	this.unresRelQualID = ""
}

func (this MainObjCommon) IsUnresolved() bool {
	return this.unresRelQualID != ""
}

func (this *MainObjCommon) GetIntrinsicID() string {
	if this.GetKind() == MOK_Unknown {
		unresRelQualID := this.unresRelQualID
		lastidx := strings.LastIndex(unresRelQualID, "/")

		if lastidx == -1 {
			return unresRelQualID
		} else if len(unresRelQualID) < lastidx + 1 {
			return ""
		} else {
			return unresRelQualID[lastidx+1:]
		}
	} else {
		return convIDInternalToString(this.GetKind(), this.numID)
	}
}

func (this *MainObjCommon) GetRelQualifiedID(baseChart *Chart) (ret string) {
	if this.IsUnresolved() {
		if this.parent != nil {
			absQualID := this.GetAbsQualifiedID()
			return strings.TrimPrefix(this.parent.GetAbsQualifiedID(), absQualID)
		} else {
			// THIS SHOULDN'T HAPPEN.
			return this.unresRelQualID
		}
	} else {
		if this.parent != nil && (this.GetChart() != baseChart || 
			this.parent.GetChart() == baseChart) {
			ret = this.parent.GetRelQualifiedID(baseChart) + "/"
		}
		ret += this.GetIntrinsicID()
		return
	}
}

func (this *MainObjCommon) GetAbsQualifiedID() (ret string) {
	if this.IsUnresolved() {
		unresRelQualID := this.unresRelQualID
		prefixid := ""
		if this.parent != nil {
			prefixid = this.parent.GetAbsQualifiedID() + "/"
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

func (this *MainObjCommon) GetParent() IMainObj {
	return this.parent
}

func (this *MainObjCommon) SetParent(p IMainObj) {
	this.parent = p
}

func (this *MainObjCommon) IsImported() bool {
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

func (this *MainObjCommon) IsEmbedded() bool {
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

func (this *MainObjCommon) IsIncluded() bool {
	return this.IsImported() || this.IsEmbedded()
}

func (this *MainObjCommon) IsLocallyImported() bool {
	return this.parent != nil && this.parent.GetKind() == MOK_Import 
}

func (this *MainObjCommon) GetAttrs(key string) []string {
	if attrlist, ok := this.attrs[key]; ok {
		// `attrlist` CANNOT be an empty list. (at least [""])
		return attrlist
	} else {
		return []string{}
	}
}

func (this *MainObjCommon) HasAttr(key string, value string) bool {
	if attrlist, ok := this.attrs[key]; ok {
		return slices.Contains(attrlist, value)
	}
	return false
}
	
// `{Add,Remove}Attr` are passive methods; they don't update other
// objects that may be referenced by the added/removed attribute.

func (this *MainObjCommon) AddAttr(key string, value string) {
	if _, ok := this.attrs[key]; !ok {
		this.attrs[key] = []string{}
	}
	this.attrs[key] = append(this.attrs[key], value)
}

func (this *MainObjCommon) RemoveAttr(key string, value string) {
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

func (this *MainObjCommon) convIDRelToAbs(relid string) string {
	if this.parent != nil {
		return this.parent.GetAbsQualifiedID() + "/" + relid
	} else {
		return relid
	}
}

func (this *MainObjCommon) convIDAbsToRel(absid string) string {
	if this.parent != nil {
		return strings.TrimPrefix(absid, this.parent.GetAbsQualifiedID() + "/")
	} else {
		return absid
	}
}

func (this *MainObjCommon) resolveReferenceTo(c IChart) {
	// No EPSource. Nothing to do.
	return
}

func (this *MainObjCommon) unresolveReferenceTo(c IChart) {
	// No EPSource. Nothing to do.
	return
}

func (this *MainObjCommon) initialize() {
	*this = MainObjCommon{
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

func CreateEmptyMainObjCommon() *MainObjCommon {
	obj := MainObjCommon{}
	obj.initialize()
	return &obj
}
