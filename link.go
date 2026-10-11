package libltc 

// endpoints are the "endpoints" of links. Each main object should include it
// with an appropriate field name.

type endpoint[T IMajor] struct {
	objs []T
}

type epSink[SrcT IMajor] struct {
	endpoint[SrcT]
}

type epSource[SrcT IMajor, SinkT IMajor] struct {
	endpoint[SinkT]
	base SrcT
}

func (this *endpoint[T]) Get() []T {
	return this.objs
}

func (this *endpoint[T]) Has(obj T) bool {
	for _, o := range this.objs {
		if o.Equal(obj) {
			return true
		}
	}
	return false
}

func (this *endpoint[T]) find(pred func(T) bool) int {
	for i, o := range this.objs {
		if pred(o) {
			return i
		}
	}
	return -1
}

func (this *endpoint[T]) add(obj T) {
	this.objs = append(this.objs, obj)
}

func (this *endpoint[T]) swap(idx int, newobj T) {
	this.objs[idx] = newobj
}

func (this *endpoint[T]) remove(idx int) {
	this.objs = append(this.objs[:idx], this.objs[idx+1:]...) 
}

func (this *epSource[SrcT, _]) initialize(base SrcT) {
	this.base = base
}

////

func reserveLinkImpl[SrcT IMajor, SinkT IMajor](
		osrc SrcT, getepSource func(SrcT) *epSource[SrcT, SinkT], relQualID string, getepSink func(SinkT) *epSink[SrcT], createUnres func(string, SrcT) SinkT) {
	epSrc := getepSource(osrc)

	absQualID := osrc.convIDRelToAbs(relQualID)
	idx := epSrc.find(func (elem SinkT) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	// If the same ID was not reserved, add unresolved.
	if idx == -1 {
		unresOsink := createUnres(relQualID, osrc)
		linkImpl[SrcT, SinkT](osrc, getepSource, unresOsink, getepSink, createUnres, true)
	}
}

// force: create if not found (prev. reserved or linked)
func linkImpl[SrcT IMajor, SinkT IMajor](
		osrc SrcT, getepSource func(SrcT) *epSource[SrcT, SinkT], osink SinkT, getepSink func(SinkT) *epSink[SrcT], createUnres func(string, SrcT) SinkT, force bool) bool {
	epSrc := getepSource(osrc)
	epSink := getepSink(osink)

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
			unlinkImpl[SrcT, SinkT](osrc, getepSource, oldOsink, getepSink, createUnres)
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
		osrc SrcT, getepSource func(SrcT) *epSource[SrcT, SinkT], osink SinkT, getepSink func(SinkT) *epSink[SrcT], createUnres func(string, SrcT) SinkT) {
	epSrc := getepSource(osrc)
	epSink := getepSink(osink)

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
		osrc SrcT, getepSource func(SrcT) *epSource[SrcT, SinkT], absQualID string, getepSink func(SinkT) *epSink[SrcT], createUnres func(string, SrcT) SinkT) {
	epSrc := getepSource(osrc)

	idx := epSrc.find(func (elem SinkT) bool {
		return elem.GetAbsQualifiedID() == absQualID
	})

	// If the same ID was reserved, (break and) remove it.
	if idx != -1 {
		osink := epSrc.Get()[idx]
		if !osink.IsUnresolved() {
			unlinkImpl[SrcT, SinkT](osrc, getepSource, osink, getepSink, createUnres)
		}
		epSrc.remove(idx)
	}
}

// Here are some specialized links.

type AttachTo struct { epSource[*Annex, IMajor] }
type Attached struct { epSink[*Annex] }

func (this AttachTo) Reserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		reserveLinkImpl[*Annex, IMajor](aobj, this.getepSource, absQualID, this.getepSink, this.createUnresolvedSink)
	}
}
func (this AttachTo) Link(obj IMajor, force bool) {
	if aobj := this.base; aobj != nil {
		linkImpl[*Annex, IMajor](aobj, this.getepSource, obj, this.getepSink, this.createUnresolvedSink, force)
	}
}
func (this AttachTo) Unlink(obj IMajor) {
	if aobj := this.base; aobj != nil {
		unlinkImpl[*Annex, IMajor](aobj, this.getepSource, obj, this.getepSink, this.createUnresolvedSink)
	}
}
func (this AttachTo) Unreserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		unreserveLinkImpl[*Annex, IMajor](aobj, this.getepSource, absQualID, this.getepSink, this.createUnresolvedSink)
	}
}

