package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"time"}}
import "time"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"github.com/livecodelife/sedum/internal/growlog"}}
import "github.com/livecodelife/sedum/internal/growlog"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"github.com/livecodelife/sedum/internal/pipeline"}}
import "github.com/livecodelife/sedum/internal/pipeline"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"github.com/livecodelife/sedum/internal/recording"}}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire.go","path":"github.com/livecodelife/sedum/internal/selection"}}
import "github.com/livecodelife/sedum/internal/selection"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// defaultGrowLogPath is where sedum grow appends its persistent grow log.\nconst defaultGrowLogPath = \".sedum/grow.jsonl\"","file":"internal/cli/growlogwire.go","name":"defaultGrowLogPath"}}
// defaultGrowLogPath is where sedum grow appends its persistent grow log.
const defaultGrowLogPath = ".sedum/grow.jsonl"

// /sedum:setConst
// sedum:setFunction {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// growLogSettings reads the grow log's switches from the environment:\n// SEDUM_GROW_LOG names the file, or is \"off\" to write no log, and\n// SEDUM_GROW_LOG_SKIPS=1 also logs records skipped as current.\nfunc growLogSettings(getenv func(string) string) (path string, on, skips bool) {\n\tpath = getenv(\"SEDUM_GROW_LOG\")\n\tif path == \"off\" {\n\t\treturn \"\", false, false\n\t}\n\tif path == \"\" {\n\t\tpath = defaultGrowLogPath\n\t}\n\treturn path, true, getenv(\"SEDUM_GROW_LOG_SKIPS\") == \"1\"\n}","file":"internal/cli/growlogwire.go","name":"growLogSettings"}}
// growLogSettings reads the grow log's switches from the environment:
// SEDUM_GROW_LOG names the file, or is "off" to write no log, and
// SEDUM_GROW_LOG_SKIPS=1 also logs records skipped as current.
func growLogSettings(getenv func(string) string) (path string, on, skips bool) {
	path = getenv("SEDUM_GROW_LOG")
	if path == "off" {
		return "", false, false
	}
	if path == "" {
		path = defaultGrowLogPath
	}
	return path, true, getenv("SEDUM_GROW_LOG_SKIPS") == "1"
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// growLogModel names the model a grow used, as the log records it.\nfunc growLogModel(cfg GrowConfig) string {\n\tswitch {\n\tcase cfg.LocalModel != \"\":\n\t\treturn \"local:\" + filepath.Base(cfg.LocalModel)\n\tcase cfg.Model != \"\":\n\t\treturn cfg.Model\n\t}\n\treturn \"bundled\"\n}","file":"internal/cli/growlogwire.go","name":"growLogModel"}}
// growLogModel names the model a grow used, as the log records it.
func growLogModel(cfg GrowConfig) string {
	switch {
	case cfg.LocalModel != "":
		return "local:" + filepath.Base(cfg.LocalModel)
	case cfg.Model != "":
		return cfg.Model
	}
	return "bundled"
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// growLogMode maps a pipeline mode to the log's mode names.\nfunc growLogMode(m pipeline.Mode) string {\n\tswitch m {\n\tcase pipeline.ModeSkip:\n\t\treturn \"skip\"\n\tcase pipeline.ModeRerender:\n\t\treturn \"rerender\"\n\t}\n\treturn \"select\"\n}","file":"internal/cli/growlogwire.go","name":"growLogMode"}}
// growLogMode maps a pipeline mode to the log's mode names.
func growLogMode(m pipeline.Mode) string {
	switch m {
	case pipeline.ModeSkip:
		return "skip"
	case pipeline.ModeRerender:
		return "rerender"
	}
	return "select"
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// literalInvocation reports whether inv carries model-written code: a\n// non-empty kwarg that the action's package declares as type literal.\nfunc literalInvocation(set *genpkg.Set, inv recording.Invocation) bool {\n\tif set == nil {\n\t\treturn false\n\t}\n\tfor _, p := range set.Packages {\n\t\ta, ok := p.Actions[inv.Action]\n\t\tif !ok {\n\t\t\tcontinue\n\t\t}\n\t\tfor name, k := range a.Kwargs {\n\t\t\tif v, ok := inv.Kwargs[name]; ok && k.Type == \"literal\" && v != nil && v != \"\" {\n\t\t\t\treturn true\n\t\t\t}\n\t\t}\n\t\treturn false\n\t}\n\treturn false\n}","file":"internal/cli/growlogwire.go","name":"literalInvocation"}}
// literalInvocation reports whether inv carries model-written code: a
// non-empty kwarg that the action's package declares as type literal.
func literalInvocation(set *genpkg.Set, inv recording.Invocation) bool {
	if set == nil {
		return false
	}
	for _, p := range set.Packages {
		a, ok := p.Actions[inv.Action]
		if !ok {
			continue
		}
		for name, k := range a.Kwargs {
			if v, ok := inv.Kwargs[name]; ok && k.Type == "literal" && v != nil && v != "" {
				return true
			}
		}
		return false
	}
	return false
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// growLogEntries builds one entry per record a finished grow handled. A record\n// skipped as current gets an entry only when skips is set. ms is the wall time\n// of the whole run, which every entry of the run shares.\nfunc growLogEntries(res *pipeline.Result, model string, at time.Time, elapsed time.Duration, skips bool) []growlog.Entry {\n\tif res == nil {\n\t\treturn nil\n\t}\n\tselections := map[string]pipeline.Selection{}\n\tfor _, s := range res.Selections {\n\t\tselections[s.RecordID] = s\n\t}\n\ttype tally struct {\n\t\tregions growlog.Regions\n\t\tfiles   map[string]bool\n\t}\n\ttallies := map[string]*tally{}\n\tfor _, in := range res.Injections {\n\t\tt := tallies[in.Record]\n\t\tif t == nil {\n\t\t\tt = &tally{files: map[string]bool{}}\n\t\t\ttallies[in.Record] = t\n\t\t}\n\t\tswitch {\n\t\tcase in.Pruned:\n\t\t\tt.regions.Pruned++\n\t\tcase in.Replaced:\n\t\t\tt.regions.Replaced++\n\t\tcase in.Skipped, in.Unselected, in.Held:\n\t\t\tt.regions.Unchanged++\n\t\tdefault:\n\t\t\tt.regions.Injected++\n\t\t}\n\t\tt.files[in.Path] = true\n\t}\n\tvar out []growlog.Entry\n\tfor _, c := range res.Classifications {\n\t\tif c.Mode == pipeline.ModeNone || (c.Mode == pipeline.ModeSkip && !skips) {\n\t\t\tcontinue\n\t\t}\n\t\te := growlog.Entry{\n\t\t\tTS:      at.UTC().Format(time.RFC3339),\n\t\t\tSedum:   recording.Version,\n\t\t\tRecord:  c.RecordID,\n\t\t\tMode:    growLogMode(c.Mode),\n\t\t\tMS:      elapsed.Milliseconds(),\n\t\t\tOutcome: \"ok\",\n\t\t}\n\t\tif s, ok := selections[c.RecordID]; ok {\n\t\t\te.Calls, e.Rejected, e.Completeness, e.Repairs = s.Calls, s.Rejected, s.Completeness, s.Repairs\n\t\t\te.PromptTokens, e.CompletionTokens = s.PromptTokens, s.CompletionTokens\n\t\t\tfor _, inv := range s.Invocations {\n\t\t\t\te.Invocations++\n\t\t\t\tif literalInvocation(res.Packages, inv) {\n\t\t\t\t\te.Literal++\n\t\t\t\t} else {\n\t\t\t\t\te.Declarative++\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t\tif e.Calls > 0 {\n\t\t\te.Model = model\n\t\t}\n\t\tif t := tallies[c.RecordID]; t != nil {\n\t\t\te.Regions, e.Files = t.regions, len(t.files)\n\t\t}\n\t\tout = append(out, e)\n\t}\n\treturn out\n}","file":"internal/cli/growlogwire.go","name":"growLogEntries"}}
// growLogEntries builds one entry per record a finished grow handled. A record
// skipped as current gets an entry only when skips is set. ms is the wall time
// of the whole run, which every entry of the run shares.
func growLogEntries(res *pipeline.Result, model string, at time.Time, elapsed time.Duration, skips bool) []growlog.Entry {
	if res == nil {
		return nil
	}
	selections := map[string]pipeline.Selection{}
	for _, s := range res.Selections {
		selections[s.RecordID] = s
	}
	type tally struct {
		regions growlog.Regions
		files   map[string]bool
	}
	tallies := map[string]*tally{}
	for _, in := range res.Injections {
		t := tallies[in.Record]
		if t == nil {
			t = &tally{files: map[string]bool{}}
			tallies[in.Record] = t
		}
		switch {
		case in.Pruned:
			t.regions.Pruned++
		case in.Replaced:
			t.regions.Replaced++
		case in.Skipped, in.Unselected, in.Held:
			t.regions.Unchanged++
		default:
			t.regions.Injected++
		}
		t.files[in.Path] = true
	}
	var out []growlog.Entry
	for _, c := range res.Classifications {
		if c.Mode == pipeline.ModeNone || (c.Mode == pipeline.ModeSkip && !skips) {
			continue
		}
		e := growlog.Entry{
			TS:      at.UTC().Format(time.RFC3339),
			Sedum:   recording.Version,
			Record:  c.RecordID,
			Mode:    growLogMode(c.Mode),
			MS:      elapsed.Milliseconds(),
			Outcome: "ok",
		}
		if s, ok := selections[c.RecordID]; ok {
			e.Calls, e.Rejected, e.Completeness, e.Repairs = s.Calls, s.Rejected, s.Completeness, s.Repairs
			e.PromptTokens, e.CompletionTokens = s.PromptTokens, s.CompletionTokens
			for _, inv := range s.Invocations {
				e.Invocations++
				if literalInvocation(res.Packages, inv) {
					e.Literal++
				} else {
					e.Declarative++
				}
			}
		}
		if e.Calls > 0 {
			e.Model = model
		}
		if t := tallies[c.RecordID]; t != nil {
			e.Regions, e.Files = t.regions, len(t.files)
		}
		out = append(out, e)
	}
	return out
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// growLogFailure builds the entry for a grow that returned err. When the error\n// is a selection.Rejection it names the record and carries what the refused\n// calls cost.\nfunc growLogFailure(err error, model string, at time.Time, elapsed time.Duration) growlog.Entry {\n\te := growlog.Entry{\n\t\tTS:      at.UTC().Format(time.RFC3339),\n\t\tSedum:   recording.Version,\n\t\tMode:    \"select\",\n\t\tMS:      elapsed.Milliseconds(),\n\t\tOutcome: \"failed\",\n\t\tError:   err.Error(),\n\t}\n\tvar rej *selection.Rejection\n\tif errors.As(err, &rej) {\n\t\te.Record, e.Calls, e.Rejected, e.Completeness = rej.RecordID, rej.Calls, rej.Rejected, rej.Completeness\n\t\te.PromptTokens, e.CompletionTokens = rej.PromptTokens, rej.CompletionTokens\n\t}\n\tvar ver *selection.VerifyError\n\tif errors.As(err, &ver) {\n\t\te.Record, e.Calls, e.Rejected, e.Completeness, e.Repairs = ver.RecordID, ver.Calls, ver.Rejected, ver.Completeness, ver.Repairs\n\t\te.PromptTokens, e.CompletionTokens = ver.PromptTokens, ver.CompletionTokens\n\t}\n\tif e.Calls > 0 {\n\t\te.Model = model\n\t}\n\treturn e\n}","file":"internal/cli/growlogwire.go","name":"growLogFailure"}}
// growLogFailure builds the entry for a grow that returned err. When the error
// is a selection.Rejection it names the record and carries what the refused
// calls cost.
func growLogFailure(err error, model string, at time.Time, elapsed time.Duration) growlog.Entry {
	e := growlog.Entry{
		TS:      at.UTC().Format(time.RFC3339),
		Sedum:   recording.Version,
		Mode:    "select",
		MS:      elapsed.Milliseconds(),
		Outcome: "failed",
		Error:   err.Error(),
	}
	var rej *selection.Rejection
	if errors.As(err, &rej) {
		e.Record, e.Calls, e.Rejected, e.Completeness = rej.RecordID, rej.Calls, rej.Rejected, rej.Completeness
		e.PromptTokens, e.CompletionTokens = rej.PromptTokens, rej.CompletionTokens
	}
	var ver *selection.VerifyError
	if errors.As(err, &ver) {
		e.Record, e.Calls, e.Rejected, e.Completeness, e.Repairs = ver.RecordID, ver.Calls, ver.Rejected, ver.Completeness, ver.Repairs
		e.PromptTokens, e.CompletionTokens = ver.PromptTokens, ver.CompletionTokens
	}
	if e.Calls > 0 {
		e.Model = model
	}
	return e
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// writeGrowLog appends the entries of a finished grow to the grow log. It\n// writes nothing for a run that changes nothing (--dry-run, --check, --stamp,\n// --stop-after) or when SEDUM_GROW_LOG is off, and a log that cannot be\n// written is a warning on errOut, never a failed grow.\nfunc writeGrowLog(errOut io.Writer, cfg GrowConfig, res *pipeline.Result, runErr error, started time.Time) {\n\tpath, on, skips := growLogSettings(os.Getenv)\n\tif !on || cfg.DryRun || cfg.Check || cfg.Stamp || cfg.StopAfter != \"\" {\n\t\treturn\n\t}\n\telapsed := time.Since(started)\n\tmodel := growLogModel(cfg)\n\tentries := growLogEntries(res, model, started, elapsed, skips)\n\tif runErr != nil {\n\t\tentries = []growlog.Entry{growLogFailure(runErr, model, started, elapsed)}\n\t}\n\tfor _, e := range entries {\n\t\tif err := growlog.Append(path, e); err != nil {\n\t\t\tfmt.Fprintf(errOut, \"warning: %v\\n\", err)\n\t\t\treturn\n\t\t}\n\t}\n}","file":"internal/cli/growlogwire.go","name":"writeGrowLog"}}
// writeGrowLog appends the entries of a finished grow to the grow log. It
// writes nothing for a run that changes nothing (--dry-run, --check, --stamp,
// --stop-after) or when SEDUM_GROW_LOG is off, and a log that cannot be
// written is a warning on errOut, never a failed grow.
func writeGrowLog(errOut io.Writer, cfg GrowConfig, res *pipeline.Result, runErr error, started time.Time) {
	path, on, skips := growLogSettings(os.Getenv)
	if !on || cfg.DryRun || cfg.Check || cfg.Stamp || cfg.StopAfter != "" {
		return
	}
	elapsed := time.Since(started)
	model := growLogModel(cfg)
	entries := growLogEntries(res, model, started, elapsed, skips)
	if runErr != nil {
		entries = []growlog.Entry{growLogFailure(runErr, model, started, elapsed)}
	}
	for _, e := range entries {
		if err := growlog.Append(path, e); err != nil {
			fmt.Fprintf(errOut, "warning: %v\n", err)
			return
		}
	}
}

// /sedum:setFunction
