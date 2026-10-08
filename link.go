package ltc 

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

func (this *Endpoint[T]) find(pred func(T) bool) int {
	return slices.IndexFunc(this.objs, pred)
}

func (this *Endpoint[T]) add(obj T) {
	this.objs = append(this.objs, obj)
}

func (this *Endpoint[T]) swap(idx int, newobj T) {
	this.objs[idx] = newobj
}

func (this *Endpoint[T]) remove(idx int) {
	this.objs = append(this.objs[:idx], this.objs[idx+1:]...) 
}

// Links are dummy structs that provide the appropriate endpoints in the
// source/sink objects. It should define the following methods.

type ILinkBase[SrcT any, SinkT any, SrcTPtr IMajorPtr[SrcT], SinkTPtr IMajorPtr[SinkT]] interface {
	GetSourceEndpoint(SrcTPtr) *EPSource[SinkTPtr]
	GetSinkEndpoint(SinkTPtr) *EPSink[SrcTPtr]
}

func ReserveLink[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT any, SinkT any, SrcTPtr IMajorPtr[SrcT], SinkTPtr IMajorPtr[SinkT]](
		osrc SrcTPtr, relQualID string) {
	var l LinkT
	epSrc := l.GetSourceEndpoint(osrc)

	absQualID := osrc.convIDRelToAbs(relQualID)
	idx := epSrc.find(func (elem SinkTPtr) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	// If the same ID was not reserved, add unresolved.
	if idx == -1 {
		var unresOsink SinkTPtr = new(SinkT)
		unresOsink.initialize()
		unresOsink.markUnresolved(relQualID, osrc)
		CreateLink[LinkT](osrc, unresOsink)
	}
}

func CreateLink[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT any, SinkT any, SrcTPtr IMajorPtr[SrcT], SinkTPtr IMajorPtr[SinkT]](
		osrc SrcTPtr, osink SinkTPtr) bool {
	var l LinkT
	epSrc := l.GetSourceEndpoint(osrc)
	epSink := l.GetSinkEndpoint(osink)

    absQualID := osink.GetAbsQualifiedID() 
	idx := epSrc.find(func (elem SinkTPtr) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	if idx != -1 {
		oldOsink := epSrc.Get()[idx]

		// Swap an already registered (different) object.
		// Update sink endpoint(s).
		if oldOsink != osink {
			epSrc.swap(idx, osink)
			BreakLink[LinkT](osrc, oldOsink)
			if !osink.IsUnresolved() {
				// Add 'osrc'.
				epSink.add(osrc)
			}
		}
	} else {
		// Otherwise, add it.
		epSrc.add(osink)
		if !osink.IsUnresolved() {
			epSink.add(osrc)
		}
	}

	return idx != -1 
}

func BreakLink[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT any, SinkT any, SrcTPtr IMajorPtr[SrcT], SinkTPtr IMajorPtr[SinkT]](
		osrc SrcTPtr, osink SinkTPtr) {
	var l LinkT
	epSrc := l.GetSourceEndpoint(osrc)
	epSink := l.GetSinkEndpoint(osink)

    absQualID := osink.GetAbsQualifiedID() 
	idx := epSrc.find(func (elem SinkTPtr) bool {
		return elem == osink 
	})

	// Break the link between the two objects, if exist.
	if idx != -1 {
		var unresOsink SinkTPtr = new(SinkT)
		unresOsink.initialize()
		unresOsink.markUnresolved(osrc.convIDAbsToRel(absQualID), osrc)
		epSrc.swap(idx, unresOsink)
		if !osink.IsUnresolved() {
			srcIdx := epSink.find(func (elem SrcTPtr) bool {
				return elem == osrc
			})
			if srcIdx != -1 {
				epSink.remove(srcIdx)
			}
		}
	}
}

func UnreserveLink[LinkT ILinkBase[SrcT, SinkT, SrcTPtr, SinkTPtr], SrcT any, SinkT any, SrcTPtr IMajorPtr[SrcT], SinkTPtr IMajorPtr[SinkT]](
		osrc SrcTPtr, absQualID string) {
	var l LinkT
	epSrc := l.GetSourceEndpoint(osrc)

	idx := epSrc.find(func (elem SinkTPtr) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	// If the same ID was reserved, (break and) remove it.
	if idx != -1 {
		osink := epSrc.Get()[idx]
		if !osink.IsUnresolved() {
			BreakLink[LinkT](osrc, osink)
		}
		epSrc.remove(idx)
	}
}

// Here are some specialized links.

type AttachTo struct {}
func (this AttachTo) GetSourceEndpoint(aobj *Annex) *EPSource[*MajorCommon] {
	return &aobj.AttachTo
}
func (this AttachTo) GetSinkEndpoint(obj *MajorCommon) *EPSink[*Annex] {
	return &obj.Attached
}

type ExtraNoteOf struct {}
func (this ExtraNoteOf) GetSourceEndpoint(aobj *Annex) *EPSource[*MajorCommon] {
	return &aobj.ExtraNoteOf
}
func (this ExtraNoteOf) GetSinkEndpoint(obj *MajorCommon) *EPSink[*Annex] {
	return &obj.ExtraNote
}

type ContinuedFrom struct {}
func (this ContinuedFrom) GetSourceEndpoint(aobj *Event) *EPSource[*Event] {
	return &aobj.ContinuedFrom
}
func (this ContinuedFrom) GetSinkEndpoint(obj *Event) *EPSink[*Event] {
	return &obj.ContinuedTo
}
