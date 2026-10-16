package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build.go","path":"github.com/livecodelife/sedum/internal/recording"}}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Adopted is one adopted file together with its split spans.\ntype Adopted struct {\n\tCandidate Candidate\n\tSpans     []Span\n}","file":"internal/adopt/build.go","name":"Adopted"}}
// Adopted is one adopted file together with its split spans.
type Adopted struct {
	Candidate Candidate
	Spans     []Span
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// BuildRecording turns adopted files into a recording. Files are ordered by\n// path, and records by the first appearance of their ID in that order.\nfunc BuildRecording(items []Adopted) (recording.Recording, error) {\n\tsorted := make([]Adopted, len(items))\n\tcopy(sorted, items)\n\tsort.SliceStable(sorted, func(i, j int) bool {\n\t\treturn sorted[i].Candidate.Path < sorted[j].Candidate.Path\n\t})\n\trec := recording.Recording{\n\t\tSedumVersion: recording.Version,\n\t\tPackages:     map[string]recording.Package{},\n\t}\n\tindex := map[string]int{}\n\tfor i, it := range sorted {\n\t\tc := it.Candidate\n\t\tif i > 0 && sorted[i-1].Candidate.Path == c.Path {\n\t\t\treturn recording.Recording{}, fmt.Errorf(\"duplicate path %s\", c.Path)\n\t\t}\n\t\tif c.Record == nil {\n\t\t\treturn recording.Recording{}, fmt.Errorf(\"%s: no record\", c.Path)\n\t\t}\n\t\tvar invs []recording.Invocation\n\t\tfor _, s := range it.Spans {\n\t\t\tif s.Kind != KindUnit {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tkw := make(map[string]any, len(s.Kwargs))\n\t\t\tfor k, v := range s.Kwargs {\n\t\t\t\tkw[k] = v\n\t\t\t}\n\t\t\tinvs = append(invs, recording.Invocation{Action: s.Action, Kwargs: kw})\n\t\t}\n\t\tif len(invs) == 0 {\n\t\t\treturn recording.Recording{}, fmt.Errorf(\"%s: no unit spans\", c.Path)\n\t\t}\n\t\tpkgName := \"\"\n\t\tif c.Package != nil {\n\t\t\tpkgName = c.Package.Name\n\t\t\tif _, ok := rec.Packages[pkgName]; !ok {\n\t\t\t\trec.Packages[pkgName] = recording.Package{\n\t\t\t\t\tExtensions: append([]string(nil), c.Package.Extensions...),\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t\tvar caps map[string]string\n\t\tif len(c.Captures) > 0 {\n\t\t\tcaps = make(map[string]string, len(c.Captures))\n\t\t\tfor k, v := range c.Captures {\n\t\t\t\tcaps[k] = v\n\t\t\t}\n\t\t}\n\t\tidx, ok := index[c.Record.ID]\n\t\tif !ok {\n\t\t\tidx = len(rec.Records)\n\t\t\tindex[c.Record.ID] = idx\n\t\t\trec.Records = append(rec.Records, recording.Record{\n\t\t\t\tRecordID:  c.Record.ID,\n\t\t\t\tGen:       c.Record.Gen,\n\t\t\t\tCreatedAt: c.Record.CreatedAt,\n\t\t\t})\n\t\t}\n\t\tr := &rec.Records[idx]\n\t\tr.Files = append(r.Files, recording.File{\n\t\t\tPath:     c.Path,\n\t\t\tPackage:  pkgName,\n\t\t\tTemplate: c.Template,\n\t\t\tCaptures: caps,\n\t\t})\n\t\tif len(r.Phases) == 0 {\n\t\t\tr.Phases = []recording.Phase{{Name: recording.DefaultPhase}}\n\t\t}\n\t\tr.Phases[0].Invocations = append(r.Phases[0].Invocations, invs...)\n\t}\n\treturn rec, nil\n}","file":"internal/adopt/build.go","name":"BuildRecording"}}
// BuildRecording turns adopted files into a recording. Files are ordered by
// path, and records by the first appearance of their ID in that order.
func BuildRecording(items []Adopted) (recording.Recording, error) {
	sorted := make([]Adopted, len(items))
	copy(sorted, items)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Candidate.Path < sorted[j].Candidate.Path
	})
	rec := recording.Recording{
		SedumVersion: recording.Version,
		Packages:     map[string]recording.Package{},
	}
	index := map[string]int{}
	for i, it := range sorted {
		c := it.Candidate
		if i > 0 && sorted[i-1].Candidate.Path == c.Path {
			return recording.Recording{}, fmt.Errorf("duplicate path %s", c.Path)
		}
		if c.Record == nil {
			return recording.Recording{}, fmt.Errorf("%s: no record", c.Path)
		}
		var invs []recording.Invocation
		for _, s := range it.Spans {
			if s.Kind != KindUnit {
				continue
			}
			kw := make(map[string]any, len(s.Kwargs))
			for k, v := range s.Kwargs {
				kw[k] = v
			}
			invs = append(invs, recording.Invocation{Action: s.Action, Kwargs: kw})
		}
		if len(invs) == 0 {
			return recording.Recording{}, fmt.Errorf("%s: no unit spans", c.Path)
		}
		pkgName := ""
		if c.Package != nil {
			pkgName = c.Package.Name
			if _, ok := rec.Packages[pkgName]; !ok {
				rec.Packages[pkgName] = recording.Package{
					Extensions: append([]string(nil), c.Package.Extensions...),
				}
			}
		}
		var caps map[string]string
		if len(c.Captures) > 0 {
			caps = make(map[string]string, len(c.Captures))
			for k, v := range c.Captures {
				caps[k] = v
			}
		}
		idx, ok := index[c.Record.ID]
		if !ok {
			idx = len(rec.Records)
			index[c.Record.ID] = idx
			rec.Records = append(rec.Records, recording.Record{
				RecordID:  c.Record.ID,
				Gen:       c.Record.Gen,
				CreatedAt: c.Record.CreatedAt,
			})
		}
		r := &rec.Records[idx]
		r.Files = append(r.Files, recording.File{
			Path:     c.Path,
			Package:  pkgName,
			Template: c.Template,
			Captures: caps,
		})
		if len(r.Phases) == 0 {
			r.Phases = []recording.Phase{{Name: recording.DefaultPhase}}
		}
		r.Phases[0].Invocations = append(r.Phases[0].Invocations, invs...)
	}
	return rec, nil
}

// /sedum:setFunction
