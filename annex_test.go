package ltc 

import (
	"testing"
	"libltc/internal/file"
	"libltc/internal/test"
)

func TestSingleAnnex_Annex(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'single_annex.ltc'...")
	file, wf, err := file.LoadFromURI[file.File]("single_annex.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}
	test.PrintAllWarnings(wf)

	t.Logf("Loading annex...")
	annex, wa := createAnnexFromParsed(file.Annex[0])
	test.PrintAllWarnings(wa)

	t.Logf("Investigating loaded annex...")
	test.AssertEQ(annex.GetIntrinsicID(), "a0")
	test.AssertEQ(annex.Title, "Whatever")

	t.Log("End test.")
}
