package ltc 

import (
	"testing"
	"libltc/internal/file"
	"libltc/internal/test"
)

func TestSingleEvent_Event(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'single_event.ltc'...")
	file, wf, err := file.LoadFromURI[file.File]("single_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}
	test.PrintAllWarnings(wf)

	t.Logf("Loading event...")
	event, we := createEventFromParsed(file.Event[0])
	test.PrintAllWarnings(we)

	t.Logf("Investigating loaded event...")
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

	t.Log("End test.")
}

func TestContinuedEvent_Event(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'cont_events.ltc'...")
	file, wf, err := file.LoadFromURI[file.File]("cont_events.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}
	test.PrintAllWarnings(wf)

	t.Logf("Loading events...")
	e0, w0 := createEventFromParsed(file.Event[0])
	test.PrintAllWarnings(w0)
	e1, w1 := createEventFromParsed(file.Event[1])
	test.PrintAllWarnings(w1)

	t.Logf("Linking events...")

	// 'e0' should have been reserved at event creation.
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Get()[0].GetIntrinsicID(), "e0")
	test.AssertEQ(e1.ContinuedFrom.Get()[0].GetAbsQualifiedID(), "e0")
	test.AssertEQ(e1.ContinuedFrom.Get()[0].IsUnresolved(), true)
	test.AssertEQ(e1.ContinuedFrom.Get()[0].getNumberID(), NID_Invalid)

	// Creating and re-reserving the same object shouldn't change anything.
	e1.ContinuedFrom.Link(e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	e1.ContinuedFrom.Link(e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	e1.ContinuedFrom.Reserve("e0")
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	e1.ContinuedFrom.Unlink(e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 0)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), false)
	test.AssertEQ(e0.ContinuedTo.Has(e1), false)
	test.AssertEQ(e0.IsUnresolved(), false)
	e1.ContinuedFrom.Unreserve("e0")
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 0)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 0)

	// Re-creating should also be possible. 
	e1.ContinuedFrom.Link(e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	e1.ContinuedFrom.Unlink(e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 0)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), false)
	test.AssertEQ(e0.ContinuedTo.Has(e1), false)
	test.AssertEQ(e0.IsUnresolved(), false)
	e1.ContinuedFrom.Link(e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)

	t.Log("End test.")
}

func TestEventEmbeddingEvent_Event(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'event_embedding_event.ltc'...")
	file, wf, err := file.LoadFromURI[file.File]("event_embedding_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}
	test.PrintAllWarnings(wf)

	t.Logf("Loading event (all embed)...")
	e2, w2 := createEventFromParsed(file.Event[2])
	test.PrintAllWarnings(w2)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e2.IsEventEmbedding(), true)
	test.AssertEQ(e2.GetTitle(), "Wherever")
	test.AssertEQ(e2.GetStartDate().String(), "1991-03-29 (Attr0)")
	test.AssertEQ(e2.GetEndDate().String(), "1991-04-** (Attr1)")
	test.AssertEQ(e2.GetEmbeddedEvent().Note.Title, "")
	test.AssertEQ(e2.GetEmbeddedEvent().Note.Snippets[0].Time.String(), "(unknown)")
	test.AssertEQ(e2.GetEmbeddedEvent().Note.Snippets[0].Text, "Some note")

	t.Logf("Loading event (partial overriding)...")
	e3, w3 := createEventFromParsed(file.Event[3])
	test.PrintAllWarnings(w3)
	e4, w4 := createEventFromParsed(file.Event[4])
	test.PrintAllWarnings(w4)
	e5, w5 := createEventFromParsed(file.Event[5])
	test.PrintAllWarnings(w5)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e3.GetTitle(), "Overridden title")
	test.AssertEQ(e4.GetStartDate().String(), "0000-01-02")
	test.AssertEQ(e5.GetEndDate().String(), "9999-01-02")

	t.Log("End test.")
}

//func TestChartEmbeddingEvent_Event(t *testing.T) {
// TODO
//}

