package pipeline

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup.go","path":"path"}}
import "path"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup.go","path":"github.com/livecodelife/sedum/internal/resolve"}}
import "github.com/livecodelife/sedum/internal/resolve"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup.go","path":"github.com/livecodelife/sedum/internal/runlog"}}
import "github.com/livecodelife/sedum/internal/runlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// runCleanup records what a run created so that a failed run can remove it.\ntype runCleanup struct {\n\toutput string\n\tdirs   []string\n\tfiles  []string\n}","file":"internal/pipeline/cleanup.go","name":"runCleanup"}}
// runCleanup records what a run created so that a failed run can remove it.
type runCleanup struct {
	output string
	dirs   []string
	files  []string
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// newRunCleanup records the ancestor directories of paths that do not exist yet, deepest first.\nfunc newRunCleanup(output string, paths []string) *runCleanup {\n\tc := &runCleanup{output: output}\n\tseen := map[string]bool{}\n\tfor _, p := range paths {\n\t\tfor dir := path.Dir(p); dir != \".\" && dir != \"/\" && !seen[dir]; dir = path.Dir(dir) {\n\t\t\tseen[dir] = true\n\t\t\tif _, err := os.Stat(filepath.Join(output, filepath.FromSlash(dir))); os.IsNotExist(err) {\n\t\t\t\tc.dirs = append(c.dirs, dir)\n\t\t\t}\n\t\t}\n\t}\n\tsort.SliceStable(c.dirs, func(i, j int) bool {\n\t\treturn strings.Count(c.dirs[i], \"/\") > strings.Count(c.dirs[j], \"/\")\n\t})\n\treturn c\n}","file":"internal/pipeline/cleanup.go","name":"newRunCleanup"}}
// newRunCleanup records the ancestor directories of paths that do not exist yet, deepest first.
func newRunCleanup(output string, paths []string) *runCleanup {
	c := &runCleanup{output: output}
	seen := map[string]bool{}
	for _, p := range paths {
		for dir := path.Dir(p); dir != "." && dir != "/" && !seen[dir]; dir = path.Dir(dir) {
			seen[dir] = true
			if _, err := os.Stat(filepath.Join(output, filepath.FromSlash(dir))); os.IsNotExist(err) {
				c.dirs = append(c.dirs, dir)
			}
		}
	}
	sort.SliceStable(c.dirs, func(i, j int) bool {
		return strings.Count(c.dirs[i], "/") > strings.Count(c.dirs[j], "/")
	})
	return c
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// created records the files this run wrote that did not exist before it.\nfunc (c *runCleanup) created(files []resolve.File) {\n\tfor _, f := range files {\n\t\tif !f.Existed && !f.Unmanaged {\n\t\t\tc.files = append(c.files, f.Path)\n\t\t}\n\t}\n}","file":"internal/pipeline/cleanup.go","name":"created","receiver":"runCleanup"}}
// created records the files this run wrote that did not exist before it.
func (c *runCleanup) created(files []resolve.File) {
	for _, f := range files {
		if !f.Existed && !f.Unmanaged {
			c.files = append(c.files, f.Path)
		}
	}
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// undo removes the recorded files, then the recorded directories deepest first.\nfunc (c *runCleanup) undo(log *runlog.Log) {\n\tif log == nil {\n\t\tlog = runlog.Discard()\n\t}\n\tfor _, f := range c.files {\n\t\tif err := os.Remove(filepath.Join(c.output, filepath.FromSlash(f))); err == nil {\n\t\t\tlog.Info(\"removed file created by a failed run\", \"path\", f)\n\t\t}\n\t}\n\tfor _, d := range c.dirs {\n\t\t_ = os.Remove(filepath.Join(c.output, filepath.FromSlash(d)))\n\t}\n}","file":"internal/pipeline/cleanup.go","name":"undo","receiver":"runCleanup"}}
// undo removes the recorded files, then the recorded directories deepest first.
func (c *runCleanup) undo(log *runlog.Log) {
	if log == nil {
		log = runlog.Discard()
	}
	for _, f := range c.files {
		if err := os.Remove(filepath.Join(c.output, filepath.FromSlash(f))); err == nil {
			log.Info("removed file created by a failed run", "path", f)
		}
	}
	for _, d := range c.dirs {
		_ = os.Remove(filepath.Join(c.output, filepath.FromSlash(d)))
	}
}

// /sedum:setMethod
