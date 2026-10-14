package libltc 

import (
	"slices"
	"testing"
	"github.com/gwangmu/libltc/internal/test"
)

func TestSingleEvent_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'single_event.ltc'...")
	chart, wc := CreateChart("single_event.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	event := chart.GetObjectByIntrinsicID("e0").(*Event)
	test.AssertEQ(event, chart.GetObjectByQualifiedID("e0").(*Event))

	test.AssertEQ(event.GetTitle(), "Wherever")
	test.AssertEQ(event.GetCategory(), "Residence")
	test.AssertEQ(event.GetStartDate().String(), "1991-03-29 (Attr0)")
	test.AssertEQ(event.GetEndDate().String(), "1991-04-** (Attr1)")
	test.AssertEQ(event.HasAttrKey("Attr2"), true)
	test.AssertEQ(event.GetAttrs("Attr2")[0], "")
	test.AssertEQ(event.HasAttrKey("Attr3"), true)
	test.AssertEQ(event.GetAttrs("Attr3")[0], "")
	test.AssertEQ(event.HasAttrKey("Attr4"), true)
	test.AssertEQ(event.GetAttrs("Attr4")[0], "Val")
	test.AssertEQ(event.HasAttrKey("Attr5"), true)
	test.AssertEQ(event.GetAttrs("Attr5")[0], "Val0")
	test.AssertEQ(event.GetAttrs("Attr5")[1], "Val1")
	test.AssertEQ(event.Note.Title, "")
	test.AssertEQ(event.Note.Snippets[0].Time.String(), "(unknown)")
	test.AssertEQ(event.Note.Snippets[0].Text, "Some note")

	t.Logf("End test.")
}

func TestDuplicateIDEvents_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'dup_id_events.ltc'...")
	chart, wc := CreateChart("dup_id_events.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	test.AssertWarnContains(wc, "duplicate")
	test.AssertEQ(chart.GetEventsInCategory("")[0].GetTitle(), "Event0")
	test.AssertEQ(chart.GetEventsInCategory("")[0].GetIntrinsicID(), "e1")
	test.AssertEQ(chart.GetEventsInCategory("")[0].GetAttrs("OldID")[0], "e0")
	test.AssertEQ(chart.GetEventsInCategory("")[1].GetTitle(), "Event1")
	test.AssertEQ(chart.GetEventsInCategory("")[1].GetIntrinsicID(), "e2")
	test.AssertEQ(chart.GetEventsInCategory("")[1].GetAttrs("OldID")[0], "e0")
	test.AssertEQ(chart.GetEventsInCategory("")[2].GetTitle(), "Event2")
	test.AssertEQ(chart.GetEventsInCategory("")[2].GetIntrinsicID(), "e0")

	t.Logf("End test.")
}

func TestBrokenEvents_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'broken_event.ltc'...")
	chart, wc:= CreateChart("broken_event.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	test.AssertWarnContains(wc, "Broken ID.*non-standard.*ID")
	test.AssertWarnContains(wc, "Broken ID.*OldID:borked.*corrected")
	test.AssertWarnContains(wc, "Infinite past date.*OldID:e1.*corrected")
	test.AssertWarnContains(wc, "Infinite future date.*OldID:e1.*corrected")
	test.AssertWarnContains(wc, "Explicit infinite.*OldID:e1.*corrected")
	test.AssertWarnContains(wc, "Empty 'ContinuedFrom'.*OldID:e1.*corrected")
	test.AssertWarnContains(wc, "Swapped dates \\(1\\).*OldID:e1.*corrected")
	test.AssertWarnContains(wc, "Ambiguous.*\\(1\\).*OldID:e1.*corrected")
	test.AssertWarnContains(wc, "Ambiguous.*\\(2\\).*OldID:e1.*corrected")
	test.AssertWarnContains(wc, "Ambiguous.*\\(3\\).*OldID:e1.*corrected")
	test.AssertWarnContains(wc, "Broken date \\(2\\).*OldID:e1.*corrected")
	test.AssertEQ(chart.GetObjectByQualifiedID("e0").(*Event).IsChartEmbedding(), false)

	// Ensure no duplicate IDs.
	aevIDs := []string{}
	for _, ev := range chart.GetEventsInCategory("") {
		aevIDs = append(aevIDs, ev.GetIntrinsicID())
	}
	slices.Sort(aevIDs)
	naevIDs := slices.Compact(aevIDs)
	test.AssertEQ(len(aevIDs), len(naevIDs))

	t.Logf("End test.")
}

