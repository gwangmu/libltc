package file

import (
	"testing"
	"libltc/internal/test"
)

func TestSingleEvent_File(t *testing.T) {
	test.Initialize_File(t, "../../assets")

	t.Log("Loading 'single_event.ltc'...")
	file, warns, err := LoadFromURI[File]("single_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	} else {
		test.PrintAllWarnings(warns)
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(file.Version, "26.09.1")
	test.AssertEQ(file.Setting.CalendarSystem, "Gregorian")
	test.AssertEQ(file.Subject.Name.Summary(), "Gwangmu Lee (Nickname:Lee)")

	test.AssertEQ(len(file.Event), 1)
	test.AssertEQ(file.Event[0].ID, "e0")
	test.AssertEQ(len(file.Event[0].Attrs), 5)
	test.AssertEQ(file.Event[0].Attrs[0], "Attr2")
	test.AssertEQ(*file.Event[0].StartDate.Year, 1991)
	test.AssertEQ(*file.Event[0].StartDate.Month, 3)
	test.AssertEQ(*file.Event[0].StartDate.Day, 29)
	test.AssertEQ(file.Event[0].StartDate.Attrs[0], "Attr0")
	test.AssertEQ(file.Event[0].EndDate.Attrs[0], "Attr1")
	test.AssertEQ(file.Event[0].EndDate.Day, nil)

	t.Log("End test.")
}

func TestSimpleEmbedEvent_File(t *testing.T) {
	test.Initialize_File(t, "../../assets")

	t.Log("Loading 'simple_embed_event.ltc'...")
	file, warns, err := LoadFromURI[File]("simple_embed_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	} else {
		test.PrintAllWarnings(warns)
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(len(file.Event), 2)
	test.AssertEQ(file.Event[1].ID, "e1")
	test.AssertEQ(*file.Event[1].EmbedEvent, "./simple_embed_event.ltc?id=e0")

	t.Log("Loading embedded event manually...")
	event, ewarns, eerr := LoadFromURI[Event](*file.Event[1].EmbedEvent)
	if eerr != nil {
		t.Fatal(eerr)
		return
	} else {
		test.PrintAllWarnings(ewarns)
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(event.ID, "e0")
	test.AssertEQ(len(event.Attrs), 1)
	test.AssertEQ(event.Attrs[0], "AmbiguousPeriod")
	test.AssertEQ(*event.StartDate.Year, 1991)
	test.AssertEQ(*event.StartDate.Month, 3)
	test.AssertEQ(*event.StartDate.Day, 29)
	test.AssertEQ(event.StartDate.Attrs[0], "Approx")
	test.AssertEQ(event.EndDate.Attrs[0], "Approx")
	test.AssertEQ(event.EndDate.Day, nil)

	t.Log("End test.")
}

func TestNestedEmbedEvent_File(t *testing.T) {
	test.Initialize_File(t, "../../assets")

	t.Log("Loading 'nested_embed_event.ltc'...")
	file, warns, err := LoadFromURI[File]("nested_embed_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	} else {
		test.PrintAllWarnings(warns)
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(len(file.Event), 3)
	test.AssertEQ(file.Event[2].ID, "e2")
	test.AssertEQ(*file.Event[2].EmbedEvent, "./nested_embed_event.ltc?id=e1")

	t.Log("Loading embedded event manually...")
	event, ewarns, eerr := LoadFromURI[Event](*file.Event[2].EmbedEvent)
	if eerr != nil {
		t.Fatal(eerr)
		return
	} else {
		test.PrintAllWarnings(ewarns)
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(event.ID, "e0")
	test.AssertEQ(len(event.Attrs), 1)
	test.AssertEQ(event.Attrs[0], "AmbiguousPeriod")
	test.AssertEQ(*event.StartDate.Year, 1991)
	test.AssertEQ(*event.StartDate.Month, 3)
	test.AssertEQ(*event.StartDate.Day, 29)
	test.AssertEQ(event.StartDate.Attrs[0], "Approx")
	test.AssertEQ(event.EndDate.Attrs[0], "Approx")
	test.AssertEQ(event.EndDate.Day, nil)

	t.Log("End test.")
}

func TestNonexistentEmbedEvent_File(t *testing.T) {
	test.Initialize_File(t, "../../assets")

	t.Log("Loading 'nonexistent_embed_event.ltc'...")
	file, warns, err := LoadFromURI[File]("nonexistent_embed_event.ltc")
	if err != nil {
		t.Fatal(err)
		return
	} else {
		test.PrintAllWarnings(warns)
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(len(file.Event), 3)
	test.AssertEQ(file.Event[1].ID, "e1")
	test.AssertEQ(*file.Event[1].EmbedEvent, "nonexistent")
	test.AssertEQ(file.Event[2].ID, "e2")
	test.AssertEQ(*file.Event[2].EmbedEvent, "./nonexistent_embed_event.ltc?id=e999")

	t.Log("Loading embedded event manually...")
	_, ewarns, eerr := LoadFromURI[Event](*file.Event[1].EmbedEvent)
	test.AssertNE(eerr, nil)
	test.AssertContains(eerr.Error(), "no such file")

	for _, warn := range ewarns {
		t.Log(" - " + warn.GetDesc())
	}

	t.Log("Loading embedded event manually...")
	_, e2warns, e2err := LoadFromURI[Event](*file.Event[2].EmbedEvent)
	test.AssertNE(e2err, nil)
	test.AssertContains(e2err.Error(), "Object")
	test.AssertContains(e2err.Error(), "not found")

	for _, warn := range e2warns {
		t.Log(" - " + warn.GetDesc())
	}

	t.Log("End test.")
}

func TestQualIDEmbedEvent_File(t *testing.T) {
	test.Initialize_File(t, "../../assets")

	t.Log("Loading 'qualid_embed_event_2.ltc'...")
	file, warns, err := LoadFromURI[File]("qualid_embed_event_2.ltc")
	if err != nil {
		t.Fatal(err)
		return
	} else {
		test.PrintAllWarnings(warns)
	}

	t.Log("Loading embedded event manually...")
	event, ewarns, eerr := LoadFromURI[Event](*file.Event[0].EmbedEvent)
	if eerr != nil {
		t.Fatal(eerr)
		return
	} else {
		test.PrintAllWarnings(ewarns)
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(event.ID, "e0")
	test.AssertEQ(len(event.Attrs), 5)
	test.AssertEQ(event.Attrs[0], "Attr2")
	test.AssertEQ(*event.StartDate.Year, 1991)
	test.AssertEQ(*event.StartDate.Month, 3)
	test.AssertEQ(*event.StartDate.Day, 29)
	test.AssertEQ(event.StartDate.Attrs[0], "Attr0")
	test.AssertEQ(event.EndDate.Attrs[0], "Attr1")
	test.AssertEQ(event.EndDate.Day, nil)

	t.Log("End test.")
}

func TestReadFile_File(t *testing.T) {
	test.Initialize_File(t, "../../assets")

	t.Log("Loading 'binary.bin'...")
	bs, err := ReadFileFromURI("binary.bin")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Investigating loaded binary...")
	test.AssertEQ(bs[0], 0x00)
	test.AssertEQ(bs[1], 0x31)
	test.AssertEQ(bs[2], 0xf8)
	test.AssertEQ(bs[3], 0xe4)

	t.Log("Loading 'hello.txt'...")
	bs2, err2 := ReadFileFromURI("hello.txt")
	if err2 != nil {
		t.Fatal(err2)
	}

	t.Log("Investigating loaded text...")
	test.AssertEQ(string(bs2), "hello, world")

	t.Log("End test.")
}

