package link

import (
	"slices"
)

type IElem interface {
	GetQualifiedID() string
	IsUnresolved() bool
}

//-- struct Sink[ElemT]

type ISink[ElemT IElem] interface {
	Get() []*ElemT
	Has(*ElemT) bool
	add(*ElemT)
	remove(*ElemT)
}

type Sink[ElemT IElem] struct {
	objs []*ElemT
}

func (this *Sink[ElemT]) Get() []*ElemT {
	return this.objs
}

func (this *Sink[ElemT]) Has(obj *ElemT) bool {
	return slices.Contains(this.objs, obj)
}

func (this *Sink[ElemT]) add(obj *ElemT) {
	if !this.Has(obj) {
		this.objs = append(this.objs, obj)
	}
}

func (this *Sink[ElemT]) remove(obj *ElemT) {
	for i, elem := range this.objs {
		if elem == obj {
			this.objs = append(this.objs[:i], this.objs[i+1:]...) 
			return
		}
	}
}

//-- struct Source[T]

type Source[ElemT IElem, SinkT ISink[ElemT]] struct {
	objs []*ElemT
}

func (this *Source[ElemT, _]) Get() []*ElemT {
	return this.objs
}

func (this *Source[ElemT, _]) Has(obj *ElemT) bool {
	return slices.Contains(this.objs, obj)
}

func (this *Source[ElemT, SinkT]) add(obj *ElemT) {
    if !this.Has(obj) {
        this.objs = append(this.objs, obj)
        if !obj.IsUnresolved() {
            obj.SinkT.add(this)
        }
    }
}

func (this *Source[ElemT, SinkT]) resolve(newobj *ElemT) bool {
    absQualID := newobj.GetQualifiedID() 
	for i, elem := range this.objs {
		if elem.GetQualifiedID() == absQualID {
			if elem.IsUnresolved() {
				this.objs[i] = newobj
				newobj.SinkT.add(this)
			}
			return true
		}
	}
    return false
}

func (this *Source[ElemT, SinkT]) unresolve(absQualID string) {
	for i, elem := range this.objs {
		if elem.GetQualifiedID() == absQualID {
			if !elem.IsUnresolved() {
				elem.markUnresolved(this.convIDAbsToRel(absQualID))
			}
			return
		}
	}
}

func (this *Source[ElemT, SinkT]) remove(absQualID string) {
	for i, elem := range this.objs {
		if elem.GetQualifiedID() == absQualID {
			if !elem.IsUnresolved() {
				elem.SinkT.remove(this)
			}
			this.objs = append(this.objs[:i], this.objs[i+1:]...) 
			return
		}
	}
}