func TestBrokenIDAnnex_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'broken_id_annex.ltc'...")
	chart, wc:= CreateChart("broken_id_annex.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	test.AssertWarnContains(wc, "Broken ID.*non-standard.*ID")
	test.AssertWarnContains(wc, "Broken ID.*OldID:borked.*corrected")
	test.AssertWarnContains(wc, "Wrong ID.*non-standard.*ID")
	test.AssertWarnContains(wc, "Wrong ID.*OldID:e0.*corrected")
	test.AssertWarnContains(wc, "Annex1.*OldID:a1.*corrected")
	test.AssertWarnContains(wc, "Annex2.*OldID:a1.*corrected")
	test.AssertWarnContains(wc, "Annex3.*OldID:a1.*corrected")

	// Ensure no duplicate IDs.
	aaIDs := []string{}
	for _, a := range chart.GetAnnexs() {
		aaIDs = append(aaIDs, a.GetIntrinsicID())
	}
	slices.Sort(aaIDs)
	naaIDs := slices.Compact(aaIDs)
	test.AssertEQ(len(aaIDs), len(naaIDs))

	t.Logf("End test.")
}

func TestLinkedObjects_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'linked_objs.ltc'...")
	chart, wc := CreateChart("linked_objs.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	test.AssertWarnContains(wc, "Event1.*dangling 'ContinuedFrom'")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("e0").(*Event).ContinuedFrom.Get()), 1)
	test.AssertEQ(chart.GetObjectByQualifiedID("e0").(*Event).ContinuedFrom.Get()[0].GetTitle(), "Event1")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("e1").(*Event).ContinuedTo.Get()), 1)
	test.AssertEQ(chart.GetObjectByQualifiedID("e1").(*Event).ContinuedTo.Get()[0].GetTitle(), "Event0")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("e0").(*Event).Attached.Get()), 1)
	test.AssertEQ(chart.GetObjectByQualifiedID("e0").(*Event).Attached.Get()[0].Title, "Annex0")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("e0").(*Event).ExtraNote.Get()), 1)
	test.AssertEQ(chart.GetObjectByQualifiedID("e0").(*Event).ExtraNote.Get()[0].Title, "Annex1")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("a0").(*Annex).AttachTo.Get()), 1)
	test.AssertEQ(chart.GetObjectByQualifiedID("a0").(*Annex).AttachTo.Get()[0].(*Event).GetTitle(), "Event0")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("a1").(*Annex).ExtraNoteOf.Get()), 2)
	test.AssertEQ(chart.GetObjectByQualifiedID("a1").(*Annex).ExtraNoteOf.Get()[0].(*Event).GetTitle(), "Event0")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("e1").(*Event).ContinuedFrom.Get()), 1)
	test.AssertEQ(chart.GetObjectByQualifiedID("e1").(*Event).ContinuedFrom.Get()[0].IsUnresolved(), true) 
	test.AssertEQ(chart.GetObjectByQualifiedID("e1").(*Event).ContinuedFrom.Get()[0].GetQualifiedID(), "e4") 

	// Mimic runtime chart manipulation
	// E.g., add event/annex, create/remove links, change dates, ...
	t.Logf("Mimicking runtime chart manipulation...")
	nev := CreateEmptyEvent()
	test.AssertEQ(nev.GetIntrinsicID(), "?")

	// After adding it to the chart, it should assign the ID.
	if err := chart.AddObject(nev); err != nil {
		t.Fatal(err)
		return
	}
	test.AssertEQ(nev.GetIntrinsicID(), "e2")

	// Setting a new title
	nev.SetTitle("New event")
	test.AssertEQ(nev.GetTitle(), "New event")

	// Removing and (re-)adding the event.
	chart.RemoveObject(nev)
	test.AssertEQ(nev.GetIntrinsicID(), "?")
	test.AssertEQ(nev.GetTitle(), "New event")
	if err := chart.AddObject(nev); err != nil {
		t.Fatal(err)
		return
	}
	test.AssertEQ(nev.GetIntrinsicID(), "e2")
	test.AssertEQ(nev.GetTitle(), "New event")

	// Assigning a start date.
	chart.RemoveObject(nev)
	nev.SetStartDate(TimeAt(2026, 1, 11))
	if err := chart.AddObject(nev); err != nil {
		t.Fatal(err)
		return
	}
	eidsInit := []string{"e0", "e2", "e1"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eidsInit[i])
	}

	// Changing a start date.
	nev.SetStartDate(TimeAt(2026, 1, 14))
	eidsC0 := []string{"e0", "e1", "e2"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eidsC0[i])
	}

	nev.SetStartDate(UnknownTime)
	eidsC1 := []string{"e2", "e0", "e1"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eidsC1[i])
	}

	nev.SetStartDate(InfiniteFuture)
	eidsC2 := []string{"e0", "e1", "e2"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eidsC2[i])
	}

	nev.SetEndDate(InfinitePast)
	eidsCx := []string{"e2", "e0", "e1"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eidsCx[i])
	}
	test.AssertEQ(nev.GetStartDate().String(), "(-inf)")
	test.AssertEQ(nev.GetEndDate().String(), "(inf)")

	nev.SetStartDate(TimeAt(2026, 1, 11))
	eidsCy := []string{"e0", "e2", "e1"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eidsCy[i])
	}
	test.AssertEQ(nev.GetStartDate().String(), "2026-01-11")
	test.AssertEQ(nev.GetEndDate().String(), "(inf)")

	nev.SetEndDate(TimeAt(2026, 1, 8))
	eidsCz := []string{"e2", "e0", "e1"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eidsCz[i])
	}
	test.AssertEQ(nev.GetStartDate().String(), "2026-01-08")
	test.AssertEQ(nev.GetEndDate().String(), "2026-01-11")

	nev.SetEndDate(TimeAt(2026, 1, 8))
	eidsCw := []string{"e2", "e0", "e1"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eidsCw[i])
	}
	test.AssertEQ(nev.GetStartDate().String(), "2026-01-08")
	test.AssertEQ(nev.GetEndDate().String(), "2026-01-08")
	test.AssertEQ(TimeAt(2026, 1, 8, 0, 0, 0).IsAfter(nev.GetStartDate()), true)
	test.AssertEQ(nev.GetEndDate().IsAfter(TimeAt(2026, 1, 8, 23, 59, 59)), true)

	nev.SetStartDate(InfinitePast)
	eidsC3 := []string{"e2", "e0", "e1"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eidsC3[i])
	}

	// Event starting from infinite future.
	nev2 := CreateEmptyEvent()
	nev2.SetStartDate(InfiniteFuture)
	if err := chart.AddObject(nev2); err != nil {
		t.Fatal(err)
		return
	}
	test.AssertEQ(nev2.GetStartDate().String(), "(inf)")
	test.AssertEQ(nev2.GetEndDate().String(), "(inf)")

	eids2 := []string{"e2", "e0", "e1", "e3"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eids2[i])
	}

	// Move category.
	e1 := chart.GetObjectByQualifiedID("e1").(*Event)
	e1.SetCategory("Some Category")
	eids3 := []string{"e2", "e0", "e3"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eids3[i])
	}
	test.AssertEQ(chart.GetEventsInCategory("Some Category")[0].GetTitle(), "Event1")

	chart.AddObject(e1)
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eids3[i])
	}
	test.AssertEQ(chart.GetEventsInCategory("Some Category")[0].GetTitle(), "Event1")

	e1.SetStartDate(TimeAt(2026, 1, 5))
	e1.SetCategory("")
	eids4 := []string{"e2", "e1", "e0", "e3"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eids4[i])
	}

	chart.AddObject(e1)
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eids4[i])
	}

	e1.SetCategory("Another Category")
	eids5 := []string{"e2", "e0", "e3"}
	for i, eobj := range chart.GetEventsInCategory("") {
		test.AssertEQ(eobj.GetIntrinsicID(), eids5[i])
	}

	// Create ContinuedFrom.
	nev3 := CreateEmptyEvent()
	nev3.SetTitle("New event 2")
	if err := chart.AddObject(nev3); err != nil {
		t.Fatal(err)
		return
	}
	test.AssertEQ(nev3.GetIntrinsicID(), "e4")
	test.AssertEQ(chart.GetObjectByQualifiedID("e1").(*Event).ContinuedFrom.Get()[0].GetTitle(), "New event 2")
	test.AssertEQ(nev3.ContinuedTo.Get()[0].GetTitle(), "Event1")

	e0 := chart.GetObjectByQualifiedID("e0").(*Event)
	nev3.ContinuedFrom.Link(e0, true)
	test.AssertEQ(nev3.ContinuedFrom.Get()[0].GetTitle(), "Event0")
	test.AssertEQ(e0.ContinuedTo.Get()[0].GetTitle(), "New event 2")

	e1.ContinuedFrom.Unlink(nev3)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(nev3.ContinuedTo.Get()), 0)
	test.AssertEQ(e1.ContinuedFrom.Get()[0].IsUnresolved(), true)

	e1.ContinuedFrom.Link(e0, false)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(nev3.ContinuedTo.Get()), 0)
	test.AssertEQ(e1.ContinuedFrom.Get()[0].IsUnresolved(), true)

	e1.ContinuedFrom.Link(nev3, false)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(nev3.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Get()[0].GetTitle(), "New event 2")

	e1.ContinuedFrom.Unreserve("e4")
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 0)
	test.AssertEQ(len(nev3.ContinuedTo.Get()), 0)

	// Create AttachTo.
	a0 := chart.GetObjectByQualifiedID("a0").(*Annex)
	a0.AttachTo.Link(nev3, true)
	test.AssertEQ(len(a0.AttachTo.Get()), 2)
	test.AssertEQ(len(nev3.Attached.Get()), 1)
	test.AssertEQ(nev3.Attached.Get()[0].Title, "Annex0")
	test.AssertEQ(a0.AttachTo.Get()[1].(*Event).GetTitle(), "New event 2")

	a0.AttachTo.Unreserve("e0")
	test.AssertEQ(len(a0.AttachTo.Get()), 1)
	test.AssertEQ(len(e0.Attached.Get()), 0)

	// Create ExtraNoteOf.
	a1 := chart.GetObjectByQualifiedID("a1").(*Annex)
	nax1 := CreateEmptyAnnex()
	nax1.Title = "New annex"
	if err := chart.AddObject(nax1); err != nil {
		t.Fatal(err)
		return
	}
	test.AssertEQ(nax1.GetIntrinsicID(), "a2")
	test.AssertEQ(len(a1.ExtraNoteOf.Get()), 2)
	test.AssertEQ(e0.ExtraNote.Get()[0].Title, "Annex1")
	test.AssertEQ(a1.ExtraNoteOf.Get()[0].(*Event).GetTitle(), "Event0")
	test.AssertEQ(a1.ExtraNoteOf.Get()[1].(*Annex).Title, "New annex")
	test.AssertEQ(nax1.ExtraNote.Get()[0].Title, "Annex1")

	chart.RemoveObject(nax1)
	test.AssertEQ(a1.ExtraNoteOf.Get()[1].IsUnresolved(), true)

	t.Logf("End test.")
}

