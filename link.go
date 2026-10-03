package libltc 

import (
	"slices"
)

//-- struct Endpoint[T]

type ILinkElem interface {
	GetQualifiedID() string
	IsUnresolved() bool

	markUnresolved(unresRelQualID string)
	unmarkUnresolved()
	initialize()
}

type ILinkElemPtr[T any] interface {
	*T
	ILinkElem
}

type IEndpoint[T any, TPtr ILinkElemPtr[T]] interface {
	Get() []TPtr
	Has(TPtr) bool
	add(TPtr)
	remove(TPtr)
}

type Endpoint[T any, TPtr ILinkElemPtr[T]] struct {
	objs []TPtr
}

func (this *Endpoint[_, TPtr]) Get() []TPtr {
	return this.objs
}

func (this *Endpoint[_, TPtr]) Has(obj TPtr) bool {
	return slices.Contains(this.objs, obj)
}

func (this *Endpoint[_, TPtr]) add(obj TPtr) {
	if !this.Has(obj) {
		this.objs = append(this.objs, obj)
	}
}

func (this *Endpoint[_, TPtr]) remove(obj TPtr) {
	for i, elem := range this.objs {
		if elem == obj {
			this.objs = append(this.objs[:i], this.objs[i+1:]...) 
			return
		}
	}
}

//type ILinkBase[SrcT ILinkElem, SinkT ILinkElem] interface {
//	GetSource(*SrcT) *Endpoint[SinkT]
//	GetSink(*SinkT) *Endpoint[SrcT]
//}
//type AttachTo struct {}
//func (this AttachTo) GetSource(aobj *Annex) *Endpoint[IMainObj] {
//	return aobj.epAttachTo
//}
//func (this AttachTo) GetSink(obj IMainObj) *Endpoint[Annex] {
// 	return obj.epAttached
//}
//
//Create[AttachTo](aobj, obj)
//Resolve[AttachTo](aobj, "e000/e001", obj)
//Unresolve[AttachTo](aobj, "e000/e001")
//Remove[AttachTo](aobj, "e000/e001")

type ILinkBase[SrcT ILinkElem, SinkT ILinkElem, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]] interface {
	GetSourceEndpoint(SrcTPtr) *Endpoint[SinkT, SinkTPtr]
	GetSinkEndpoint(SinkTPtr) *Endpoint[SrcT, SrcTPtr]
}

func Create[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT ILinkElem, SinkT ILinkElem, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]](
		osrc SrcTPtr, osink SinkTPtr) {
	var l LinkT
	epSrc := l.GetSourceEndpoint(osrc)
	epSink := l.GetSinkEndpoint(osink)

    if !epSrc.Has(osink) {
        epSrc.add(osink)
        if !osink.IsUnresolved() {
            epSink.add(osrc)
        }
    }
}

func Resolve[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT ILinkElem, SinkT ILinkElem, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]](
		osrc SrcT, absQualID string, osink SinkT) {
	// TODO
	panic("Unimplmemented")
}

func Unresolve[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT ILinkElem, SinkT ILinkElem, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]](
		osrc SrcT, absQualID string) {
	// TODO
	panic("Unimplmemented")
}

func Remove[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT ILinkElem, SinkT ILinkElem, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]](
		osrc SrcT, absQualID string) {
	// TODO
	panic("Unimplmemented")
}

/*
type ILinkElemPtr[T any] interface {
	*T
	ILinkElem
}

func Create[ElemT ILinkElem, ElemTPtr ILinkElemPtr[ILinkElem]](fromObj ElemTPtr, from *Endpoint, toObj ElemTPtr, to *Endpoint) {
    if !fromobjs.Has(toObj) {
        from.objs = append(from.objs, toObj)
        if !toObj.IsUnresolved() {
            toObj.add(fromObj)
        }
    }
}

func Create[ElemT ILinkElem, ElemTPtr ILinkElemPtr[ILinkElem]](fromObj ElemTPtr, from *Endpoint, toObj ElemTPtr, to *Endpoint) {
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

func (this *Source[ElemT, SinkT, ElemTPtr]) unresolve(absQualID string) {
	for i, elem := range this.objs {
		if elem.GetQualifiedID() == absQualID {
			if !elem.IsUnresolved() {
				elem.markUnresolved(this.convIDAbsToRel(absQualID))
			}
			return
		}
	}
}

func (this *Source[ElemT, SinkT, ElemTPtr]) remove(absQualID string) {
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
*/
