package libltc

import (
	"errors"
	"slices"
	"strings"
	"maps"
	"regexp"

	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

type IChart interface {
	GetEmbeddingEvent() *Event
	GetStringifiedSubjectName() string

	GetEventCategories() []string
	GetEventsInCategory(category string) []*Event
	GetAnnexs() []*Annex
	GetImports() []*Import
	GetAllMainObjects(inclLoadedSubchartObjs bool) []IMainObj

	GetObjectByAbsQualID(qualid string) IMainObj 
 	GetObjectByRelQualID(iqualid string, baseChart *Chart) IMainObj 
	GetObjectsByIntrinsicID(locid string) []IMainObj 

	getNextNumberID(kind MainObjKind) (NumberID, error) 

	asChart() *Chart
	resolveReferenceTo(IMainObj)
	unresolveReferenceTo(IMainObj)
	changeEventStartDate(*Event)
	changeEventCategory(*Event, string)
}

type Chart struct {
	Filepath string
	Version Version
	Note Note

	Setting Setting
	Subject Subject

	// `events`, `annexs`, and `imports` will contain imported objects, too.
	events map[string][]*Event			// By category, sorted by StartDate
	annexs []*Annex
	imports []*Import

	embeddingEvent *Event	// Only if this chart embedded as a subchart. 
}

//-- interface IChart (package internal)

func (this *Chart) GetEmbeddingEvent() *Event {
	return this.embeddingEvent
}

func (this *Chart) GetStringifiedSubjectName() string {
	return this.Subject.Name.String()
}

func (this *Chart) asChart() *Chart {
	return this
}

func (this *Chart) resolveReferenceTo(obj IMainObj) {
	// Ensure it's the main object itself, not just the common part. 
	obj = obj.GetEnclosingObject()

	// If any chart event points to 'obj', resolve 'ContinuedFrom'.
	if eobj, ok := obj.(*Event); ok {
		for _, evs := range this.events {
			for _, ceobj := range evs {
				ResolveLink[ContinuedFrom](ceobj, eobj)
			}
		}
		// NOTE: subchart objects should be resolved when it's loaded.
	}

	// Resolve the references from existing annexs to 'obj'.
	for _, caobj := range this.annexs {
		ResolveLink[AttachTo](caobj, obj.GetCommon())
		ResolveLink[ExtraNoteOf](caobj, obj.GetCommon())
	}
}

func (this *Chart) unresolveReferenceTo(obj IMainObj) {
	// Ensure it's the main object itself, not just the common part. 
	obj = obj.GetEnclosingObject()

	// If 'obj' points to an event, unresolve 'ContinuedFrom'.
	if eobj, ok := obj.(*Event); ok {
		if eobj.IsSubchartLoaded() {
			subchart, _ := eobj.GetSubchart()
			for _, scobj := range subchart.GetAllMainObjects(true) {
				this.unresolveReferenceTo(scobj)
			}
		}

		for _, evs := range this.events {
			for _, event := range evs {
				UnresolveLink[ContinuedFrom](event, eobj)
			}
		}
	}

	// Unresolve the references from existing annexs to 'obj'.
	for _, annex := range this.annexs {
		UnresolveLink[AttachTo](annex, obj.GetCommon())
		UnresolveLink[ExtraNoteOf](annex, obj.GetCommon())
	}
}

func (this *Chart) changeEventStartDate(eobj *Event) {
	category := eobj.GetCategory()
	tmpevs := removeMainObjectFromSlice[Event](this.events[category], eobj)
	this.events[category] = insertEventToSlice(tmpevs, eobj)
}

func (this *Chart) changeEventCategory(eobj *Event, oldc string) {
	newc := eobj.GetCategory()
	tmpevs := removeMainObjectFromSlice[Event](this.events[oldc], eobj)
	this.events[newc] = insertEventToSlice(tmpevs, eobj)
}

//-- interface IWarningObj

func (this *Chart) Summary() (ret string) {
	ret = "the LTC chart"
	if !this.Subject.IsUnknown() {
		ret += " for " + this.Subject.Summary()
	} else if this.Filepath != "" {
		ret += " from the file at '" + this.Filepath + "'"
	}
	return
}

func (this Chart) IsUnknown() bool {
	return this.Subject.IsUnknown() && this.Filepath == ""
}

//-- method (main object association)

func (this *Chart) getNextNumberID(kind MainObjKind) (NumberID, error) {
	mainobjArr := []IMainObj{}

	switch kind {
	case MOK_Event:
		for _, evs := range this.events {
			for _, o := range evs {
				mainobjArr = append(mainobjArr, o)
			}	
		}
	case MOK_Annex:
		for _, o := range this.annexs {
			mainobjArr = append(mainobjArr, o)
		}
	case MOK_Import:
		for _, o := range this.imports {
			mainobjArr = append(mainobjArr, o)
		}
	default:
		return NID_Invalid, errors.New("Unrecognized main object kind")
	}

	// Next next ID = Max existing numID + 1
	// FIX: error out if 'maxNumID' == UINT_MAX - 1. Unlikely, but still.
	var maxNumID NumberID = 0
	for _, mainobj := range mainobjArr {
		maxNumID = max(maxNumID, mainobj.getNumberID())
	}
	if (maxNumID == NID_Max) {
		return 0, errors.New("Cannot get the next number ID")
	} else {
		return maxNumID + 1, nil
	}
}

//-- method (getters)

func (this *Chart) GetAnnexs() []*Annex {
	return this.annexs
}

func (this *Chart) GetImports() []*Import {
	return this.imports
}

//-- method (setters)

func (this *Chart) setEmbeddingEvent(o *Event) {
	this.embeddingEvent = o
}

//-- method (main object manipulation)

func (this *Chart) GetObjectByAbsQualID(absQualID string) IMainObj {
	inChartQualIDs := splitIntoInChartQualIDs(absQualID)
	curChart := this
	for i, inChartQualID := range inChartQualIDs {
		curObj := curChart.GetObjectByRelQualID(inChartQualID, curChart)

		if curObj != nil {
			if i == len(inChartQualIDs) - 1 {
				return curObj
			} else if eobj, ok := curObj.(*Event); ok {
				if eobj.IsChartEmbedding() {
					curChart, _ = eobj.GetSubchart()
					continue
				} else {
					break
				}
			}
		}
	}
	return nil
}

func (this *Chart) GetObjectByRelQualID(relQualID string, baseChart *Chart) IMainObj {
	objs := this.getObjectsWithIDGetter(relQualID, func (obj IMainObj) string {
		return obj.GetRelQualifiedID(baseChart)
	}, true)

	if len(objs) > 0 {
		return objs[0]
	} else {
		return nil
	}
}

func (this *Chart) GetObjectsByIntrinsicID(intrID string) []IMainObj {
	return this.getObjectsWithIDGetter(intrID, func (obj IMainObj) string {
		return obj.GetIntrinsicID()
	}, false)
}

func (this *Chart) GetAllMainObjects(inclLoadedSubchartObjs bool) (ret []IMainObj) {
	for _, evs := range this.events {
		for _, eobj := range evs {
			ret = append(ret, eobj)
			if inclLoadedSubchartObjs && eobj.IsSubchartLoaded() {
				schart, _ := eobj.GetSubchart()
				for _, seobj := range schart.GetAllMainObjects(true) {
					ret = append(ret, seobj)
				}
			}
		}
	}
	for _, aobj := range this.annexs {
		ret = append(ret, aobj)
	}
	for _, iobj := range this.imports {
		ret = append(ret, iobj)
	}
	return
}

func (this *Chart) GetEventCategories() []string {
	return slices.Collect(maps.Keys(this.events))
}

func (this *Chart) GetEventsInCategory(category string) []*Event {
	if evs, ok := this.events[category]; ok {
		return evs
	} else {
		return []*Event{}
	}
}

// start, end: all inclusive.
func (this *Chart) GetEventsInCategoryBetween(category string, start Time, end Time) (ret []*Event) {
	if evs, ok := this.events[category]; !ok {
		return
	} else {
		inbound := false 
		for _, eobj := range evs {
			if !inbound {
				// Check inbound (based on StartDate).
				inbound = eobj.GetStartDate().IsAfter(start)
			}
			
			if inbound {
				// Check out-of-bound (based on StartDate).
				if eobj.GetStartDate().IsAfter(end) {
					break
				}
			}

			if inbound {
				// Add if it's actually inbound.
				if end.IsAfter(eobj.GetEndDate()) {
					ret = append(ret, eobj)
				}
			}
		}
		return
	}
}

func (this *Chart) AddObject(obj IMainObj) error {
	if obj == nil || obj.GetKind() == MOK_Unknown {
		return errors.New("Bogus object. Cannot add to the chart.")
	}

	obj = obj.GetEnclosingObject()
	objKind := obj.GetKind()

	// Associate 'obj' to this chart.
	if nnid, err := this.getNextNumberID(objKind); err != nil {
		obj.setChart(this)
		obj.setNumberID(nnid)
	} else {
		return err
	}

	switch objKind {
	case MOK_Event:
		eobj, _ := obj.(*Event)

		// Resolve references.
		this.resolveReferenceTo(eobj)
		eobj.resolveReferenceTo(this)

		// Insert 'obj' to the slice.
		category := eobj.GetCategory()
		this.ReserveEventCategory(category)
		this.events[category] = insertEventToSlice(this.events[category], eobj)
	case MOK_Annex:
		aobj, _ := obj.(*Annex)

		// Resolve references.
		this.resolveReferenceTo(aobj)
		aobj.resolveReferenceTo(this)

		// Insert 'obj' to the slice.
		this.annexs = append(this.annexs, aobj)
	case MOK_Import:
		iobj, _ := obj.(*Import)

		// Resolve references.
		this.resolveReferenceTo(iobj)
		iobj.resolveReferenceTo(this)
		for _, imobj := range iobj.GetImportedMainObjects() {
			this.resolveReferenceTo(imobj)
			// NOTE: Not the other way round. No back reference.
		}

		// Insert 'obj' to the slice.
		this.imports = append(this.imports, iobj)
		for _, ieobj := range iobj.GetImportedEvents() {
			category := ieobj.GetCategory()
			this.ReserveEventCategory(category)
			this.events[category] = insertEventToSlice(this.events[category], ieobj)
		}
		for _, iaobj := range iobj.GetImportedAnnexs() {
			this.annexs = append(this.annexs, iaobj)
		}
		for _, iiobj := range iobj.GetImportedImports() {
			this.imports = append(this.imports, iiobj)
		}
	}

	return nil
}

func (this *Chart) RemoveObject(obj IMainObj) {
	if obj == nil || obj.GetKind() == MOK_Unknown {
		return
	}

	objKind := obj.GetKind()

	// Dissociate 'obj' to this chart.
	obj.setChart(nil)
	obj.setNumberID(NID_Invalid)

	switch objKind {
	case MOK_Event:
		eobj, _ := obj.(*Event)

		// Remove 'obj' from the slice.
		category := eobj.GetCategory()
		this.events[category] = removeMainObjectFromSlice[Event](this.events[category], eobj)

		// Unresolve references.
		this.unresolveReferenceTo(eobj)
		eobj.unresolveReferenceTo(this)
	case MOK_Annex:
		aobj, _ := obj.(*Annex)

		// Remove 'obj' from the slice.
		this.annexs = removeMainObjectFromSlice[Annex](this.annexs, aobj)

		// Unresolve references.
		this.unresolveReferenceTo(aobj)
		aobj.unresolveReferenceTo(this)
	case MOK_Import:
		iobj, _ := obj.(*Import)

		// Remove 'obj' from the slice.
		this.imports = removeMainObjectFromSlice[Import](this.imports, iobj)
		for _, ieobj := range iobj.GetImportedEvents() {
			category := ieobj.GetCategory()
			this.events[category] = removeMainObjectFromSlice[Event](this.events[category], ieobj)
		}
		for _, iaobj := range iobj.GetImportedAnnexs() {
			this.annexs = removeMainObjectFromSlice[Annex](this.annexs, iaobj)
		}
		for _, iiobj := range iobj.GetImportedImports() {
			this.imports = removeMainObjectFromSlice[Import](this.imports, iiobj)
		}

		// Unresolve references.
		this.unresolveReferenceTo(iobj)
		iobj.unresolveReferenceTo(this)
		for _, imobj := range iobj.GetImportedMainObjects() {
			this.unresolveReferenceTo(imobj)
			// NOTE: Not the other way round. No back reference.
		}
	}
}

func (this *Chart) ReserveEventCategory(category string) {
	if _, ok := this.events[category]; !ok {
		this.events[category] = []*Event{}
	}
}

func (this *Chart) RemoveEventCategory(category string) error {
	// TODO: remove events if it's only empty.
	if evs, ok := this.events[category]; ok {
		if len(evs) != 0 {
			return errors.New("Not an empty category")
		} else {
			delete(this.events, category)
		}
	}
	return nil
}

//-- method (creation)

func createChartFromParsed(o *file.File) (*Chart, warning.Warnings) {
	chart := &Chart{}
	warns := warning.Warnings{}

	// Load small preamble objects.
	chart.Filepath = o.Filepath
	chart.Version = GetVersion(o.Version)

	note, moreWarns := CreateNoteFromString(o.Note)
	warns.Concat(moreWarns)
	chart.Note = note

	// Load setting object.
	setting, moreWarns := createSettingFromParsed(o.Setting)
	warns.Concat(moreWarns)
	chart.Setting = setting

	// Pre-register categories.
	// We need this to save (potentially) empty categories, too.
	for _, category := range o.Setting.Categories {
		chart.ReserveEventCategory(category)
	}

	// Load subject object.
	subject, moreWarns := createSubjectFromParsed(o.Subject)
	warns.Concat(moreWarns)
	chart.Subject = subject

	// Load event objects.
	for _, fevent := range o.Event {
		cevent, moreWarns := createEventFromParsed(fevent)
		warns.Concat(moreWarns)
		chart.AddObject(cevent)
	}
	
	// Load annex objects.
	for _, fannex := range o.Annex {
		cannex, moreWarns := createAnnexFromParsed(fannex)
		warns.Concat(moreWarns)
		chart.AddObject(cannex)
	}

	// Load import objects.
	for _, fimport := range o.Import {
		cimport, moreWarns := createImportFromParsed(fimport)
		warns.Concat(moreWarns)
		chart.AddObject(cimport)
	}

	// Diagnose main objects. (non-local: comfirmed after adding to chart)
	for _, cevs := range chart.events {
		for _, cevent := range cevs {
			warns.Concat(cevent.DiagnoseNonLocal())
		}
	}
	for _, cannex := range chart.annexs {
		warns.Concat(cannex.DiagnoseNonLocal())
	}
	for _, cimport := range chart.imports {
		warns.Concat(cimport.DiagnoseNonLocal())
	}

	return chart, warns
}

func CreateChart(uri string) (*Chart, warning.Warnings) {
	// Parse LTC file in TOML format.
	ltcf, warns, err := file.LoadFromURI[file.File](uri)
	if err != nil {
		return CreateEmptyChart(), warns
	}

	// Recursively convert (TOML-format) file to (in-memory) chart.
	chart, warnsMore := createChartFromParsed(ltcf)
	warns = append(warns, warnsMore...)

	return chart, warns 
}

func CreateEmptyChart() *Chart {
	return &Chart{
		Filepath: "",
		Version: VER_26_09_1,
		Note: Note{},

		Setting: Setting{},
		Subject: Subject{},

		events: map[string][]*Event{},
		annexs: []*Annex{},
		imports: []*Import{},

		embeddingEvent: nil,
	}
}

//-- method (util)

func (this *Chart) getObjectsWithIDGetter(id string, idGetter func(IMainObj) string, onlyFirst bool) (ret []IMainObj) {
	objKind := findObjectKindByID(id)
	switch objKind {
	case MOK_Event:
		for _, evs := range this.events {
			for _, event := range evs {
				if idGetter(event) == id {
					ret = append(ret, event)
					if onlyFirst {
						return
					}
				}
			}
		}
	case MOK_Annex:
		for _, annex := range this.annexs {
			if idGetter(annex) == id {
				ret = append(ret, annex)
				if onlyFirst {
					return
				}
			}
		}
	case MOK_Import:
		for _, cimport := range this.imports {
			if idGetter(cimport) == id {
				ret = append(ret, cimport)
				if onlyFirst {
					return
				}
			}
		}
	}
	return 
}

func splitIntoInChartQualIDs(absQualID string) (ret []string) {
	re := regexp.MustCompile(`e[0-9]+/`)

	prevI := 0
	for prevI < len(absQualID) {
		if idxs := re.FindStringIndex(absQualID[prevI:]); idxs != nil {
			ret = append(ret, absQualID[prevI:prevI+idxs[1]-1])
			prevI += idxs[1]
		} else {
			break
		}
	}
	if prevI < len(absQualID) {
		ret = append(ret, absQualID[prevI:])
	}

	return 
}

func findObjectKindByID(id string) MainObjKind {
	ids := strings.Split(id, "/")
	if len(ids) != 0 {
		lastid := ids[len(ids)-1]
		re := regexp.MustCompile(`^([a-z])[0-9]+$`)
		if match := re.FindStringSubmatch(lastid); match != nil {
			kind, _ := GetMainObjKind(match[1])
			return kind
		}
	}
	return MOK_Unknown
}

func insertEventToSlice(evs []*Event, eobj *Event) (ret []*Event) {
	// Events are inserted by its start date (ascending).
	inserted := false
	for i, oeobj := range evs {
		if oeobj.GetStartDate().IsAfter(eobj.GetStartDate()) {
			ret = slices.Insert(evs, i, eobj)
			inserted = true
			break
		}
	}
	if !inserted {
		ret = append(evs, eobj)
	}
	return
}

func removeMainObjectFromSlice[T any, TPtr IMainObjPtr[T]](objs []TPtr, obj TPtr) []TPtr {
	idx := slices.Index(objs, obj)
	if idx != -1 {
		return slices.Delete(objs, idx, idx+1)
	} else {
		return objs
	}
}

//-- struct Setting 

type Setting struct {
	CalendarSystem string
	NoteFormat string
	Attrs map[string][]string
}

//-- struct Setting: interface IWarningObj

func (this *Setting) Summary() string {
	return warning.BuildSummaryString(
		WSE{"calendar", this.CalendarSystem},
		WSE{"note", this.NoteFormat},
	)
}

func (this Setting) IsUnknown() bool {
	return false
}

//-- struct Setting: interface IDiagnosable

func (this *Setting) DiagnoseLocal() (warns warning.Warnings) {
	// For now, any `CalendarSystem` or `NoteFormat` is valid,
	// assuming that the front-end tool will generate warnings if they
	// don't support it. Since any `CalendarSystem` is fine in libLTC,
	// the time object should implement local logic to compare times,
	// not relying on the GO's time package.

	if this.CalendarSystem == "" {
		this.CalendarSystem = "Gregorian"
		warns.Add("unspecified calendar system in @0@. setting to 'Gregorian'...", this)
	}

	if this.NoteFormat == "" {
		this.NoteFormat = "Markdown"
		warns.Add("unspecified note format in @0@. setting to 'Markdown'...", this)
	}
	
	if len(this.Attrs) != 0 {
		for akey, _ := range this.Attrs {
			warns.Add("unrecognized attribute '%s' in @0@.", akey, this)
		}
	}
	return
}

func (this *Setting) DiagnoseNonLocal() warning.Warnings {
	return warning.Warnings{}
}

//-- struct Setting: method (creation)

func createSettingFromParsed(o file.Setting) (Setting, warning.Warnings) {
	setting := Setting{
		CalendarSystem: "",
		NoteFormat: "",
		Attrs: map[string][]string{},
	}

	setting.CalendarSystem = o.CalendarSystem
	setting.NoteFormat = o.NoteFormat
	setting.Attrs = convAttrsFileToChart(o.Attrs)

	warns := setting.DiagnoseLocal()

	return setting, warns
}

//-- struct Subject

type Subject struct {
	Name Name
	StartDate Time
	EndDate Time
	Sex string
}

//-- struct Subject: interface IWarningObj

func (this *Subject) Summary() (ret string) {
	ret = this.Name.Summary()
	
	extraStr := warning.BuildSummaryString(
		WSE{"", this.Sex}, 
		WSE{"started", this.StartDate},
		WSE{"ended", this.EndDate},
	)
	if (len(extraStr) > 0) {
		ret += " (" + extraStr + ")"
	}
	return
}

func (this Subject) IsUnknown() bool {
	return this.Name.IsUnknown() && this.StartDate.IsUnknown() &&
			this.EndDate.IsUnknown() && this.Sex == ""
}

//-- struct Subject: interface IDiagnosable

func (this *Subject) DiagnoseLocal() (warns warning.Warnings) {
	if this.Name.IsUnknown() {
		warns.Add("unknown subject.")
	}
	return
}

func (this *Subject) DiagnoseNonLocal() warning.Warnings {
	return warning.Warnings{}
}

//-- struct Subject: method (creation) 

func createSubjectFromParsed(o file.Subject) (Subject, warning.Warnings) {
	subject := Subject{
		Name: Name{},
		StartDate: CreateEmptyTime(TUK_Start),
		EndDate: CreateEmptyTime(TUK_End),
		Sex: "",
	}

	subject.Name = createNameFromParsed(o.Name)
	subject.StartDate = createTimeFromParsed(o.StartDate, TUK_Start)
	subject.EndDate = createTimeFromParsed(o.EndDate, TUK_End)
	subject.Sex = o.Sex

	warns := subject.DiagnoseLocal()

	return subject, warns
}
