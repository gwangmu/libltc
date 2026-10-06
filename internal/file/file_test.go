package file

import (
	"testing"
	"libltc/internal/test"
)

var TestAssetDirPath string = "../../test_asset"

func TestSingleEvent(t *testing.T) {
	t.Chdir(TestAssetDirPath)

	t.Log("Loading 'single_event.ltc'...")
	file, warns, err := LoadFromURI[File]("single_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("Printing all warnings...")
	for _, warn := range warns {
		t.Log(warn.GetDesc())
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(t, file.Version, "26.09.1")
	test.AssertEQ(t, file.Setting.CalendarSystem, "Gregorian")
	test.AssertEQ(t, file.Subject.Name.Summary(), "Gwangmu Lee (Nickname:Lee)")

	test.AssertEQ(t, len(file.Event), 1)
	test.AssertEQ(t, file.Event[0].ID, "e0")
	test.AssertEQ(t, len(file.Event[0].Attrs), 1)
	test.AssertEQ(t, file.Event[0].Attrs[0], "AmbiguousPeriod")
	test.AssertEQ(t, *file.Event[0].StartDate.Year, 1991)
	test.AssertEQ(t, *file.Event[0].StartDate.Month, 3)
	test.AssertEQ(t, *file.Event[0].StartDate.Day, 29)
	test.AssertEQ(t, file.Event[0].StartDate.Attrs[0], "Approx")
	test.AssertEQ(t, file.Event[0].EndDate.Attrs[0], "Approx")
	test.AssertEQ(t, file.Event[0].EndDate.Day, nil)

	t.Log("End test.")
}

func TestSimpleEmbedEvent(t *testing.T) {
	t.Chdir(TestAssetDirPath)

	t.Log("Loading 'simple_embed_event.ltc'...")
	file, warns, err := LoadFromURI[File]("simple_embed_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("Printing all warnings...")
	for _, warn := range warns {
		t.Log(warn.GetDesc())
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(t, len(file.Event), 2)
	test.AssertEQ(t, file.Event[1].ID, "e1")
	test.AssertEQ(t, *file.Event[1].EmbedEvent, "./simple_embed_event.ltc?id=e0")

	t.Log("Loading embedded event manually...")
	event, ewarns, eerr := LoadFromURI[Event](*file.Event[1].EmbedEvent)
	if eerr != nil {
		t.Fatal(eerr)
		return
	}

	t.Log("Printing all warnings...")
	for _, warn := range ewarns {
		t.Log(warn.GetDesc())
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(t, event.ID, "e0")
	test.AssertEQ(t, len(event.Attrs), 1)
	test.AssertEQ(t, event.Attrs[0], "AmbiguousPeriod")
	test.AssertEQ(t, *event.StartDate.Year, 1991)
	test.AssertEQ(t, *event.StartDate.Month, 3)
	test.AssertEQ(t, *event.StartDate.Day, 29)
	test.AssertEQ(t, event.StartDate.Attrs[0], "Approx")
	test.AssertEQ(t, event.EndDate.Attrs[0], "Approx")
	test.AssertEQ(t, event.EndDate.Day, nil)

	t.Log("End test.")
}

func TestNestedEmbedEvent(t *testing.T) {
	t.Chdir(TestAssetDirPath)

	t.Log("Loading 'nested_embed_event.ltc'...")
	file, warns, err := LoadFromURI[File]("nested_embed_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("Printing all warnings...")
	for _, warn := range warns {
		t.Log(warn.GetDesc())
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(t, len(file.Event), 3)
	test.AssertEQ(t, file.Event[2].ID, "e2")
	test.AssertEQ(t, *file.Event[2].EmbedEvent, "./nested_embed_event.ltc?id=e1")

	t.Log("Loading embedded event manually...")
	event, ewarns, eerr := LoadFromURI[Event](*file.Event[2].EmbedEvent)
	if eerr != nil {
		t.Fatal(eerr)
		return
	}

	t.Log("Printing all warnings...")
	for _, warn := range ewarns {
		t.Log(warn.GetDesc())
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(t, event.ID, "e0")
	test.AssertEQ(t, len(event.Attrs), 1)
	test.AssertEQ(t, event.Attrs[0], "AmbiguousPeriod")
	test.AssertEQ(t, *event.StartDate.Year, 1991)
	test.AssertEQ(t, *event.StartDate.Month, 3)
	test.AssertEQ(t, *event.StartDate.Day, 29)
	test.AssertEQ(t, event.StartDate.Attrs[0], "Approx")
	test.AssertEQ(t, event.EndDate.Attrs[0], "Approx")
	test.AssertEQ(t, event.EndDate.Day, nil)

	t.Log("End test.")
}

func TestNonexistentEmbedEvent(t *testing.T) {
	t.Chdir(TestAssetDirPath)

	t.Log("Loading 'nonexistent_embed_event.ltc'...")
	file, warns, err := LoadFromURI[File]("nonexistent_embed_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("Printing all warnings...")
	for _, warn := range warns {
		t.Log(warn.GetDesc())
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(t, len(file.Event), 3)
	test.AssertEQ(t, file.Event[1].ID, "e1")
	test.AssertEQ(t, *file.Event[1].EmbedEvent, "nonexistent")
	test.AssertEQ(t, file.Event[2].ID, "e2")
	test.AssertEQ(t, *file.Event[2].EmbedEvent, "./nonexistent_embed_event.ltc?id=e999")

	t.Log("Loading embedded event manually...")
	_, ewarns, eerr := LoadFromURI[Event](*file.Event[1].EmbedEvent)
	test.AssertNE(t, eerr, nil)

	t.Log("Printing all warnings...")
	for _, warn := range ewarns {
		t.Log(warn.GetDesc())
	}

	t.Log("End test.")
}
