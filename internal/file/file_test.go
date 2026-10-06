package file

import (
	"testing"
)

var TestAssetDirPath string = "../../test_asset"

func AssertEQ[T comparable](t *testing.T, entity string, got T, want T) {
	t.Logf(" - %s: %v", entity, got)
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
	AssertEQ(t, "file.Version", file.Version, "26.09.1")
	AssertEQ(t, "file.Setting.CalendarSystem", file.Setting.CalendarSystem, "Gregorian")
	AssertEQ(t, "file.Subject.Name.Summary()", file.Subject.Name.Summary(), "Gwangmu Lee (Nickname:Lee)")

	AssertEQ(t, "len(file.Event)", len(file.Event), 1)
	AssertEQ(t, "file.Event[0].ID", file.Event[0].ID, "e0")
	AssertEQ(t, "len(file.Event[0].Attrs)", len(file.Event[0].Attrs), 1)
	AssertEQ(t, "file.Event[0].Attrs[0]", file.Event[0].Attrs[0], "AmbiguousPeriod")
	AssertEQ(t, "*file.Event[0].StartDate.Year", *file.Event[0].StartDate.Year, 1991)
	AssertEQ(t, "*file.Event[0].StartDate.Month", *file.Event[0].StartDate.Month, 3)
	AssertEQ(t, "*file.Event[0].StartDate.Day", *file.Event[0].StartDate.Day, 29)
	AssertEQ(t, "file.Event[0].StartDate.Attrs[0]", file.Event[0].StartDate.Attrs[0], "Approx")
	AssertEQ(t, "file.Event[0].EndDate.Attrs[0]", file.Event[0].EndDate.Attrs[0], "Approx")
	AssertEQ(t, "file.Event[0].EndDate.Day", file.Event[0].EndDate.Day, nil)

	t.Log("End test.")
}

func TestSimpleEmbedEvent(t *testing.T) {
	t.Chdir(TestAssetDirPath)

	t.Log("Loading 'simple_embed_event.ltc'...")
	file, warns, err := LoadFromURI[File]("simple_embed_event.ltc")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Printing all warnings...")
	for _, warn := range warns {
		t.Log(warn.GetDesc())
	}

	t.Log("Investigating loaded file...")
	AssertEQ(t, "len(file.Event)", len(file.Event), 2)
	AssertEQ(t, "file.Event[1].ID", file.Event[1].ID, "e1")
	AssertEQ(t, "*file.Event[1].EmbedEvent", *file.Event[1].EmbedEvent, "./simple_embed_event.ltc?id=e0")

	t.Log("Loading embedded event manually...")
	event, ewarns, eerr := LoadFromURI[Event](*file.Event[1].EmbedEvent)
	if eerr != nil {
		t.Error(eerr)
		return
	}

	t.Log("Printing all warnings...")
	for _, warn := range ewarns {
		t.Log(warn.GetDesc())
	}

	t.Log("Investigating loaded file...")
	AssertEQ(t, "event.ID", event.ID, "e0")
	AssertEQ(t, "len(event.Attrs)", len(event.Attrs), 1)
	AssertEQ(t, "event.Attrs[0]", event.Attrs[0], "AmbiguousPeriod")
	AssertEQ(t, "*event.StartDate.Year", *event.StartDate.Year, 1991)
	AssertEQ(t, "*event.StartDate.Month", *event.StartDate.Month, 3)
	AssertEQ(t, "*event.StartDate.Day", *event.StartDate.Day, 29)
	AssertEQ(t, "event.StartDate.Attrs[0]", event.StartDate.Attrs[0], "Approx")
	AssertEQ(t, "event.EndDate.Attrs[0]", event.EndDate.Attrs[0], "Approx")
	AssertEQ(t, "event.EndDate.Day", event.EndDate.Day, nil)

	t.Log("End test.")
}
