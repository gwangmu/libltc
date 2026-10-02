package libltc

import (
	"slices"
	"strings"
)

type IMainObj interface {
	GetKind() MainObjKind

	GetChart() *Chart
	setChart(c *Chart)

	getNumberID() NumberID 
	setNumberID(nid NumberID)

	getUnresRelQualID() string
	setUnresRelQualID(unresRelQualID string)
	IsUnresolved() bool

	GetLocalID() string 
	GetQualifiedID() string

	GetParent() IMainObj
	SetParent(p IMainObj)

	GetExtraNotes() (ret []*Note)
	HasExtraNoteAnnex(aobj *Annex) bool
	AddExtraNoteAnnex(aobj *Annex)
	RemoveExtraNoteAnnex(aobj *Annex)

	GetAttachedAnnexs() []*Annex
	HasAttachedAnnex(aobj *Annex) bool
	AddAttachedAnnex(aobj *Annex)
	RemoveAttachedAnnex(aobj *Annex)

	GetAttrs(key string) []string
	HasAttr(key string, value string) bool
	AddAttr(key string, value string)
	RemoveAttr(key string, value string)

	convIDRelToAbs(relid string) string
}

type MainObjCommon struct {
	kind MainObjKind			// Main object kind

	chart IChart				// Associated chart
	parent IMainObj				// Included by... (nil: chart-local)
	numID NumberID 				// Numeric part of ID
	unresRelQualID string 				// Full ID (ONLY FOR UNRESOLVED! Relative to parent)

	extraNoteAnnexs []*Annex	// Annexs as extra notes
	attachedAnnexs []*Annex		// Annexs as attachments
	attrs map[string][]string	// Attributes
}

//-- struct MainObjCommon: method

func (this *MainObjCommon) GetKind() MainObjKind {
	return this.kind
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

func (this *MainObjCommon) setUnresRelQualID(unresRelQualID string) {
	this.unresRelQualID = unresRelQualID
}

func (this MainObjCommon) IsUnresolved() bool {
	return this.unresRelQualID != ""
}

func (this *MainObjCommon) GetLocalID() string {
	if this.kind == MOK_Unknown {
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
		return convIDInternalToString(this.kind, this.numID)
	}
}

func (this *MainObjCommon) GetQualifiedID() (ret string) {
	if this.IsUnresolved() {
		unresRelQualID := this.unresRelQualID
		prefixid := ""
		if this.parent != nil {
			prefixid = this.parent.GetQualifiedID() + "/"
		}
		return prefixid + unresRelQualID
	} else {
		if this.parent != nil {
			ret = this.parent.GetQualifiedID() + "/"
		}

		ret += this.GetLocalID()
		return
	}
}

func (this *MainObjCommon) GetParent() IMainObj {
	return this.parent
}

func (this *MainObjCommon) SetParent(p IMainObj) {
	this.parent = p
}

func (this *MainObjCommon) GetExtraNotes() (ret []*Note) {
	for _, aobj := range this.extraNoteAnnexs {
		ret = append(ret, &aobj.Note)
	}
	return
}

func (this *MainObjCommon) HasExtraNoteAnnex(aobj *Annex) bool {
	return slices.Contains(this.extraNoteAnnexs, aobj)
}

func (this *MainObjCommon) AddExtraNoteAnnex(aobj *Annex) {
	if !this.HasExtraNoteAnnex(aobj) {
		this.extraNoteAnnexs = append(this.extraNoteAnnexs, aobj)
	}
}

func (this *MainObjCommon) RemoveExtraNoteAnnex(aobj *Annex) {
	for i, elem := range this.extraNoteAnnexs {
		if elem == aobj {
			this.extraNoteAnnexs = append(this.extraNoteAnnexs[:i], this.extraNoteAnnexs[i+1:]...) 
			break
		}
	}
}

func (this *MainObjCommon) GetAttachedAnnexs() []*Annex {
	return this.attachedAnnexs
}

func (this *MainObjCommon) HasAttachedAnnex(aobj *Annex) bool {
	return slices.Contains(this.attachedAnnexs, aobj)
}

func (this *MainObjCommon) AddAttachedAnnex(aobj *Annex) {
	if !this.HasAttachedAnnex(aobj) {
		this.attachedAnnexs = append(this.attachedAnnexs, aobj)
	}
}

func (this *MainObjCommon) RemoveAttachedAnnex(aobj *Annex) {
	for i, elem := range this.attachedAnnexs {
		if elem == aobj {
			this.attachedAnnexs = append(this.attachedAnnexs[:i], this.attachedAnnexs[i+1:]...) 
			break
		}
	}
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
		return this.parent.GetQualifiedID() + "/" + relid
	} else {
		return relid
	}
}

func (this *MainObjCommon) convIDAbsToRel(absid string) string {
	if this.parent != nil {
		return strings.TrimPrefix(absid, this.parent.GetQualifiedID() + "/")
	} else {
		return absid
	}
}

// Unresolved = non-empty `unresRelQualID`
// `unresRelQualID` should be a "qualified" ID, meaning if an object was included,
// the `unresRelQualID` here should prepend the qualified ID of the chart-embedding event object.
func createMainObjCommon(kind MainObjKind, unresRelQualID string) *MainObjCommon {
	return &MainObjCommon{
		kind: kind,

		chart: nil,
		parent: nil,
		numID: NID_Invalid,
		unresRelQualID: unresRelQualID,

		extraNoteAnnexs: []*Annex{},
		attachedAnnexs: []*Annex{},
		attrs: map[string][]string{},
	}
}

func CreateEmptyMainObjCommon(kind MainObjKind) *MainObjCommon {
	return createMainObjCommon(kind, "")
}

func CreateUnresolvedMainObjCommon(kind MainObjKind, relQualID string) *MainObjCommon {
	return createMainObjCommon(kind, relQualID)
}
