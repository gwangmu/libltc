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

	t.Logf("Linking events... (indirect)")
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Get()[0].GetIntrinsicID(), "e0")
	test.AssertEQ(e1.ContinuedFrom.Get()[0].GetAbsQualifiedID(), "e0")
	test.AssertEQ(e1.ContinuedFrom.Get()[0].IsUnresolved(), true)
	test.AssertEQ(e1.ContinuedFrom.Get()[0].getNumberID(), NID_Invalid)
	CreateLink[ContinuedFrom](e1, e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), true)
	test.AssertEQ(e0.ContinuedTo.Has(e1), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	BreakLink[ContinuedFrom](e1, e0)
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 1)
	test.AssertEQ(e1.ContinuedFrom.Has(e0), false)
	test.AssertEQ(e0.ContinuedTo.Has(e1), false)
	test.AssertEQ(e0.IsUnresolved(), false)
	UnreserveLink[ContinuedFrom](e1, "e0")
	test.AssertEQ(len(e1.ContinuedFrom.Get()), 0)

	t.Log("End test.")
}
