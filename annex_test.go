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
	test.AssertEQ(annex.GetFormat(), "txt")

	rawData, derr := annex.GetRawData()
	if derr != nil {
		t.Fatal(derr)
		return
	}

	test.AssertEQ(string(rawData), "hello, world")

	t.Log("End test.")
}

func TestIndirectAnnex_Annex(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'indirect_annex.ltc'...")
	file, wf, err := file.LoadFromURI[file.File]("indirect_annex.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}
	test.PrintAllWarnings(wf)

	t.Logf("Loading annex...")
	a0, w0 := createAnnexFromParsed(file.Annex[0])
	test.PrintAllWarnings(w0)

	t.Logf("Investigating loaded annex...")
	r0, e0 := a0.GetRawData()
	if e0 != nil {
		t.Fatal(e0)
		return
	}

	test.AssertEQ(string(r0), "hello, world")

	t.Logf("Loading annex...")
	a1, w1 := createAnnexFromParsed(file.Annex[1])
	test.PrintAllWarnings(w1)

	t.Logf("Investigating loaded annex...")
	test.AssertWarnContains(w1, "both")
	r1, e1 := a1.GetRawData()
	if e1 != nil {
		t.Fatal(e1)
		return
	}

	test.AssertEQ(string(r1), "hi, globe")

	t.Logf("Loading annex...")
	a2, w2 := createAnnexFromParsed(file.Annex[2])
	test.PrintAllWarnings(w2)

	t.Logf("Investigating loaded annex...")
	r2, e2 := a2.GetRawData()
	if e2 != nil {
		t.Fatal(e2)
		return
	}

	test.AssertEQ(r2[0], 0x00)
	test.AssertEQ(r2[1], 0x31)
	test.AssertEQ(r2[2], 0xf8)
	test.AssertEQ(r2[3], 0xe4)

	t.Log("End test.")
}

func TestEncodedAnnex_Annex(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'encoded_annex.ltc'...")
	file, wf, err := file.LoadFromURI[file.File]("encoded_annex.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}
	test.PrintAllWarnings(wf)

	t.Logf("Loading annex...")
	a0, w0 := createAnnexFromParsed(file.Annex[0])
	test.PrintAllWarnings(w0)

	t.Logf("Investigating loaded annex...")
	r0, e0 := a0.GetRawData()
	if e0 != nil {
		t.Fatal(e0)
		return
	}

	test.AssertEQ(string(r0), "hello")

	t.Log("End test.")
}

func TestAttachedAnnex_Annex(t *testing.T) {
	test.Initialize(t, "assets")

	t.Logf("Loading 'attached_annex.ltc'...")
	file, wf, err := file.LoadFromURI[file.File]("attached_annex.ltc")
	if err != nil {
		t.Fatal(err)
		return
	}
	test.PrintAllWarnings(wf)

	t.Logf("Loading event/annex...")
	e0, _ := createEventFromParsed(file.Event[0])
	a0, w0 := createAnnexFromParsed(file.Annex[0])
	a1, w1 := createAnnexFromParsed(file.Annex[1])
	test.PrintAllWarnings(w0)
	test.PrintAllWarnings(w1)

	t.Logf("Investigating loaded annex...")

	// 'AttachTo:e0' should have been reserved at annex creation.
	test.AssertEQ(len(a0.AttachTo.Get()), 1)
	test.AssertEQ(a0.AttachTo.Get()[0].GetIntrinsicID(), "e0")
	test.AssertEQ(a0.AttachTo.Get()[0].GetAbsQualifiedID(), "e0")
	test.AssertEQ(a0.AttachTo.Get()[0].IsUnresolved(), true)
	test.AssertEQ(a0.AttachTo.Get()[0].getNumberID(), NID_Invalid)

	// 'ExtraNoteOf:e0' should have been reserved at annex creation.
	test.AssertEQ(len(a1.ExtraNoteOf.Get()), 1)
	test.AssertEQ(a1.ExtraNoteOf.Get()[0].GetIntrinsicID(), "e0")
	test.AssertEQ(a1.ExtraNoteOf.Get()[0].GetAbsQualifiedID(), "e0")
	test.AssertEQ(a1.ExtraNoteOf.Get()[0].IsUnresolved(), true)
	test.AssertEQ(a1.ExtraNoteOf.Get()[0].getNumberID(), NID_Invalid)

	// Creating and re-reserving the same object shouldn't change anything.
	a0.AttachTo.Create(e0)
	test.AssertEQ(len(a0.AttachTo.Get()), 1)
	test.AssertEQ(len(e0.Attached.Get()), 1)
	test.AssertEQ(a0.AttachTo.Has(e0), true)
	test.AssertEQ(e0.Attached.Has(a0), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	a0.AttachTo.Create(e0)
	test.AssertEQ(len(a0.AttachTo.Get()), 1)
	test.AssertEQ(len(e0.Attached.Get()), 1)
	test.AssertEQ(a0.AttachTo.Has(e0), true)
	test.AssertEQ(e0.Attached.Has(a0), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	a0.AttachTo.Reserve("e0")
	test.AssertEQ(len(a0.AttachTo.Get()), 1)
	test.AssertEQ(len(e0.Attached.Get()), 1)
	test.AssertEQ(a0.AttachTo.Has(e0), true)
	test.AssertEQ(e0.Attached.Has(a0), true)
	test.AssertEQ(e0.IsUnresolved(), false)
	a0.AttachTo.Break(e0)
	test.AssertEQ(len(a0.AttachTo.Get()), 1)
	test.AssertEQ(len(e0.Attached.Get()), 0)
	test.AssertEQ(a0.AttachTo.Has(e0), false)
	test.AssertEQ(e0.Attached.Has(a0), false)
	test.AssertEQ(e0.IsUnresolved(), false)
	a0.AttachTo.Unreserve("e0")
	test.AssertEQ(len(a0.AttachTo.Get()), 0)
	test.AssertEQ(len(e0.Attached.Get()), 0)

	t.Log("End test.")
}
