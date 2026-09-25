package libltc

import (
	"errors"
)

type ltcMainObjCommon struct {
	chart *LTCChart
	parent LTCMainObj
	numID LTCNumberID 

	extraNoteAnnexs []*LTCAnnex
	attachedAnnexs []*LTCAnnex
	attrs map[string][]string
}

type LTCMainObj interface {
	GetChart() *LTCChart
	GetLocalID() string
	GetQualifiedID() string

	GetExtraNotes() []*LTCNote
	AddExtraNoteAnnex(*LTCAnnex)
	RemoveExtraNoteAnnex(*LTCAnnex)

	GetAttachedAnnexs() []*LTCAnnex
	AddAttachedAnnex(*LTCAnnex)
	RemoveAttachedAnnex(*LTCAnnex)

	// `{Add,Remove}Attr` are passive methods; they don't update other
	// objects that may be referenced by the added/removed attribute.
	GetAttrs(key string) []string 
	HasAttr(key string, value string) bool
	AddAttr(key string, value string)
	RemoveAttr(key string, value string)

	getNumberID() LTCNumberID
}

//-- struct ltcMainObjCommon: interface LTCMainObj

func (this *ltcMainObjCommon) GetChart() *LTCChart {
	return this.chart
}

func (this *ltcMainObjCommon) GetLocalID(prefix string) string {
	return prefix + strconv.Itoa(this.numID)
}

func (this *ltcMainObjCommon) GetQualifiedID(prefix string) (ret string) {
	if this.parent != nil {
		ret = this.parent.GetQualifiedID() + '/'
	}

	ret += this.GetLocalID(prefix)
	return
}

func (this *ltcMainObjCommon) GetExtraNotes() (ret []*LTCNote) {
	for _, aobj := range this.extraNoteAnnexs {
		append(ret, &aobj.Note)
	}
	return
}

func (this *ltcMainObjCommon) AddExtraNoteAnnex(aobj *LTCAnnex) {
	append(this.extraNoteAnnexs, aobj)
}

func (this *ltcMainObjCommon) RemoveExtraNoteAnnex(aobj *LTCAnnex) {
	for i, elem := range this.extraNoteAnnexs {
		if elem == aobj {
			this.extraNoteAnnexs = append(this.extraNoteAnnexs[:i], this.extraNoteAnnexs[i+1:]...) 
			return
		}
	}
}

func (this *ltcMainObjCommon) GetAttachedAnnexs() []*LTCAnnex {
	return this.attachedAnnexs
}

func (this *ltcMainObjCommon) AddAttachedAnnex(aobj *LTCAnnex) {
	append(this.attachedAnnexs, aobj)
}

func (this *ltcMainObjCommon) RemoveAttachedAnnex(aobj *LTCAnnex) {
	for i, elem := range this.attachedAnnexs {
		if elem == aobj {
			this.attachedAnnexs = append(this.attachedAnnexs[:i], this.attachedAnnexs[i+1:]...) 
			return
		}
	}
}

func (this *ltcMainObjCommon) GetAttrs(key string) []string {
	if attrlist, ok := this.attrs[key]; ok {
		// `attrlist` CANNOT be an empty list. (at least [""])
		return attrlist
	} else {
		return []string{}
	}
}

func (this *ltcMainObjCommon) HasAttr(key string, value string) bool {
	if attrlist, ok := this.attrs[key]; ok {
		for _, elem := range attrlist {
			if elem == value {
				return true
			}
		}
	}
	return false
}

func (this *ltcMainObjCommon) AddAttr(key string, value string) {
	if _, ok := this.attrs[key]; !ok {
		this.attrs[key] = []string{}
	}
	append(this.attrs[key], value)
}

func (this *ltcMainObjCommon) RemoveAttr(key string, value string) {
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

func (this *ltcMainObjCommon) getNumberID() LTCNumberID {
	return this.numID
}
