package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/build.go","path":"fmt"},"kwargs":"sha256:a7518489b293b2e15f2149b6bf8c569372aa58200c6d391f4cb654c0b98ac1a4","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/build.go","path":"sort"},"kwargs":"sha256:9456642b00d4845b75f6842365a401285b386a527484b61c9360cfa1eeb31455","body":"sha256:f197172617ae21ef23b1855b09aafb14571020ae81bff877d63d3ff81f4bf63c"}
import "sort"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/build.go","path":"github.com/livecodelife/sedum/internal/recording"},"kwargs":"sha256:772270b8be800280a1450542a55c6fd02aa19269af5fed58743e2dda6be30a45","body":"sha256:77eda6b9872215a21291b42d59b9ec4b021705be869f51a3d61adc6dcb120e03"}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/build.go","name":"Adopted"},"kwargs":"sha256:ec0a32cb6cd210949419d5bc82018b6a8b125f6daf67351aa6eb9f6a7a8dc97d","body":"sha256:078db400b3ab62b3ff24310e57bd1d1e7d51513d2672bd62b3bc9fd6f33c66d6"}
// Adopted is one candidate file together with the unit spans adopted from it.
type Adopted struct {
	Candidate Candidate
	Spans     []Span
}

// /sedum:setType
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/build.go","name":"BuildRecording"},"kwargs":"sha256:cd4999a818453c3c4fb7bf456ec8ecd3a8912db565d4313f52b264704793d235","body":"sha256:f0c1ef1d42397fc520f4df27b63f52e4add32afb75daf3438d80e7f3e7c1179a"}
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
