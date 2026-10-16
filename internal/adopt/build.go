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
// sedum:setType {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Adopted is one candidate file together with the unit spans adopted from it.\ntype Adopted struct {\n\tCandidate Candidate\n\tSpans     []Span\n}","file":"internal/adopt/build.go","name":"Adopted"}}
// Adopted is one candidate file together with the unit spans adopted from it.
type Adopted struct {
	Candidate Candidate
	Spans     []Span
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// BuildRecording turns adopted files into a recording, one record per adoption record ID.\nfunc BuildRecording(items []Adopted) (recording.Recording, error) {\n\tsorted := make([]Adopted, len(items))\n\tcopy(sorted, items)\n\tsort.SliceStable(sorted, func(i, j int) bool {\n\t\treturn sorted[i].Candidate.Path < sorted[j].Candidate.Path\n\t})\n\tout := recording.Recording{\n\t\tSedumVersion: recording.Version,\n\t\tPackages:     map[string]recording.Package{},\n\t}\n\tseen := map[string]bool{}\n\tbyRecord := map[string][]Adopted{}\n\tvar order []string\n\tfor _, it := range sorted {\n\t\tc := it.Candidate\n\t\tif seen[c.Path] {\n\t\t\treturn recording.Recording{}, fmt.Errorf(\"duplicate path %s\", c.Path)\n\t\t}\n\t\tseen[c.Path] = true\n\t\tif c.Record == nil {\n\t\t\treturn recording.Recording{}, fmt.Errorf(\"%s: no record\", c.Path)\n\t\t}\n\t\tunits := 0\n\t\tfor _, s := range it.Spans {\n\t\t\tif s.Kind == KindUnit {\n\t\t\t\tunits++\n\t\t\t}\n\t\t}\n\t\tif units == 0 {\n\t\t\treturn recording.Recording{}, fmt.Errorf(\"%s: no unit spans\", c.Path)\n\t\t}\n\t\tif c.Package != nil {\n\t\t\tout.Packages[c.Package.Name] = recording.Package{\n\t\t\t\tExtensions: append([]string(nil), c.Package.Extensions...),\n\t\t\t}\n\t\t}\n\t\tid := c.Record.ID\n\t\tif _, ok := byRecord[id]; !ok {\n\t\t\torder = append(order, id)\n\t\t}\n\t\tbyRecord[id] = append(byRecord[id], it)\n\t}\n\tfor _, id := range order {\n\t\tgroup := byRecord[id]\n\t\tfirst := group[0].Candidate.Record\n\t\trec := recording.Record{RecordID: id, Gen: first.Gen, CreatedAt: first.CreatedAt}\n\t\tphase := recording.Phase{Name: recording.DefaultPhase}\n\t\tfor _, it := range group {\n\t\t\tc := it.Candidate\n\t\t\tfile := recording.File{Path: c.Path, Template: c.Template}\n\t\t\tif c.Package != nil {\n\t\t\t\tfile.Package = c.Package.Name\n\t\t\t}\n\t\t\tif c.Captures != nil {\n\t\t\t\tfile.Captures = make(map[string]string, len(c.Captures))\n\t\t\t\tfor k, v := range c.Captures {\n\t\t\t\t\tfile.Captures[k] = v\n\t\t\t\t}\n\t\t\t}\n\t\t\trec.Files = append(rec.Files, file)\n\t\t\tfor _, s := range it.Spans {\n\t\t\t\tif s.Kind != KindUnit {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tkw := make(map[string]any, len(s.Kwargs))\n\t\t\t\tfor k, v := range s.Kwargs {\n\t\t\t\t\tkw[k] = v\n\t\t\t\t}\n\t\t\t\tphase.Invocations = append(phase.Invocations, recording.Invocation{Action: s.Action, Kwargs: kw})\n\t\t\t}\n\t\t}\n\t\trec.Phases = []recording.Phase{phase}\n\t\tout.Records = append(out.Records, rec)\n\t}\n\treturn out, nil\n}","file":"internal/adopt/build.go","name":"BuildRecording"}}
// BuildRecording turns adopted files into a recording, one record per adoption record ID.
func BuildRecording(items []Adopted) (recording.Recording, error) {
	sorted := make([]Adopted, len(items))
	copy(sorted, items)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Candidate.Path < sorted[j].Candidate.Path
	})
	out := recording.Recording{
		SedumVersion: recording.Version,
		Packages:     map[string]recording.Package{},
	}
	seen := map[string]bool{}
	byRecord := map[string][]Adopted{}
	var order []string
	for _, it := range sorted {
		c := it.Candidate
		if seen[c.Path] {
			return recording.Recording{}, fmt.Errorf("duplicate path %s", c.Path)
		}
		seen[c.Path] = true
		if c.Record == nil {
			return recording.Recording{}, fmt.Errorf("%s: no record", c.Path)
		}
		units := 0
		for _, s := range it.Spans {
			if s.Kind == KindUnit {
				units++
			}
		}
		if units == 0 {
			return recording.Recording{}, fmt.Errorf("%s: no unit spans", c.Path)
		}
		if c.Package != nil {
			out.Packages[c.Package.Name] = recording.Package{
				Extensions: append([]string(nil), c.Package.Extensions...),
			}
		}
		id := c.Record.ID
		if _, ok := byRecord[id]; !ok {
			order = append(order, id)
		}
		byRecord[id] = append(byRecord[id], it)
	}
	for _, id := range order {
		group := byRecord[id]
		first := group[0].Candidate.Record
		rec := recording.Record{RecordID: id, Gen: first.Gen, CreatedAt: first.CreatedAt}
		phase := recording.Phase{Name: recording.DefaultPhase}
		for _, it := range group {
			c := it.Candidate
			file := recording.File{Path: c.Path, Template: c.Template}
			if c.Package != nil {
				file.Package = c.Package.Name
			}
			if c.Captures != nil {
				file.Captures = make(map[string]string, len(c.Captures))
				for k, v := range c.Captures {
					file.Captures[k] = v
				}
			}
			rec.Files = append(rec.Files, file)
			for _, s := range it.Spans {
				if s.Kind != KindUnit {
					continue
				}
				kw := make(map[string]any, len(s.Kwargs))
				for k, v := range s.Kwargs {
					kw[k] = v
				}
				phase.Invocations = append(phase.Invocations, recording.Invocation{Action: s.Action, Kwargs: kw})
			}
		}
		rec.Phases = []recording.Phase{phase}
		out.Records = append(out.Records, rec)
	}
	return out, nil
}

// /sedum:setFunction
