package libltc

import (
	"errors"
)

type IMainObj interface {
	Common() *MainObjCommon
	IsUnresolved() bool
}

type MainObjCommon struct {
	kind MainObjKind			// Main object kind

	chart *Chart				// Linked chart
	parent *IMainObj				// Included by... (nil: chart-local)
	numID NumberID 				// Numeric part of ID
	fullID string 				// Full ID (ONLY FOR MOK_Unknown!)

	extraNoteAnnexs []*Annex	// Annexs as extra notes
	attachedAnnexs []*Annex		// Annexs as attachments
	attrs map[string][]string	// Attributes
}

//-- struct MainObjCommon: method

func (this *MainObjCommon) GetKind() MainObjKind {
	return this.kind
}

func (this *MainObjCommon) GetChart() *Chart {
	return this.chart
}

func (this *MainObjCommon) GetLocalID(prefix string) string {
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

func (this *MainObjCommon) GetQualifiedID() (ret string) {
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

func (this *MainObjCommon) GetParent() IMainObj {
	return this.parent
}

func (this *MainObjCommon) SetParent(p IMainObj) {
	this.parent = p
}

func (this *MainObjCommon) GetExtraNotes() (ret []*Note) {
	for _, aobj := range this.extraNoteAnnexs {
		append(ret, &aobj.Note)
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
		append(this.extraNoteAnnexs, aobj)
	}
}

func (this *MainObjCommon) RemoveExtraNoteAnnex(aobj *Annex) {
	for _, elem := range this.extraNoteAnnexs {
		if elem == aobj {
			this.extraNoteAnnexs = append(this.extraNoteAnnexs[:i], this.extraNoteAnnexs[i+1:]...) 
			break
		}
	}
}

func (this *MainObjCommon) GetAttachedAnnexs() []*Annex {
	return this.attachedAnnexs
}

func (this *MainObjCommon) AddAttachedAnnex(aobj *Annex) {
	if !this.HasAttachedAnnex(aobj) {
		append(this.attachedAnnexs, aobj)
	}
}

func (this *MainObjCommon) RemoveAttachedAnnex(aobj *Annex) {
	for _, elem := range this.attachedAnnexs {
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
	append(this.attrs[key], value)
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

func (this *MainObjCommon) getNumberID() NumberID {
	return this.numID
}

func (this *MainObjCommon) IsUnresolved() bool {
	return fullID != ""
}

func CreateUnresolvedMainObj(kind MainObjectKind, id string) *MainObjCommon {
	return &MainObjCommon{
		kind: kind,

		chart: nil,
		parent: nil,
		numID: NID_Invalid,
		fullID: id,

		extraNoteAnnexs: []*Annex{},
		attachedAnnexs: []*Annex{},
		attrs: map[string][]string{},
	}
}
