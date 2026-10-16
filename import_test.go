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
