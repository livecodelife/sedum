package cli

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans.go","path":"fmt"},"kwargs":"sha256:5f2c923a65ab1787b05f7c4bbbfda12f0ea4a17cc87d08887358cbba19542962","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans.go","path":"io"},"kwargs":"sha256:2d29522d94cec7fde6f6ac091495195e442084ebd6c34214fb03daca3a8f6de2","body":"sha256:ef8e8abf4ef5ee40a19dc85b49bf09c5df1fe32abfe834e0ab4c2175baf80ddf"}
import "io"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans.go","path":"path"},"kwargs":"sha256:aef622eb800177bc7e3d52345bc5a274381520d7575c223dffb15a299ea0eb77","body":"sha256:c8c46e2b4136d27afcfddd73e765db239e056b6cba04d989689431e735645f13"}
import "path"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans.go","path":"sort"},"kwargs":"sha256:357e3ae0bb84c39f4796c9ea552591b9eeb97e54f8320588dc9525baa1a3bbf1","body":"sha256:f197172617ae21ef23b1855b09aafb14571020ae81bff877d63d3ff81f4bf63c"}
import "sort"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans.go","path":"strings"},"kwargs":"sha256:f4f10706bfb5a3d5ca22f5214af2fb9db687a24a014f0f47217b9c2fe16d86f8","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans.go","path":"github.com/livecodelife/sedum/internal/pipeline"},"kwargs":"sha256:0bdc28ece2aa5ce074ab235a4194f27bf021f64078be909665f03281711fee88","body":"sha256:4dfadc322ca3e8f8dc06b68db404b84dcb9cf1254d835d4ed3eedc1206db5fe8"}
import "github.com/livecodelife/sedum/internal/pipeline"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/cli/orphans.go","name":"orphanDirs"},"kwargs":"sha256:d334327ad82d1a8bdd7b561cafc13c9a4e3e26a87322b6860a60e5a979480f78","body":"sha256:ea3b24a15847eef783ba37536ebc3277ef36463c0862482148eb371b128fd87f"}
// orphanDirs is how many directories a summary names.
const orphanDirs = 3

// /sedum:setConst
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/cli/orphans.go","name":"orphanGroup"},"kwargs":"sha256:c6f8a3258e22406fc3293a33932feb77bfe9e5be0bc5abbdb2ce999bab6c786e","body":"sha256:ff88d4a81171609b23558e49e29ce4f53086f6689504ce67708bfcb470dc8ccc"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/orphans.go","name":"newOrphanGroup"},"kwargs":"sha256:73f05df4a5e3bfd99affe0c6a07264abd655aa5b58aaf3e1c2043864e228a5eb","body":"sha256:c742b6d5753c3d1bc36773b40e73599c68c3bd3d12ad4230cb86001c10710773"}
func newOrphanGroup() *orphanGroup {
	return &orphanGroup{files: map[string]bool{}, records: map[string]bool{}, dirs: map[string]int{}}
}

// /sedum:setFunction
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/cli/orphans.go","name":"add","receiver":"orphanGroup"},"kwargs":"sha256:5d37fd10c0e364105a67db72691810dc4bda81d8ef71aa2bf76808878ccb5341","body":"sha256:5f2499c30d4f945c205e8f3c0890f5d93425bd21a4b6a514f9a5ad154642f10a"}
func (g *orphanGroup) add(o pipeline.Orphan) {
	g.regions++
	g.files[o.Path] = true
	g.records[o.Record] = true
	g.dirs[path.Dir(o.Path)]++
}

// /sedum:setMethod
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/cli/orphans.go","name":"topDirs","receiver":"orphanGroup"},"kwargs":"sha256:48be128a91c1b2f87cca2e957b68bcc4f15e61319ad43cca17f3d7aa50424ae9","body":"sha256:85d195f5cdd7c1479a1f68092c4eb676a6917c33cac7318d9076645e1dc8de7c"}
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
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/cli/orphans.go","name":"line","receiver":"orphanGroup"},"kwargs":"sha256:3cd25cbeb84a1470de2add81429a19488cda4eb9cc2da917afd8649d121e677f","body":"sha256:b74fb2648e2dc67ffd536a97f538cb33535c94fb06d7a0da240209cf2c4786f4"}
// line is the group's one-line summary: the verb and the closing clause come
// from the caller, the counts and the busiest directories from the group.
func (g *orphanGroup) line(verb, rest string) string {
	return fmt.Sprintf("%s%d orphan region(s) in %d file(s) from %d gone or deprecated record(s)%s; top directories: %s",
		verb, g.regions, len(g.files), len(g.records), rest, g.topDirs())
}

// /sedum:setMethod
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5cfbb898","gen":"sha256:31991879433df3bbf93581e15551477f28ecfbbabbfc061e2403d1dd4302478c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/orphans.go","name":"printOrphans"},"kwargs":"sha256:d9d0f1edcd5fe04d2d6cf8b98f8ef410ea004ce44b9a8e4c911fcd602516930f","body":"sha256:c25c16d5fb9614b79addf6142dcaf89f3d6a012a60123600f1ba24c5b4b8141d"}
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
