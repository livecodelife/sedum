package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/revise_run.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/revise_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/revise_run.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/revise_run.go","path":"os/exec"}}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/revise_run.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/revise_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/revise_run.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/revise_run.go","path":"io"}}
import "io"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func recordArg(cmd *cobra.Command, args []string) (string, error) {\n\tflag, _ := cmd.Flags().GetString(\"record\")\n\tif len(args) > 0 && flag != \"\" && args[0] != flag {\n\t\treturn \"\", fmt.Errorf(\"--record %s and argument %s name different records\", flag, args[0])\n\t}\n\tif len(args) > 0 {\n\t\treturn args[0], nil\n\t}\n\tif flag != \"\" {\n\t\treturn flag, nil\n\t}\n\treturn \"\", errors.New(\"a record id is required (as the argument or --record)\")\n}","file":"cmd/provenance/revise_run.go","name":"recordArg"}}
func recordArg(cmd *cobra.Command, args []string) (string, error) {
	flag, _ := cmd.Flags().GetString("record")
	if len(args) > 0 && flag != "" && args[0] != flag {
		return "", fmt.Errorf("--record %s and argument %s name different records", flag, args[0])
	}
	if len(args) > 0 {
		return args[0], nil
	}
	if flag != "" {
		return flag, nil
	}
	return "", errors.New("a record id is required (as the argument or --record)")
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func versionable(r *repo, rec *record.Record) error {\n\tif rec.Schema != 2 {\n\t\tr.fail(\"Cannot change %s: it is a v1 record; run 'provenance migrate' first.\", rec.ID)\n\t\treturn fmt.Errorf(\"%s is a v1 record\", rec.ID)\n\t}\n\tif rec.Status == record.StatusDeprecated {\n\t\tr.fail(\"Cannot change %s: it is deprecated and takes no new versions.\", rec.ID)\n\t\treturn fmt.Errorf(\"%s is deprecated\", rec.ID)\n\t}\n\treturn nil\n}","file":"cmd/provenance/revise_run.go","name":"versionable"}}
func versionable(r *repo, rec *record.Record) error {
	if rec.Schema != 2 {
		r.fail("Cannot change %s: it is a v1 record; run 'provenance migrate' first.", rec.ID)
		return fmt.Errorf("%s is a v1 record", rec.ID)
	}
	if rec.Status == record.StatusDeprecated {
		r.fail("Cannot change %s: it is deprecated and takes no new versions.", rec.ID)
		return fmt.Errorf("%s is deprecated", rec.ID)
	}
	return nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func saveVersion(out io.Writer, r *repo, rec *record.Record, doc *record.Document, reason string, fix bool) error {\n\tabs := filepath.Join(r.Root, r.recordPath(rec.ID))\n\toriginal, err := os.ReadFile(abs)\n\tif err != nil {\n\t\treturn err\n\t}\n\tencoded, err := doc.Encode(rec)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif string(encoded) == string(original) {\n\t\tfmt.Fprintf(out, \"\\n✓ No change: %s is unchanged; nothing written.\\n\\n\", rec.ID)\n\t\treturn nil\n\t}\n\tversions, err := r.Store.Versions(rec.ID)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif len(versions) == 0 {\n\t\tif err := r.Store.Save(rec, doc); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tfmt.Fprintf(out, \"\\n✓ Updated %s (a draft: open writes its version 1)\\n\\n\", rec.ID)\n\t\treturn nil\n\t}\n\tlatest, err := r.Store.Snapshot(rec.ID, versions[len(versions)-1])\n\tif err != nil {\n\t\treturn err\n\t}\n\thash, _ := rec.ContentHash(r.Store.Specs())\n\tif hash == latest.ContentHash {\n\t\tif err := r.Store.Save(rec, doc); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tfmt.Fprintf(out, \"\\n✓ Updated %s (no content change; still version %d)\\n\\n\", rec.ID, rec.Version)\n\t\treturn nil\n\t}\n\tif reason == \"\" {\n\t\tr.fail(\"A new version of %s needs a reason: pass --reason \\\"...\\\".\", rec.ID)\n\t\treturn fmt.Errorf(\"%s: no reason given\", rec.ID)\n\t}\n\treopened := rec.Status == record.StatusImplemented\n\tkind := \"\"\n\tif fix {\n\t\tkind = record.ReasonKindFix\n\t}\n\tstates := snapshotFiles(abs, filepath.Join(r.Root, manifestPath))\n\tsnap, err := record.Commit(r.Store, rec, doc, record.Bump{Reason: reason, ReasonKind: kind, Author: r.author()})\n\tif err != nil {\n\t\trestoreFiles(states)\n\t\tr.fail(\"Cannot save %s: %v\", rec.ID, err)\n\t\treturn err\n\t}\n\tif err := r.writeManifest(); err != nil {\n\t\trestoreFiles(states)\n\t\tos.Remove(filepath.Join(r.Root, r.versionsDir(rec.ID), snap.FileName()))\n\t\treturn err\n\t}\n\tfmt.Fprintf(out, \"\\n✓ %s is now version %d\", rec.ID, snap.Version)\n\tif reopened {\n\t\tfmt.Fprint(out, \" (reopened)\")\n\t}\n\tfmt.Fprintf(out, \"\\n\\n  Commit it with the tag %s\\n\\n\", tag(rec.ID, snap.Version))\n\treturn nil\n}","file":"cmd/provenance/revise_run.go","name":"saveVersion"}}
func saveVersion(out io.Writer, r *repo, rec *record.Record, doc *record.Document, reason string, fix bool) error {
	abs := filepath.Join(r.Root, r.recordPath(rec.ID))
	original, err := os.ReadFile(abs)
	if err != nil {
		return err
	}
	encoded, err := doc.Encode(rec)
	if err != nil {
		return err
	}
	if string(encoded) == string(original) {
		fmt.Fprintf(out, "\n✓ No change: %s is unchanged; nothing written.\n\n", rec.ID)
		return nil
	}
	versions, err := r.Store.Versions(rec.ID)
	if err != nil {
		return err
	}
	if len(versions) == 0 {
		if err := r.Store.Save(rec, doc); err != nil {
			return err
		}
		fmt.Fprintf(out, "\n✓ Updated %s (a draft: open writes its version 1)\n\n", rec.ID)
		return nil
	}
	latest, err := r.Store.Snapshot(rec.ID, versions[len(versions)-1])
	if err != nil {
		return err
	}
	hash, _ := rec.ContentHash(r.Store.Specs())
	if hash == latest.ContentHash {
		if err := r.Store.Save(rec, doc); err != nil {
			return err
		}
		fmt.Fprintf(out, "\n✓ Updated %s (no content change; still version %d)\n\n", rec.ID, rec.Version)
		return nil
	}
	if reason == "" {
		r.fail("A new version of %s needs a reason: pass --reason \"...\".", rec.ID)
		return fmt.Errorf("%s: no reason given", rec.ID)
	}
	reopened := rec.Status == record.StatusImplemented
	kind := ""
	if fix {
		kind = record.ReasonKindFix
	}
	states := snapshotFiles(abs, filepath.Join(r.Root, manifestPath))
	snap, err := record.Commit(r.Store, rec, doc, record.Bump{Reason: reason, ReasonKind: kind, Author: r.author()})
	if err != nil {
		restoreFiles(states)
		r.fail("Cannot save %s: %v", rec.ID, err)
		return err
	}
	if err := r.writeManifest(); err != nil {
		restoreFiles(states)
		os.Remove(filepath.Join(r.Root, r.versionsDir(rec.ID), snap.FileName()))
		return err
	}
	fmt.Fprintf(out, "\n✓ %s is now version %d", rec.ID, snap.Version)
	if reopened {
		fmt.Fprint(out, " (reopened)")
	}
	fmt.Fprintf(out, "\n\n  Commit it with the tag %s\n\n", tag(rec.ID, snap.Version))
	return nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runRevise(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\tid, err := recordArg(cmd, args)\n\tif err != nil {\n\t\treturn err\n\t}\n\trec, doc, err := r.resolve(id)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif err := versionable(r, rec); err != nil {\n\t\treturn err\n\t}\n\treason, _ := cmd.Flags().GetString(\"reason\")\n\tfix, _ := cmd.Flags().GetBool(\"fix\")\n\n\tabs := filepath.Join(r.Root, r.recordPath(rec.ID))\n\toriginal, err := os.ReadFile(abs)\n\tif err != nil {\n\t\treturn err\n\t}\n\ttmp, err := os.CreateTemp(\"\", rec.ID+\"-*.yml\")\n\tif err != nil {\n\t\treturn err\n\t}\n\tdefer os.Remove(tmp.Name())\n\tif _, err := tmp.Write(original); err != nil {\n\t\ttmp.Close()\n\t\treturn err\n\t}\n\tif err := tmp.Close(); err != nil {\n\t\treturn err\n\t}\n\n\teditor := os.Getenv(\"EDITOR\")\n\tif editor == \"\" {\n\t\teditor = \"vi\"\n\t}\n\ted := exec.Command(editor, tmp.Name())\n\ted.Stdin = os.Stdin\n\ted.Stdout = os.Stdout\n\ted.Stderr = os.Stderr\n\tif err := ed.Run(); err != nil {\n\t\treturn err\n\t}\n\n\tedited, err := os.ReadFile(tmp.Name())\n\tif err != nil {\n\t\treturn err\n\t}\n\tedDoc, err := record.Parse(edited)\n\tif err != nil {\n\t\treturn fmt.Errorf(\"%s: %w\", tmp.Name(), err)\n\t}\n\tchanged := edDoc.Record()\n\n\tchangedFields := []struct {\n\t\tkey     string\n\t\tdiffers bool\n\t}{\n\t\t{\"id\", changed.ID != rec.ID},\n\t\t{\"status\", changed.Status != rec.Status},\n\t\t{\"version\", changed.Version != rec.Version},\n\t\t{\"reason\", changed.Reason != rec.Reason},\n\t\t{\"reason_kind\", changed.ReasonKind != rec.ReasonKind},\n\t\t{\"sealed_at_sha\", changed.SealedAtSHA != rec.SealedAtSHA},\n\t\t{\"created_at\", fmt.Sprint(changed.CreatedAt) != fmt.Sprint(rec.CreatedAt)},\n\t\t{\"author\", changed.Author != rec.Author},\n\t\t{\"implements_version\", changed.ImplementsVersion != rec.ImplementsVersion},\n\t\t{\"aliases\", fmt.Sprint(changed.Aliases) != fmt.Sprint(rec.Aliases)},\n\t}\n\tfor _, c := range changedFields {\n\t\tif c.differs {\n\t\t\tr.fail(\"revise changes content fields only; %s cannot be edited here.\", c.key)\n\t\t\treturn fmt.Errorf(\"%s cannot be edited here\", c.key)\n\t\t}\n\t}\n\treturn saveVersion(out, r, changed, doc, reason, fix)\n}","file":"cmd/provenance/revise_run.go","name":"runRevise"}}
func runRevise(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}
	id, err := recordArg(cmd, args)
	if err != nil {
		return err
	}
	rec, doc, err := r.resolve(id)
	if err != nil {
		return err
	}
	if err := versionable(r, rec); err != nil {
		return err
	}
	reason, _ := cmd.Flags().GetString("reason")
	fix, _ := cmd.Flags().GetBool("fix")

	abs := filepath.Join(r.Root, r.recordPath(rec.ID))
	original, err := os.ReadFile(abs)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp("", rec.ID+"-*.yml")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(original); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	ed := exec.Command(editor, tmp.Name())
	ed.Stdin = os.Stdin
	ed.Stdout = os.Stdout
	ed.Stderr = os.Stderr
	if err := ed.Run(); err != nil {
		return err
	}

	edited, err := os.ReadFile(tmp.Name())
	if err != nil {
		return err
	}
	edDoc, err := record.Parse(edited)
	if err != nil {
		return fmt.Errorf("%s: %w", tmp.Name(), err)
	}
	changed := edDoc.Record()

	changedFields := []struct {
		key     string
		differs bool
	}{
		{"id", changed.ID != rec.ID},
		{"status", changed.Status != rec.Status},
		{"version", changed.Version != rec.Version},
		{"reason", changed.Reason != rec.Reason},
		{"reason_kind", changed.ReasonKind != rec.ReasonKind},
		{"sealed_at_sha", changed.SealedAtSHA != rec.SealedAtSHA},
		{"created_at", fmt.Sprint(changed.CreatedAt) != fmt.Sprint(rec.CreatedAt)},
		{"author", changed.Author != rec.Author},
		{"implements_version", changed.ImplementsVersion != rec.ImplementsVersion},
		{"aliases", fmt.Sprint(changed.Aliases) != fmt.Sprint(rec.Aliases)},
	}
	for _, c := range changedFields {
		if c.differs {
			r.fail("revise changes content fields only; %s cannot be edited here.", c.key)
			return fmt.Errorf("%s cannot be edited here", c.key)
		}
	}
	return saveVersion(out, r, changed, doc, reason, fix)
}

// /sedum:setFunction
