package libltc

import (
	"strings"
)

type IMainObj interface {
	GetKind() MainObjKind

	GetChart() *Chart
	setChart(c *Chart)

	getNumberID() NumberID 
	setNumberID(nid NumberID)

	getFullID() string
	setFullID(fullid string)
	IsUnresolved() bool

	GetLocalID() string 
	GetQualifiedIDFrom(stopAt IMainObj) string
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
}

type MainObjCommon struct {
	kind MainObjKind			// Main object kind

	chart IChart				// Associated chart
	parent IMainObj				// Included by... (nil: chart-local)
	numID NumberID 				// Numeric part of ID
	fullID string 				// Full ID (ONLY FOR UNRESOLVED!)

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

func (this *MainObjCommon) getFullID() string {
	return this.fullID
}

func (this *MainObjCommon) setFullID(fullid string) {
	this.fullID = fullid
}

func (this *MainObjCommon) IsUnresolved() bool {
	return this.fullID != ""
}

func (this *MainObjCommon) GetLocalID() string {
	if this.kind == MOK_Unknown {
		fullid := this.fullID
		lastidx := strings.LastIndex(fullid, "/")

		if lastidx == -1 {
			return fullid
		} else if len(fullid) < lastidx + 1 {
			return ""
		} else {
			return fullid[lastidx+1:]
		}
	} else {
		return convIDInternalToString(this.kind, this.numID)
	}
}

func (this *MainObjCommon) GetQualifiedIDFrom(stopAt IMainObj) (ret string) {
	if this.kind == MOK_Unknown {
		fullid := this.fullID
		baseid := stopAt.GetQualifiedID()
		return strings.TrimLeft(strings.TrimPrefix(fullid, baseid), "/")
	} else {
		if this.parent != nil && this.parent != stopAt {
			ret = this.parent.GetQualifiedIDFrom(stopAt) + "/"
		}

		ret += this.GetLocalID()
		return
	}
}

func (this *MainObjCommon) GetQualifiedID() string {
	return this.GetQualifiedIDFrom(nil)
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
	for _, elem := range this.extraNoteAnnexs {
		if elem == aobj {
			return true
		}
	}
	return false
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
	for _, elem := range this.attachedAnnexs {
		if elem == aobj {
			return true
		}
	}
	return false
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
		for _, elem := range attrlist {
			if elem == value {
				return true
			}
		}
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

// Unresolved = non-empty `fullid`
// `fullid` should be a "qualified" ID, meaning if an object was included,
// the `fullid` here should prepend the qualified ID of the chart-embedding event object.
func createMainObjCommon(kind MainObjKind, fullid string) MainObjCommon {
	return MainObjCommon{
		kind: kind,

		chart: nil,
		parent: nil,
		numID: NID_Invalid,
		fullID: fullid,

		extraNoteAnnexs: []*Annex{},
		attachedAnnexs: []*Annex{},
		attrs: map[string][]string{},
	}
}

func CreateEmptyMainObjCommon(kind MainObjKind) MainObjCommon {
	return createMainObjCommon(kind, "")
}

func CreateUnresolvedMainObjCommon(kind MainObjKind, fullid string) MainObjCommon {
	return createMainObjCommon(kind, fullid)
}
