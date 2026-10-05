package file

import (
	"testing"
)

var TestAssetDirPath string = "../../test_asset"

func AssertEQ[T comparable](t *testing.T, got T, want T) {
	if got != want {
		t.Fatalf("assertion failed: %v != %v", got, want)
	}
}

func TestSingleEvent(t *testing.T) {
	t.Chdir(TestAssetDirPath)

	t.Log("Loading 'single_event.ltc'...")
	file, warns, err := LoadFromURI[File]("single_event.ltc")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Printing all warnings...")
	for _, warn := range warns {
		t.Log(warn.GetDesc())
	}

	t.Log("Investigating loaded file...")
	AssertEQ(t, file.Version, "26.09.1")
	AssertEQ(t, file.Setting.CalendarSystem, "Gregorian")
	AssertEQ(t, file.Subject.Name.Summary(), "Gwangmu Lee (Nickname:Lee)")

	AssertEQ(t, len(file.Event), 1)
	AssertEQ(t, file.Event[0].ID, "e0")
	AssertEQ(t, len(file.Event[0].Attrs), 1)
	AssertEQ(t, file.Event[0].Attrs[0], "AmbiguousPeriod")
	AssertEQ(t, *file.Event[0].StartDate.Year, 1991)
	AssertEQ(t, *file.Event[0].StartDate.Month, 3)
	AssertEQ(t, *file.Event[0].StartDate.Day, 29)
	AssertEQ(t, file.Event[0].EndDate.Day, nil)

	t.Log("End test.")
}
