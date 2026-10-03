package libltc 

import (
	"slices"
)

type ILinkElem interface {
	GetQualifiedID() string
	IsUnresolved() bool

	markUnresolved(unresRelQualID string)
	unmarkUnresolved()

	convIDRelToAbs(relid string) string
	convIDAbsToRel(absid string) string

	initialize()
}

type ILinkElemPtr[T any] interface {
	*T
	ILinkElem
}

// Endpoints are the "endpoints" of links. Each main object should include it
// with an appropriate field name.

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

func (this *Endpoint[_, TPtr]) swap(oldobj TPtr, newobj TPtr) {
	if i := slices.Index(this.objs, oldobj); i != -1 {
		this.objs[i] = newobj
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

// Links are dummy structs that provide the appropriate endpoints in the
// source/sink objects. It should define the following methods.

type ILinkBase[SrcT any, SinkT any, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]] interface {
	GetSourceEndpoint(SrcTPtr) *Endpoint[SinkT, SinkTPtr]
	GetSinkEndpoint(SinkTPtr) *Endpoint[SrcT, SrcTPtr]
}

func CreateLink[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT any, SinkT any, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]](
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

func ResolveLink[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT any, SinkT any, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]](
		osrc SrcTPtr, osink SinkTPtr) bool {
	var l LinkT
	epSrc := l.GetSourceEndpoint(osrc)
	epSink := l.GetSinkEndpoint(osink)

	if osink.IsUnresolved() {
		return false
	}

    absQualID := osink.GetQualifiedID() 
	for _, elem := range epSrc.Get() {
		if elem.GetQualifiedID() == absQualID {
			if elem.IsUnresolved() {
				epSrc.swap(elem, osink)
				epSink.add(osrc)
			}
			return true
		}
	}
    return false
}

func UnresolveLink[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT any, SinkT any, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]](
		osrc SrcTPtr, osink SinkTPtr) {
	var l LinkT
	epSrc := l.GetSourceEndpoint(osrc)
	epSink := l.GetSinkEndpoint(osink)

	absQualID := osink.GetQualifiedID()
	epSink.remove(osrc)

	for _, elem := range epSrc.Get() {
		if elem == osink {
			if !elem.IsUnresolved() {
				elem.initialize()
				elem.markUnresolved(osrc.convIDAbsToRel(absQualID))
			}
			return
		}
	}
}

func RemoveLink[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT any, SinkT any, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]](
		osrc SrcTPtr, absQualID string) {
	var l LinkT
	epSrc := l.GetSourceEndpoint(osrc)

	for _, elem := range epSrc.Get() {
		if elem.GetQualifiedID() == absQualID {
			if !elem.IsUnresolved() {
				epSink := l.GetSinkEndpoint(elem)
				epSink.remove(osrc)
			}
			epSrc.remove(elem)
			return
		}
	}
}

// Here are some specialized links.

type AttachTo struct {}
func (this AttachTo) GetSourceEndpoint(aobj *Annex) *Endpoint[MainObjCommon, *MainObjCommon] {
	return &aobj.AttachTo
}
func (this AttachTo) GetSinkEndpoint(obj *MainObjCommon) *Endpoint[Annex, *Annex] {
	return &obj.Attached
}

type ExtraNoteOf struct {}
func (this ExtraNoteOf) GetSourceEndpoint(aobj *Annex) *Endpoint[MainObjCommon, *MainObjCommon] {
	return &aobj.ExtraNoteOf
}
func (this ExtraNoteOf) GetSinkEndpoint(obj *MainObjCommon) *Endpoint[Annex, *Annex] {
	return &obj.ExtraNote
}

type ContinuedFrom struct {}
func (this ContinuedFrom) GetSourceEndpoint(aobj *Event) *Endpoint[Event, *Event] {
	return &aobj.ContinuedFrom
}
func (this ContinuedFrom) GetSinkEndpoint(obj *Event) *Endpoint[Event, *Event] {
	return &obj.ContinuedTo
}
