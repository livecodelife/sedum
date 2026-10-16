package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"github.com/livecodelife/provenance/lint"}}
import "github.com/livecodelife/provenance/lint"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/complete_run.go","path":"github.com/livecodelife/provenance/scope"}}
import "github.com/livecodelife/provenance/scope"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func commitsForRecord(r *repo, rec *record.Record) []string {\n\tvar commits []string\n\tseen := map[string]bool{}\n\tids := append([]string{rec.ID}, rec.Aliases...)\n\tfor _, id := range ids {\n\t\tout, err := r.git(\"log\", \"--all\", \"-E\", \"--grep\", `\\[`+id+`(@[0-9]+)?\\]`, \"--format=%H\")\n\t\tif err != nil {\n\t\t\tcontinue\n\t\t}\n\t\tfor _, line := range strings.Split(out, \"\\n\") {\n\t\t\tline = strings.TrimSpace(line)\n\t\t\tif line == \"\" || seen[line] {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tseen[line] = true\n\t\t\tcommits = append(commits, line)\n\t\t}\n\t}\n\treturn commits\n}","file":"cmd/provenance/complete_run.go","name":"commitsForRecord"}}
func commitsForRecord(r *repo, rec *record.Record) []string {
	var commits []string
	seen := map[string]bool{}
	ids := append([]string{rec.ID}, rec.Aliases...)
	for _, id := range ids {
		out, err := r.git("log", "--all", "-E", "--grep", `\[`+id+`(@[0-9]+)?\]`, "--format=%H")
		if err != nil {
			continue
		}
		for _, line := range strings.Split(out, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || seen[line] {
				continue
			}
			seen[line] = true
			commits = append(commits, line)
		}
	}
	return commits
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func filesInCommits(r *repo, commits []string) []string {\n\tset := map[string]bool{}\n\tfor _, c := range commits {\n\t\tout, err := r.git(\"show\", \"--name-only\", \"--format=\", c)\n\t\tif err != nil {\n\t\t\tcontinue\n\t\t}\n\t\tfor _, line := range strings.Split(out, \"\\n\") {\n\t\t\tline = strings.TrimSpace(line)\n\t\t\tif line != \"\" {\n\t\t\t\tset[line] = true\n\t\t\t}\n\t\t}\n\t}\n\tfiles := make([]string, 0, len(set))\n\tfor f := range set {\n\t\tfiles = append(files, f)\n\t}\n\tsort.Strings(files)\n\treturn files\n}","file":"cmd/provenance/complete_run.go","name":"filesInCommits"}}
func filesInCommits(r *repo, commits []string) []string {
	set := map[string]bool{}
	for _, c := range commits {
		out, err := r.git("show", "--name-only", "--format=", c)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(out, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				set[line] = true
			}
		}
	}
	files := make([]string, 0, len(set))
	for f := range set {
		files = append(files, f)
	}
	sort.Strings(files)
	return files
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func overlappingSealed(self *record.Record, all []*record.Record, files []string) []*record.Record {\n\tvar out []*record.Record\n\tfor _, rec := range all {\n\t\tif rec == self {\n\t\t\tcontinue\n\t\t}\n\t\tif rec.Status != record.StatusImplemented || rec.SealedAtSHA == \"\" {\n\t\t\tcontinue\n\t\t}\n\t\tfor _, f := range files {\n\t\t\tok, err := scope.InScope(rec, f)\n\t\t\tif err == nil && ok {\n\t\t\t\tout = append(out, rec)\n\t\t\t\tbreak\n\t\t\t}\n\t\t}\n\t}\n\treturn out\n}","file":"cmd/provenance/complete_run.go","name":"overlappingSealed"}}
func overlappingSealed(self *record.Record, all []*record.Record, files []string) []*record.Record {
	var out []*record.Record
	for _, rec := range all {
		if rec == self {
			continue
		}
		if rec.Status != record.StatusImplemented || rec.SealedAtSHA == "" {
			continue
		}
		for _, f := range files {
			ok, err := scope.InScope(rec, f)
			if err == nil && ok {
				out = append(out, rec)
				break
			}
		}
	}
	return out
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func overlapMode(v string) string {\n\tswitch v {\n\tcase \"warn\", \"off\":\n\t\treturn v\n\t}\n\treturn \"block\"\n}","file":"cmd/provenance/complete_run.go","name":"overlapMode"}}
func overlapMode(v string) string {
	switch v {
	case "warn", "off":
		return v
	}
	return "block"
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func recordList(recs []*record.Record) string {\n\tvar b strings.Builder\n\tfor _, rec := range recs {\n\t\tfmt.Fprintf(&b, \"    · %s — %s\\n\", rec.ID, rec.Title)\n\t}\n\treturn b.String()\n}","file":"cmd/provenance/complete_run.go","name":"recordList"}}
func recordList(recs []*record.Record) string {
	var b strings.Builder
	for _, rec := range recs {
		fmt.Fprintf(&b, "    · %s — %s\n", rec.ID, rec.Title)
	}
	return b.String()
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runComplete(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\tid, _ := cmd.Flags().GetString(\"record\")\n\tif id == \"\" {\n\t\treturn errors.New(`required flag(s) \"record\" not set`)\n\t}\n\tforce, _ := cmd.Flags().GetBool(\"force\")\n\trec, doc, err := r.resolve(id)\n\tif err != nil {\n\t\treturn err\n\t}\n\tids, err := r.Store.List()\n\tif err != nil {\n\t\treturn err\n\t}\n\tvar all []*record.Record\n\tfor _, other := range ids {\n\t\to, _, err := r.Store.Load(other)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tall = append(all, o)\n\t}\n\tif rec.Schema != 2 {\n\t\tr.fail(\"Cannot complete %s: it is a v1 record; run 'provenance migrate' first.\", rec.ID)\n\t\treturn errors.New(\"v1 record\")\n\t}\n\tif rec.Status == record.StatusImplemented {\n\t\tr.fail(\"Record %s is already implemented\", rec.ID)\n\t\treturn errors.New(\"already implemented\")\n\t}\n\tif err := lint.CheckTransition(rec.Status, record.StatusImplemented, false); err != nil {\n\t\tr.fail(\"Cannot complete %s: %v\", rec.ID, err)\n\t\treturn err\n\t}\n\tif !force {\n\t\tvar missing strings.Builder\n\t\tfor _, s := range rec.AssociatedSpecs {\n\t\t\tif _, err := os.Stat(filepath.Join(r.Root, s.Path)); err != nil {\n\t\t\t\tfmt.Fprintf(&missing, \"    · %s  ✗ not found\\n\", s.Path)\n\t\t\t}\n\t\t}\n\t\tif missing.Len() > 0 {\n\t\t\tr.fail(\"Cannot mark %s as implemented\\n\\n  The following associated specs do not exist on disk:\\n%s\\n  Create the missing spec files or remove them from\\n  associated_specs before completing this record.\", rec.ID, missing.String())\n\t\t\treturn errors.New(\"missing specs\")\n\t\t}\n\t}\n\thead, err := r.git(\"rev-parse\", \"HEAD\")\n\tif err != nil {\n\t\tr.fail(\"Failed to get HEAD SHA: %v\", err)\n\t\treturn err\n\t}\n\trel := r.recordPath(rec.ID)\n\tstates := snapshotFiles(filepath.Join(r.Root, rel))\n\torig := rec.Status\n\trollback := func(unstage bool) {\n\t\trestoreFiles(states)\n\t\tif unstage {\n\t\t\tr.unstage(rel)\n\t\t}\n\t}\n\trolled := fmt.Sprintf(\"  The transition was rolled back — %s is unchanged (still %s).\", rec.ID, orig)\n\trec.Status = record.StatusImplemented\n\trec.SealedAtSHA = head\n\tif err := r.Store.Save(rec, doc); err != nil {\n\t\trollback(false)\n\t\tr.fail(\"Failed to save record: %v\", err)\n\t\treturn err\n\t}\n\tres, err := lintRepo(r, \"\")\n\tif err != nil {\n\t\trollback(false)\n\t\treturn err\n\t}\n\trres := recordLint(res, rec.ID)\n\tif rres.ErrorCount > 0 {\n\t\tvar issues strings.Builder\n\t\tfor _, is := range rres.Issues {\n\t\t\tif is.Severity == lint.SeverityError {\n\t\t\t\tfmt.Fprintf(&issues, \"    · %s (%s)\\n\", is.Message, is.Field)\n\t\t\t}\n\t\t}\n\t\trollback(false)\n\t\tr.fail(\"Cannot complete %s: sealing it would make it implemented, but it does not pass validation:\\n\\n%s\\n%s\\n  Fix the errors above and run 'provenance complete' again.\", rec.ID, issues.String(), rolled)\n\t\treturn errors.New(\"record failed validation\")\n\t}\n\tvar outcomes []specOutcome\n\tif !force && r.Config.RunAssociatedSpecsOnComplete && len(rec.AssociatedSpecs) > 0 {\n\t\tfmt.Fprintf(out, \"\\nRunning %d associated spec(s) for %s...\\n\\n\", len(rec.AssociatedSpecs), rec.ID)\n\t\t_, failed, oc := runRecordSpecs(out, r, rec, nil)\n\t\toutcomes = oc\n\t\tif failed {\n\t\t\trollback(false)\n\t\t\tr.fail(\"Cannot complete %s: one or more of its own associated specs failed.\\n\\n%s\\n  Fix the failing spec(s) above and run 'provenance complete' again.\", rec.ID, rolled)\n\t\t\treturn errors.New(\"associated spec failed\")\n\t\t}\n\t}\n\tmode := overlapMode(r.Config.OverlapSpecsOnComplete)\n\tif mode != \"off\" {\n\t\ttouched := overlappingSealed(rec, all, filesInCommits(r, commitsForRecord(r, rec)))\n\t\tif len(touched) > 0 {\n\t\t\tseen := map[string]bool{}\n\t\t\tvar failedRecs, unverified []*record.Record\n\t\t\tfor _, t := range touched {\n\t\t\t\tif !r.Config.RunAssociatedSpecsOnComplete || len(t.AssociatedSpecs) == 0 {\n\t\t\t\t\tunverified = append(unverified, t)\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tfmt.Fprintf(out, \"\\nVerifying sealed record %s touched by this change (%s)...\\n\", t.ID, t.Title)\n\t\t\t\tran, failed, _ := runRecordSpecs(out, r, t, seen)\n\t\t\t\tif failed {\n\t\t\t\t\tfailedRecs = append(failedRecs, t)\n\t\t\t\t}\n\t\t\t\tif !ran {\n\t\t\t\t\tunverified = append(unverified, t)\n\t\t\t\t}\n\t\t\t}\n\t\t\tif len(failedRecs) > 0 && mode == \"block\" {\n\t\t\t\trollback(false)\n\t\t\t\tr.fail(\"Cannot complete %s: it changes files governed by a sealed record whose specs now fail.\\n\\n%s\\n  A sealed record's proven behavior appears broken by this change, so completion is blocked.\\n%s\\n  Fix the failing spec(s) above, or revise the affected record, then run 'provenance complete' again.\\n  (To downgrade this gate, set overlap_specs_on_complete: warn|off in .provenance.yml.)\", rec.ID, recordList(failedRecs), rolled)\n\t\t\t\treturn errors.New(\"overlapping sealed record spec failed\")\n\t\t\t}\n\t\t\tif len(failedRecs) > 0 {\n\t\t\t\tfmt.Fprintf(out, \"\\nℹ Governance overlap (non-blocking, warn mode): specs of sealed record(s) you touched failed:\\n\\n%s\\n  Completion proceeds because overlap_specs_on_complete is 'warn'. If the failure is real, fix it or revise the affected record.\\n\\n\", recordList(failedRecs))\n\t\t\t}\n\t\t\tif len(unverified) > 0 {\n\t\t\t\tfmt.Fprintf(out, \"\\nℹ Governance overlap (non-blocking): this change also touches files governed by sealed record(s):\\n\\n%s\\n  No runnable specs were available to auto-verify them. If your change conflicts with their intent, revise them.\\n\\n\", recordList(unverified))\n\t\t\t}\n\t\t}\n\t}\n\tif r.Config.CommitOnStatusChange {\n\t\tmsg := \"Complete provenance record \" + tag(rec.ID, rec.Version)\n\t\tif err := r.commit(msg, rel); err != nil {\n\t\t\trollback(true)\n\t\t\tr.fail(\"Cannot complete %s: the auto-commit was rejected:\\n\\n    %v\\n\\n%s Nothing was staged.\\n  Resolve the issue above (e.g. a failing pre-commit hook) and run 'provenance complete' again.\", rec.ID, err, rolled)\n\t\t\treturn err\n\t\t}\n\t}\n\tfmt.Fprintf(out, \"\\n✓ %s marked as implemented\\n\\n\", rec.ID)\n\tif outcomes != nil {\n\t\tfmt.Fprint(out, \"  Associated specs verified:\\n\")\n\t\tfor _, o := range outcomes {\n\t\t\ts := \"(skipped — no type or run_command)\"\n\t\t\tswitch o.Status {\n\t\t\tcase \"passed\":\n\t\t\t\ts = \"✓ passed\"\n\t\t\tcase \"failed\":\n\t\t\t\ts = \"✗ failed\"\n\t\t\t}\n\t\t\tfmt.Fprintf(out, \"    · %-50s %s\\n\", o.Path, s)\n\t\t}\n\t\tfmt.Fprint(out, \"\\n\")\n\t}\n\treturn nil\n}","file":"cmd/provenance/complete_run.go","name":"runComplete"}}
func runComplete(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}
	id, _ := cmd.Flags().GetString("record")
	if id == "" {
		return errors.New(`required flag(s) "record" not set`)
	}
	force, _ := cmd.Flags().GetBool("force")
	rec, doc, err := r.resolve(id)
	if err != nil {
		return err
	}
	ids, err := r.Store.List()
	if err != nil {
		return err
	}
	var all []*record.Record
	for _, other := range ids {
		o, _, err := r.Store.Load(other)
		if err != nil {
			return err
		}
		all = append(all, o)
	}
	if rec.Schema != 2 {
		r.fail("Cannot complete %s: it is a v1 record; run 'provenance migrate' first.", rec.ID)
		return errors.New("v1 record")
	}
	if rec.Status == record.StatusImplemented {
		r.fail("Record %s is already implemented", rec.ID)
		return errors.New("already implemented")
	}
	if err := lint.CheckTransition(rec.Status, record.StatusImplemented, false); err != nil {
		r.fail("Cannot complete %s: %v", rec.ID, err)
		return err
	}
	if !force {
		var missing strings.Builder
		for _, s := range rec.AssociatedSpecs {
			if _, err := os.Stat(filepath.Join(r.Root, s.Path)); err != nil {
				fmt.Fprintf(&missing, "    · %s  ✗ not found\n", s.Path)
			}
		}
		if missing.Len() > 0 {
			r.fail("Cannot mark %s as implemented\n\n  The following associated specs do not exist on disk:\n%s\n  Create the missing spec files or remove them from\n  associated_specs before completing this record.", rec.ID, missing.String())
			return errors.New("missing specs")
		}
	}
	head, err := r.git("rev-parse", "HEAD")
	if err != nil {
		r.fail("Failed to get HEAD SHA: %v", err)
		return err
	}
	rel := r.recordPath(rec.ID)
	states := snapshotFiles(filepath.Join(r.Root, rel))
	orig := rec.Status
	rollback := func(unstage bool) {
		restoreFiles(states)
		if unstage {
			r.unstage(rel)
		}
	}
	rolled := fmt.Sprintf("  The transition was rolled back — %s is unchanged (still %s).", rec.ID, orig)
	rec.Status = record.StatusImplemented
	rec.SealedAtSHA = head
	if err := r.Store.Save(rec, doc); err != nil {
		rollback(false)
		r.fail("Failed to save record: %v", err)
		return err
	}
	res, err := lintRepo(r, "")
	if err != nil {
		rollback(false)
		return err
	}
	rres := recordLint(res, rec.ID)
	if rres.ErrorCount > 0 {
		var issues strings.Builder
		for _, is := range rres.Issues {
			if is.Severity == lint.SeverityError {
				fmt.Fprintf(&issues, "    · %s (%s)\n", is.Message, is.Field)
			}
		}
		rollback(false)
		r.fail("Cannot complete %s: sealing it would make it implemented, but it does not pass validation:\n\n%s\n%s\n  Fix the errors above and run 'provenance complete' again.", rec.ID, issues.String(), rolled)
		return errors.New("record failed validation")
	}
	var outcomes []specOutcome
	if !force && r.Config.RunAssociatedSpecsOnComplete && len(rec.AssociatedSpecs) > 0 {
		fmt.Fprintf(out, "\nRunning %d associated spec(s) for %s...\n\n", len(rec.AssociatedSpecs), rec.ID)
		_, failed, oc := runRecordSpecs(out, r, rec, nil)
		outcomes = oc
		if failed {
			rollback(false)
			r.fail("Cannot complete %s: one or more of its own associated specs failed.\n\n%s\n  Fix the failing spec(s) above and run 'provenance complete' again.", rec.ID, rolled)
			return errors.New("associated spec failed")
		}
	}
	mode := overlapMode(r.Config.OverlapSpecsOnComplete)
	if mode != "off" {
		touched := overlappingSealed(rec, all, filesInCommits(r, commitsForRecord(r, rec)))
		if len(touched) > 0 {
			seen := map[string]bool{}
			var failedRecs, unverified []*record.Record
			for _, t := range touched {
				if !r.Config.RunAssociatedSpecsOnComplete || len(t.AssociatedSpecs) == 0 {
					unverified = append(unverified, t)
					continue
				}
				fmt.Fprintf(out, "\nVerifying sealed record %s touched by this change (%s)...\n", t.ID, t.Title)
				ran, failed, _ := runRecordSpecs(out, r, t, seen)
				if failed {
					failedRecs = append(failedRecs, t)
				}
				if !ran {
					unverified = append(unverified, t)
				}
			}
			if len(failedRecs) > 0 && mode == "block" {
				rollback(false)
				r.fail("Cannot complete %s: it changes files governed by a sealed record whose specs now fail.\n\n%s\n  A sealed record's proven behavior appears broken by this change, so completion is blocked.\n%s\n  Fix the failing spec(s) above, or revise the affected record, then run 'provenance complete' again.\n  (To downgrade this gate, set overlap_specs_on_complete: warn|off in .provenance.yml.)", rec.ID, recordList(failedRecs), rolled)
				return errors.New("overlapping sealed record spec failed")
			}
			if len(failedRecs) > 0 {
				fmt.Fprintf(out, "\nℹ Governance overlap (non-blocking, warn mode): specs of sealed record(s) you touched failed:\n\n%s\n  Completion proceeds because overlap_specs_on_complete is 'warn'. If the failure is real, fix it or revise the affected record.\n\n", recordList(failedRecs))
			}
			if len(unverified) > 0 {
				fmt.Fprintf(out, "\nℹ Governance overlap (non-blocking): this change also touches files governed by sealed record(s):\n\n%s\n  No runnable specs were available to auto-verify them. If your change conflicts with their intent, revise them.\n\n", recordList(unverified))
			}
		}
	}
	if r.Config.CommitOnStatusChange {
		msg := "Complete provenance record " + tag(rec.ID, rec.Version)
		if err := r.commit(msg, rel); err != nil {
			rollback(true)
			r.fail("Cannot complete %s: the auto-commit was rejected:\n\n    %v\n\n%s Nothing was staged.\n  Resolve the issue above (e.g. a failing pre-commit hook) and run 'provenance complete' again.", rec.ID, err, rolled)
			return err
		}
	}
	fmt.Fprintf(out, "\n✓ %s marked as implemented\n\n", rec.ID)
	if outcomes != nil {
		fmt.Fprint(out, "  Associated specs verified:\n")
		for _, o := range outcomes {
			s := "(skipped — no type or run_command)"
			switch o.Status {
			case "passed":
				s = "✓ passed"
			case "failed":
				s = "✗ failed"
			}
			fmt.Fprintf(out, "    · %-50s %s\n", o.Path, s)
		}
		fmt.Fprint(out, "\n")
	}
	return nil
}

// /sedum:setFunction
