package ltc 

import (
	"testing"
	"libltc/internal/file"
	"libltc/internal/test"
)

var TestAssetDirPath string = "assets"

func TestNotes(t *testing.T) {
	test.Initialize(t, TestAssetDirPath)

	t.Logf("Loading 'notes.ltc'...")
	file, _, err := file.LoadFromURI[file.File]("notes.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("Loading the note '%s'...", *file.Event[0].Title)
	n0, w0 := CreateNoteFromString(file.Event[0].Note)
	test.PrintAllWarnings(w0)

	t.Log("Investigating loaded file...")
	test.AssertEQ(n0.Title, "")
	test.AssertEQ(len(n0.Snippets), 2)
	test.AssertEQ(n0.Snippets[0].Time.String(), "(unknown)")
	test.AssertEQ(n0.Snippets[0].Text, "aaa")
	test.AssertEQ(n0.Snippets[1].Time.String(), "2026-09-01 00:01:02 (UTC)")
	test.AssertEQ(n0.Snippets[1].Text, "bbb\nccc")

	t.Logf("Loading the note '%s'...", *file.Event[1].Title)
	n1, w1 := CreateNoteFromString(file.Event[1].Note)
	test.PrintAllWarnings(w1)

	t.Log("Investigating loaded file...")
	test.AssertEQ(n1.Title, "")
	test.AssertEQ(len(n1.Snippets), 2)
	test.AssertEQ(n1.Snippets[0].Time.String(), "2026-09-01 00:01:02 (UTC)")
	test.AssertEQ(n1.Snippets[0].Text, "aaa")
	test.AssertEQ(n1.Snippets[1].Time.String(), "2026-09-02 03:04:05 (UTC)")
	test.AssertEQ(n1.Snippets[1].Text, "bbb\nccc")

	t.Logf("Loading the note '%s'...", *file.Event[2].Title)
	n2, w2 := CreateNoteFromString(file.Event[2].Note)
	test.PrintAllWarnings(w2)

	t.Log("Investigating loaded file...")
	test.AssertEQ(n2.Title, "")
	test.AssertEQ(len(n2.Snippets), 2)
	test.AssertEQ(n2.Snippets[0].Time.String(), "(unknown)")
	test.AssertEQ(n2.Snippets[0].Text, "aaa")
	test.AssertEQ(n2.Snippets[1].Time.String(), "2026-09-02 11:22:33 (UTC)")
	test.AssertEQ(n2.Snippets[1].Text, "bbb\nccc")

	t.Logf("Loading the note '%s'...", *file.Event[3].Title)
	n3, w3 := CreateNoteFromString(file.Event[3].Note)
	test.PrintAllWarnings(w3)

	t.Log("Investigating loaded file...")
	test.AssertEQ(n3.Title, "")
	test.AssertEQ(len(n3.Snippets), 8)
	test.AssertEQ(n3.Snippets[0].Time.String(), "(unknown)")
	test.AssertEQ(n3.Snippets[0].Text, "aaa")
	test.AssertEQ(n3.Snippets[1].Time.String(), "2026-09-02 11:22:33 (UTC)")
	test.AssertEQ(n3.Snippets[1].Text, "bbb\nccc")
	test.AssertEQ(n3.Snippets[2].Time.String(), "2026-09-02 11:22:00")
	test.AssertEQ(n3.Snippets[2].Text, "ddd\neee")
	test.AssertEQ(n3.Snippets[3].Time.String(), "2026-09-02 (UTC)")
	test.AssertEQ(n3.Snippets[3].Text, "fff\nggg")
	test.AssertEQ(n3.Snippets[4].Time.String(), "2026-09-02")
	test.AssertEQ(n3.Snippets[4].Text, "hhh\niii")
	test.AssertEQ(n3.Snippets[5].Time.String(), "2026-09-??")
	test.AssertEQ(n3.Snippets[5].Text, "jjj\nkkk")
	test.AssertEQ(n3.Snippets[6].Time.String(), "2026-09-??")
	test.AssertEQ(n3.Snippets[6].Text, "lll\nmmm")
	test.AssertEQ(n3.Snippets[7].Time.String(), "(unknown)")
	test.AssertEQ(n3.Snippets[7].Text, "nnn\nooo")

	t.Logf("Loading the note '%s'...", *file.Event[4].Title)
	n4, w4 := CreateNoteFromString(file.Event[4].Note)
	test.PrintAllWarnings(w4)

	t.Log("Investigating loaded file...")
	test.AssertEQ(n4.Title, "whatever")
	test.AssertEQ(len(n4.Snippets), 2)
	test.AssertEQ(n4.Snippets[0].Time.String(), "2026-09-01 00:01:02 (UTC)")
	test.AssertEQ(n4.Snippets[0].Text, "aaa")
	test.AssertEQ(n4.Snippets[1].Time.String(), "2026-09-01 00:01:02 (UTC)")
	test.AssertEQ(n4.Snippets[1].Text, "bbb\nccc")

	t.Logf("Loading the note '%s'...", *file.Event[5].Title)
	n5, w5 := CreateNoteFromString(file.Event[5].Note)
	test.PrintAllWarnings(w5)

	t.Log("Investigating loaded file...")
	test.AssertEQ(n5.Title, "")
	test.AssertEQ(len(n5.Snippets), 3)
	test.AssertEQ(n5.Snippets[0].Time.String(), "2026-09-01 00:01:02 (CET)")
	test.AssertEQ(n5.Snippets[0].Text, "aaa")
	test.AssertEQ(n5.Snippets[1].Time.String(), "2026-09-01 (Asia/Seoul)")
	test.AssertEQ(n5.Snippets[1].Text, "aaa")
	test.AssertEQ(n5.Snippets[2].Time.String(), "2026-09-01 00:01:02")
	test.AssertEQ(n5.Snippets[2].Text, "bbb\nccc")

	t.Log("End test.")
}
