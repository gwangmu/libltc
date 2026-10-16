package libltc 

import (
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
		for _, iobj := range chart.GetImports() {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		for _, eobj := range chart.GetEvents() {
			gotIDs = append(gotIDs, eobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	ieobj := chart.GetObjectByQualifiedID("i0/e0").(*Event)
	test.AssertEQ(ieobj.GetTitle(), "Wherever")
	test.AssertEQ(ieobj.GetStartDate().String(), "1991-03-29 (Attr0)")
	test.AssertEQ(ieobj.GetCategory(), "Residence")
	test.AssertEQ(chart.GetEventsInCategory("Residence")[0], ieobj)

	t.Logf("Changing import target...")
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
		for _, iobj := range chart.GetImports() {
			gotIDs = append(gotIDs, iobj.GetQualifiedID())
		}
		for _, eobj := range chart.GetEvents() {
			gotIDs = append(gotIDs, eobj.GetQualifiedID())
		}
		test.AssertSetEQ(gotIDs, wantIDs)
	}

	t.Logf("End test.")
}
