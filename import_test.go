package libltc 

import (
	"slices"
	"testing"
	"github.com/gwangmu/libltc/internal/test"
)

func TestSingleImport_Import(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'single_import.ltc'...")
	chart, wc := CreateChart("single_import.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded chart...")
	{
		gotIDs := []string{}
		wantIDs := []string{ "i0", "i0/e0" }
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	ieobj := chart.GetObjectByQualifiedID("i0/e0").(*Event)
	test.AssertEQ(ieobj.GetTitle(), "Wherever")
	test.AssertEQ(ieobj.GetStartDate().String(), "1991-03-29 (Attr0)")
	test.AssertEQ(ieobj.GetCategory(), "Residence")
	test.AssertEQ(chart.GetEventsInCategory("Residence")[0], ieobj)

	t.Logf("Changing import target to 'single_annex.ltc'...")
	iobj := chart.GetImports()[0]
	wc2 := iobj.SetLink("./single_annex.ltc")
	test.PrintAllWarnings(wc2)

	test.AssertEQ(len(chart.GetEvents()), 0)
	test.AssertEQ(len(chart.GetAnnexs()), 1)
	test.AssertEQ(len(chart.GetImports()), 1)

	t.Logf("Investigating loaded chart...")
	{
		gotIDs := []string{}
		wantIDs := []string{ "i0", "i0/a0" }
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("Changing import target to 'nonexistent'...")
	wc3 := iobj.SetLink("nonexistent")
	test.PrintAllWarnings(wc3)

	t.Logf("Investigating loaded chart...")
	test.AssertWarnContains(wc3, "Cannot read an object")
	{
		gotIDs := []string{}
		wantIDs := []string{ "i0" }
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("Changing import target to 'single_annex.ltc'...")
	wc4 := iobj.SetLink("./single_annex.ltc")
	test.PrintAllWarnings(wc4)

	test.AssertEQ(len(chart.GetEvents()), 0)
	test.AssertEQ(len(chart.GetAnnexs()), 1)
	test.AssertEQ(len(chart.GetImports()), 1)

	t.Logf("Investigating loaded chart...")
	{
		gotIDs := []string{}
		wantIDs := []string{ "i0", "i0/a0" }
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("Removing import...")
	chart.RemoveObject(iobj)

	test.AssertEQ(len(chart.GetEvents()), 0)
	test.AssertEQ(len(chart.GetAnnexs()), 0)
	test.AssertEQ(len(chart.GetImports()), 0)

	t.Logf("Investigating loaded chart...")
	{
		gotIDs := []string{}
		wantIDs := []string{}
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("End test.")
}

func TestNestedImport_Import(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'nested_import.ltc'...")
	chart, wc := CreateChart("nested_import.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded chart...")
	{
		gotIDs := []string{}
		wantIDs := []string{ "e0", "e1", "i0", "i0/i0", "i0/i0/e0" }
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("Adding import './single_import.ltc'...")
	i1 := CreateEmptyImport()
	i1.SetLink("./single_import.ltc")
	chart.AddObject(i1)
	{
		gotIDs := []string{}
		wantIDs := []string{ "e0", "e1", "i0", "i0/i0", "i0/i0/e0", "i1", "i1/i0", "i1/i0/e0" }
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("Adding import './nested_import.ltc'...")
	i2 := CreateEmptyImport()
	i2.SetLink("./nested_import.ltc")
	chart.AddObject(i2)
	{
		gotIDs := []string{}
		wantIDs := []string{ "e0", "e1", "i0", "i0/i0", "i0/i0/e0", "i1", "i1/i0", "i1/i0/e0", "i2", "i2/e0", "i2/e1", "i2/i0", "i2/i0/i0", "i2/i0/i0/e0" }
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("Changing import from './nested_import.ltc' to './single_event.ltc...")
	i2.SetLink("./single_event.ltc")
	{
		gotIDs := []string{}
		wantIDs := []string{ "e0", "e1", "i0", "i0/i0", "i0/i0/e0", "i1", "i1/i0", "i1/i0/e0", "i2", "i2/e0" }
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("Removing import './single_import.ltc'...")
	chart.RemoveObject(i1)
	{
		gotIDs := []string{}
		wantIDs := []string{ "e0", "e1", "i0", "i0/i0", "i0/i0/e0", "i2", "i2/e0" }
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("End test.")
}

func TestLinkedObjectsAcrossImport_Import(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'linked_objs_2.ltc'...")
	chart, wc := CreateChart("linked_objs_2.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded chart...")
	e0 := chart.GetObjectByQualifiedID("e0").(*Event)
	test.AssertEQ(e0.ContinuedFrom.Get()[0].GetTitle(), "Event1")
	test.AssertEQ(e0.ContinuedFrom.Get()[0].GetQualifiedID(), "i0/e1")

	a0 := chart.GetObjectByQualifiedID("a0").(*Annex)
	test.AssertEQ(a0.AttachTo.Get()[0].(*Import).GetLink(), "./single_event.ltc")
	test.AssertEQ(a0.AttachTo.Get()[0].GetQualifiedID(), "i0/i0/i0")

	a1 := chart.GetObjectByQualifiedID("a1").(*Annex)
	test.AssertEQ(a1.ExtraNoteOf.Get()[0].(*Event).GetTitle(), "Event0")
	test.AssertEQ(a1.ExtraNoteOf.Get()[1].(*Event).GetTitle(), "Wherever")
	test.AssertEQ(a1.ExtraNoteOf.Get()[1].GetQualifiedID(), "i0/i0/i0/e0")

	t.Logf("Adding import './single_event.ltc'...")
	i1 := CreateEmptyImport()
	i1.SetLink("./single_event.ltc")
	chart.AddObject(i1)
	test.AssertEQ(a0.AttachTo.Get()[1].IsUnresolved(), false)
	test.AssertEQ(a0.AttachTo.Get()[1].(*Event).GetTitle(), "Wherever")

	t.Logf("Creating link to import (then removing and re-adding import)...")
	a1.AttachTo.Link(chart.GetObjectByQualifiedID("i1/e0"), true)
	test.AssertEQ(len(a1.AttachTo.Get()), 1)
	test.AssertEQ(a1.AttachTo.Get()[0].(*Event).GetTitle(), "Wherever")

	chart.RemoveObject(i1)
	test.AssertEQ(len(a1.AttachTo.Get()), 1)
	test.AssertEQ(a1.AttachTo.Get()[0].IsUnresolved(), true)
	test.AssertEQ(a1.AttachTo.Get()[0].GetQualifiedID(), "i1/e0")

	i1 = CreateEmptyImport()
	i1.SetLink("./single_event.ltc")
	chart.AddObject(i1)
	test.AssertEQ(a1.AttachTo.Get()[0].(*Event).GetTitle(), "Wherever")

	t.Logf("End test.")
}

func TestImportStartEndDates_Import(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'start_end_import.ltc'...")
	chart, wc := CreateChart("start_end_import.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded chart...")
	test.AssertWarnContains(wc, "inverted")
	wantIDs := []string{}
 	wantIDs = append(wantIDs, "i0", "i0/a0", "i0/e0", "i0/a1", "i0/e1")
	wantIDs = append(wantIDs, "i1", "i1/a0", "i1/e0", "i1/a1", "i1/e1")
 	wantIDs = append(wantIDs, "i2", "i2/a1", "i2/e1") 
 	wantIDs = append(wantIDs, "i3", "i3/a1", "i3/e1") 
 	wantIDs = append(wantIDs, "i4") 
 	wantIDs = append(wantIDs, "i5", "i5/a0", "i5/e0", "i5/a1", "i5/e1")
 	wantIDs = append(wantIDs, "i6", "i6/a0", "i6/e0", "i6/a1", "i6/e1")
 	wantIDs = append(wantIDs, "i7", "i7/a0", "i7/e0", "i7/a1", "i7/e1") 
 	wantIDs = append(wantIDs, "i8", "i8/a0", "i8/a1", "i8/e0") 
 	wantIDs = append(wantIDs, "i9", "i9/a0", "i9/e0") 
 	wantIDs = append(wantIDs, "i10") 
 	wantIDs = append(wantIDs, "i11", "i11/a0", "i11/e0", "i11/a1", "i11/e1") 
 	wantIDs = append(wantIDs, "i12", "i12/a0", "i12/e0", "i12/a1", "i12/e1") 
 	wantIDs = append(wantIDs, "i13", "i13/a1") 
 	wantIDs = append(wantIDs, "i14", "i14/a0", "i14/a1", "i14/e0") 
 	wantIDs = append(wantIDs, "i15", "i15/a0", "i15/a1", "i15/e0") 
	{
		gotIDs := []string{}
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		slices.Sort(gotIDs)
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("Adjusting import start/end...")
	i14 := chart.GetObjectByQualifiedID("i14").(*Import)
	i14.SetEndDate(TimeAt(2026, 1, 11))
	idx := slices.Index(wantIDs, "i14/a1")
	wantIDs = append(wantIDs[:idx], wantIDs[idx+1:]...)
	{
		gotIDs := []string{}
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		slices.Sort(gotIDs)
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	i13 := chart.GetObjectByQualifiedID("i13").(*Import)
	i13.SetEndDate(TimeAt(2026, 1, 10))
	wantIDs = append(wantIDs, "i13/e0", "i13/a0")
	{
		gotIDs := []string{}
		for _, iobj := range chart.GetAllMajors(false) {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		slices.Sort(gotIDs)
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("End test.")
}

func TestImportOffsetDate_Import(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'offset_import.ltc'...")
	chart, wc := CreateChart("offset_import.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded chart...")
	test.AssertEQ(chart.GetObjectByQualifiedID("i0/e0").(*Event).GetStartDate().String(), "2026-02-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i0/e0").(*Event).GetEndDate().String(), "2026-02-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i0/e1").(*Event).GetStartDate().String(), "2026-02-12")
	test.AssertEQ(chart.GetObjectByQualifiedID("i0/e1").(*Event).GetEndDate().String(), "2026-02-13")
	test.AssertEQ(chart.GetObjectByQualifiedID("i0/a0").(*Annex).GetCreateDate().String(), "2026-02-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i0/a1").(*Annex).GetCreateDate().String(), "2026-02-12")

	test.AssertEQ(chart.GetObjectByQualifiedID("i1/e0").(*Event).GetStartDate().String(), "2026-01-10 01:__:__")
	test.AssertEQ(chart.GetObjectByQualifiedID("i1/e0").(*Event).GetEndDate().String(), "2026-01-11 00:**:**")
	test.AssertEQ(chart.GetObjectByQualifiedID("i1/e1").(*Event).GetStartDate().String(), "2026-01-12 01:__:__")
	test.AssertEQ(chart.GetObjectByQualifiedID("i1/e1").(*Event).GetEndDate().String(), "2026-01-14 00:**:**")
	test.AssertEQ(chart.GetObjectByQualifiedID("i1/a0").(*Annex).GetCreateDate().String(), "2026-01-10 01:??:??")
	test.AssertEQ(chart.GetObjectByQualifiedID("i1/a1").(*Annex).GetCreateDate().String(), "2026-01-12 01:??:??")

	test.AssertEQ(chart.GetObjectByQualifiedID("i2/e0").(*Event).GetStartDate().String(), "2026-01-11")
	test.AssertEQ(chart.GetObjectByQualifiedID("i2/e0").(*Event).GetEndDate().String(), "2026-01-11")
	test.AssertEQ(chart.GetObjectByQualifiedID("i2/e1").(*Event).GetStartDate().String(), "2026-01-13")
	test.AssertEQ(chart.GetObjectByQualifiedID("i2/e1").(*Event).GetEndDate().String(), "2026-01-14")
	test.AssertEQ(chart.GetObjectByQualifiedID("i2/a0").(*Annex).GetCreateDate().String(), "2026-01-11")
	test.AssertEQ(chart.GetObjectByQualifiedID("i2/a1").(*Annex).GetCreateDate().String(), "2026-01-13")

	test.AssertEQ(chart.GetObjectByQualifiedID("i3/e0").(*Event).GetStartDate().String(), "2026-01-11 01:__:__")
	test.AssertEQ(chart.GetObjectByQualifiedID("i3/e0").(*Event).GetEndDate().String(), "2026-01-12 00:**:**")
	test.AssertEQ(chart.GetObjectByQualifiedID("i3/e1").(*Event).GetStartDate().String(), "2026-01-13 01:__:__")
	test.AssertEQ(chart.GetObjectByQualifiedID("i3/e1").(*Event).GetEndDate().String(), "2026-01-15 00:**:**")
	test.AssertEQ(chart.GetObjectByQualifiedID("i3/a0").(*Annex).GetCreateDate().String(), "2026-01-11 01:??:??")
	test.AssertEQ(chart.GetObjectByQualifiedID("i3/a1").(*Annex).GetCreateDate().String(), "2026-01-13 01:??:??")

	test.AssertEQ(chart.GetObjectByQualifiedID("i4/e0").(*Event).GetStartDate().String(), "2026-01-09 23:__:__")
	test.AssertEQ(chart.GetObjectByQualifiedID("i4/e0").(*Event).GetEndDate().String(), "2026-01-10 22:**:**")
	test.AssertEQ(chart.GetObjectByQualifiedID("i4/e1").(*Event).GetStartDate().String(), "2026-01-11 23:__:__")
	test.AssertEQ(chart.GetObjectByQualifiedID("i4/e1").(*Event).GetEndDate().String(), "2026-01-13 22:**:**")
	test.AssertEQ(chart.GetObjectByQualifiedID("i4/a0").(*Annex).GetCreateDate().String(), "2026-01-09 23:??:??")
	test.AssertEQ(chart.GetObjectByQualifiedID("i4/a1").(*Annex).GetCreateDate().String(), "2026-01-11 23:??:??")

	test.AssertEQ(chart.GetObjectByQualifiedID("i5/e0").(*Event).GetStartDate().String(), "2026-01-09 23:__:__")
	test.AssertEQ(chart.GetObjectByQualifiedID("i5/e0").(*Event).GetEndDate().String(), "2026-01-10 22:**:**")
	test.AssertEQ(chart.GetObjectByQualifiedID("i5/e1").(*Event).GetStartDate().String(), "2026-01-11 23:__:__")
	test.AssertEQ(chart.GetObjectByQualifiedID("i5/e1").(*Event).GetEndDate().String(), "2026-01-13 22:**:**")
	test.AssertEQ(chart.GetObjectByQualifiedID("i5/a0").(*Annex).GetCreateDate().String(), "2026-01-09 23:??:??")
	test.AssertEQ(chart.GetObjectByQualifiedID("i5/a1").(*Annex).GetCreateDate().String(), "2026-01-11 23:??:??")

	test.AssertWarnContains(wc, "wrong offset date.*\\(unknown\\)")
	test.AssertEQ(chart.GetObjectByQualifiedID("i6/e0").(*Event).GetStartDate().String(), "2026-01-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i6/e0").(*Event).GetEndDate().String(), "2026-01-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i6/e1").(*Event).GetStartDate().String(), "2026-01-12")
	test.AssertEQ(chart.GetObjectByQualifiedID("i6/e1").(*Event).GetEndDate().String(), "2026-01-13")
	test.AssertEQ(chart.GetObjectByQualifiedID("i6/a0").(*Annex).GetCreateDate().String(), "2026-01-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i6/a1").(*Annex).GetCreateDate().String(), "2026-01-12")

	test.AssertEQ(chart.GetObjectByQualifiedID("i7/e0").(*Event).GetStartDate().String(), "2026-01-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i7/e0").(*Event).GetEndDate().String(), "2026-01-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i7/e1").(*Event).GetStartDate().String(), "2026-01-12")
	test.AssertEQ(chart.GetObjectByQualifiedID("i7/e1").(*Event).GetEndDate().String(), "2026-01-13")
	test.AssertEQ(chart.GetObjectByQualifiedID("i7/a0").(*Annex).GetCreateDate().String(), "2026-01-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i7/a1").(*Annex).GetCreateDate().String(), "2026-01-12")

	test.AssertWarnContains(wc, "wrong offset date.*\\(inf\\)")
	test.AssertEQ(chart.GetObjectByQualifiedID("i8/e0").(*Event).GetStartDate().String(), "2026-01-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i8/e0").(*Event).GetEndDate().String(), "2026-01-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i8/e1").(*Event).GetStartDate().String(), "2026-01-12")
	test.AssertEQ(chart.GetObjectByQualifiedID("i8/e1").(*Event).GetEndDate().String(), "2026-01-13")
	test.AssertEQ(chart.GetObjectByQualifiedID("i8/a0").(*Annex).GetCreateDate().String(), "2026-01-10")
	test.AssertEQ(chart.GetObjectByQualifiedID("i8/a1").(*Annex).GetCreateDate().String(), "2026-01-12")

	t.Logf("End test.")
}
