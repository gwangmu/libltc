package ltc 

// Endpoints are the "endpoints" of links. Each main object should include it
// with an appropriate field name.

type Endpoint[T IMajor] struct {
	objs []T
}

type EPSource[T IMajor] struct {
	Endpoint[T]
}

type EPSink[T IMajor] struct {
	Endpoint[T]
}

func (this *Endpoint[T]) Get() []T {
	return this.objs
}

func (this *Endpoint[T]) Has(obj T) bool {
	for _, o := range this.objs {
		if o.Equal(obj) {
			return true
		}
	}
	return false
}

func (this *Endpoint[T]) find(pred func(T) bool) int {
	for i, o := range this.objs {
		if pred(o) {
			return i
		}
	}
	return -1
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

///////////

func reserveLinkImpl[SrcT IMajor, SinkT IMajor](
		osrc SrcT, getEPSource func(SrcT) *EPSource[SinkT], relQualID string, getEPSink func(SinkT) *EPSink[SrcT], createUnres func(string, SrcT) SinkT) {
	epSrc := getEPSource(osrc)

	absQualID := osrc.convIDRelToAbs(relQualID)
	idx := epSrc.find(func (elem SinkT) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	// If the same ID was not reserved, add unresolved.
	if idx == -1 {
		unresOsink := createUnres(relQualID, osrc)
		createLinkImpl[SrcT, SinkT](osrc, getEPSource, unresOsink, getEPSink, createUnres)
	}
}

func createLinkImpl[SrcT IMajor, SinkT IMajor](
		osrc SrcT, getEPSource func(SrcT) *EPSource[SinkT], osink SinkT, getEPSink func(SinkT) *EPSink[SrcT], createUnres func(string, SrcT) SinkT) bool {
	epSrc := getEPSource(osrc)
	epSink := getEPSink(osink)

    absQualID := osink.GetAbsQualifiedID() 
	idx := epSrc.find(func (elem SinkT) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	if idx != -1 {
		oldOsink := epSrc.Get()[idx]

		// Swap an already registered (different) object.
		// Update sink endpoint(s).
		if !oldOsink.Equal(osink) {
			epSrc.swap(idx, osink)
			breakLinkImpl[SrcT, SinkT](osrc, getEPSource, osink, getEPSink, createUnres)
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

func breakLinkImpl[SrcT IMajor, SinkT IMajor](
		osrc SrcT, getEPSource func(SrcT) *EPSource[SinkT], osink SinkT, getEPSink func(SinkT) *EPSink[SrcT], createUnres func(string, SrcT) SinkT) {
	epSrc := getEPSource(osrc)
	epSink := getEPSink(osink)

    absQualID := osink.GetAbsQualifiedID() 
	idx := epSrc.find(func (elem SinkT) bool {
		return elem.Equal(osink)
	})

	// Break the link between the two objects, if exist.
	if idx != -1 {
		unresOsink := createUnres(osrc.convIDAbsToRel(absQualID), osrc)
		epSrc.swap(idx, unresOsink)
		if !osink.IsUnresolved() {
			srcIdx := epSink.find(func (elem SrcT) bool {
				return elem.Equal(osrc)
			})
			if srcIdx != -1 {
				epSink.remove(srcIdx)
			}
		}
	}
}

func unreserveLinkImpl[SrcT IMajor, SinkT IMajor](
		osrc SrcT, getEPSource func(SrcT) *EPSource[SinkT], absQualID string, getEPSink func(SinkT) *EPSink[SrcT], createUnres func(string, SrcT) SinkT) {
	epSrc := getEPSource(osrc)

	idx := epSrc.find(func (elem SinkT) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	// If the same ID was reserved, (break and) remove it.
	if idx != -1 {
		osink := epSrc.Get()[idx]
		if !osink.IsUnresolved() {
			breakLinkImpl[SrcT, SinkT](osrc, getEPSource, osink, getEPSink, createUnres)
		}
		epSrc.remove(idx)
	}
}

/////////////

// Here are some specialized links.

type AttachToLink struct {}
var AttachTo AttachToLink
func (AttachToLink) GetSourceEndpoint(aobj *Annex) *EPSource[*MajorCommon] {
	return &aobj.AttachTo
}
func (AttachToLink) GetSinkEndpoint(obj *MajorCommon) *EPSink[*Annex] {
	return &obj.Attached
}
func (this AttachToLink) ReserveLink(aobj *Annex, absQualID string) {
	reserveLinkImpl[*Annex, IMajor](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
}
func (this AttachToLink) CreateLink(aobj *Annex, obj IMajor) {
	createLinkImpl[*Annex, IMajor](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink)
}
func (this AttachToLink) BreakLink(aobj *Annex, obj IMajor) {
	breakLinkImpl[*Annex, IMajor](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink)
}
func (this AttachToLink) UnreserveLink(aobj *Annex, absQualID string) {
	unreserveLinkImpl[*Annex, IMajor](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
}
func (AttachToLink) getEPSource(aobj *Annex) *EPSource[IMajor] {
	return &aobj.AttachToTest
}
func (AttachToLink) getEPSink(obj IMajor) *EPSink[*Annex] {
	return &obj.GetCommon().Attached
}
func (AttachToLink) createUnresolvedSink(relQualID string, refcer *Annex) IMajor {
	ret := CreateEmptyMajorCommon()
	ret.markUnresolved(relQualID, refcer)
	return ret
}

type ExtraNoteOf struct {}
func (ExtraNoteOf) GetSourceEndpoint(aobj *Annex) *EPSource[*MajorCommon] {
	return &aobj.ExtraNoteOf
}
func (ExtraNoteOf) GetSinkEndpoint(obj *MajorCommon) *EPSink[*Annex] {
	return &obj.ExtraNote
}

type ContinuedFrom struct {}
func (ContinuedFrom) GetSourceEndpoint(aobj *Event) *EPSource[*Event] {
	return &aobj.ContinuedFrom
}
func (ContinuedFrom) GetSinkEndpoint(obj *Event) *EPSink[*Event] {
	return &obj.ContinuedTo
}
