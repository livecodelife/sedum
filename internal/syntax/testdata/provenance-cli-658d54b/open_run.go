package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/open_run.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/open_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/open_run.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/open_run.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/open_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/open_run.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/open_run.go","path":"io"}}
import "io"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runOpen(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tid, err := cmd.Flags().GetString(\"record\")\n\tif err != nil {\n\t\treturn err\n\t}\n\tif id == \"\" {\n\t\treturn errors.New(`required flag(s) \"record\" not set`)\n\t}\n\n\trec, doc, err := r.resolve(id)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tif rec.Schema != 2 {\n\t\tr.fail(\"Cannot open %s: it is a v1 record; run 'provenance migrate' first.\", rec.ID)\n\t\treturn errors.New(\"record is a v1 record\")\n\t}\n\n\tif rec.Status != record.StatusDraft {\n\t\tr.fail(\"Cannot open %s: record status is %q (open only transitions from draft).\\n\\n  Use 'provenance complete' to mark an open record as implemented.\", rec.ID, rec.Status)\n\t\treturn errors.New(\"record is not a draft\")\n\t}\n\n\tstates := snapshotFiles(filepath.Join(r.Root, r.recordPath(rec.ID)), filepath.Join(r.Root, manifestPath))\n\n\trec.Status = record.StatusOpen\n\tsnap, err := record.Commit(r.Store, rec, doc, record.Bump{Author: r.author()})\n\tif err != nil {\n\t\trestoreFiles(states)\n\t\tr.fail(\"Cannot open %s: %v\", rec.ID, err)\n\t\treturn err\n\t}\n\n\trollback := func() {\n\t\trestoreFiles(states)\n\t\tos.Remove(filepath.Join(r.Root, r.versionsDir(rec.ID), snap.FileName()))\n\t\tos.Remove(filepath.Join(r.Root, r.versionsDir(rec.ID)))\n\t}\n\n\tif err := r.writeManifest(); err != nil {\n\t\trollback()\n\t\treturn err\n\t}\n\n\tif r.Config.CommitOnStatusChange {\n\t\tpaths := []string{r.recordPath(rec.ID), r.versionsDir(rec.ID), manifestPath}\n\t\tmsg := \"Open provenance record \" + rec.ID + \" \" + tag(rec.ID, snap.Version)\n\t\tif err := r.commit(msg, paths...); err != nil {\n\t\t\tr.unstage(paths...)\n\t\t\trollback()\n\t\t\tr.fail(\"Cannot open %s: the auto-commit was rejected:\\n\\n    %v\\n\\n  The transition was rolled back: %s is unchanged (still draft) and nothing was staged.\\n  Resolve the issue above (e.g. a failing pre-commit hook) and run 'provenance open' again.\", rec.ID, err, rec.ID)\n\t\t\treturn err\n\t\t}\n\t}\n\n\tfmt.Fprintf(out, \"\\n✓ %s transitioned from draft → open\\n\\n\", rec.ID)\n\tfmt.Fprint(out, \"  Scope and spec enforcement now apply to this record.\\n\\n\")\n\treturn nil\n}","file":"cmd/provenance/open_run.go","name":"runOpen"}}
func runOpen(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}

	id, err := cmd.Flags().GetString("record")
	if err != nil {
		return err
	}
	if id == "" {
		return errors.New(`required flag(s) "record" not set`)
	}

	rec, doc, err := r.resolve(id)
	if err != nil {
		return err
	}

	if rec.Schema != 2 {
		r.fail("Cannot open %s: it is a v1 record; run 'provenance migrate' first.", rec.ID)
		return errors.New("record is a v1 record")
	}

	if rec.Status != record.StatusDraft {
		r.fail("Cannot open %s: record status is %q (open only transitions from draft).\n\n  Use 'provenance complete' to mark an open record as implemented.", rec.ID, rec.Status)
		return errors.New("record is not a draft")
	}

	states := snapshotFiles(filepath.Join(r.Root, r.recordPath(rec.ID)), filepath.Join(r.Root, manifestPath))

	rec.Status = record.StatusOpen
	snap, err := record.Commit(r.Store, rec, doc, record.Bump{Author: r.author()})
	if err != nil {
		restoreFiles(states)
		r.fail("Cannot open %s: %v", rec.ID, err)
		return err
	}

	rollback := func() {
		restoreFiles(states)
		os.Remove(filepath.Join(r.Root, r.versionsDir(rec.ID), snap.FileName()))
		os.Remove(filepath.Join(r.Root, r.versionsDir(rec.ID)))
	}

	if err := r.writeManifest(); err != nil {
		rollback()
		return err
	}

	if r.Config.CommitOnStatusChange {
		paths := []string{r.recordPath(rec.ID), r.versionsDir(rec.ID), manifestPath}
		msg := "Open provenance record " + rec.ID + " " + tag(rec.ID, snap.Version)
		if err := r.commit(msg, paths...); err != nil {
			r.unstage(paths...)
			rollback()
			r.fail("Cannot open %s: the auto-commit was rejected:\n\n    %v\n\n  The transition was rolled back: %s is unchanged (still draft) and nothing was staged.\n  Resolve the issue above (e.g. a failing pre-commit hook) and run 'provenance open' again.", rec.ID, err, rec.ID)
			return err
		}
	}

	fmt.Fprintf(out, "\n✓ %s transitioned from draft → open\n\n", rec.ID)
	fmt.Fprint(out, "  Scope and spec enforcement now apply to this record.\n\n")
	return nil
}

// /sedum:setFunction
