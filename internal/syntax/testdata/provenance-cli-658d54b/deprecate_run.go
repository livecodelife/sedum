package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/deprecate_run.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/deprecate_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/deprecate_run.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/deprecate_run.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/deprecate_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/deprecate_run.go","path":"github.com/livecodelife/provenance/lint"}}
import "github.com/livecodelife/provenance/lint"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/deprecate_run.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runDeprecate(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\tid, _ := cmd.Flags().GetString(\"record\")\n\tif id == \"\" {\n\t\treturn errors.New(`required flag(s) \"record\" not set`)\n\t}\n\treason, _ := cmd.Flags().GetString(\"reason\")\n\tif reason == \"\" {\n\t\treturn errors.New(`required flag(s) \"reason\" not set`)\n\t}\n\trec, doc, err := r.resolve(id)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif rec.Schema != 2 {\n\t\tr.fail(\"Cannot deprecate %s: it is a v1 record; run 'provenance migrate' first.\", rec.ID)\n\t\treturn errors.New(\"not a v2 record\")\n\t}\n\tif rec.Status == record.StatusDeprecated {\n\t\tr.fail(\"Record %s is already deprecated\", rec.ID)\n\t\treturn errors.New(\"already deprecated\")\n\t}\n\tif err := lint.CheckTransition(rec.Status, record.StatusDeprecated, false); err != nil {\n\t\tr.fail(\"Cannot deprecate %s: %v\", rec.ID, err)\n\t\treturn err\n\t}\n\trel := r.recordPath(rec.ID)\n\tstates := snapshotFiles(filepath.Join(r.Root, rel))\n\torig := rec.Status\n\trec.Status = record.StatusDeprecated\n\tif err := r.Store.Save(rec, doc); err != nil {\n\t\trestoreFiles(states)\n\t\tr.fail(\"Failed to save record: %v\", err)\n\t\treturn err\n\t}\n\tif r.Config.CommitOnStatusChange {\n\t\tmsg := \"Deprecate provenance record \" + tag(rec.ID, rec.Version) + \": \" + reason\n\t\tif err := r.commit(msg, rel); err != nil {\n\t\t\trestoreFiles(states)\n\t\t\tr.unstage(rel)\n\t\t\tr.fail(\"Cannot deprecate %s: the auto-commit was rejected:\\n\\n    %v\\n\\n  The transition was rolled back — %s is unchanged (still %s). Nothing was staged.\\n  Resolve the issue above (e.g. a failing pre-commit hook)\\n  and run 'provenance deprecate' again.\", rec.ID, err, rec.ID, orig)\n\t\t\treturn err\n\t\t}\n\t}\n\tfmt.Fprintf(out, \"\\n✓ %s marked as deprecated\\n\\n\", rec.ID)\n\treturn nil\n}","file":"cmd/provenance/deprecate_run.go","name":"runDeprecate"}}
func runDeprecate(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}
	id, _ := cmd.Flags().GetString("record")
	if id == "" {
		return errors.New(`required flag(s) "record" not set`)
	}
	reason, _ := cmd.Flags().GetString("reason")
	if reason == "" {
		return errors.New(`required flag(s) "reason" not set`)
	}
	rec, doc, err := r.resolve(id)
	if err != nil {
		return err
	}
	if rec.Schema != 2 {
		r.fail("Cannot deprecate %s: it is a v1 record; run 'provenance migrate' first.", rec.ID)
		return errors.New("not a v2 record")
	}
	if rec.Status == record.StatusDeprecated {
		r.fail("Record %s is already deprecated", rec.ID)
		return errors.New("already deprecated")
	}
	if err := lint.CheckTransition(rec.Status, record.StatusDeprecated, false); err != nil {
		r.fail("Cannot deprecate %s: %v", rec.ID, err)
		return err
	}
	rel := r.recordPath(rec.ID)
	states := snapshotFiles(filepath.Join(r.Root, rel))
	orig := rec.Status
	rec.Status = record.StatusDeprecated
	if err := r.Store.Save(rec, doc); err != nil {
		restoreFiles(states)
		r.fail("Failed to save record: %v", err)
		return err
	}
	if r.Config.CommitOnStatusChange {
		msg := "Deprecate provenance record " + tag(rec.ID, rec.Version) + ": " + reason
		if err := r.commit(msg, rel); err != nil {
			restoreFiles(states)
			r.unstage(rel)
			r.fail("Cannot deprecate %s: the auto-commit was rejected:\n\n    %v\n\n  The transition was rolled back — %s is unchanged (still %s). Nothing was staged.\n  Resolve the issue above (e.g. a failing pre-commit hook)\n  and run 'provenance deprecate' again.", rec.ID, err, rec.ID, orig)
			return err
		}
	}
	fmt.Fprintf(out, "\n✓ %s marked as deprecated\n\n", rec.ID)
	return nil
}

// /sedum:setFunction
