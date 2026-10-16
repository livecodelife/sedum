package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build.go","path":"github.com/livecodelife/sedum/internal/recording"}}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type Adopted struct {\n\tCandidate Candidate\n\tSpans     []Span\n}","file":"internal/adopt/build.go","name":"Adopted"}}
type Adopted struct {
	Candidate Candidate
	Spans     []Span
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func BuildRecording(items []Adopted) (recording.Recording, error) {\n\tsorted := make([]Adopted, len(items))\n\tcopy(sorted, items)\n\tsort.SliceStable(sorted, func(i, j int) bool {\n\t\treturn sorted[i].Candidate.Path < sorted[j].Candidate.Path\n\t})\n\tout := recording.Recording{\n\t\tSedumVersion: recording.Version,\n\t\tPackages:     make(map[string]recording.Package),\n\t}\n\tseen := make(map[string]bool)\n\tindex := make(map[string]int)\n\tfor _, it := range sorted {\n\t\tc := it.Candidate\n\t\tif seen[c.Path] {\n\t\t\treturn recording.Recording{}, fmt.Errorf(\"adopt: duplicate path %s\", c.Path)\n\t\t}\n\t\tseen[c.Path] = true\n\t\tif c.Record == nil {\n\t\t\treturn recording.Recording{}, fmt.Errorf(\"adopt: %s has no record\", c.Path)\n\t\t}\n\t\tvar invs []recording.Invocation\n\t\tfor _, sp := range it.Spans {\n\t\t\tif sp.Kind != KindUnit {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tkw := make(map[string]any, len(sp.Kwargs))\n\t\t\tfor k, v := range sp.Kwargs {\n\t\t\t\tkw[k] = v\n\t\t\t}\n\t\t\tinvs = append(invs, recording.Invocation{Action: sp.Action, Kwargs: kw})\n\t\t}\n\t\tif len(invs) == 0 {\n\t\t\treturn recording.Recording{}, fmt.Errorf(\"adopt: %s has no unit spans\", c.Path)\n\t\t}\n\t\tpkgName := \"\"\n\t\tif c.Package != nil {\n\t\t\tpkgName = c.Package.Name\n\t\t\tout.Packages[pkgName] = recording.Package{Extensions: append([]string(nil), c.Package.Extensions...)}\n\t\t}\n\t\tvar caps map[string]string\n\t\tif len(c.Captures) > 0 {\n\t\t\tcaps = make(map[string]string, len(c.Captures))\n\t\t\tfor k, v := range c.Captures {\n\t\t\t\tcaps[k] = v\n\t\t\t}\n\t\t}\n\t\tfor _, sp := range it.Spans {\n\t\t\tif sp.Kind != KindHeader {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tfor k, v := range sp.Kwargs {\n\t\t\t\tif caps == nil {\n\t\t\t\t\tcaps = make(map[string]string)\n\t\t\t\t}\n\t\t\t\tcaps[k] = fmt.Sprint(v)\n\t\t\t}\n\t\t}\n\t\tidx, ok := index[c.Record.ID]\n\t\tif !ok {\n\t\t\tidx = len(out.Records)\n\t\t\tindex[c.Record.ID] = idx\n\t\t\tout.Records = append(out.Records, recording.Record{\n\t\t\t\tRecordID:  c.Record.ID,\n\t\t\t\tGen:       c.Record.Gen,\n\t\t\t\tCreatedAt: c.Record.CreatedAt,\n\t\t\t})\n\t\t}\n\t\trec := &out.Records[idx]\n\t\trec.Files = append(rec.Files, recording.File{\n\t\t\tPath:     c.Path,\n\t\t\tPackage:  pkgName,\n\t\t\tTemplate: c.Template,\n\t\t\tCaptures: caps,\n\t\t})\n\t\tif len(rec.Phases) == 0 {\n\t\t\trec.Phases = append(rec.Phases, recording.Phase{Name: recording.DefaultPhase})\n\t\t}\n\t\trec.Phases[0].Invocations = append(rec.Phases[0].Invocations, invs...)\n\t}\n\treturn out, nil\n}","file":"internal/adopt/build.go","name":"BuildRecording"}}
func BuildRecording(items []Adopted) (recording.Recording, error) {
	sorted := make([]Adopted, len(items))
	copy(sorted, items)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Candidate.Path < sorted[j].Candidate.Path
	})
	out := recording.Recording{
		SedumVersion: recording.Version,
		Packages:     make(map[string]recording.Package),
	}
	seen := make(map[string]bool)
	index := make(map[string]int)
	for _, it := range sorted {
		c := it.Candidate
		if seen[c.Path] {
			return recording.Recording{}, fmt.Errorf("adopt: duplicate path %s", c.Path)
		}
		seen[c.Path] = true
		if c.Record == nil {
			return recording.Recording{}, fmt.Errorf("adopt: %s has no record", c.Path)
		}
		var invs []recording.Invocation
		for _, sp := range it.Spans {
			if sp.Kind != KindUnit {
				continue
			}
			kw := make(map[string]any, len(sp.Kwargs))
			for k, v := range sp.Kwargs {
				kw[k] = v
			}
			invs = append(invs, recording.Invocation{Action: sp.Action, Kwargs: kw})
		}
		if len(invs) == 0 {
			return recording.Recording{}, fmt.Errorf("adopt: %s has no unit spans", c.Path)
		}
		pkgName := ""
		if c.Package != nil {
			pkgName = c.Package.Name
			out.Packages[pkgName] = recording.Package{Extensions: append([]string(nil), c.Package.Extensions...)}
		}
		var caps map[string]string
		if len(c.Captures) > 0 {
			caps = make(map[string]string, len(c.Captures))
			for k, v := range c.Captures {
				caps[k] = v
			}
		}
		for _, sp := range it.Spans {
			if sp.Kind != KindHeader {
				continue
			}
			for k, v := range sp.Kwargs {
				if caps == nil {
					caps = make(map[string]string)
				}
				caps[k] = fmt.Sprint(v)
			}
		}
		idx, ok := index[c.Record.ID]
		if !ok {
			idx = len(out.Records)
			index[c.Record.ID] = idx
			out.Records = append(out.Records, recording.Record{
				RecordID:  c.Record.ID,
				Gen:       c.Record.Gen,
				CreatedAt: c.Record.CreatedAt,
			})
		}
		rec := &out.Records[idx]
		rec.Files = append(rec.Files, recording.File{
			Path:     c.Path,
			Package:  pkgName,
			Template: c.Template,
			Captures: caps,
		})
		if len(rec.Phases) == 0 {
			rec.Phases = append(rec.Phases, recording.Phase{Name: recording.DefaultPhase})
		}
		rec.Phases[0].Invocations = append(rec.Phases[0].Invocations, invs...)
	}
	return out, nil
}

// /sedum:setFunction