func TestSimpleSubchartEmbedding_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'subchart_embed_simple.ltc'...")
	chart, wc := CreateChart("subchart_embed_simple.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Loading single_event subchart...")
	sc1, wsc1 := chart.GetObjectByQualifiedID("e2").(*Event).GetSubchart()
	test.PrintAllWarnings(wsc1)

	t.Logf("Investigating loaded subchart...")
	test.AssertEQ(sc1.GetObjectByIntrinsicID("e0").(*Event).GetTitle(), "Wherever")
	test.AssertEQ(chart.GetObjectByQualifiedID("e2/e0").(*Event).GetTitle(), "Wherever")

	e2e0 := chart.GetObjectByQualifiedID("e2/e0").(*Event)
	test.AssertEQ(e2e0.IsUnresolved(), false)
	test.AssertEQ(e2e0.GetIntrinsicID(), "e0")
	test.AssertEQ(e2e0.GetQualifiedID(), "e2/e0")
	test.AssertEQ(e2e0.GetRelQualifiedID(sc1), "e0")
	test.AssertEQ(e2e0.GetRelQualifiedID(chart), "e2/e0")

	t.Logf("Loading single_annex subchart...")
	sc2, wsc2 := chart.GetObjectByQualifiedID("e3").(*Event).GetSubchart()
	test.PrintAllWarnings(wsc2)

	t.Logf("Investigating loaded subchart...")
	test.AssertEQ(sc2.GetObjectByIntrinsicID("a0").(*Annex).Title, "Whatever")
	test.AssertEQ(chart.GetObjectByQualifiedID("e3/a0").(*Annex).Title, "Whatever")

	e3a0 := chart.GetObjectByQualifiedID("e3/a0").(*Annex)
	test.AssertEQ(e3a0.IsUnresolved(), false)
	test.AssertEQ(e3a0.GetIntrinsicID(), "a0")
	test.AssertEQ(e3a0.GetQualifiedID(), "e3/a0")
	test.AssertEQ(e3a0.GetRelQualifiedID(sc2), "a0")
	test.AssertEQ(e3a0.GetRelQualifiedID(chart), "e3/a0")

	t.Logf("End test.")
}

