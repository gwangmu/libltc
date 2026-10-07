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
	t.Log("Title: " + e0.Note.Title)
	for _, snippet := range e0.Note.Snippets {
		t.Log("Time: " + snippet.Time.String())
		t.Logf("Text: %q\n", []rune(snippet.Text))
	}

	t.Log("End test.")
}
