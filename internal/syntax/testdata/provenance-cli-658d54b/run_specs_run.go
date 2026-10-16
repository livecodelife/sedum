package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/run_specs_run.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/run_specs_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/run_specs_run.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/run_specs_run.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/run_specs_run.go","path":"os/exec"}}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/run_specs_run.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/run_specs_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/run_specs_run.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type specOutcome struct {\n\tPath   string\n\tStatus string\n}","file":"cmd/provenance/run_specs_run.go","name":"specOutcome"}}
type specOutcome struct {
	Path   string
	Status string
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func specCommandBase(spec record.AssociatedSpec) (base string, templated bool, skip bool) {\n\tif spec.RunCommand != \"\" {\n\t\tif strings.Contains(spec.RunCommand, \"{{path}}\") {\n\t\t\treturn \"\", true, false\n\t\t}\n\t\treturn spec.RunCommand, false, false\n\t}\n\tswitch spec.Type {\n\tcase \"linespec\":\n\t\tif _, err := os.Stat(\"./linespec\"); err == nil {\n\t\t\treturn \"./linespec test\", false, false\n\t\t}\n\t\treturn \"linespec test\", false, false\n\tcase \"rspec\":\n\t\treturn \"bundle exec rspec\", false, false\n\tcase \"pytest\":\n\t\treturn \"pytest\", false, false\n\tcase \"jest\":\n\t\treturn \"npx jest\", false, false\n\t}\n\treturn \"\", false, true\n}","file":"cmd/provenance/run_specs_run.go","name":"specCommandBase"}}
func specCommandBase(spec record.AssociatedSpec) (base string, templated bool, skip bool) {
	if spec.RunCommand != "" {
		if strings.Contains(spec.RunCommand, "{{path}}") {
			return "", true, false
		}
		return spec.RunCommand, false, false
	}
	switch spec.Type {
	case "linespec":
		if _, err := os.Stat("./linespec"); err == nil {
			return "./linespec test", false, false
		}
		return "linespec test", false, false
	case "rspec":
		return "bundle exec rspec", false, false
	case "pytest":
		return "pytest", false, false
	case "jest":
		return "npx jest", false, false
	}
	return "", false, true
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func buildSpecCommand(spec record.AssociatedSpec) (string, bool) {\n\tbase, templated, skip := specCommandBase(spec)\n\tif skip {\n\t\treturn \"\", true\n\t}\n\tif templated {\n\t\treturn strings.ReplaceAll(spec.RunCommand, \"{{path}}\", spec.Path), false\n\t}\n\treturn base + \" \" + spec.Path, false\n}","file":"cmd/provenance/run_specs_run.go","name":"buildSpecCommand"}}
func buildSpecCommand(spec record.AssociatedSpec) (string, bool) {
	base, templated, skip := specCommandBase(spec)
	if skip {
		return "", true
	}
	if templated {
		return strings.ReplaceAll(spec.RunCommand, "{{path}}", spec.Path), false
	}
	return base + " " + spec.Path, false
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func specEnv() []string {\n\tvar env []string\n\tfor _, kv := range os.Environ() {\n\t\tif strings.HasPrefix(kv, \"GIT_\") {\n\t\t\tcontinue\n\t\t}\n\t\tenv = append(env, kv)\n\t}\n\treturn env\n}","file":"cmd/provenance/run_specs_run.go","name":"specEnv"}}
func specEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "GIT_") {
			continue
		}
		env = append(env, kv)
	}
	return env
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runShell(out io.Writer, dir, command string) error {\n\tcmd := exec.Command(\"sh\", \"-c\", command)\n\tcmd.Dir = dir\n\tcmd.Env = specEnv()\n\tcmd.Stdout = out\n\tcmd.Stderr = out\n\treturn cmd.Run()\n}","file":"cmd/provenance/run_specs_run.go","name":"runShell"}}
func runShell(out io.Writer, dir, command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	cmd.Env = specEnv()
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runRecordSpecs(out io.Writer, r *repo, rec *record.Record, seen map[string]bool) (ran bool, failed bool, outcomes []specOutcome) {\n\tspecs := rec.AssociatedSpecs\n\tfor i := 0; i < len(specs); {\n\t\tspec := specs[i]\n\t\t_, skip := buildSpecCommand(spec)\n\t\tif skip {\n\t\t\tfmt.Fprintf(out, \"  · %s  (skipped — no type or run_command)\\n\", spec.Path)\n\t\t\toutcomes = append(outcomes, specOutcome{Path: spec.Path, Status: \"skipped\"})\n\t\t\ti++\n\t\t\tcontinue\n\t\t}\n\t\tbase, templated, _ := specCommandBase(spec)\n\t\tgroup := []record.AssociatedSpec{spec}\n\t\tj := i + 1\n\t\tif !templated {\n\t\t\tfor j < len(specs) {\n\t\t\t\tnext := specs[j]\n\t\t\t\tnb, nt, ns := specCommandBase(next)\n\t\t\t\tif ns || nt || nb != base {\n\t\t\t\t\tbreak\n\t\t\t\t}\n\t\t\t\tgroup = append(group, next)\n\t\t\t\tj++\n\t\t\t}\n\t\t}\n\t\tran = true\n\t\tvar kept []record.AssociatedSpec\n\t\tvar keptCmds []string\n\t\tfor _, m := range group {\n\t\t\tc, _ := buildSpecCommand(m)\n\t\t\tif seen != nil && seen[c] {\n\t\t\t\tfmt.Fprintf(out, \"  · %s  (already verified)\\n\", c)\n\t\t\t\toutcomes = append(outcomes, specOutcome{Path: m.Path, Status: \"passed\"})\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tkept = append(kept, m)\n\t\t\tkeptCmds = append(keptCmds, c)\n\t\t}\n\t\tif seen != nil {\n\t\t\tfor _, c := range keptCmds {\n\t\t\t\tseen[c] = true\n\t\t\t}\n\t\t}\n\t\tswitch len(kept) {\n\t\tcase 0:\n\t\tcase 1:\n\t\t\tfmt.Fprintf(out, \"  · %s\\n    %s\\n\", kept[0].Path, keptCmds[0])\n\t\t\tif err := runShell(out, r.Root, keptCmds[0]); err != nil {\n\t\t\t\tfmt.Fprint(out, \"    ✗ failed\\n\\n\")\n\t\t\t\tfailed = true\n\t\t\t\toutcomes = append(outcomes, specOutcome{Path: kept[0].Path, Status: \"failed\"})\n\t\t\t} else {\n\t\t\t\tfmt.Fprint(out, \"    ✓ passed\\n\\n\")\n\t\t\t\toutcomes = append(outcomes, specOutcome{Path: kept[0].Path, Status: \"passed\"})\n\t\t\t}\n\t\tdefault:\n\t\t\tpaths := make([]string, len(kept))\n\t\t\tfor k, m := range kept {\n\t\t\t\tpaths[k] = m.Path\n\t\t\t}\n\t\t\tbatch := base + \" \" + strings.Join(paths, \" \")\n\t\t\tfmt.Fprintf(out, \"  · %s\\n    %s\\n\", strings.Join(paths, \", \"), batch)\n\t\t\tif err := runShell(out, r.Root, batch); err == nil {\n\t\t\t\tfor _, p := range paths {\n\t\t\t\t\tfmt.Fprintf(out, \"    ✓ %s passed\\n\", p)\n\t\t\t\t\toutcomes = append(outcomes, specOutcome{Path: p, Status: \"passed\"})\n\t\t\t\t}\n\t\t\t\tfmt.Fprint(out, \"\\n\")\n\t\t\t} else {\n\t\t\t\tfailed = true\n\t\t\t\tfmt.Fprint(out, \"    ✗ batch failed — localizing per spec:\\n\")\n\t\t\t\tfor k, m := range kept {\n\t\t\t\t\tif err := runShell(out, r.Root, keptCmds[k]); err != nil {\n\t\t\t\t\t\tfmt.Fprintf(out, \"    ✗ %s failed\\n\", m.Path)\n\t\t\t\t\t\toutcomes = append(outcomes, specOutcome{Path: m.Path, Status: \"failed\"})\n\t\t\t\t\t} else {\n\t\t\t\t\t\tfmt.Fprintf(out, \"    ✓ %s passed\\n\", m.Path)\n\t\t\t\t\t\toutcomes = append(outcomes, specOutcome{Path: m.Path, Status: \"passed\"})\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t\tfmt.Fprint(out, \"\\n\")\n\t\t\t}\n\t\t}\n\t\ti = j\n\t}\n\treturn ran, failed, outcomes\n}","file":"cmd/provenance/run_specs_run.go","name":"runRecordSpecs"}}
func runRecordSpecs(out io.Writer, r *repo, rec *record.Record, seen map[string]bool) (ran bool, failed bool, outcomes []specOutcome) {
	specs := rec.AssociatedSpecs
	for i := 0; i < len(specs); {
		spec := specs[i]
		_, skip := buildSpecCommand(spec)
		if skip {
			fmt.Fprintf(out, "  · %s  (skipped — no type or run_command)\n", spec.Path)
			outcomes = append(outcomes, specOutcome{Path: spec.Path, Status: "skipped"})
			i++
			continue
		}
		base, templated, _ := specCommandBase(spec)
		group := []record.AssociatedSpec{spec}
		j := i + 1
		if !templated {
			for j < len(specs) {
				next := specs[j]
				nb, nt, ns := specCommandBase(next)
				if ns || nt || nb != base {
					break
				}
				group = append(group, next)
				j++
			}
		}
		ran = true
		var kept []record.AssociatedSpec
		var keptCmds []string
		for _, m := range group {
			c, _ := buildSpecCommand(m)
			if seen != nil && seen[c] {
				fmt.Fprintf(out, "  · %s  (already verified)\n", c)
				outcomes = append(outcomes, specOutcome{Path: m.Path, Status: "passed"})
				continue
			}
			kept = append(kept, m)
			keptCmds = append(keptCmds, c)
		}
		if seen != nil {
			for _, c := range keptCmds {
				seen[c] = true
			}
		}
		switch len(kept) {
		case 0:
		case 1:
			fmt.Fprintf(out, "  · %s\n    %s\n", kept[0].Path, keptCmds[0])
			if err := runShell(out, r.Root, keptCmds[0]); err != nil {
				fmt.Fprint(out, "    ✗ failed\n\n")
				failed = true
				outcomes = append(outcomes, specOutcome{Path: kept[0].Path, Status: "failed"})
			} else {
				fmt.Fprint(out, "    ✓ passed\n\n")
				outcomes = append(outcomes, specOutcome{Path: kept[0].Path, Status: "passed"})
			}
		default:
			paths := make([]string, len(kept))
			for k, m := range kept {
				paths[k] = m.Path
			}
			batch := base + " " + strings.Join(paths, " ")
			fmt.Fprintf(out, "  · %s\n    %s\n", strings.Join(paths, ", "), batch)
			if err := runShell(out, r.Root, batch); err == nil {
				for _, p := range paths {
					fmt.Fprintf(out, "    ✓ %s passed\n", p)
					outcomes = append(outcomes, specOutcome{Path: p, Status: "passed"})
				}
				fmt.Fprint(out, "\n")
			} else {
				failed = true
				fmt.Fprint(out, "    ✗ batch failed — localizing per spec:\n")
				for k, m := range kept {
					if err := runShell(out, r.Root, keptCmds[k]); err != nil {
						fmt.Fprintf(out, "    ✗ %s failed\n", m.Path)
						outcomes = append(outcomes, specOutcome{Path: m.Path, Status: "failed"})
					} else {
						fmt.Fprintf(out, "    ✓ %s passed\n", m.Path)
						outcomes = append(outcomes, specOutcome{Path: m.Path, Status: "passed"})
					}
				}
				fmt.Fprint(out, "\n")
			}
		}
		i = j
	}
	return ran, failed, outcomes
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-c31c9c59","gen":"sha256:9a18a4eb62cedc8b94c8ee92552042184ff66c2dcd6e7dd094c735af31f5de23","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runRunSpecs(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\tid, _ := cmd.Flags().GetString(\"record\")\n\tif id == \"\" {\n\t\treturn errors.New(`required flag(s) \"record\" not set`)\n\t}\n\trec, _, err := r.resolve(id)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif !r.Config.RunAssociatedSpecsOnComplete || len(rec.AssociatedSpecs) == 0 {\n\t\treturn nil\n\t}\n\tfmt.Fprintf(out, \"\\nRunning %d associated spec(s) for %s...\\n\\n\", len(rec.AssociatedSpecs), rec.ID)\n\t_, failed, _ := runRecordSpecs(out, r, rec, nil)\n\tif failed {\n\t\treturn errors.New(\"one or more associated specs failed — commit blocked\")\n\t}\n\tfmt.Fprint(out, \"All associated specs passed.\\n\\n\")\n\treturn nil\n}","file":"cmd/provenance/run_specs_run.go","name":"runRunSpecs"}}
func runRunSpecs(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}
	id, _ := cmd.Flags().GetString("record")
	if id == "" {
		return errors.New(`required flag(s) "record" not set`)
	}
	rec, _, err := r.resolve(id)
	if err != nil {
		return err
	}
	if !r.Config.RunAssociatedSpecsOnComplete || len(rec.AssociatedSpecs) == 0 {
		return nil
	}
	fmt.Fprintf(out, "\nRunning %d associated spec(s) for %s...\n\n", len(rec.AssociatedSpecs), rec.ID)
	_, failed, _ := runRecordSpecs(out, r, rec, nil)
	if failed {
		return errors.New("one or more associated specs failed — commit blocked")
	}
	fmt.Fprint(out, "All associated specs passed.\n\n")
	return nil
}

// /sedum:setFunction