func TestComplexSubchartEmbedding_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'subchart_embed_complex.ltc'...")
	chart, wc := CreateChart("subchart_embed_complex.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Loading subchart_embed_simple subchart...")
	sc1, wsc1 := chart.GetObjectByQualifiedID("e0").(*Event).GetSubchart()
	ssc1, wssc1 := chart.GetObjectByQualifiedID("e0/e2").(*Event).GetSubchart()
	ssc2, wssc2 := chart.GetObjectByQualifiedID("e0/e3").(*Event).GetSubchart()
	test.PrintAllWarnings(wsc1)
	test.PrintAllWarnings(wssc1)
	test.PrintAllWarnings(wssc2)

	t.Logf("Investigating loaded subchart...")
	test.AssertEQ(ssc1.GetObjectByIntrinsicID("e0").(*Event).GetTitle(), "Wherever")
	test.AssertEQ(sc1.GetObjectByQualifiedID("e2/e0").(*Event).GetTitle(), "Wherever")
	test.AssertEQ(chart.GetObjectByQualifiedID("e0/e2/e0").(*Event).GetTitle(), "Wherever")

	e0e2e0 := chart.GetObjectByQualifiedID("e0/e2/e0").(*Event)
	test.AssertEQ(e0e2e0.IsUnresolved(), false) 
	test.AssertEQ(e0e2e0.GetIntrinsicID(), "e0")
	test.AssertEQ(e0e2e0.GetQualifiedID(), "e0/e2/e0")
	test.AssertEQ(e0e2e0.GetRelQualifiedID(ssc1), "e0")
	test.AssertEQ(e0e2e0.GetRelQualifiedID(sc1), "e2/e0")
	test.AssertEQ(e0e2e0.GetRelQualifiedID(chart), "e0/e2/e0")

	test.AssertEQ(ssc2.GetObjectByIntrinsicID("a0").(*Annex).Title, "Whatever")
	test.AssertEQ(sc1.GetObjectByQualifiedID("e3/a0").(*Annex).Title, "Whatever")
	test.AssertEQ(chart.GetObjectByQualifiedID("e0/e3/a0").(*Annex).Title, "Whatever")

	e0e3a0 := chart.GetObjectByQualifiedID("e0/e3/a0").(*Annex)
	test.AssertEQ(e0e3a0.IsUnresolved(), false) 
	test.AssertEQ(e0e3a0.GetIntrinsicID(), "a0")
	test.AssertEQ(e0e3a0.GetQualifiedID(), "e0/e3/a0")
	test.AssertEQ(e0e3a0.GetRelQualifiedID(ssc2), "a0")
	test.AssertEQ(e0e3a0.GetRelQualifiedID(sc1), "e3/a0")
	test.AssertEQ(e0e3a0.GetRelQualifiedID(chart), "e0/e3/a0")

	test.AssertEQ(chart.GetObjectByIntrinsicID("a0").(*Annex).AttachTo.Get()[0].IsUnresolved(), false)
	test.AssertEQ(chart.GetObjectByIntrinsicID("a0").(*Annex).AttachTo.Get()[0].(*Event).GetTitle(), "Dummy1")
	
	t.Logf("Loading linked_objs subchart...")
	sc2, wsc2 := chart.GetObjectByQualifiedID("e1").(*Event).GetSubchart()
	test.PrintAllWarnings(wsc1)

	t.Logf("Investigating loaded subchart...")
	test.AssertWarnContains(wsc2, "Event1.*dangling 'ContinuedFrom'")
	test.AssertEQ(len(sc2.GetObjectByIntrinsicID("e0").(*Event).ContinuedFrom.Get()), 1)
	test.AssertEQ(sc2.GetObjectByIntrinsicID("e0").(*Event).ContinuedFrom.Get()[0].GetTitle(), "Event1")
	test.AssertEQ(len(sc2.GetObjectByQualifiedID("e1").(*Event).ContinuedTo.Get()), 1)
	test.AssertEQ(sc2.GetObjectByQualifiedID("e1").(*Event).ContinuedTo.Get()[0].GetTitle(), "Event0")
	test.AssertEQ(len(sc2.GetObjectByQualifiedID("e0").(*Event).Attached.Get()), 1)
	test.AssertEQ(sc2.GetObjectByQualifiedID("e0").(*Event).Attached.Get()[0].Title, "Annex0")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("e1/e0").(*Event).ExtraNote.Get()), 1)
	test.AssertEQ(chart.GetObjectByQualifiedID("e1/e0").(*Event).ExtraNote.Get()[0].Title, "Annex1")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("e1/a0").(*Annex).AttachTo.Get()), 1)
	test.AssertEQ(chart.GetObjectByQualifiedID("e1/a0").(*Annex).AttachTo.Get()[0].(*Event).GetTitle(), "Event0")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("e1/a1").(*Annex).ExtraNoteOf.Get()), 2)
	test.AssertEQ(chart.GetObjectByQualifiedID("e1/a1").(*Annex).ExtraNoteOf.Get()[0].(*Event).GetTitle(), "Event0")
	test.AssertEQ(len(chart.GetObjectByQualifiedID("e1/e1").(*Event).ContinuedFrom.Get()), 1)
	test.AssertEQ(sc2.GetObjectByQualifiedID("e1").(*Event).ContinuedFrom.Get()[0].IsUnresolved(), true) 
	test.AssertEQ(sc2.GetObjectByQualifiedID("e1").(*Event).ContinuedFrom.Get()[0].GetQualifiedID(), "e4") 

	nax1 := CreateEmptyAnnex()
	e1e1 := sc2.GetObjectByIntrinsicID("e1").(*Event)
	nax1.Title = "New annex"
	chart.AddObject(nax1)

	// Cross-subchart 'forward' refs should be possible.
	// NOTE: 'back' refs are not possible. (i.e., subchart objs refing the main chart objs)
	nax1.AttachTo.Link(e1e1, true)
	
	test.AssertEQ(len(nax1.AttachTo.Get()), 1)
	test.AssertEQ(nax1.AttachTo.Get()[0].(*Event).GetTitle(), "Event1")
	test.AssertEQ(len(e1e1.Attached.Get()), 1)
	test.AssertEQ(e1e1.Attached.Get()[0].Title, "New annex")

	// Unreserving an ID should erase the ref altogether. 
	nax1.AttachTo.Unreserve("e1/e1")

	test.AssertEQ(len(nax1.AttachTo.Get()), 0)
	test.AssertEQ(len(e1e1.Attached.Get()), 0)

	nax1.AttachTo.Link(e1e1, true)

	test.AssertEQ(len(nax1.AttachTo.Get()), 1)
	test.AssertEQ(len(e1e1.Attached.Get()), 1)

	// Removing an object 'unresolve' refs, not 'unreserve' them.
	chart.RemoveObject(nax1)

	test.AssertEQ(len(nax1.AttachTo.Get()), 1)
	test.AssertEQ(len(e1e1.Attached.Get()), 0)

	t.Logf("End test.")
}

