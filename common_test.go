package ltc 

import (
	"testing"
	"libltc/internal/file"
	"libltc/internal/test"
)

var TestAssetDirPath string = "assets"

func TestNotes(t *testing.T) {
	test.Initialize(t, TestAssetDirPath)

	t.Log("Loading 'notes.ltc'...")
	file, _, err := file.LoadFromURI[file.File]("notes.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("Loading the first event...")
	e0, w0 := createEventFromParsed(file.Event[0])

	t.Log("Printing all warnings...")
	for _, w := range w0 {
		t.Log(" - " + w.GetDesc())
	}

	t.Log("Investigating loaded file...")
	test.AssertEQ(e0.Note.Title, "")
	test.AssertEQ(len(e0.Note.Snippets), 2)
	test.AssertEQ(e0.Note.Snippets[0].Time.String(), "(unknown)")
	test.AssertEQ(e0.Note.Snippets[0].Text, "aaa")
	test.AssertEQ(e0.Note.Snippets[1].Time.String(), "2026-09-01 00:01:02 (UTC)")
	test.AssertEQ(e0.Note.Snippets[1].Text, "bbb\nccc")

	t.Log("End test.")
}
