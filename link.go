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

////

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
		oldOsink := epSrc.objs[idx]

		// Swap an already registered (different) object.
		// Update sink endpoint(s).
		if !oldOsink.Equal(osink) {
			epSrc.swap(idx, osink)
			breakLinkImpl[SrcT, SinkT](osrc, getEPSource, oldOsink, getEPSink, createUnres)
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

// Here are some specialized links.

type AttachToLink struct{}
var AttachTo AttachToLink

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
	return &aobj.AttachTo
}
func (AttachToLink) getEPSink(obj IMajor) *EPSink[*Annex] {
	return &obj.GetCommon().Attached
}
func (AttachToLink) createUnresolvedSink(relQualID string, refcer *Annex) IMajor {
	ret := CreateEmptyMajorCommon()
	ret.markUnresolved(relQualID, refcer)
	return ret
}

type ExtraNoteOfLink struct{}
var ExtraNoteOf ExtraNoteOfLink

func (this ExtraNoteOfLink) ReserveLink(aobj *Annex, absQualID string) {
	reserveLinkImpl[*Annex, IMajor](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
}
func (this ExtraNoteOfLink) CreateLink(aobj *Annex, obj IMajor) {
	createLinkImpl[*Annex, IMajor](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink)
}
func (this ExtraNoteOfLink) BreakLink(aobj *Annex, obj IMajor) {
	breakLinkImpl[*Annex, IMajor](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink)
}
func (this ExtraNoteOfLink) UnreserveLink(aobj *Annex, absQualID string) {
	unreserveLinkImpl[*Annex, IMajor](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
}

func (ExtraNoteOfLink) getEPSource(aobj *Annex) *EPSource[IMajor] {
	return &aobj.ExtraNoteOf
}
func (ExtraNoteOfLink) getEPSink(obj IMajor) *EPSink[*Annex] {
	return &obj.GetCommon().ExtraNote
}
func (ExtraNoteOfLink) createUnresolvedSink(relQualID string, refcer *Annex) IMajor {
	ret := CreateEmptyMajorCommon()
	ret.markUnresolved(relQualID, refcer)
	return ret
}

type ContinuedFromLink struct{}
var ContinuedFrom ContinuedFromLink

func (this ContinuedFromLink) ReserveLink(aobj *Event, absQualID string) {
	reserveLinkImpl[*Event, *Event](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
}
func (this ContinuedFromLink) CreateLink(aobj *Event, obj *Event) {
	createLinkImpl[*Event, *Event](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink)
}
func (this ContinuedFromLink) BreakLink(aobj *Event, obj *Event) {
	breakLinkImpl[*Event, *Event](aobj, this.getEPSource, obj, this.getEPSink, this.createUnresolvedSink)
}
func (this ContinuedFromLink) UnreserveLink(aobj *Event, absQualID string) {
	unreserveLinkImpl[*Event, *Event](aobj, this.getEPSource, absQualID, this.getEPSink, this.createUnresolvedSink)
}

func (ContinuedFromLink) getEPSource(obj *Event) *EPSource[*Event] {
	return &obj.ContinuedFrom
}
func (ContinuedFromLink) getEPSink(obj *Event) *EPSink[*Event] {
	return &obj.ContinuedTo
}
func (ContinuedFromLink) createUnresolvedSink(relQualID string, refcer *Event) *Event {
	ret := CreateEmptyEvent()
	ret.markUnresolved(relQualID, refcer)
	return ret
}
