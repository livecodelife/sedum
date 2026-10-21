package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/pipeline"}}
import "github.com/livecodelife/sedum/internal/pipeline"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/recording"}}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/runlog"}}
import "github.com/livecodelife/sedum/internal/runlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonCoverage marks a file whose spans do not cover its bytes.\nconst ReasonCoverage Reason = \"coverage\"","file":"internal/adopt/run.go","name":"ReasonCoverage"}}
// ReasonCoverage marks a file whose spans do not cover its bytes.
const ReasonCoverage Reason = "coverage"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonRoundTrip marks a file whose replayed output failed verification.\nconst ReasonRoundTrip Reason = \"round_trip\"","file":"internal/adopt/run.go","name":"ReasonRoundTrip"}}
// ReasonRoundTrip marks a file whose replayed output failed verification.
const ReasonRoundTrip Reason = "round_trip"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonAdopterFailed marks a file whose adopter returned an error.\nconst ReasonAdopterFailed Reason = \"adopter_failed\"","file":"internal/adopt/run.go","name":"ReasonAdopterFailed"}}
// ReasonAdopterFailed marks a file whose adopter returned an error.
const ReasonAdopterFailed Reason = "adopter_failed"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonNoUnits marks a file with no unit spans.\nconst ReasonNoUnits Reason = \"no_units\"","file":"internal/adopt/run.go","name":"ReasonNoUnits"}}
// ReasonNoUnits marks a file with no unit spans.
const ReasonNoUnits Reason = "no_units"

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Config holds the inputs to Run.\ntype Config struct {\n\t// Root is the project root that Paths are relative to.\n\tRoot string\n\t// Generators is the directory of generator packages.\n\tGenerators string\n\t// Records is the directory of provenance records.\n\tRecords string\n\t// Paths are the files or directories to select.\n\tPaths []string\n\t// Tests includes test files when true.\n\tTests bool\n\t// AllowCommand permits command adopters when true.\n\tAllowCommand bool\n\t// DryRun reports what would change without writing files under Root.\n\tDryRun bool\n\t// RecordingPath is where the recording is written. Empty means the default.\n\tRecordingPath string\n}","file":"internal/adopt/run.go","name":"Config"}}
// Config holds the inputs to Run.
type Config struct {
	// Root is the project root that Paths are relative to.
	Root string
	// Generators is the directory of generator packages.
	Generators string
	// Records is the directory of provenance records.
	Records string
	// Paths are the files or directories to select.
	Paths []string
	// Tests includes test files when true.
	Tests bool
	// AllowCommand permits command adopters when true.
	AllowCommand bool
	// DryRun reports what would change without writing files under Root.
	DryRun bool
	// RecordingPath is where the recording is written. Empty means the default.
	RecordingPath string
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Run adopts the files selected by cfg. It loads the generator packages and\n// records, selects candidate files, splits each into spans, and checks that\n// the spans cover the file and leave no residual text. Files that pass are\n// rebuilt by replaying a recording, verified against their original text, and\n// written back under Root unless cfg.DryRun is set.\nfunc Run(cfg Config) (*Report, error) {\n\tset, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{})\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"load generators: %w\", err)\n\t}\n\tvar problems []string\n\tfor _, f := range findings {\n\t\tif f.Kind == genpkg.KindError {\n\t\t\tproblems = append(problems, f.String())\n\t\t}\n\t}\n\tif len(problems) > 0 {\n\t\treturn nil, fmt.Errorf(\"generator packages did not load:\\n  %s\", strings.Join(problems, \"\\n  \"))\n\t}\n\trecords, _, err := record.Load(cfg.Records, record.Options{})\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"load records: %w\", err)\n\t}\n\tcands, skips, err := Select(SelectConfig{Root: cfg.Root, Set: set, Records: records, Tests: cfg.Tests}, cfg.Paths)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"select files: %w\", err)\n\t}\n\trep := NewReport()\n\tfor _, s := range skips {\n\t\trep.Add(FileReport{Path: s.Path, Status: StatusSkipped, Reason: s.Reason, Detail: s.Detail})\n\t}\n\n\tadopters := make(map[string]Adopter)\n\tvar items []Adopted\n\tvar oks []FileReport\n\tfor _, c := range cands {\n\t\tfr := FileReport{Path: c.Path, Package: c.Package.Name}\n\t\tif c.Record != nil {\n\t\t\tfr.Record = c.Record.ID\n\t\t}\n\t\tad, ok := adopters[c.Package.Name]\n\t\tif !ok {\n\t\t\tad, err = ForPackage(c.Package, cfg.AllowCommand)\n\t\t\tif err != nil {\n\t\t\t\treturn nil, fmt.Errorf(\"adopter for package %s: %w\", c.Package.Name, err)\n\t\t\t}\n\t\t\tif ca, isCmd := ad.(*CommandAdopter); isCmd {\n\t\t\t\tca.Dir = cfg.Root\n\t\t\t}\n\t\t\tadopters[c.Package.Name] = ad\n\t\t}\n\t\tspans, err := ad.Split(c.Path, c.Src)\n\t\tif err != nil {\n\t\t\tfr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonAdopterFailed, err.Error()\n\t\t\trep.Add(fr)\n\t\t\tcontinue\n\t\t}\n\t\tif err := CheckSpans(c.Src, spans); err != nil {\n\t\t\tfr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonCoverage, err.Error()\n\t\t\trep.Add(fr)\n\t\t\tcontinue\n\t\t}\n\t\tif res := ResidualsOf(c.Src, spans, c.Package.Comment()); len(res) > 0 {\n\t\t\tfr.Status, fr.Reason, fr.Residuals = StatusSkipped, res[0].Reason, res\n\t\t\trep.Add(fr)\n\t\t\tcontinue\n\t\t}\n\t\tfor _, sp := range spans {\n\t\t\tif sp.Kind == KindUnit {\n\t\t\t\tfr.Units++\n\t\t\t}\n\t\t}\n\t\tif fr.Units == 0 {\n\t\t\tfr.Status, fr.Reason = StatusSkipped, ReasonNoUnits\n\t\t\trep.Add(fr)\n\t\t\tcontinue\n\t\t}\n\t\titems = append(items, Adopted{Candidate: c, Spans: spans})\n\t\toks = append(oks, fr)\n\t}\n\n\tfinish := func() (*Report, error) {\n\t\tsort.SliceStable(rep.Files, func(i, j int) bool {\n\t\t\treturn rep.Files[i].Path < rep.Files[j].Path\n\t\t})\n\t\treturn rep, nil\n\t}\n\tif len(items) == 0 {\n\t\treturn finish()\n\t}\n\n\trecPath := cfg.RecordingPath\n\tif recPath == \"\" {\n\t\trecPath = filepath.Join(cfg.Root, \".sedum\", \"adopt\", \"adopt.json\")\n\t}\n\trec, err := BuildRecording(items)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"build recording: %w\", err)\n\t}\n\tif err := os.MkdirAll(filepath.Dir(recPath), 0o755); err != nil {\n\t\treturn nil, fmt.Errorf(\"create recording directory: %w\", err)\n\t}\n\tf, err := os.Create(recPath)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"create recording: %w\", err)\n\t}\n\tif err := recording.Write(f, rec); err != nil {\n\t\tf.Close()\n\t\treturn nil, fmt.Errorf(\"write recording: %w\", err)\n\t}\n\tif err := f.Close(); err != nil {\n\t\treturn nil, fmt.Errorf(\"close recording: %w\", err)\n\t}\n\trep.Recording = recPath\n\n\tif cfg.DryRun {\n\t\tfor _, fr := range oks {\n\t\t\tfr.Status, fr.Detail = StatusAdopted, \"dry run\"\n\t\t\trep.Add(fr)\n\t\t}\n\t\treturn finish()\n\t}\n\n\tout, err := os.MkdirTemp(\"\", \"sedum-adopt-\")\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"create output directory: %w\", err)\n\t}\n\tdefer os.RemoveAll(out)\n\tif _, err := pipeline.Replay(rec, pipeline.ReplayConfig{\n\t\tGenerators:     cfg.Generators,\n\t\tRecords:        cfg.Records,\n\t\tOutput:         out,\n\t\tLog:            runlog.Discard(),\n\t\tAllowUnexposed: true,\n\t}); err != nil {\n\t\treturn nil, fmt.Errorf(\"replay recording: %w\", err)\n\t}\n\n\ttype pending struct {\n\t\tc    Candidate\n\t\tfr   FileReport\n\t\ttext string\n\t}\n\tvar ready []pending\n\tfor i, it := range items {\n\t\tc := it.Candidate\n\t\tfr := oks[i]\n\t\tdata, rerr := os.ReadFile(filepath.Join(out, c.Path))\n\t\tif rerr == nil {\n\t\t\trerr = Verify(c.Package.Comment(), string(c.Src), string(data))\n\t\t}\n\t\tif rerr != nil {\n\t\t\tfr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonRoundTrip, rerr.Error()\n\t\t\trep.Add(fr)\n\t\t\tcontinue\n\t\t}\n\t\tready = append(ready, pending{c: c, fr: fr, text: string(data)})\n\t}\n\n\treplace := func(rel string, data []byte) error {\n\t\tp := filepath.Join(cfg.Root, rel)\n\t\tinfo, err := os.Stat(p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\ttmp, err := os.CreateTemp(filepath.Dir(p), \".sedum-adopt-*\")\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\ttmpName := tmp.Name()\n\t\tif _, err := tmp.Write(data); err != nil {\n\t\t\ttmp.Close()\n\t\t\tos.Remove(tmpName)\n\t\t\treturn err\n\t\t}\n\t\tif err := tmp.Close(); err != nil {\n\t\t\tos.Remove(tmpName)\n\t\t\treturn err\n\t\t}\n\t\tif err := os.Chmod(tmpName, info.Mode().Perm()); err != nil {\n\t\t\tos.Remove(tmpName)\n\t\t\treturn err\n\t\t}\n\t\tif err := os.Rename(tmpName, p); err != nil {\n\t\t\tos.Remove(tmpName)\n\t\t\treturn err\n\t\t}\n\t\treturn nil\n\t}\n\tvar done []Candidate\n\tfor _, p := range ready {\n\t\tif err := replace(p.c.Path, []byte(p.text)); err != nil {\n\t\t\tfor _, d := range done {\n\t\t\t\t_ = replace(d.Path, d.Src)\n\t\t\t}\n\t\t\treturn nil, fmt.Errorf(\"write %s: %w\", p.c.Path, err)\n\t\t}\n\t\tdone = append(done, p.c)\n\t\tp.fr.Status = StatusAdopted\n\t\trep.Add(p.fr)\n\t}\n\treturn finish()\n}","file":"internal/adopt/run.go","name":"Run"}}
// Run adopts the files selected by cfg. It loads the generator packages and
// records, selects candidate files, splits each into spans, and checks that
// the spans cover the file and leave no residual text. Files that pass are
// rebuilt by replaying a recording, verified against their original text, and
// written back under Root unless cfg.DryRun is set.
func Run(cfg Config) (*Report, error) {
	set, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{})
	if err != nil {
		return nil, fmt.Errorf("load generators: %w", err)
	}
	var problems []string
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			problems = append(problems, f.String())
		}
	}
	if len(problems) > 0 {
		return nil, fmt.Errorf("generator packages did not load:\n  %s", strings.Join(problems, "\n  "))
	}
	records, _, err := record.Load(cfg.Records, record.Options{})
	if err != nil {
		return nil, fmt.Errorf("load records: %w", err)
	}
	cands, skips, err := Select(SelectConfig{Root: cfg.Root, Set: set, Records: records, Tests: cfg.Tests}, cfg.Paths)
	if err != nil {
		return nil, fmt.Errorf("select files: %w", err)
	}
	rep := NewReport()
	for _, s := range skips {
		rep.Add(FileReport{Path: s.Path, Status: StatusSkipped, Reason: s.Reason, Detail: s.Detail})
	}

	adopters := make(map[string]Adopter)
	var items []Adopted
	var oks []FileReport
	for _, c := range cands {
		fr := FileReport{Path: c.Path, Package: c.Package.Name}
		if c.Record != nil {
			fr.Record = c.Record.ID
		}
		ad, ok := adopters[c.Package.Name]
		if !ok {
			ad, err = ForPackage(c.Package, cfg.AllowCommand)
			if err != nil {
				return nil, fmt.Errorf("adopter for package %s: %w", c.Package.Name, err)
			}
			if ca, isCmd := ad.(*CommandAdopter); isCmd {
				ca.Dir = cfg.Root
			}
			adopters[c.Package.Name] = ad
		}
		spans, err := ad.Split(c.Path, c.Src)
		if err != nil {
			fr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonAdopterFailed, err.Error()
			rep.Add(fr)
			continue
		}
		if err := CheckSpans(c.Src, spans); err != nil {
			fr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonCoverage, err.Error()
			rep.Add(fr)
			continue
		}
		if res := ResidualsOf(c.Src, spans, c.Package.Comment()); len(res) > 0 {
			fr.Status, fr.Reason, fr.Residuals = StatusSkipped, res[0].Reason, res
			rep.Add(fr)
			continue
		}
		for _, sp := range spans {
			if sp.Kind == KindUnit {
				fr.Units++
			}
		}
		if fr.Units == 0 {
			fr.Status, fr.Reason = StatusSkipped, ReasonNoUnits
			rep.Add(fr)
			continue
		}
		items = append(items, Adopted{Candidate: c, Spans: spans})
		oks = append(oks, fr)
	}

	finish := func() (*Report, error) {
		sort.SliceStable(rep.Files, func(i, j int) bool {
			return rep.Files[i].Path < rep.Files[j].Path
		})
		return rep, nil
	}
	if len(items) == 0 {
		return finish()
	}

	recPath := cfg.RecordingPath
	if recPath == "" {
		recPath = filepath.Join(cfg.Root, ".sedum", "adopt", "adopt.json")
	}
	rec, err := BuildRecording(items)
	if err != nil {
		return nil, fmt.Errorf("build recording: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(recPath), 0o755); err != nil {
		return nil, fmt.Errorf("create recording directory: %w", err)
	}
	f, err := os.Create(recPath)
	if err != nil {
		return nil, fmt.Errorf("create recording: %w", err)
	}
	if err := recording.Write(f, rec); err != nil {
		f.Close()
		return nil, fmt.Errorf("write recording: %w", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("close recording: %w", err)
	}
	rep.Recording = recPath

	if cfg.DryRun {
		for _, fr := range oks {
			fr.Status, fr.Detail = StatusAdopted, "dry run"
			rep.Add(fr)
		}
		return finish()
	}

	out, err := os.MkdirTemp("", "sedum-adopt-")
	if err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}
	defer os.RemoveAll(out)
	if _, err := pipeline.Replay(rec, pipeline.ReplayConfig{
		Generators:     cfg.Generators,
		Records:        cfg.Records,
		Output:         out,
		Log:            runlog.Discard(),
		AllowUnexposed: true,
	}); err != nil {
		return nil, fmt.Errorf("replay recording: %w", err)
	}

	type pending struct {
		c    Candidate
		fr   FileReport
		text string
	}
	var ready []pending
	for i, it := range items {
		c := it.Candidate
		fr := oks[i]
		data, rerr := os.ReadFile(filepath.Join(out, c.Path))
		if rerr == nil {
			rerr = Verify(c.Package.Comment(), string(c.Src), string(data))
		}
		if rerr != nil {
			fr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonRoundTrip, rerr.Error()
			rep.Add(fr)
			continue
		}
		ready = append(ready, pending{c: c, fr: fr, text: string(data)})
	}

	replace := func(rel string, data []byte) error {
		p := filepath.Join(cfg.Root, rel)
		info, err := os.Stat(p)
		if err != nil {
			return err
		}
		tmp, err := os.CreateTemp(filepath.Dir(p), ".sedum-adopt-*")
		if err != nil {
			return err
		}
		tmpName := tmp.Name()
		if _, err := tmp.Write(data); err != nil {
			tmp.Close()
			os.Remove(tmpName)
			return err
		}
		if err := tmp.Close(); err != nil {
			os.Remove(tmpName)
			return err
		}
		if err := os.Chmod(tmpName, info.Mode().Perm()); err != nil {
			os.Remove(tmpName)
			return err
		}
		if err := os.Rename(tmpName, p); err != nil {
			os.Remove(tmpName)
			return err
		}
		return nil
	}
	var done []Candidate
	for _, p := range ready {
		if err := replace(p.c.Path, []byte(p.text)); err != nil {
			for _, d := range done {
				_ = replace(d.Path, d.Src)
			}
			return nil, fmt.Errorf("write %s: %w", p.c.Path, err)
		}
		done = append(done, p.c)
		p.fr.Status = StatusAdopted
		rep.Add(p.fr)
	}
	return finish()
}

// /sedum:setFunction
