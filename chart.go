package libltc

import (
	"errors"
	"slices"
	"maps"

	"github.com/gwangmu/libltc/warning"
	"github.com/gwangmu/libltc/internal/file"
)

type IChart interface {
	GetStringifiedSubjectName() string
	GetAnnexs() []*Annex
	GetImports() []*Import
	GetObject(qualid string) IMainObj 
	GetObjectByLocalID(locid string) IMainObj 
	GetAllMainObjects() []IMainObj
	GetEmbeddingEvent() *Event
	GetEventCategories() []string
	GetEventsInCategory(category string) []*Event
	HasObject(absQualID string) bool

	getNextNumberID(kind MainObjKind) (NumberID, error) 

	asChart() *Chart
	resolveReferenceTo(IMainObj)
	unresolveReferenceTo(IMainObj)
	changeEventStartDate(*Event)
	changeEventCategory(*Event)
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

func (this *Chart) asChart() *Chart {
	return this
}

func (this *Chart) resolveReferenceTo(obj IMainObj) {
	for _, evs := range this.events {
		for _, event := range evs {
			switch eobj := obj.(type) {
			case *Event:
				event.resolveContinuedFromEvent(eobj)
			}
		}
	}

	for _, annex := range this.annexs {
		annex.resolveAttachToObject(obj)
		annex.resolveExtraNoteOfObject(obj)
	}
}

func (this *Chart) unresolveReferenceTo(obj IMainObj) {
	objID := obj.GetQualifiedID()
	
	for _, evs := range this.events {
		for _, event := range evs {
			event.unresolveContinuedFromEvent(objID)
		}
	}

	for _, annex := range this.annexs {
		annex.unresolveAttachToObject(objID)
		annex.unresolveExtraNoteOfObject(objID)
	}
}

func (this *Chart) changeEventStartDate(eobj *Event) {
	// TODO: find the previous index of `eobj` in `this.events`.
	// TODO: should find! remove this index in `this.events`.
	// TODO: scanning `this.events` from the beginning,
	// TODO: find the index where `eobj`'s StartDate is !IsAfter() first.
	// TODO: insert `eobj` in that index.
	panic("Unimplmemented")
}

func (this *Chart) changeEventCategory(eobj *Event) {
	// TODO: find the previous index of `eobj` in `this.events`.
	// TODO: should find! remove this index in `this.events`.
	// TODO: scanning `this.events` from the beginning,
	// TODO: find the index where `eobj`'s StartDate is !IsAfter() first.
	// TODO: insert `eobj` in that index.
	panic("Unimplemented")
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

func (this *Chart) GetStringifiedSubjectName() string {
	return this.Subject.Name.String()
}

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

func (this *Chart) GetObject(qualid string) IMainObj {
	// TODO: return obj by qualified id.
	panic("Unimplemented")

	// TODO: cut at the earliest 'e*' => before, after (no / inbetween).
	// TODO: find 'before' by the last-ID's kind. (matching qualid in events, annex, imports)
	// TODO: if the found object is event and embedding a chart and after != "",
	// TODO: call 'GetObject(after)' of the embedded chart. return that.
}

func (this *Chart) GetObjectByLocalID(locid string) IMainObj {
	// TODO: return obj by local id.
	panic("Unimplemented")
}

func (this *Chart) GetAllMainObjects() (ret []IMainObj) {
	for _, evs := range this.events {
		for _, eobj := range evs {
			ret = append(ret, eobj)
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

func (this *Chart) GetEventsInCategoryBetween(category string, start Time, end Time, inclusive bool) []*Event {
	// TODO: inclusive = false: start < estart && eend < end
	// TODO: inclusive = true: start <= eend || estart <= end
	panic("Unimplemented")
}

func (this *Chart) HasObject(absQualID string) bool {
	// TODO
	panic("Unimplemented")
}

func (this *Chart) AddObject(obj IMainObj) {
	// TODO: add already-created obj. update chart and id.
	// TODO: sort events by StartDate.
	// TODO: resolveReferenceTo 'obj'
	// TODO: resolve the references of 'obj' itself.
	// TODO: for import objects, add `imported*` to the chart, too.
	// TODO: for import objects, resorveReferenceTo all imported objs.
	panic("Unimplemented")
}

func (this *Chart) RemoveObject(obj IMainObj) {
	// TODO: dispose of any possible links to other objs.
	// TODO: for import objects, unresorveReferenceTo all imported objs.
	// TODO: for import objects, remove `imported*` from the chart, too.
	// TODO: unresolveReferenceTo 'obj'
	// TODO: unresolve the references of 'obj' itself.
	// TODO: remove itself from the chart.
	panic("Unimplemented")
}

func (this *Chart) ReserveEventCategory(category string) {
	if _, ok := this.events[category]; !ok {
		this.events[category] = []*Event{}
	}
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
