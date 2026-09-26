package libltc

import (
	"errors"
)

type MainObj interface {
	Common() *MainObjCommon
	IsUnresolved() bool
}

type MainObjCommon struct {
	kind MainObjKind			// Main object kind

	chart *Chart				// Linked chart
	parent *MainObj				// Included by... (nil: chart-local)
	numID NumberID 				// Numeric part of ID
	fullID string 				// Full ID (ONLY FOR MOK_Unknown!)

	extraNoteAnnexs []*Annex	// Annexs as extra notes
	attachedAnnexs []*Annex		// Annexs as attachments
	attrs map[string][]string	// Attributes
}

//-- struct MainObj: method

func (this *MainObj) GetKind() MainObjKind {
	return this.kind
}

func (this *MainObj) GetChart() *Chart {
	return this.chart
}

func (this *MainObj) GetLocalID(prefix string) string {
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

func (this *MainObj) GetQualifiedID() (ret string) {
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

func (this *MainObj) GetParent() *MainObj {
	return this.parent
}

func (this *MainObj) SetParent(p *MainObj) {
	this.parent = p
}

func (this *MainObj) GetExtraNotes() (ret []*Note) {
	for _, aobj := range this.extraNoteAnnexs {
		append(ret, &aobj.Note)
	}
	return
}

func (this *MainObj) HasExtraNoteAnnex(aobj *Annex) bool {
	for _, elem := range this.extraNoteAnnexs {
		if elem == aobj {
			return true
		}
	}
	return false
}

func (this *MainObj) AddExtraNoteAnnex(aobj *Annex) {
	if !this.HasExtraNoteAnnex(aobj) {
		append(this.extraNoteAnnexs, aobj)
	}
}

func (this *MainObj) RemoveExtraNoteAnnex(aobj *Annex) {
	for _, elem := range this.extraNoteAnnexs {
		if elem == aobj {
			this.extraNoteAnnexs = append(this.extraNoteAnnexs[:i], this.extraNoteAnnexs[i+1:]...) 
			break
		}
	}
}

func (this *MainObj) GetAttachedAnnexs() []*Annex {
	return this.attachedAnnexs
}

func (this *MainObj) AddAttachedAnnex(aobj *Annex) {
	if !this.HasAttachedAnnex(aobj) {
		append(this.attachedAnnexs, aobj)
	}
}

func (this *MainObj) RemoveAttachedAnnex(aobj *Annex) {
	for _, elem := range this.attachedAnnexs {
		if elem == aobj {
			this.attachedAnnexs = append(this.attachedAnnexs[:i], this.attachedAnnexs[i+1:]...) 
			break
		}
	}
}

func (this *MainObj) GetAttrs(key string) []string {
	if attrlist, ok := this.attrs[key]; ok {
		// `attrlist` CANNOT be an empty list. (at least [""])
		return attrlist
	} else {
		return []string{}
	}
}

func (this *MainObj) HasAttr(key string, value string) bool {
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

func (this *MainObj) AddAttr(key string, value string) {
	if _, ok := this.attrs[key]; !ok {
		this.attrs[key] = []string{}
	}
	append(this.attrs[key], value)
}

func (this *MainObj) RemoveAttr(key string, value string) {
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

func (this *MainObj) getNumberID() NumberID {
	return this.numID
}

func (this *MainObj) IsUnresolved() bool {
	return fullID != ""
}

func CreateUnresolvedMainObj(kind MainObjectKind, id string) *MainObj {
	return &MainObj{
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
