package ltc 

import (
	"slices"
	"testing"
	"libltc/internal/test"
)

func TestSingleEvent_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'single_event.ltc'...")
	chart, wc := CreateChart("single_event.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	event := chart.GetObjectsByIntrinsicID("e0")[0].(*Event)
	test.AssertEQ(event, chart.GetObjectByAbsQualID("e0").(*Event))
	test.AssertEQ(event, chart.GetObjectByRelQualID("e0", chart).(*Event))
	test.AssertEQ(event, chart.GetObjectByRelQualID("e0", nil).(*Event))

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
	test.AssertEQ(chart.GetObjectByAbsQualID("e0").(*Event).IsChartEmbedding(), false)

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
	chart, wc:= CreateChart("linked_objs.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	test.AssertWarnContains(wc, "Event1.*dangling 'ContinuedFrom'")
	test.AssertEQ(len(chart.GetObjectByAbsQualID("e0").(*Event).ContinuedFrom.Get()), 1)
	test.AssertEQ(chart.GetObjectByAbsQualID("e0").(*Event).ContinuedFrom.Get()[0].GetTitle(), "Event1")
	test.AssertEQ(len(chart.GetObjectByAbsQualID("e1").(*Event).ContinuedTo.Get()), 1)
	test.AssertEQ(chart.GetObjectByAbsQualID("e1").(*Event).ContinuedTo.Get()[0].GetTitle(), "Event0")
	test.AssertEQ(len(chart.GetObjectByAbsQualID("e0").(*Event).Attached.Get()), 1)
	test.AssertEQ(chart.GetObjectByAbsQualID("e0").(*Event).Attached.Get()[0].Title, "Annex0")
	test.AssertEQ(len(chart.GetObjectByAbsQualID("e0").(*Event).ExtraNote.Get()), 1)
	test.AssertEQ(chart.GetObjectByAbsQualID("e0").(*Event).ExtraNote.Get()[0].Title, "Annex1")

	t.Logf("End test.")
}

// TODO: mimic runtime chart manipulation (add event/annex, create/remove links, change dates, ...)
// TODO: check subchart functionality (load, change link, reference inside)
