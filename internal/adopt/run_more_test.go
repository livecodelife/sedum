package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_more_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_more_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_more_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunDryRunTouchesNothing(t *testing.T) {\n\tcfg := newRunFixture(t, map[string]string{\"a.py\": \"def f():\\n    return 1\\n\"})\n\tpath := filepath.Join(cfg.Root, \"a.py\")\n\tbefore, err := os.ReadFile(path)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tcfg.DryRun = true\n\tcfg.RecordingPath = filepath.Join(t.TempDir(), \"rec.json\")\n\trep, err := Run(cfg)\n\tif err != nil {\n\t\tt.Fatalf(\"Run: %v\", err)\n\t}\n\tafter, err := os.ReadFile(path)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif string(before) != string(after) {\n\t\tt.Errorf(\"dry run changed file bytes\")\n\t}\n\tif _, err := os.Stat(cfg.RecordingPath); err != nil {\n\t\tt.Errorf(\"recording not written: %v\", err)\n\t}\n\tfor _, fr := range rep.Files {\n\t\tif fr.Detail != \"dry run\" {\n\t\t\tt.Errorf(\"%s: Detail = %q, want %q\", fr.Path, fr.Detail, \"dry run\")\n\t\t}\n\t}\n}","file":"internal/adopt/run_more_test.go","name":"TestRunDryRunTouchesNothing"}}
func TestRunDryRunTouchesNothing(t *testing.T) {
	cfg := newRunFixture(t, map[string]string{"a.py": "def f():\n    return 1\n"})
	path := filepath.Join(cfg.Root, "a.py")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cfg.DryRun = true
	cfg.RecordingPath = filepath.Join(t.TempDir(), "rec.json")
	rep, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Errorf("dry run changed file bytes")
	}
	if _, err := os.Stat(cfg.RecordingPath); err != nil {
		t.Errorf("recording not written: %v", err)
	}
	for _, fr := range rep.Files {
		if fr.Detail != "dry run" {
			t.Errorf("%s: Detail = %q, want %q", fr.Path, fr.Detail, "dry run")
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunSkipsResidualFiles(t *testing.T) {\n\tbSrc := \"import os\\n\\ndef f():\\n    return 1\\n\"\n\tcSrc := \"# licence\\n\\ndef f():\\n    return 1\\n\"\n\tcfg := newRunFixture(t, map[string]string{\"b.py\": bSrc, \"c.py\": cSrc})\n\trep, err := Run(cfg)\n\tif err != nil {\n\t\tt.Fatalf(\"Run: %v\", err)\n\t}\n\tfor name, want := range map[string]string{\"b.py\": bSrc, \"c.py\": cSrc} {\n\t\tgot, err := os.ReadFile(filepath.Join(cfg.Root, name))\n\t\tif err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t\tif string(got) != want {\n\t\t\tt.Errorf(\"%s changed on disk\", name)\n\t\t}\n\t}\n\treasons := map[string]Reason{}\n\tfor _, fr := range rep.Files {\n\t\treasons[filepath.Base(fr.Path)] = fr.Reason\n\t}\n\tif reasons[\"b.py\"] != ReasonOtherResidual {\n\t\tt.Errorf(\"b.py reason = %v, want %v\", reasons[\"b.py\"], ReasonOtherResidual)\n\t}\n\tif reasons[\"c.py\"] != ReasonHeader {\n\t\tt.Errorf(\"c.py reason = %v, want %v\", reasons[\"c.py\"], ReasonHeader)\n\t}\n\tif rep.Totals.SkippedByReason[ReasonOtherResidual] != 1 {\n\t\tt.Errorf(\"SkippedByReason[ReasonOtherResidual] = %d, want 1\", rep.Totals.SkippedByReason[ReasonOtherResidual])\n\t}\n\tif rep.Totals.SkippedByReason[ReasonHeader] != 1 {\n\t\tt.Errorf(\"SkippedByReason[ReasonHeader] = %d, want 1\", rep.Totals.SkippedByReason[ReasonHeader])\n\t}\n}","file":"internal/adopt/run_more_test.go","name":"TestRunSkipsResidualFiles"}}
func TestRunSkipsResidualFiles(t *testing.T) {
	bSrc := "import os\n\ndef f():\n    return 1\n"
	cSrc := "# licence\n\ndef f():\n    return 1\n"
	cfg := newRunFixture(t, map[string]string{"b.py": bSrc, "c.py": cSrc})
	rep, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	for name, want := range map[string]string{"b.py": bSrc, "c.py": cSrc} {
		got, err := os.ReadFile(filepath.Join(cfg.Root, name))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Errorf("%s changed on disk", name)
		}
	}
	reasons := map[string]Reason{}
	for _, fr := range rep.Files {
		reasons[filepath.Base(fr.Path)] = fr.Reason
	}
	if reasons["b.py"] != ReasonOtherResidual {
		t.Errorf("b.py reason = %v, want %v", reasons["b.py"], ReasonOtherResidual)
	}
	if reasons["c.py"] != ReasonHeader {
		t.Errorf("c.py reason = %v, want %v", reasons["c.py"], ReasonHeader)
	}
	if rep.Totals.SkippedByReason[ReasonOtherResidual] != 1 {
		t.Errorf("SkippedByReason[ReasonOtherResidual] = %d, want 1", rep.Totals.SkippedByReason[ReasonOtherResidual])
	}
	if rep.Totals.SkippedByReason[ReasonHeader] != 1 {
		t.Errorf("SkippedByReason[ReasonHeader] = %d, want 1", rep.Totals.SkippedByReason[ReasonHeader])
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunTwiceSkipsAdopted(t *testing.T) {\n\tcfg := newRunFixture(t, map[string]string{\"a.py\": \"def f():\\n    return 1\\n\"})\n\tif _, err := Run(cfg); err != nil {\n\t\tt.Fatalf(\"first Run: %v\", err)\n\t}\n\trep, err := Run(cfg)\n\tif err != nil {\n\t\tt.Fatalf(\"second Run: %v\", err)\n\t}\n\tif rep.Totals.Adopted != 0 {\n\t\tt.Errorf(\"second run Adopted = %d, want 0\", rep.Totals.Adopted)\n\t}\n\tfound := false\n\tfor _, fr := range rep.Files {\n\t\tif filepath.Base(fr.Path) == \"a.py\" {\n\t\t\tfound = true\n\t\t\tif fr.Reason != ReasonAlreadyAdopted {\n\t\t\t\tt.Errorf(\"a.py reason = %v, want %v\", fr.Reason, ReasonAlreadyAdopted)\n\t\t\t}\n\t\t}\n\t}\n\tif !found {\n\t\tt.Errorf(\"a.py missing from second report\")\n\t}\n}","file":"internal/adopt/run_more_test.go","name":"TestRunTwiceSkipsAdopted"}}
func TestRunTwiceSkipsAdopted(t *testing.T) {
	cfg := newRunFixture(t, map[string]string{"a.py": "def f():\n    return 1\n"})
	if _, err := Run(cfg); err != nil {
		t.Fatalf("first Run: %v", err)
	}
	rep, err := Run(cfg)
	if err != nil {
		t.Fatalf("second Run: %v", err)
	}
	if rep.Totals.Adopted != 0 {
		t.Errorf("second run Adopted = %d, want 0", rep.Totals.Adopted)
	}
	found := false
	for _, fr := range rep.Files {
		if filepath.Base(fr.Path) == "a.py" {
			found = true
			if fr.Reason != ReasonAlreadyAdopted {
				t.Errorf("a.py reason = %v, want %v", fr.Reason, ReasonAlreadyAdopted)
			}
		}
	}
	if !found {
		t.Errorf("a.py missing from second report")
	}
}

// /sedum:setFunction