func TestBrokenEvent_Event(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'broken_event.ltc'...")
	file, wf, err := file.LoadFromURI[file.File]("broken_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}
	test.PrintAllWarnings(wf)

	t.Logf("Investigating loaded file...")
	test.AssertWarnContains(wf, "non-standard event ID")

	t.Logf("Loading event '%s'...", *file.Event[0].Title)
	e0, w0 := createEventFromParsed(file.Event[0])
	test.PrintAllWarnings(w0)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e0.GetAttrs("OldID")[0], "borked")
	test.AssertEQ(e0.GetIntrinsicID(), "?")
	test.AssertEQ(e0.GetAbsQualifiedID(), "?")

	t.Logf("Loading event '%s'...", *file.Event[1].Title)
	e1, w1 := createEventFromParsed(file.Event[1])
	test.PrintAllWarnings(w1)

	t.Logf("Investigating loaded event...")
	test.AssertWarnContains(w1, "cannot load an embedded event")
	test.AssertWarnContains(w1, "attempted to embed both")
	test.AssertEQ(e1.GetEmbedEventLink(), "/borked/path")
	test.AssertEQ(e1.GetEmbedChartLink(), "/borked/path")
	test.AssertEQ(e1.IsEventEmbedding(), true)
	test.AssertNE(e1.GetEmbeddedEvent(), nil)

	t.Logf("Loading event '%s'...", *file.Event[2].Title)
	e2, w2 := createEventFromParsed(file.Event[2])
	test.PrintAllWarnings(w2)

	t.Logf("Investigating loaded event...")
	test.AssertWarnContains(w2, "inverted")
	test.AssertEQ(e2.GetStartDate().String(), "0000-01-02")
	test.AssertEQ(e2.GetEndDate().String(), "9999-01-02")

	t.Logf("Loading event '%s'...", *file.Event[3].Title)
	e3, w3 := createEventFromParsed(file.Event[3])
	test.PrintAllWarnings(w3)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e3.GetEndDate().IsInfiniteFuture(), true)
	test.AssertEQ(e3.GetEndDate().String(), "(inf)")

	t.Logf("Loading event '%s'...", *file.Event[4].Title)
	e4, w4 := createEventFromParsed(file.Event[4])
	test.PrintAllWarnings(w4)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e4.GetStartDate().IsInfinitePast(), true)
	test.AssertEQ(e4.GetStartDate().String(), "(-inf)")

	t.Logf("Loading event '%s'...", *file.Event[5].Title)
	e5, w5 := createEventFromParsed(file.Event[5])
	test.PrintAllWarnings(w5)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e5.GetStartDate().String(), "9999-01-__")
	test.AssertEQ(e5.GetEndDate().String(), "9999-02-**")

	t.Logf("Loading event '%s'...", *file.Event[6].Title)
	e6, w6 := createEventFromParsed(file.Event[6])
	test.PrintAllWarnings(w6)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e6.GetStartDate().String(), "9999-01-__")
	test.AssertEQ(e6.GetEndDate().String(), "9999-01-**")

	t.Logf("Loading event '%s'...", *file.Event[7].Title)
	e7, w7 := createEventFromParsed(file.Event[7])
	test.PrintAllWarnings(w7)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e7.GetStartDate().String(), "9998-03-__")
	test.AssertEQ(e7.GetEndDate().String(), "9999-02-**")

	t.Logf("Loading event '%s'...", *file.Event[8].Title)
	e8, w8 := createEventFromParsed(file.Event[8])
	test.PrintAllWarnings(w8)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e8.GetStartDate().String(), "9999-02-00")
	test.AssertEQ(e8.GetEndDate().String(), "9999-02-**")

	t.Logf("Loading event '%s'...", *file.Event[9].Title)
	e9, w9 := createEventFromParsed(file.Event[9])
	test.PrintAllWarnings(w9)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e9.GetStartDate().String(), "9999-02-00")
	test.AssertEQ(e9.GetEndDate().String(), "9999-02-08")

	t.Logf("Loading event '%s'...", *file.Event[10].Title)
	e10, w10 := createEventFromParsed(file.Event[10])
	test.PrintAllWarnings(w10)

	t.Logf("Investigating loaded event...")
	test.AssertEQ(e10.GetStartDate().String(), "(-inf)")
	test.AssertEQ(e10.GetEndDate().String(), "0000-01-02")

	t.Logf("Loading event '%s'...", *file.Event[11].Title)
	_, w11 := createEventFromParsed(file.Event[11])
	test.PrintAllWarnings(w11)

	t.Logf("Investigating loaded event...")
	test.AssertWarnContains(w11, "'ContinuedFrom' attribute")
	test.AssertWarnContains(w11, "empty.")

	t.Log("End test.")
}
