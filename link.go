package libltc 

import (
	"slices"
)

// Endpoints are the "endpoints" of links. Each main object should include it
// with an appropriate field name.

type Endpoint[T comparable] struct {
	objs []T
}

type EPSource[T comparable] struct {
	Endpoint[T]
}

type EPSink[T comparable] struct {
	Endpoint[T]
}

func (this *Endpoint[T]) Get() []T {
	return this.objs
}

func (this *Endpoint[T]) Has(obj T) bool {
	return slices.Contains(this.objs, obj)
}

func (this *Endpoint[T]) add(obj T) {
	if !this.Has(obj) {
		this.objs = append(this.objs, obj)
	}
}

func (this *Endpoint[T]) swap(oldobj T, newobj T) {
	if i := slices.Index(this.objs, oldobj); i != -1 {
		this.objs[i] = newobj
	}
}

func (this *Endpoint[T]) remove(obj T) {
	for i, elem := range this.objs {
		if elem == obj {
			this.objs = append(this.objs[:i], this.objs[i+1:]...) 
			return
		}
	}
}

// Links are dummy structs that provide the appropriate endpoints in the
// source/sink objects. It should define the following methods.

type ILinkElem interface {
	GetFullQualifiedID() string
	IsUnresolved() bool

	markUnresolved(unresRelQualID string)
	unmarkUnresolved()

	convIDLocalToFull(relid string) string
	convIDFullToLocal(absid string) string

	initialize()
}

type ILinkElemPtr[T any] interface {
	*T
	ILinkElem
}

type ILinkBase[SrcT any, SinkT any, SrcTPtr ILinkElemPtr[SrcT], SinkTPtr ILinkElemPtr[SinkT]] interface {
	GetSourceEndpoint(SrcTPtr) *EPSource[SinkTPtr]
	GetSinkEndpoint(SinkTPtr) *EPSink[SrcTPtr]
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

    absQualID := osink.GetFullQualifiedID() 
	for _, elem := range epSrc.Get() {
		if elem.GetFullQualifiedID() == absQualID {
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

	absQualID := osink.GetFullQualifiedID()
	epSink.remove(osrc)

	for _, elem := range epSrc.Get() {
		if elem == osink {
			if !elem.IsUnresolved() {
				elem.initialize()
				elem.markUnresolved(osrc.convIDFullToLocal(absQualID))
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
		if elem.GetFullQualifiedID() == absQualID {
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
func (this AttachTo) GetSourceEndpoint(aobj *Annex) *EPSource[*MainObjCommon] {
	return &aobj.AttachTo
}
func (this AttachTo) GetSinkEndpoint(obj *MainObjCommon) *EPSink[*Annex] {
	return &obj.Attached
}

type ExtraNoteOf struct {}
func (this ExtraNoteOf) GetSourceEndpoint(aobj *Annex) *EPSource[*MainObjCommon] {
	return &aobj.ExtraNoteOf
}
func (this ExtraNoteOf) GetSinkEndpoint(obj *MainObjCommon) *EPSink[*Annex] {
	return &obj.ExtraNote
}

type ContinuedFrom struct {}
func (this ContinuedFrom) GetSourceEndpoint(aobj *Event) *EPSource[*Event] {
	return &aobj.ContinuedFrom
}
func (this ContinuedFrom) GetSinkEndpoint(obj *Event) *EPSink[*Event] {
	return &obj.ContinuedTo
}
