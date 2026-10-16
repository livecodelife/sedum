package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"github.com/livecodelife/provenance/lint"}}
import "github.com/livecodelife/provenance/lint"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/check_run.go","path":"github.com/livecodelife/provenance/scope"}}
import "github.com/livecodelife/provenance/scope"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type violation struct {\n\tRecordID string\n\tFile     string\n\tMessage  string\n}","file":"cmd/provenance/check_run.go","name":"violation"}}
type violation struct {
	RecordID string
	File     string
	Message  string
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func stagedTree(r *repo) (string, func(), error) {\n\tdir, err := os.MkdirTemp(\"\", \"provenance-staged-\")\n\tif err != nil {\n\t\treturn \"\", nil, err\n\t}\n\tif _, err := r.git(\"checkout-index\", \"-a\", \"--prefix=\"+dir+\"/\"); err != nil {\n\t\tos.RemoveAll(dir)\n\t\treturn \"\", nil, err\n\t}\n\treturn dir, func() { os.RemoveAll(dir) }, nil\n}","file":"cmd/provenance/check_run.go","name":"stagedTree"}}
func stagedTree(r *repo) (string, func(), error) {
	dir, err := os.MkdirTemp("", "provenance-staged-")
	if err != nil {
		return "", nil, err
	}
	if _, err := r.git("checkout-index", "-a", "--prefix="+dir+"/"); err != nil {
		os.RemoveAll(dir)
		return "", nil, err
	}
	return dir, func() { os.RemoveAll(dir) }, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func stagedPaths(r *repo) ([]string, error) {\n\tout, err := r.git(\"diff\", \"--cached\", \"--name-only\", \"--no-renames\")\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tvar files []string\n\tfor _, line := range strings.Split(out, \"\\n\") {\n\t\tif line != \"\" {\n\t\t\tfiles = append(files, line)\n\t\t}\n\t}\n\treturn files, nil\n}","file":"cmd/provenance/check_run.go","name":"stagedPaths"}}
func stagedPaths(r *repo) ([]string, error) {
	out, err := r.git("diff", "--cached", "--name-only", "--no-renames")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, line := range strings.Split(out, "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func headRecord(r *repo, path string) (*record.Record, bool) {\n\tdata, err := r.git(\"show\", \"HEAD:\"+path)\n\tif err != nil {\n\t\treturn nil, false\n\t}\n\tdoc, err := record.Parse([]byte(data + \"\\n\"))\n\tif err != nil {\n\t\treturn nil, false\n\t}\n\treturn doc.Record(), true\n}","file":"cmd/provenance/check_run.go","name":"headRecord"}}
func headRecord(r *repo, path string) (*record.Record, bool) {
	data, err := r.git("show", "HEAD:"+path)
	if err != nil {
		return nil, false
	}
	doc, err := record.Parse([]byte(data + "\n"))
	if err != nil {
		return nil, false
	}
	return doc.Record(), true
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func headHas(r *repo, path string) bool {\n\t_, err := r.git(\"cat-file\", \"-e\", \"HEAD:\"+path)\n\treturn err == nil\n}","file":"cmd/provenance/check_run.go","name":"headHas"}}
func headHas(r *repo, path string) bool {
	_, err := r.git("cat-file", "-e", "HEAD:"+path)
	return err == nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func snapshotOwner(r *repo, path string) string {\n\tprefix := filepath.ToSlash(r.Config.Dir) + \"/.versions/\"\n\tif !strings.HasPrefix(path, prefix) {\n\t\treturn \"\"\n\t}\n\trest := strings.TrimPrefix(path, prefix)\n\tif i := strings.Index(rest, \"/\"); i >= 0 {\n\t\treturn rest[:i]\n\t}\n\treturn \"\"\n}","file":"cmd/provenance/check_run.go","name":"snapshotOwner"}}
func snapshotOwner(r *repo, path string) string {
	prefix := filepath.ToSlash(r.Config.Dir) + "/.versions/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	rest := strings.TrimPrefix(path, prefix)
	if i := strings.Index(rest, "/"); i >= 0 {
		return rest[:i]
	}
	return ""
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func recordFileID(r *repo, path string) string {\n\tprefix := r.Config.Dir + \"/\"\n\tif !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, \".yml\") {\n\t\treturn \"\"\n\t}\n\tname := strings.TrimSuffix(strings.TrimPrefix(path, prefix), \".yml\")\n\tif name == \"\" || strings.Contains(name, \"/\") {\n\t\treturn \"\"\n\t}\n\treturn name\n}","file":"cmd/provenance/check_run.go","name":"recordFileID"}}
func recordFileID(r *repo, path string) string {
	prefix := r.Config.Dir + "/"
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, ".yml") {
		return ""
	}
	name := strings.TrimSuffix(strings.TrimPrefix(path, prefix), ".yml")
	if name == "" || strings.Contains(name, "/") {
		return ""
	}
	return name
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func checkStaged(r *repo, message string) ([]violation, error) {\n\tfiles, err := stagedPaths(r)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\ttree, cleanup, err := stagedTree(r)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tdefer cleanup()\n\tstaged := &record.FilesStore{Root: tree, Dir: r.Config.Dir, SpecRoot: tree}\n\tvar vs []violation\n\n\tfor _, f := range files {\n\t\tif owner := snapshotOwner(r, f); owner != \"\" && headHas(r, f) {\n\t\t\tvs = append(vs, violation{RecordID: owner, Message: \"snapshot files are append-only: \" + f + \" already exists in HEAD\"})\n\t\t}\n\t}\n\n\tfor _, f := range files {\n\t\tid := recordFileID(r, f)\n\t\tif id == \"\" {\n\t\t\tcontinue\n\t\t}\n\t\tif _, err := os.Stat(filepath.Join(tree, f)); err != nil {\n\t\t\tcontinue\n\t\t}\n\t\tnow, _, err := staged.Load(id)\n\t\tif err != nil {\n\t\t\tcontinue\n\t\t}\n\t\tif prev, ok := headRecord(r, f); ok && prev.Status != now.Status {\n\t\t\tif err := lint.CheckTransition(prev.Status, now.Status, now.Version > prev.Version); err != nil {\n\t\t\t\tvs = append(vs, violation{RecordID: id, Message: fmt.Sprintf(\"%s: %v\", id, err)})\n\t\t\t}\n\t\t}\n\t}\n\n\tenforcement := r.Config.Enforcement\n\tif enforcement == \"\" {\n\t\tenforcement = \"none\"\n\t}\n\tres, err := lint.Lint(staged, lint.Options{Root: tree, Enforcement: enforcement, ExcludePaths: r.Config.ExcludePaths, ManifestPath: manifestPath})\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\ttags := lint.ParseCommitTags(message)\n\trelevant := map[string]bool{}\n\tfor _, t := range tags {\n\t\trelevant[t.ID] = true\n\t}\n\tfor _, f := range files {\n\t\tif id := recordFileID(r, f); id != \"\" {\n\t\t\trelevant[id] = true\n\t\t}\n\t}\n\tfor _, is := range res.Issues {\n\t\tif is.Severity == lint.SeverityError && (is.RecordID == \"\" || relevant[is.RecordID]) {\n\t\t\tvs = append(vs, violation{RecordID: is.RecordID, Message: \"lint \" + is.Rule + \": \" + is.Message})\n\t\t}\n\t}\n\n\tvar nonExcluded []string\n\tfor _, f := range files {\n\t\tif !scope.Excluded(f, r.Config.ExcludePaths) {\n\t\t\tnonExcluded = append(nonExcluded, f)\n\t\t}\n\t}\n\n\tif len(tags) == 0 {\n\t\tif r.Config.CommitTagRequired && len(nonExcluded) > 0 {\n\t\t\tvs = append(vs, violation{Message: \"Commit tag required but no provenance ID found in message\"})\n\t\t}\n\t\treturn vs, nil\n\t}\n\tif len(nonExcluded) == 0 {\n\t\treturn vs, nil\n\t}\n\n\tfor _, tag := range tags {\n\t\tresolved, err := lint.ResolveTag(staged, tag)\n\t\tif err != nil {\n\t\t\tif tag.Version == 0 {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tvs = append(vs, violation{RecordID: tag.ID, Message: err.Error()})\n\t\t\tcontinue\n\t\t}\n\t\trec, _, err := staged.Load(resolved.ID)\n\t\tif err != nil {\n\t\t\tcontinue\n\t\t}\n\t\tif tag.Version != 0 && tag.Version != rec.Version {\n\t\t\tvs = append(vs, violation{RecordID: rec.ID, Message: fmt.Sprintf(\"%s is at version %d in this commit: tag it [%s@%d], not %s\", rec.ID, rec.Version, rec.ID, rec.Version, tag.String())})\n\t\t\tcontinue\n\t\t}\n\t\town := r.Config.Dir + \"/\" + rec.ID + \".yml\"\n\t\tvar enforce []string\n\t\tfor _, f := range nonExcluded {\n\t\t\tif f == manifestPath || snapshotOwner(r, f) == rec.ID {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tif id := recordFileID(r, f); id != \"\" {\n\t\t\t\tif other, _, err := staged.Load(id); err == nil && other.Implements == rec.ID {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t}\n\t\t\tenforce = append(enforce, f)\n\t\t}\n\t\tif len(enforce) == 0 {\n\t\t\tcontinue\n\t\t}\n\t\tif rec.Status == record.StatusImplemented {\n\t\t\townStaged := false\n\t\t\tfor _, f := range files {\n\t\t\t\tif f == own {\n\t\t\t\t\townStaged = true\n\t\t\t\t}\n\t\t\t}\n\t\t\tcompletion := false\n\t\t\tif ownStaged {\n\t\t\t\tif prev, ok := headRecord(r, own); ok && prev.Status != record.StatusImplemented {\n\t\t\t\t\tcompletion = true\n\t\t\t\t}\n\t\t\t}\n\t\t\tif !completion {\n\t\t\t\tvs = append(vs, violation{RecordID: rec.ID, Message: rec.ID + \" is implemented: run 'provenance revise \" + rec.ID + \"' to start a new version, then tag the commit with it.\"})\n\t\t\t\tcontinue\n\t\t\t}\n\t\t}\n\t\tfor _, f := range enforce {\n\t\t\tif f == own {\n\t\t\t\tfor _, pat := range rec.ForbiddenScope {\n\t\t\t\t\tm, err := scope.Match(f, pat)\n\t\t\t\t\tif err != nil {\n\t\t\t\t\t\treturn nil, err\n\t\t\t\t\t}\n\t\t\t\t\tif m {\n\t\t\t\t\t\tvs = append(vs, violation{RecordID: rec.ID, File: f, Message: rec.ID + \" forbids changes to \" + f})\n\t\t\t\t\t\tbreak\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tok, err := scope.InScope(rec, f)\n\t\t\tif err != nil {\n\t\t\t\treturn nil, err\n\t\t\t}\n\t\t\tif !ok {\n\t\t\t\tvs = append(vs, violation{RecordID: rec.ID, File: f, Message: rec.ID + \" forbids changes to \" + f})\n\t\t\t}\n\t\t}\n\t}\n\treturn vs, nil\n}","file":"cmd/provenance/check_run.go","name":"checkStaged"}}
func checkStaged(r *repo, message string) ([]violation, error) {
	files, err := stagedPaths(r)
	if err != nil {
		return nil, err
	}
	tree, cleanup, err := stagedTree(r)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	staged := &record.FilesStore{Root: tree, Dir: r.Config.Dir, SpecRoot: tree}
	var vs []violation

	for _, f := range files {
		if owner := snapshotOwner(r, f); owner != "" && headHas(r, f) {
			vs = append(vs, violation{RecordID: owner, Message: "snapshot files are append-only: " + f + " already exists in HEAD"})
		}
	}

	for _, f := range files {
		id := recordFileID(r, f)
		if id == "" {
			continue
		}
		if _, err := os.Stat(filepath.Join(tree, f)); err != nil {
			continue
		}
		now, _, err := staged.Load(id)
		if err != nil {
			continue
		}
		if prev, ok := headRecord(r, f); ok && prev.Status != now.Status {
			if err := lint.CheckTransition(prev.Status, now.Status, now.Version > prev.Version); err != nil {
				vs = append(vs, violation{RecordID: id, Message: fmt.Sprintf("%s: %v", id, err)})
			}
		}
	}

	enforcement := r.Config.Enforcement
	if enforcement == "" {
		enforcement = "none"
	}
	res, err := lint.Lint(staged, lint.Options{Root: tree, Enforcement: enforcement, ExcludePaths: r.Config.ExcludePaths, ManifestPath: manifestPath})
	if err != nil {
		return nil, err
	}
	tags := lint.ParseCommitTags(message)
	relevant := map[string]bool{}
	for _, t := range tags {
		relevant[t.ID] = true
	}
	for _, f := range files {
		if id := recordFileID(r, f); id != "" {
			relevant[id] = true
		}
	}
	for _, is := range res.Issues {
		if is.Severity == lint.SeverityError && (is.RecordID == "" || relevant[is.RecordID]) {
			vs = append(vs, violation{RecordID: is.RecordID, Message: "lint " + is.Rule + ": " + is.Message})
		}
	}

	var nonExcluded []string
	for _, f := range files {
		if !scope.Excluded(f, r.Config.ExcludePaths) {
			nonExcluded = append(nonExcluded, f)
		}
	}

	if len(tags) == 0 {
		if r.Config.CommitTagRequired && len(nonExcluded) > 0 {
			vs = append(vs, violation{Message: "Commit tag required but no provenance ID found in message"})
		}
		return vs, nil
	}
	if len(nonExcluded) == 0 {
		return vs, nil
	}

	for _, tag := range tags {
		resolved, err := lint.ResolveTag(staged, tag)
		if err != nil {
			if tag.Version == 0 {
				continue
			}
			vs = append(vs, violation{RecordID: tag.ID, Message: err.Error()})
			continue
		}
		rec, _, err := staged.Load(resolved.ID)
		if err != nil {
			continue
		}
		if tag.Version != 0 && tag.Version != rec.Version {
			vs = append(vs, violation{RecordID: rec.ID, Message: fmt.Sprintf("%s is at version %d in this commit: tag it [%s@%d], not %s", rec.ID, rec.Version, rec.ID, rec.Version, tag.String())})
			continue
		}
		own := r.Config.Dir + "/" + rec.ID + ".yml"
		var enforce []string
		for _, f := range nonExcluded {
			if f == manifestPath || snapshotOwner(r, f) == rec.ID {
				continue
			}
			if id := recordFileID(r, f); id != "" {
				if other, _, err := staged.Load(id); err == nil && other.Implements == rec.ID {
					continue
				}
			}
			enforce = append(enforce, f)
		}
		if len(enforce) == 0 {
			continue
		}
		if rec.Status == record.StatusImplemented {
			ownStaged := false
			for _, f := range files {
				if f == own {
					ownStaged = true
				}
			}
			completion := false
			if ownStaged {
				if prev, ok := headRecord(r, own); ok && prev.Status != record.StatusImplemented {
					completion = true
				}
			}
			if !completion {
				vs = append(vs, violation{RecordID: rec.ID, Message: rec.ID + " is implemented: run 'provenance revise " + rec.ID + "' to start a new version, then tag the commit with it."})
				continue
			}
		}
		for _, f := range enforce {
			if f == own {
				for _, pat := range rec.ForbiddenScope {
					m, err := scope.Match(f, pat)
					if err != nil {
						return nil, err
					}
					if m {
						vs = append(vs, violation{RecordID: rec.ID, File: f, Message: rec.ID + " forbids changes to " + f})
						break
					}
				}
				continue
			}
			ok, err := scope.InScope(rec, f)
			if err != nil {
				return nil, err
			}
			if !ok {
				vs = append(vs, violation{RecordID: rec.ID, File: f, Message: rec.ID + " forbids changes to " + f})
			}
		}
	}
	return vs, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func printViolations(out io.Writer, vs []violation) {\n\tif len(vs) == 0 {\n\t\tfmt.Fprint(out, \"\\n✓ No forbidden scope violations in staged\\n\\n\")\n\t\treturn\n\t}\n\tfmt.Fprint(out, \"\\n✗ Forbidden scope violation in staged\\n\\n\")\n\tvar order []string\n\tgroups := map[string][]violation{}\n\tfor _, v := range vs {\n\t\tif v.RecordID == \"\" {\n\t\t\tfmt.Fprintf(out, \"  %s\\n\", v.Message)\n\t\t\tcontinue\n\t\t}\n\t\tif _, ok := groups[v.RecordID]; !ok {\n\t\t\torder = append(order, v.RecordID)\n\t\t}\n\t\tgroups[v.RecordID] = append(groups[v.RecordID], v)\n\t}\n\tfor _, id := range order {\n\t\tfmt.Fprintf(out, \"  %s forbids changes to:\\n\", id)\n\t\tfor _, v := range groups[id] {\n\t\t\tif v.File != \"\" {\n\t\t\t\tfmt.Fprintf(out, \"    · %s\\n\", v.File)\n\t\t\t} else {\n\t\t\t\tfmt.Fprintf(out, \"    %s\\n\", v.Message)\n\t\t\t}\n\t\t}\n\t\tfmt.Fprintln(out)\n\t}\n\tfmt.Fprint(out, \"  Hint: If this change is intentional, revise the record to widen its scope (provenance edit --add-scope) or create a record that covers it.\\n\\n\")\n}","file":"cmd/provenance/check_run.go","name":"printViolations"}}
func printViolations(out io.Writer, vs []violation) {
	if len(vs) == 0 {
		fmt.Fprint(out, "\n✓ No forbidden scope violations in staged\n\n")
		return
	}
	fmt.Fprint(out, "\n✗ Forbidden scope violation in staged\n\n")
	var order []string
	groups := map[string][]violation{}
	for _, v := range vs {
		if v.RecordID == "" {
			fmt.Fprintf(out, "  %s\n", v.Message)
			continue
		}
		if _, ok := groups[v.RecordID]; !ok {
			order = append(order, v.RecordID)
		}
		groups[v.RecordID] = append(groups[v.RecordID], v)
	}
	for _, id := range order {
		fmt.Fprintf(out, "  %s forbids changes to:\n", id)
		for _, v := range groups[id] {
			if v.File != "" {
				fmt.Fprintf(out, "    · %s\n", v.File)
			} else {
				fmt.Fprintf(out, "    %s\n", v.Message)
			}
		}
		fmt.Fprintln(out)
	}
	fmt.Fprint(out, "  Hint: If this change is intentional, revise the record to widen its scope (provenance edit --add-scope) or create a record that covers it.\n\n")
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-063f1395","gen":"sha256:dedd2ea59b62f6b6f4e787c893c15e2b33a4ffffd180e1f988b67ec7a96e4fe5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runCheck(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\tstaged, _ := cmd.Flags().GetBool(\"staged\")\n\tif !staged {\n\t\treturn errors.New(\"only --staged is supported (commit range checks come with stale and check in part 3)\")\n\t}\n\tmsgFile, _ := cmd.Flags().GetString(\"message-file\")\n\tvar message string\n\tif msgFile != \"\" {\n\t\tdata, err := os.ReadFile(msgFile)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tmessage = strings.TrimSpace(string(data))\n\t} else {\n\t\tmessage, err = r.git(\"log\", \"-1\", \"--format=%B\")\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\tvs, err := checkStaged(r, message)\n\tif err != nil {\n\t\tfmt.Fprintf(out, \"\\n✗ Check failed: %v\\n\\n\", err)\n\t\treturn err\n\t}\n\tprintViolations(out, vs)\n\tif len(vs) > 0 {\n\t\treturn errors.New(\"forbidden scope violations found\")\n\t}\n\treturn nil\n}","file":"cmd/provenance/check_run.go","name":"runCheck"}}
func runCheck(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}
	staged, _ := cmd.Flags().GetBool("staged")
	if !staged {
		return errors.New("only --staged is supported (commit range checks come with stale and check in part 3)")
	}
	msgFile, _ := cmd.Flags().GetString("message-file")
	var message string
	if msgFile != "" {
		data, err := os.ReadFile(msgFile)
		if err != nil {
			return err
		}
		message = strings.TrimSpace(string(data))
	} else {
		message, err = r.git("log", "-1", "--format=%B")
		if err != nil {
			return err
		}
	}
	vs, err := checkStaged(r, message)
	if err != nil {
		fmt.Fprintf(out, "\n✗ Check failed: %v\n\n", err)
		return err
	}
	printViolations(out, vs)
	if len(vs) > 0 {
		return errors.New("forbidden scope violations found")
	}
	return nil
}

// /sedum:setFunction
