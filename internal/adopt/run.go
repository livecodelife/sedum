package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/pipeline"}}
import "github.com/livecodelife/sedum/internal/pipeline"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/recording"}}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run.go","path":"github.com/livecodelife/sedum/internal/runlog"}}
import "github.com/livecodelife/sedum/internal/runlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonCoverage means the adopted spans do not cover the file byte for byte.\nconst ReasonCoverage Reason = \"coverage\"","file":"internal/adopt/run.go","name":"ReasonCoverage"}}
// ReasonCoverage means the adopted spans do not cover the file byte for byte.
const ReasonCoverage Reason = "coverage"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonRoundTrip means the regenerated file does not verify against the original.\nconst ReasonRoundTrip Reason = \"round_trip\"","file":"internal/adopt/run.go","name":"ReasonRoundTrip"}}
// ReasonRoundTrip means the regenerated file does not verify against the original.
const ReasonRoundTrip Reason = "round_trip"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonAdopterFailed means the adopter could not split the file.\nconst ReasonAdopterFailed Reason = \"adopter_failed\"","file":"internal/adopt/run.go","name":"ReasonAdopterFailed"}}
// ReasonAdopterFailed means the adopter could not split the file.
const ReasonAdopterFailed Reason = "adopter_failed"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonNoUnits means the file contains no adoptable units.\nconst ReasonNoUnits Reason = \"no_units\"","file":"internal/adopt/run.go","name":"ReasonNoUnits"}}
// ReasonNoUnits means the file contains no adoptable units.
const ReasonNoUnits Reason = "no_units"

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Config describes one adopt run.\ntype Config struct {\n\t// Root is the project root that adopted files are read from and written under.\n\tRoot string\n\t// Generators is the directory of generator packages.\n\tGenerators string\n\t// Records is the directory of provenance records.\n\tRecords string\n\t// Paths are the paths to consider for adoption.\n\tPaths []string\n\t// Tests includes test files when true.\n\tTests bool\n\t// AllowCommand permits command adopters to run.\n\tAllowCommand bool\n\t// DryRun reports what would be adopted without touching files under Root.\n\tDryRun bool\n\t// RecordingPath is where the adoption recording is written.\n\tRecordingPath string\n}","file":"internal/adopt/run.go","name":"Config"}}
// Config describes one adopt run.
type Config struct {
	// Root is the project root that adopted files are read from and written under.
	Root string
	// Generators is the directory of generator packages.
	Generators string
	// Records is the directory of provenance records.
	Records string
	// Paths are the paths to consider for adoption.
	Paths []string
	// Tests includes test files when true.
	Tests bool
	// AllowCommand permits command adopters to run.
	AllowCommand bool
	// DryRun reports what would be adopted without touching files under Root.
	DryRun bool
	// RecordingPath is where the adoption recording is written.
	RecordingPath string
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// replaceFile atomically replaces path with data, keeping the file's permissions.\nfunc replaceFile(path string, data []byte) error {\n\tinfo, err := os.Stat(path)\n\tif err != nil {\n\t\treturn fmt.Errorf(\"stat %s: %w\", path, err)\n\t}\n\ttmp, err := os.CreateTemp(filepath.Dir(path), \".sedum-adopt-*\")\n\tif err != nil {\n\t\treturn fmt.Errorf(\"create temp file: %w\", err)\n\t}\n\ttmpName := tmp.Name()\n\tif _, err := tmp.Write(data); err != nil {\n\t\ttmp.Close()\n\t\tos.Remove(tmpName)\n\t\treturn fmt.Errorf(\"write temp file: %w\", err)\n\t}\n\tif err := tmp.Close(); err != nil {\n\t\tos.Remove(tmpName)\n\t\treturn fmt.Errorf(\"close temp file: %w\", err)\n\t}\n\tif err := os.Chmod(tmpName, info.Mode().Perm()); err != nil {\n\t\tos.Remove(tmpName)\n\t\treturn fmt.Errorf(\"chmod temp file: %w\", err)\n\t}\n\tif err := os.Rename(tmpName, path); err != nil {\n\t\tos.Remove(tmpName)\n\t\treturn fmt.Errorf(\"replace %s: %w\", path, err)\n\t}\n\treturn nil\n}","file":"internal/adopt/run.go","name":"replaceFile"}}
// replaceFile atomically replaces path with data, keeping the file's permissions.
func replaceFile(path string, data []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".sedum-adopt-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Chmod(tmpName, info.Mode().Perm()); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("chmod temp file: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Run adopts the files selected by cfg.Paths: it selects candidates, splits and\n// checks each one, records the adoption, and, unless cfg.DryRun is set, replaces\n// each file that survives a round trip through the generators.\nfunc Run(cfg Config) (*Report, error) {\n\tgens, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{})\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"load generators: %w\", err)\n\t}\n\tvar broken []string\n\tfor _, f := range findings {\n\t\tif f.Kind == genpkg.KindError {\n\t\t\tbroken = append(broken, f.String())\n\t\t}\n\t}\n\tif len(broken) > 0 {\n\t\treturn nil, fmt.Errorf(\"generator packages did not load:\\n  %s\", strings.Join(broken, \"\\n  \"))\n\t}\n\trecords, _, err := record.Load(cfg.Records, record.Options{})\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"load records: %w\", err)\n\t}\n\tcands, skips, err := Select(SelectConfig{Root: cfg.Root, Set: gens, Records: records, Tests: cfg.Tests}, cfg.Paths)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"select: %w\", err)\n\t}\n\treport := NewReport()\n\tfor _, s := range skips {\n\t\treport.Add(FileReport{Path: s.Path, Status: StatusSkipped, Reason: s.Reason, Detail: s.Detail})\n\t}\n\tadopters := make(map[string]Adopter)\n\tvar adopted []Adopted\n\tvar pending []FileReport\n\tfor _, c := range cands {\n\t\tname := c.Package.Name\n\t\tad, ok := adopters[name]\n\t\tif !ok {\n\t\t\tad, err = ForPackage(c.Package, cfg.AllowCommand)\n\t\t\tif err != nil {\n\t\t\t\treturn nil, fmt.Errorf(\"adopter for %s: %w\", name, err)\n\t\t\t}\n\t\t\tif ca, isCmd := ad.(*CommandAdopter); isCmd {\n\t\t\t\tca.Dir = cfg.Root\n\t\t\t}\n\t\t\tadopters[name] = ad\n\t\t}\n\t\tfr := FileReport{Path: c.Path, Package: name}\n\t\tif c.Record != nil {\n\t\t\tfr.Record = c.Record.ID\n\t\t}\n\t\tspans, err := ad.Split(c.Path, c.Src)\n\t\tif err != nil {\n\t\t\tfr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonAdopterFailed, err.Error()\n\t\t\treport.Add(fr)\n\t\t\tcontinue\n\t\t}\n\t\tif err := CheckSpans(c.Src, spans); err != nil {\n\t\t\tfr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonCoverage, err.Error()\n\t\t\treport.Add(fr)\n\t\t\tcontinue\n\t\t}\n\t\tif res := ResidualsOf(c.Src, spans, c.Package.Comment()); len(res) > 0 {\n\t\t\tfr.Status, fr.Reason, fr.Residuals = StatusSkipped, res[0].Reason, res\n\t\t\tfr.Detail = \"file has code outside adoptable units\"\n\t\t\treport.Add(fr)\n\t\t\tcontinue\n\t\t}\n\t\tunits := 0\n\t\tfor _, s := range spans {\n\t\t\tif s.Kind == KindUnit {\n\t\t\t\tunits++\n\t\t\t}\n\t\t}\n\t\tif units == 0 {\n\t\t\tfr.Status, fr.Reason = StatusSkipped, ReasonNoUnits\n\t\t\treport.Add(fr)\n\t\t\tcontinue\n\t\t}\n\t\tfr.Status, fr.Units = StatusAdopted, units\n\t\tpending = append(pending, fr)\n\t\tadopted = append(adopted, Adopted{Candidate: c, Spans: spans})\n\t}\n\tif len(adopted) > 0 {\n\t\trec, err := BuildRecording(adopted)\n\t\tif err != nil {\n\t\t\treturn nil, fmt.Errorf(\"build recording: %w\", err)\n\t\t}\n\t\trecPath := cfg.RecordingPath\n\t\tif recPath == \"\" {\n\t\t\trecPath = filepath.Join(cfg.Root, \".sedum\", \"adopt\", \"adopt.json\")\n\t\t}\n\t\tif err := os.MkdirAll(filepath.Dir(recPath), 0o755); err != nil {\n\t\t\treturn nil, fmt.Errorf(\"create recording directory: %w\", err)\n\t\t}\n\t\tf, err := os.Create(recPath)\n\t\tif err != nil {\n\t\t\treturn nil, fmt.Errorf(\"create recording: %w\", err)\n\t\t}\n\t\tif err := recording.Write(f, rec); err != nil {\n\t\t\tf.Close()\n\t\t\treturn nil, fmt.Errorf(\"write recording: %w\", err)\n\t\t}\n\t\tif err := f.Close(); err != nil {\n\t\t\treturn nil, fmt.Errorf(\"close recording: %w\", err)\n\t\t}\n\t\treport.Recording = recPath\n\t\tif cfg.DryRun {\n\t\t\tfor i := range pending {\n\t\t\t\tpending[i].Detail = \"dry run\"\n\t\t\t}\n\t\t} else {\n\t\t\ttmp, err := os.MkdirTemp(\"\", \"sedum-adopt-\")\n\t\t\tif err != nil {\n\t\t\t\treturn nil, fmt.Errorf(\"create output directory: %w\", err)\n\t\t\t}\n\t\t\tdefer os.RemoveAll(tmp)\n\t\t\tif _, err := pipeline.Replay(rec, pipeline.ReplayConfig{\n\t\t\t\tGenerators: cfg.Generators,\n\t\t\t\tOutput:     tmp,\n\t\t\t\tRecords:    cfg.Records,\n\t\t\t\tLog:        runlog.Discard(),\n\t\t\t}); err != nil {\n\t\t\t\treturn nil, fmt.Errorf(\"replay recording: %w\", err)\n\t\t\t}\n\t\t\ttexts := make([]string, len(adopted))\n\t\t\tfor i := range adopted {\n\t\t\t\tc := adopted[i].Candidate\n\t\t\t\tout, err := os.ReadFile(filepath.Join(tmp, c.Path))\n\t\t\t\tif err == nil {\n\t\t\t\t\terr = Verify(c.Package.Comment(), string(c.Src), string(out))\n\t\t\t\t}\n\t\t\t\tif err != nil {\n\t\t\t\t\tpending[i].Status, pending[i].Reason, pending[i].Detail = StatusRejected, ReasonRoundTrip, err.Error()\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\ttexts[i] = string(out)\n\t\t\t}\n\t\t\tvar written []int\n\t\t\tfor i := range adopted {\n\t\t\t\tif pending[i].Status != StatusAdopted {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tc := adopted[i].Candidate\n\t\t\t\tif err := replaceFile(filepath.Join(cfg.Root, c.Path), []byte(texts[i])); err != nil {\n\t\t\t\t\tfor _, j := range written {\n\t\t\t\t\t\t_ = replaceFile(filepath.Join(cfg.Root, adopted[j].Candidate.Path), adopted[j].Candidate.Src)\n\t\t\t\t\t}\n\t\t\t\t\treturn nil, fmt.Errorf(\"write %s: %w\", c.Path, err)\n\t\t\t\t}\n\t\t\t\twritten = append(written, i)\n\t\t\t}\n\t\t}\n\t\tfor _, p := range pending {\n\t\t\treport.Add(p)\n\t\t}\n\t}\n\tsort.SliceStable(report.Files, func(i, j int) bool {\n\t\treturn report.Files[i].Path < report.Files[j].Path\n\t})\n\treturn report, nil\n}","file":"internal/adopt/run.go","name":"Run"}}
// Run adopts the files selected by cfg.Paths: it selects candidates, splits and
// checks each one, records the adoption, and, unless cfg.DryRun is set, replaces
// each file that survives a round trip through the generators.
func Run(cfg Config) (*Report, error) {
	gens, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{})
	if err != nil {
		return nil, fmt.Errorf("load generators: %w", err)
	}
	var broken []string
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			broken = append(broken, f.String())
		}
	}
	if len(broken) > 0 {
		return nil, fmt.Errorf("generator packages did not load:\n  %s", strings.Join(broken, "\n  "))
	}
	records, _, err := record.Load(cfg.Records, record.Options{})
	if err != nil {
		return nil, fmt.Errorf("load records: %w", err)
	}
	cands, skips, err := Select(SelectConfig{Root: cfg.Root, Set: gens, Records: records, Tests: cfg.Tests}, cfg.Paths)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	report := NewReport()
	for _, s := range skips {
		report.Add(FileReport{Path: s.Path, Status: StatusSkipped, Reason: s.Reason, Detail: s.Detail})
	}
	adopters := make(map[string]Adopter)
	var adopted []Adopted
	var pending []FileReport
	for _, c := range cands {
		name := c.Package.Name
		ad, ok := adopters[name]
		if !ok {
			ad, err = ForPackage(c.Package, cfg.AllowCommand)
			if err != nil {
				return nil, fmt.Errorf("adopter for %s: %w", name, err)
			}
			if ca, isCmd := ad.(*CommandAdopter); isCmd {
				ca.Dir = cfg.Root
			}
			adopters[name] = ad
		}
		fr := FileReport{Path: c.Path, Package: name}
		if c.Record != nil {
			fr.Record = c.Record.ID
		}
		spans, err := ad.Split(c.Path, c.Src)
		if err != nil {
			fr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonAdopterFailed, err.Error()
			report.Add(fr)
			continue
		}
		if err := CheckSpans(c.Src, spans); err != nil {
			fr.Status, fr.Reason, fr.Detail = StatusRejected, ReasonCoverage, err.Error()
			report.Add(fr)
			continue
		}
		if res := ResidualsOf(c.Src, spans, c.Package.Comment()); len(res) > 0 {
			fr.Status, fr.Reason, fr.Residuals = StatusSkipped, res[0].Reason, res
			fr.Detail = "file has code outside adoptable units"
			report.Add(fr)
			continue
		}
		units := 0
		for _, s := range spans {
			if s.Kind == KindUnit {
				units++
			}
		}
		if units == 0 {
			fr.Status, fr.Reason = StatusSkipped, ReasonNoUnits
			report.Add(fr)
			continue
		}
		fr.Status, fr.Units = StatusAdopted, units
		pending = append(pending, fr)
		adopted = append(adopted, Adopted{Candidate: c, Spans: spans})
	}
	if len(adopted) > 0 {
		rec, err := BuildRecording(adopted)
		if err != nil {
			return nil, fmt.Errorf("build recording: %w", err)
		}
		recPath := cfg.RecordingPath
		if recPath == "" {
			recPath = filepath.Join(cfg.Root, ".sedum", "adopt", "adopt.json")
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
		report.Recording = recPath
		if cfg.DryRun {
			for i := range pending {
				pending[i].Detail = "dry run"
			}
		} else {
			tmp, err := os.MkdirTemp("", "sedum-adopt-")
			if err != nil {
				return nil, fmt.Errorf("create output directory: %w", err)
			}
			defer os.RemoveAll(tmp)
			if _, err := pipeline.Replay(rec, pipeline.ReplayConfig{
				Generators: cfg.Generators,
				Output:     tmp,
				Records:    cfg.Records,
				Log:        runlog.Discard(),
			}); err != nil {
				return nil, fmt.Errorf("replay recording: %w", err)
			}
			texts := make([]string, len(adopted))
			for i := range adopted {
				c := adopted[i].Candidate
				out, err := os.ReadFile(filepath.Join(tmp, c.Path))
				if err == nil {
					err = Verify(c.Package.Comment(), string(c.Src), string(out))
				}
				if err != nil {
					pending[i].Status, pending[i].Reason, pending[i].Detail = StatusRejected, ReasonRoundTrip, err.Error()
					continue
				}
				texts[i] = string(out)
			}
			var written []int
			for i := range adopted {
				if pending[i].Status != StatusAdopted {
					continue
				}
				c := adopted[i].Candidate
				if err := replaceFile(filepath.Join(cfg.Root, c.Path), []byte(texts[i])); err != nil {
					for _, j := range written {
						_ = replaceFile(filepath.Join(cfg.Root, adopted[j].Candidate.Path), adopted[j].Candidate.Src)
					}
					return nil, fmt.Errorf("write %s: %w", c.Path, err)
				}
				written = append(written, i)
			}
		}
		for _, p := range pending {
			report.Add(p)
		}
	}
	sort.SliceStable(report.Files, func(i, j int) bool {
		return report.Files[i].Path < report.Files[j].Path
	})
	return report, nil
}

// /sedum:setFunction
