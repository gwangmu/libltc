package ltc 

import (
	"testing"
	"libltc/internal/file"
	"libltc/internal/test"
)

func TestSingleEvent(t *testing.T) {
	test.Initialize(t, TestAssetDirPath)

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

func TestContinuedEvent(t *testing.T) {
	test.Initialize(t, TestAssetDirPath)

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
	CreateLink[ContinuedFrom](e1, e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	CreateLink[ContinuedFrom](e1, e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	ReserveLink[ContinuedFrom](e1, "e0")
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	BreakLink[ContinuedFrom](e1, e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 0)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), false)
	test.AssertEQ(e0.ContinuedTo.Has(e1), false)
	test.AssertEQ(e0.IsUnresolved(), false)
	UnreserveLink[ContinuedFrom](e1, "e0")
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 0)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 0)

	// Re-creating should also be possible. 
	CreateLink[ContinuedFrom](e1, e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	BreakLink[ContinuedFrom](e1, e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 0)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), false)
	test.AssertEQ(e0.ContinuedTo.Has(e1), false)
	test.AssertEQ(e0.IsUnresolved(), false)
	CreateLink[ContinuedFrom](e1, e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(len(e0.ContinuedTo.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)

	t.Log("End test.")
}

func TestEventEmbeddingEvent(t *testing.T) {
	test.Initialize(t, TestAssetDirPath)

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

//func TestChartEmbeddingEvent(t *testing.T) {
// TODO
//}

//func TestBrokenEvent(t *testing.T) {
// TODO: broken ID event
// TODO: both-embedding event
// TODO: swapped dates event
// TODO: empty 'ContinuedFrom' event
//}
