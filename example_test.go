package libltc_test

import (
	"testing"

	. "github.com/gwangmu/libltc"
	"github.com/gwangmu/libltc/internal/test"
)

func TestCurriculum_Example(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'student.ltc'...")
	chart, wc := CreateChart("student.ltc")
	test.PrintAllWarnings(wc)

	t.Logf("Investigating loaded event...")
	for _, eobj := range chart.GetEvents() {
		t.Log(eobj.GetTitle())
	}

	sc, wc2 := chart.GetObjectByQualifiedID("e1").(*Event).GetSubchart()
	test.PrintAllWarnings(wc2)

	//t.Log("All classes:")
	//for _, eobj := range sc.GetImports() {
	//	t.Log(" - " + eobj.GetTitle() + " @ " + eobj.GetStartDate().String() + " ~ " + eobj.GetEndDate().String())
	//}
	t.Log("All schedules:")
	for _, eobj := range sc.GetEvents() {
		t.Log(" - " + eobj.GetTitle() + " @ " + eobj.GetStartDate().String() + " ~ " + eobj.GetEndDate().String())
	}
	t.Log("Lectures:")
	for _, eobj := range sc.GetEventsInCategory("Lecture") {
		t.Log(" - " + eobj.GetTitle() + " @ " + eobj.GetStartDate().String() + " ~ " + eobj.GetEndDate().String())
	}

	t.Logf("End test.")
}
