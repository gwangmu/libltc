package ltc 

import (
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
	_, wc:= CreateChart("broken_event.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	// TODO: check duplicate warnings.
	// TODO: check duplicate ID.
	// TODO: check 'EmbedChart' is nil.
	// TODO: check date orders.

	t.Logf("End test.")
}

// TODO: check loading annex.
// TODO: check auto-resolving on load (+ new links).
