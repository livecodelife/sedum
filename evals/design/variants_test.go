package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/variants_test.go","path":"os"},"kwargs":"sha256:319aacf314f98447503e065a6987a342d9cbf106dbde37925f4e4a224c5af62a","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/variants_test.go","path":"path/filepath"},"kwargs":"sha256:86464b7d287af4f6fcfef8202c1e35147fdc21a96779ee1e1fbf863925c29016","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/variants_test.go","path":"strings"},"kwargs":"sha256:98f3f39d0bcd71d1f6fc21c5309872367e9e6c86a7e0e87ee3fe9b175e82bb69","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/variants_test.go","path":"testing"},"kwargs":"sha256:2e951aecf6cb84aadefa770768428aebb1377f8c0061b3185d7de06274cd3374","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/variants_test.go","name":"TestIntentVariants"},"kwargs":"sha256:0af2c7b59314a2b86c69a64e54a8ef828c22a24712ca787bceb8815c3f6560f9","body":"sha256:975e67d988278e4d543687f966f97398fb4c02434826c3aa3a095654dfad1d64"}
func TestIntentVariants(t *testing.T) {
	task := Task{Behaviour: "b", Signatures: "s", Code: "c"}
	for v, want := range map[Variant]string{VariantBehaviour: "b", VariantSignatures: "s", VariantCode: "c"} {
		got, err := task.Intent(v)
		if err != nil || got != want {
			t.Errorf("Intent(%s) = %q, %v, want %q", v, got, err, want)
		}
	}
	if _, err := task.Intent("nope"); err == nil {
		t.Fatal("an unknown variant must be an error")
	}
	if got := Variants(); len(got) != 3 || got[0] != VariantBehaviour || got[2] != VariantCode {
		t.Fatalf("Variants() = %v", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/variants_test.go","name":"TestIsolatedRealTasks"},"kwargs":"sha256:eb9d09a6ad849cf1c2044338f9c004be71402f3e8a1c14e1cef604540f22216c","body":"sha256:64cfdda106c93ae134e7865736931c847d03f7be7f91e78d891801dc7975e45c"}
func TestIsolatedRealTasks(t *testing.T) {
	tasks, err := LoadTasks("testdata/tasks")
	if err != nil {
		t.Fatal(err)
	}
	for _, task := range tasks {
		if err := task.Isolated(); err != nil {
			t.Errorf("%s: %v", task.ID, err)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/variants_test.go","name":"TestIsolatedCatchesLeaks"},"kwargs":"sha256:6de4eeb83bc5597d3425e313c8ca7e81c6e612496718459c3f10a850b4f6e42f","body":"sha256:5fdfe02ea45d8afad992c8a0aad23218dc66aaad7febe6f1005d78f68bdb9c9a"}
func TestIsolatedCatchesLeaks(t *testing.T) {
	root := t.TempDir()
	dir := writeTask(t, root, "ok", validTask)
	task, err := LoadTask(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := task.Isolated(); err != nil {
		t.Fatalf("a clean task: %v", err)
	}
	named := task
	named.Signatures = "see a_test.go"
	if err := named.Isolated(); err == nil || !strings.Contains(err.Error(), "a_test.go") {
		t.Fatalf("a named hidden file: %v", err)
	}
	leaked := task
	leaked.Behaviour = "copy this\nfunc TestA(t *testing.T) {}\n"
	if err := leaked.Isolated(); err == nil || !strings.Contains(err.Error(), "func TestA") {
		t.Fatalf("a leaked test line: %v", err)
	}
	if err := os.Remove(filepath.Join(dir, "hidden", "a_test.go")); err != nil {
		t.Fatal(err)
	}
	if err := task.Isolated(); err == nil {
		t.Fatal("an unreadable hidden file must be an error")
	}
}

// /sedum:setFunction
