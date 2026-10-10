package ltc 

import (
	"testing"
	"libltc/internal/file"
	"libltc/internal/test"
)

func TestSingleEvent_Chart(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'single_event.ltc'...")
	ltcf, wf, err := file.LoadFromURI[file.File]("single_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}
	test.PrintAllWarnings(wf)

	t.Logf("Loading chart...")
	chart, wc := createChartFromParsed(ltcf)
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
