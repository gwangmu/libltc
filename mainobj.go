package libltc

import (
	"errors"
)

type LTCMainObjCommon struct {
	kind LTCMainObjKind				// Main object kind

	chart *LTCChart					// Linked chart
	parent LTCMainObj				// Included by... (nil: chart-local)
	numID LTCNumberID 				// Numeric part of ID

	extraNoteAnnexs []*LTCAnnex		// Annexs as extra notes
	attachedAnnexs []*LTCAnnex		// Annexs as attachments
	attrs map[string][]string		// Attributes
}

type LTCMainObj interface {
	GetChart() *LTCChart
	GetLocalID() string
	GetQualifiedID() string
	HasParent() bool

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

//-- struct LTCMainObjCommon: interface LTCMainObj

func (this *LTCMainObjCommon) GetChart() *LTCChart {
	return this.chart
}

func (this *LTCMainObjCommon) GetLocalID(prefix string) string {
	if prefix, err := this.kind.Prefix(); err != nil {
		return prefix + strconv.Itoa(this.numID)
	} else {
		panic(err.Error())
	}
}

func (this *LTCMainObjCommon) GetQualifiedID() (ret string) {
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

func (this *LTCMainObjCommon) HasParent() bool {
	// If this is true, the object shoule be considered read-only.
	return this.parent != nil
}

func (this *LTCMainObjCommon) GetExtraNotes() (ret []*LTCNote) {
	for _, aobj := range this.extraNoteAnnexs {
		append(ret, &aobj.Note)
	}
	return
}

func (this *LTCMainObjCommon) AddExtraNoteAnnex(aobj *LTCAnnex) {
	if !aobj.HasAttr("ExtraNoteOf", this.GetQualifiedID()) {
		aobj.AddAttr("ExtraNoteOf", this.GetQualifiedID())
	}
	append(this.extraNoteAnnexs, aobj)
}

func (this *LTCMainObjCommon) RemoveExtraNoteAnnex(aobj *LTCAnnex) {
	for i, elem := range this.extraNoteAnnexs {
		if elem == aobj {
			this.extraNoteAnnexs = append(this.extraNoteAnnexs[:i], this.extraNoteAnnexs[i+1:]...) 
			break
		}
	}
	aobj.RemoveAttr("ExtraNoteOf", this.GetQualifiedID())
}

func (this *LTCMainObjCommon) GetAttachedAnnexs() []*LTCAnnex {
	return this.attachedAnnexs
}

func (this *LTCMainObjCommon) AddAttachedAnnex(aobj *LTCAnnex) {
	if !aobj.HasAttr("AttachTo", this.GetQualifiedID()) {
		aobj.AddAttr("AttachTo", this.GetQualifiedID())
	}
	append(this.attachedAnnexs, aobj)
}

func (this *LTCMainObjCommon) RemoveAttachedAnnex(aobj *LTCAnnex) {
	for i, elem := range this.attachedAnnexs {
		if elem == aobj {
			this.attachedAnnexs = append(this.attachedAnnexs[:i], this.attachedAnnexs[i+1:]...) 
			break
		}
	}
	aobj.RemoveAttr("AttachTo", this.GetQualifiedID())
}

func (this *LTCMainObjCommon) GetAttrs(key string) []string {
	if attrlist, ok := this.attrs[key]; ok {
		// `attrlist` CANNOT be an empty list. (at least [""])
		return attrlist
	} else {
		return []string{}
	}
}

func (this *LTCMainObjCommon) HasAttr(key string, value string) bool {
	if attrlist, ok := this.attrs[key]; ok {
		for _, elem := range attrlist {
			if elem == value {
				return true
			}
		}
	}
	return false
}

func (this *LTCMainObjCommon) AddAttr(key string, value string) {
	if _, ok := this.attrs[key]; !ok {
		this.attrs[key] = []string{}
	}
	append(this.attrs[key], value)
}

func (this *LTCMainObjCommon) RemoveAttr(key string, value string) {
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

func (this *LTCMainObjCommon) getNumberID() LTCNumberID {
	return this.numID
}
