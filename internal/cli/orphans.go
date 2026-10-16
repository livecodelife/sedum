package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans.go","path":"path"}}
import "path"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans.go","path":"github.com/livecodelife/sedum/internal/pipeline"}}
import "github.com/livecodelife/sedum/internal/pipeline"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// orphanDirs is how many directories a summary names.\nconst orphanDirs = 3","file":"internal/cli/orphans.go","name":"orphanDirs"}}
// orphanDirs is how many directories a summary names.
const orphanDirs = 3

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// orphanGroup is the orphans of one run that share a fate: pruned, removable\n// by --prune-orphans, or left alone because no current record authorizes\n// their file (prov-2026-3f83335f).\ntype orphanGroup struct {\n\tregions int\n\tfiles   map[string]bool\n\trecords map[string]bool\n\tdirs    map[string]int\n}","file":"internal/cli/orphans.go","name":"orphanGroup"}}
// orphanGroup is the orphans of one run that share a fate: pruned, removable
// by --prune-orphans, or left alone because no current record authorizes
// their file (prov-2026-3f83335f).
type orphanGroup struct {
	regions int
	files   map[string]bool
	records map[string]bool
	dirs    map[string]int
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func newOrphanGroup() *orphanGroup {\n\treturn &orphanGroup{files: map[string]bool{}, records: map[string]bool{}, dirs: map[string]int{}}\n}","file":"internal/cli/orphans.go","name":"newOrphanGroup"}}
func newOrphanGroup() *orphanGroup {
	return &orphanGroup{files: map[string]bool{}, records: map[string]bool{}, dirs: map[string]int{}}
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (g *orphanGroup) add(o pipeline.Orphan) {\n\tg.regions++\n\tg.files[o.Path] = true\n\tg.records[o.Record] = true\n\tg.dirs[path.Dir(o.Path)]++\n}","file":"internal/cli/orphans.go","name":"add","receiver":"orphanGroup"}}
func (g *orphanGroup) add(o pipeline.Orphan) {
	g.regions++
	g.files[o.Path] = true
	g.records[o.Record] = true
	g.dirs[path.Dir(o.Path)]++
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// topDirs names the directories holding the most orphan regions, most first,\n// ties by path, so a reader sees where to look without reading every line.\nfunc (g *orphanGroup) topDirs() string {\n\tnames := make([]string, 0, len(g.dirs))\n\tfor dir := range g.dirs {\n\t\tnames = append(names, dir)\n\t}\n\tsort.Slice(names, func(i, j int) bool {\n\t\tif g.dirs[names[i]] != g.dirs[names[j]] {\n\t\t\treturn g.dirs[names[i]] > g.dirs[names[j]]\n\t\t}\n\t\treturn names[i] < names[j]\n\t})\n\tif len(names) > orphanDirs {\n\t\tnames = names[:orphanDirs]\n\t}\n\tparts := make([]string, len(names))\n\tfor i, dir := range names {\n\t\tparts[i] = fmt.Sprintf(\"%s (%d)\", dir, g.dirs[dir])\n\t}\n\treturn strings.Join(parts, \", \")\n}","file":"internal/cli/orphans.go","name":"topDirs","receiver":"orphanGroup"}}
// topDirs names the directories holding the most orphan regions, most first,
// ties by path, so a reader sees where to look without reading every line.
func (g *orphanGroup) topDirs() string {
	names := make([]string, 0, len(g.dirs))
	for dir := range g.dirs {
		names = append(names, dir)
	}
	sort.Slice(names, func(i, j int) bool {
		if g.dirs[names[i]] != g.dirs[names[j]] {
			return g.dirs[names[i]] > g.dirs[names[j]]
		}
		return names[i] < names[j]
	})
	if len(names) > orphanDirs {
		names = names[:orphanDirs]
	}
	parts := make([]string, len(names))
	for i, dir := range names {
		parts[i] = fmt.Sprintf("%s (%d)", dir, g.dirs[dir])
	}
	return strings.Join(parts, ", ")
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// line is the group's one-line summary: the verb and the closing clause come\n// from the caller, the counts and the busiest directories from the group.\nfunc (g *orphanGroup) line(verb, rest string) string {\n\treturn fmt.Sprintf(\"%s%d orphan region(s) in %d file(s) from %d gone or deprecated record(s)%s; top directories: %s\",\n\t\tverb, g.regions, len(g.files), len(g.records), rest, g.topDirs())\n}","file":"internal/cli/orphans.go","name":"line","receiver":"orphanGroup"}}
// line is the group's one-line summary: the verb and the closing clause come
// from the caller, the counts and the busiest directories from the group.
func (g *orphanGroup) line(verb, rest string) string {
	return fmt.Sprintf("%s%d orphan region(s) in %d file(s) from %d gone or deprecated record(s)%s; top directories: %s",
		verb, g.regions, len(g.files), len(g.records), rest, g.topDirs())
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// printOrphans reports regions whose record is no longer current, and what was\n// done with each (prov-2026-cc73cf61). It prints one summary line per fate,\n// because a repository that vendors a corpus can hold thousands of orphans and\n// a line each buries the rest of the report; list restores the line each.\nfunc printOrphans(out io.Writer, result *pipeline.Result, list bool) {\n\tif len(result.Orphans) == 0 {\n\t\treturn\n\t}\n\tfmt.Fprintln(out)\n\tif list {\n\t\tfor _, o := range result.Orphans {\n\t\t\tref := describeRef(pipeline.RegionRef{Path: o.Path, Action: o.Action, Variant: o.Variant})\n\t\t\tswitch {\n\t\t\tcase o.Pruned:\n\t\t\t\tfmt.Fprintf(out, \"pruned orphan %s  %s (record %s is gone or deprecated)\\n\", o.Path, ref, o.Record)\n\t\t\tcase !o.Authorized:\n\t\t\t\tfmt.Fprintf(out, \"orphan        %s  %s (record %s is gone or deprecated; no record authorizes this file, so it is left alone)\\n\", o.Path, ref, o.Record)\n\t\t\tdefault:\n\t\t\t\tfmt.Fprintf(out, \"orphan        %s  %s (record %s is gone or deprecated; --prune-orphans removes it)\\n\", o.Path, ref, o.Record)\n\t\t\t}\n\t\t}\n\t\treturn\n\t}\n\tpruned, removable, leftAlone := newOrphanGroup(), newOrphanGroup(), newOrphanGroup()\n\tfor _, o := range result.Orphans {\n\t\tswitch {\n\t\tcase o.Pruned:\n\t\t\tpruned.add(o)\n\t\tcase !o.Authorized:\n\t\t\tleftAlone.add(o)\n\t\tdefault:\n\t\t\tremovable.add(o)\n\t\t}\n\t}\n\tif pruned.regions > 0 {\n\t\tfmt.Fprintln(out, pruned.line(\"pruned \", \"\"))\n\t}\n\tif removable.regions > 0 {\n\t\tfmt.Fprintln(out, removable.line(\"\", \"; --prune-orphans removes them\"))\n\t}\n\tif leftAlone.regions > 0 {\n\t\tfmt.Fprintln(out, leftAlone.line(\"\", \"; left alone, no current record authorizes the file\"))\n\t}\n\tfmt.Fprintln(out, \"--orphans lists every region\")\n}","file":"internal/cli/orphans.go","name":"printOrphans"}}
// printOrphans reports regions whose record is no longer current, and what was
// done with each (prov-2026-cc73cf61). It prints one summary line per fate,
// because a repository that vendors a corpus can hold thousands of orphans and
// a line each buries the rest of the report; list restores the line each.
func printOrphans(out io.Writer, result *pipeline.Result, list bool) {
	if len(result.Orphans) == 0 {
		return
	}
	fmt.Fprintln(out)
	if list {
		for _, o := range result.Orphans {
			ref := describeRef(pipeline.RegionRef{Path: o.Path, Action: o.Action, Variant: o.Variant})
			switch {
			case o.Pruned:
				fmt.Fprintf(out, "pruned orphan %s  %s (record %s is gone or deprecated)\n", o.Path, ref, o.Record)
			case !o.Authorized:
				fmt.Fprintf(out, "orphan        %s  %s (record %s is gone or deprecated; no record authorizes this file, so it is left alone)\n", o.Path, ref, o.Record)
			default:
				fmt.Fprintf(out, "orphan        %s  %s (record %s is gone or deprecated; --prune-orphans removes it)\n", o.Path, ref, o.Record)
			}
		}
		return
	}
	pruned, removable, leftAlone := newOrphanGroup(), newOrphanGroup(), newOrphanGroup()
	for _, o := range result.Orphans {
		switch {
		case o.Pruned:
			pruned.add(o)
		case !o.Authorized:
			leftAlone.add(o)
		default:
			removable.add(o)
		}
	}
	if pruned.regions > 0 {
		fmt.Fprintln(out, pruned.line("pruned ", ""))
	}
	if removable.regions > 0 {
		fmt.Fprintln(out, removable.line("", "; --prune-orphans removes them"))
	}
	if leftAlone.regions > 0 {
		fmt.Fprintln(out, leftAlone.line("", "; left alone, no current record authorizes the file"))
	}
	fmt.Fprintln(out, "--orphans lists every region")
}

// /sedum:setFunction