func TestGetEventsBetween_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'linked_objs.ltc'...")
	chart, wc := CreateChart("linked_objs.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	nev := CreateEmptyEvent()

	// Setting just one date, while the other date is unknown, makes the other date
	// "floating," meaning the event has only one specific date. 
	nev.SetDate(TimeAt(2026, 1, 11))
	if err := chart.AddObject(nev); err != nil {
		t.Fatal(err)
		return
	}
	test.AssertEQ(nev.GetIntrinsicID(), "e2")

	{
		evs := chart.GetEventsInCategoryBetween("", TimeAt(2026, 1, 8), TimeAt(2026, 1, 14))
		eids := []string{"e0", "e2", "e1"}
		test.AssertEQ(len(evs), len(eids))
		for i, eobj := range evs {
			test.AssertEQ(eobj.GetIntrinsicID(), eids[i])
		}
	}

	{
		evs := chart.GetEventsInCategoryBetween("", TimeAt(2026, 1, 10), TimeAt(2026, 1, 14))
		eids := []string{"e0", "e2", "e1"}
		test.AssertEQ(len(evs), len(eids))
		for i, eobj := range evs {
			test.AssertEQ(eobj.GetIntrinsicID(), eids[i])
		}
	}

	{
		evs := chart.GetEventsInCategoryBetween("", TimeAt(2026, 1, 10), TimeAt(2026, 1, 12, 12))
		eids := []string{"e0", "e2"}
		test.AssertEQ(len(evs), len(eids))
		for i, eobj := range evs {
			test.AssertEQ(eobj.GetIntrinsicID(), eids[i])
		}
	}

	{
		evs := chart.GetEventsInCategoryBetween("", TimeAt(2026, 1, 10, 12), TimeAt(2026, 1, 12, 12))
		eids := []string{"e2"}
		test.AssertEQ(len(evs), len(eids))
		for i, eobj := range evs {
			test.AssertEQ(eobj.GetIntrinsicID(), eids[i])
		}
	}

	// This should be the same as 'UnsetLocalDate' or 'SetDate(UnknownTime)'.
	nev.SetEndDate(UnknownTime)
	test.AssertEQ(nev.GetStartDate().String(), "(unknown)")
	test.AssertEQ(nev.GetEndDate().String(), "(unknown)")

	nev.SetStartDate(TimeAt(2026, 1, 11))
	test.AssertEQ(nev.GetStartDate().String(), "2026-01-11")
	test.AssertEQ(nev.GetEndDate().String(), "(inf)")

	{
		evs := chart.GetEventsInCategoryBetween("", TimeAt(2026, 1, 10, 12), TimeAt(2026, 1, 12, 12))
		test.AssertEQ(len(evs), 0)
	}

	t.Logf("End test.")
}
