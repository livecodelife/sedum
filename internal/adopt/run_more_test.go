package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/run_more_test.go","path":"os"},"kwargs":"sha256:4c9c2655139be0eb81e14ac1fc97b1f8eac1c8211f51da57be45652527118f52","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/run_more_test.go","path":"path/filepath"},"kwargs":"sha256:2e55528b790b0e0befccf8a24ed6e73400adea53ef88ef99d0908f76c8f0c68e","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/run_more_test.go","path":"testing"},"kwargs":"sha256:26f35404d8926345cc2249a71c22bea2d6b9b599e85394187f1eac5482569dd2","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/run_more_test.go","name":"TestRunDryRunTouchesNothing"},"kwargs":"sha256:9f950fd539c9fa1fa60f1595de8e35bcfa27af1b89b06a66645e05e06717e471","body":"sha256:2aa8c75ff19441122d4ea53c5c622087350d1e2640906b565c50ca1a8df479d7"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/run_more_test.go","name":"TestRunSkipsResidualFiles"},"kwargs":"sha256:a351aa89a5d70f9aa1e2f5b58d2025f945f1b6416ed129687b634a88bdde242c","body":"sha256:7ae1272eb3797786c7f0bce078e6217f5c135a0c426f6ca797ef2a140411bf08"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-8542d7f7","gen":"sha256:fdd3186e1634ebdb7501fed40677b2b2ac039c37932bb30690f9356f385a7e81","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/run_more_test.go","name":"TestRunTwiceSkipsAdopted"},"kwargs":"sha256:2d1cff9a4e2b3a1adaccaea046a8d752a443c2675a8f2be6b22f5d0d25ab127e","body":"sha256:ee1df3947a8e6d15c437c0dc87ada1b7ffe69d36b1da93984ec2e693bbfce0bf"}
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
