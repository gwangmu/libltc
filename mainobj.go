package libltc

import (
	"errors"
)

type LTCMainObj interface {
	Common() *LTCMainObjCommon
	IsUnresolved() bool
}

type LTCMainObjCommon struct {
	kind LTCMainObjKind				// Main object kind

	chart *LTCChart					// Linked chart
	parent *LTCMainObj				// Included by... (nil: chart-local)
	numID LTCNumberID 				// Numeric part of ID
	fullID string 					// Full ID (ONLY FOR MOK_Unknown!)

	extraNoteAnnexs []*LTCAnnex		// Annexs as extra notes
	attachedAnnexs []*LTCAnnex		// Annexs as attachments
	attrs map[string][]string		// Attributes
}

//-- struct LTCMainObj: method

func (this *LTCMainObj) GetKind() LTCMainObjKind {
	return this.kind
}

func (this *LTCMainObj) GetChart() *LTCChart {
	return this.chart
}

func (this *LTCMainObj) GetLocalID(prefix string) string {
	if this.kind == MOK_Unknown {
		return this.fullID
	} else {
		if prefix, err := this.kind.Prefix(); err != nil {
			return prefix + strconv.Itoa(this.numID)
		} else {
			panic(err.Error())
		}
	}
}

func (this *LTCMainObj) GetQualifiedID() (ret string) {
	if this.kind == MOK_Unknown {
		return this.fullID
	} else {
		if prefix, err := this.kind.Prefix(); err != nil {
			if this.parent != nil {
				ret = this.parent.GetQualifiedID() + '/'
			}

			ret += this.GetLocalID(prefix)
			return
		} else {
			panic(err.Error())
		}
	}
}

func (this *LTCMainObj) GetParent() *LTCMainObj {
	return this.parent
}

func (this *LTCMainObj) SetParent(p *LTCMainObj) {
	this.parent = p
}

func (this *LTCMainObj) GetExtraNotes() (ret []*LTCNote) {
	for _, aobj := range this.extraNoteAnnexs {
		append(ret, &aobj.Note)
	}
	return
}

func (this *LTCMainObj) HasExtraNoteAnnex(aobj *LTCAnnex) bool {
	for _, elem := range this.extraNoteAnnexs {
		if elem == aobj {
			return true
		}
	}
	return false
}

func (this *LTCMainObj) AddExtraNoteAnnex(aobj *LTCAnnex) {
	if !this.HasExtraNoteAnnex(aobj) {
		append(this.extraNoteAnnexs, aobj)
	}
}

func (this *LTCMainObj) RemoveExtraNoteAnnex(aobj *LTCAnnex) {
	for _, elem := range this.extraNoteAnnexs {
		if elem == aobj {
			this.extraNoteAnnexs = append(this.extraNoteAnnexs[:i], this.extraNoteAnnexs[i+1:]...) 
			break
		}
	}
}

func (this *LTCMainObj) GetAttachedAnnexs() []*LTCAnnex {
	return this.attachedAnnexs
}

func (this *LTCMainObj) AddAttachedAnnex(aobj *LTCAnnex) {
	if !this.HasAttachedAnnex(aobj) {
		append(this.attachedAnnexs, aobj)
	}
}

func (this *LTCMainObj) RemoveAttachedAnnex(aobj *LTCAnnex) {
	for _, elem := range this.attachedAnnexs {
		if elem == aobj {
			this.attachedAnnexs = append(this.attachedAnnexs[:i], this.attachedAnnexs[i+1:]...) 
			break
		}
	}
}

func (this *LTCMainObj) GetAttrs(key string) []string {
	if attrlist, ok := this.attrs[key]; ok {
		// `attrlist` CANNOT be an empty list. (at least [""])
		return attrlist
	} else {
		return []string{}
	}
}

func (this *LTCMainObj) HasAttr(key string, value string) bool {
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

func (this *LTCMainObj) AddAttr(key string, value string) {
	if _, ok := this.attrs[key]; !ok {
		this.attrs[key] = []string{}
	}
	append(this.attrs[key], value)
}

func (this *LTCMainObj) RemoveAttr(key string, value string) {
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

func (this *LTCMainObj) getNumberID() LTCNumberID {
	return this.numID
}

func (this *LTCMainObj) IsUnresolved() bool {
	return fullID != ""
}

func CreateUnresolvedMainObj(kind LTCMainObjectKind, id string) *LTCMainObj {
	return &LTCMainObj{
		kind: kind,

		chart: nil,
		parent: nil,
		numID: NID_Invalid,
		fullID: id,

		extraNoteAnnexs: []*LTCAnnex{},
		attachedAnnexs: []*LTCAnnex{},
		attrs: map[string][]string{},
	}
}