func (AttachTo) getepSource(aobj *Annex) *epSource[*Annex, IMajor] {
	return &aobj.AttachTo.epSource
}
func (AttachTo) getepSink(obj IMajor) *epSink[*Annex] {
	return &obj.GetCommon().Attached.epSink
}
func (AttachTo) createUnresolvedSink(relQualID string, refcer *Annex) IMajor {
	ret := CreateEmptyMajorCommon()
	ret.markUnresolved(relQualID, refcer)
	return ret
}

type ExtraNoteOf struct { epSource[*Annex, IMajor] }
type ExtraNote struct { epSink[*Annex] }

func (this ExtraNoteOf) Reserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		reserveLinkImpl[*Annex, IMajor](aobj, this.getepSource, absQualID, this.getepSink, this.createUnresolvedSink)
	}
}
func (this ExtraNoteOf) Link(obj IMajor, force bool) {
	if aobj := this.base; aobj != nil {
		linkImpl[*Annex, IMajor](aobj, this.getepSource, obj, this.getepSink, this.createUnresolvedSink, force)
	}
}
func (this ExtraNoteOf) Unlink(obj IMajor) {
	if aobj := this.base; aobj != nil {
		unlinkImpl[*Annex, IMajor](aobj, this.getepSource, obj, this.getepSink, this.createUnresolvedSink)
	}
}
func (this ExtraNoteOf) Unreserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		unreserveLinkImpl[*Annex, IMajor](aobj, this.getepSource, absQualID, this.getepSink, this.createUnresolvedSink)
	}
}

func (ExtraNoteOf) getepSource(aobj *Annex) *epSource[*Annex, IMajor] {
	return &aobj.ExtraNoteOf.epSource
}
func (ExtraNoteOf) getepSink(obj IMajor) *epSink[*Annex] {
	return &obj.GetCommon().ExtraNote.epSink
}
func (ExtraNoteOf) createUnresolvedSink(relQualID string, refcer *Annex) IMajor {
	ret := CreateEmptyMajorCommon()
	ret.markUnresolved(relQualID, refcer)
	return ret
}

type ContinuedFrom struct { epSource[*Event, *Event] }
type ContinuedTo struct { epSink[*Event] }

func (this ContinuedFrom) Reserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		reserveLinkImpl[*Event, *Event](aobj, this.getepSource, absQualID, this.getepSink, this.createUnresolvedSink)
	}
}
func (this ContinuedFrom) Link(obj *Event, force bool) {
	if aobj := this.base; aobj != nil {
		linkImpl[*Event, *Event](aobj, this.getepSource, obj, this.getepSink, this.createUnresolvedSink, force)
	}
}
func (this ContinuedFrom) Unlink(obj *Event) {
	if aobj := this.base; aobj != nil {
		unlinkImpl[*Event, *Event](aobj, this.getepSource, obj, this.getepSink, this.createUnresolvedSink)
	}
}
func (this ContinuedFrom) Unreserve(absQualID string) {
	if aobj := this.base; aobj != nil {
		unreserveLinkImpl[*Event, *Event](aobj, this.getepSource, absQualID, this.getepSink, this.createUnresolvedSink)
	}
}

func (ContinuedFrom) getepSource(eobj *Event) *epSource[*Event, *Event] {
	return &eobj.ContinuedFrom.epSource
}
func (ContinuedFrom) getepSink(eobj *Event) *epSink[*Event] {
	return &eobj.ContinuedTo.epSink
}
func (ContinuedFrom) createUnresolvedSink(relQualID string, refcer *Event) *Event {
	ret := CreateEmptyEvent()
	ret.markUnresolved(relQualID, refcer)
	return ret
}
