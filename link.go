package ltc 

// Endpoints are the "endpoints" of links. Each main object should include it
// with an appropriate field name.

type Endpoint[T IMajor] struct {
	objs []T
}

type EPSink[SrcT IMajor] struct {
	Endpoint[SrcT]
}

type EPSource[SrcT IMajor, SinkT IMajor] struct {
	Endpoint[SinkT]
	base SrcT
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

func (this *EPSource[SrcT, _]) initialize(base SrcT) {
	this.base = base
}

////

func reserveLinkImpl[SrcT IMajor, SinkT IMajor](
		osrc SrcT, getEPSource func(SrcT) *EPSource[SrcT, SinkT], relQualID string, getEPSink func(SinkT) *EPSink[SrcT], createUnres func(string, SrcT) SinkT) {
	epSrc := getEPSource(osrc)

	absQualID := osrc.convIDRelToAbs(relQualID)
	idx := epSrc.find(func (elem SinkT) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	// If the same ID was not reserved, add unresolved.
	if idx == -1 {
		unresOsink := createUnres(relQualID, osrc)
		linkImpl[SrcT, SinkT](osrc, getEPSource, unresOsink, getEPSink, createUnres, true)
	}
}

// force: create if not found (prev. reserved or linked)
func linkImpl[SrcT IMajor, SinkT IMajor](
		osrc SrcT, getEPSource func(SrcT) *EPSource[SrcT, SinkT], osink SinkT, getEPSink func(SinkT) *EPSink[SrcT], createUnres func(string, SrcT) SinkT, force bool) bool {
	epSrc := getEPSource(osrc)
	epSink := getEPSink(osink)

    absQualID := osink.GetAbsQualifiedID() 
	idx := epSrc.find(func (elem SinkT) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	if idx != -1 {
		oldOsink := epSrc.objs[idx]

		// Swap an already registered (different) object.
		// Update sink endpoint(s).
		if !oldOsink.Equal(osink) {
			epSrc.swap(idx, osink)
			unlinkImpl[SrcT, SinkT](osrc, getEPSource, oldOsink, getEPSink, createUnres)
			if !osink.IsUnresolved() {
				// Add 'osrc'.
				epSink.add(osrc)
			}
		}
	} else if force {
		// Otherwise, add it.
		epSrc.add(osink)
		if !osink.IsUnresolved() {
			epSink.add(osrc)
		}
	}

	return idx != -1 
}

func unlinkImpl[SrcT IMajor, SinkT IMajor](
		osrc SrcT, getEPSource func(SrcT) *EPSource[SrcT, SinkT], osink SinkT, getEPSink func(SinkT) *EPSink[SrcT], createUnres func(string, SrcT) SinkT) {
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
		osrc SrcT, getEPSource func(SrcT) *EPSource[SrcT, SinkT], absQualID string, getEPSink func(SinkT) *EPSink[SrcT], createUnres func(string, SrcT) SinkT) {
	epSrc := getEPSource(osrc)

	idx := epSrc.find(func (elem SinkT) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	// If the same ID was reserved, (break and) remove it.
	if idx != -1 {
		osink := epSrc.Get()[idx]
		if !osink.IsUnresolved() {
			unlinkImpl[SrcT, SinkT](osrc, getEPSource, osink, getEPSink, createUnres)
		}
		epSrc.remove(idx)
	}
}

// Here are some specialized links.

type AttachTo struct { EPSource[*Annex, IMajor] }
type Attached struct { EPSink[*Annex] }

func (this AttachTo) Reserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		reserveLinkImpl[*Annex, IMajor](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
	}
}
func (this AttachTo) Link(obj IMajor, force bool) {
	if aobj := this.base; aobj != nil {
		linkImpl[*Annex, IMajor](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink, force)
	}
}
func (this AttachTo) Unlink(obj IMajor) {
	if aobj := this.base; aobj != nil {
		unlinkImpl[*Annex, IMajor](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink)
	}
}
func (this AttachTo) Unreserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		unreserveLinkImpl[*Annex, IMajor](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
	}
}

func (AttachTo) getEPSource(aobj *Annex) *EPSource[*Annex, IMajor] {
	return &aobj.AttachTo.EPSource
}
func (AttachTo) getEPSink(obj IMajor) *EPSink[*Annex] {
	return &obj.GetCommon().Attached.EPSink
}
func (AttachTo) createUnresolvedSink(relQualID string, refcer *Annex) IMajor {
	ret := CreateEmptyMajorCommon()
	ret.markUnresolved(relQualID, refcer)
	return ret
}

type ExtraNoteOf struct { EPSource[*Annex, IMajor] }
type ExtraNote struct { EPSink[*Annex] }

func (this ExtraNoteOf) Reserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		reserveLinkImpl[*Annex, IMajor](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
	}
}
func (this ExtraNoteOf) Link(obj IMajor, force bool) {
	if aobj := this.base; aobj != nil {
		linkImpl[*Annex, IMajor](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink, force)
	}
}
func (this ExtraNoteOf) Unlink(obj IMajor) {
	if aobj := this.base; aobj != nil {
		unlinkImpl[*Annex, IMajor](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink)
	}
}
func (this ExtraNoteOf) Unreserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		unreserveLinkImpl[*Annex, IMajor](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
	}
}

func (ExtraNoteOf) getEPSource(aobj *Annex) *EPSource[*Annex, IMajor] {
	return &aobj.ExtraNoteOf.EPSource
}
func (ExtraNoteOf) getEPSink(obj IMajor) *EPSink[*Annex] {
	return &obj.GetCommon().ExtraNote.EPSink
}
func (ExtraNoteOf) createUnresolvedSink(relQualID string, refcer *Annex) IMajor {
	ret := CreateEmptyMajorCommon()
	ret.markUnresolved(relQualID, refcer)
	return ret
}

type ContinuedFrom struct { EPSource[*Event, *Event] }
type ContinuedTo struct { EPSink[*Event] }

func (this ContinuedFrom) Reserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		reserveLinkImpl[*Event, *Event](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
	}
}
func (this ContinuedFrom) Link(obj *Event, force bool) {
	if aobj := this.base; aobj != nil {
		linkImpl[*Event, *Event](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink, force)
	}
}
func (this ContinuedFrom) Unlink(obj *Event) {
	if aobj := this.base; aobj != nil {
		unlinkImpl[*Event, *Event](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink)
	}
}
func (this ContinuedFrom) Unreserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		unreserveLinkImpl[*Event, *Event](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
	}
}

func (ContinuedFrom) getEPSource(eobj *Event) *EPSource[*Event, *Event] {
	return &eobj.ContinuedFrom.EPSource
}
func (ContinuedFrom) getEPSink(eobj *Event) *EPSink[*Event] {
	return &eobj.ContinuedTo.EPSink
}
func (ContinuedFrom) createUnresolvedSink(relQualID string, refcer *Event) *Event {
	ret := CreateEmptyEvent()
	ret.markUnresolved(relQualID, refcer)
	return ret
}
