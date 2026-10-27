package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	librecord "github.com/livecodelife/provenance/record"
	"github.com/livecodelife/sedum/internal/evidence"
	"github.com/livecodelife/sedum/internal/pipeline"
	"github.com/livecodelife/sedum/internal/record"
)

// EvidenceConfig is `sedum evidence`'s flags.
type EvidenceConfig struct {
	Generators string
	Output     string
	Records    string
	Record     string
	Refresh    bool
	JSON       bool
}

func newEvidenceCommand() *cobra.Command {
	var cfg EvidenceConfig

	cmd := &cobra.Command{
		Use:   "evidence",
		Short: "Report whether each record's code is the code its latest evidence entry saw",
		Long: `Judges each record's code freshness from the regions under --output and the
record's evidence file, .sedum/evidence/<record>.jsonl, which is meant to be
committed. The code hash is the one docs/HASHING.md section 7 of the provenance
library defines, so earlier values still match.

A record's verdict is fresh (the latest entry for its current version has the
current code hash), stale (it has another), none (no entry covers the version),
unknown (an owned region has no gen or pkg: run sedum grow --stamp) or no
regions. With --refresh a new entry is appended for each stale or none record
whose hash is known; evidence files are never rewritten.

Exit 0 means no record is stale or unknown, exit 1 that one is, and exit 2 a
usage or read error. Without --refresh nothing is written; no model, no network.`,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 {
				return ExitError{Code: 2, Message: "sedum evidence takes no arguments"}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runEvidence(cmd.OutOrStdout(), cfg)
		},
	}
	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return ExitError{Code: 2, Message: err.Error()}
	})
	cmd.Flags().StringVar(&cfg.Record, "record", "", "Report only this record, by id or alias")
	cmd.Flags().StringVar(&cfg.Output, "output", defaultOutput, "The directory to read the markers from and keep .sedum/evidence in")
	cmd.Flags().StringVar(&cfg.Generators, "generators", "", "The generator packages (required)")
	cmd.Flags().StringVar(&cfg.Records, "records", "", "The provenance records directory (required)")
	cmd.Flags().BoolVar(&cfg.Refresh, "refresh", false, "Append an evidence entry for each stale or none record whose hash is known")
	cmd.Flags().BoolVar(&cfg.JSON, "json", false, "Print the sedum.evidence/1 document")
	return cmd
}

func runEvidence(out io.Writer, cfg EvidenceConfig) error {
	if cfg.Generators == "" {
		return ExitError{Code: 2, Message: "sedum evidence needs --generators"}
	}
	if cfg.Records == "" {
		return ExitError{Code: 2, Message: "sedum evidence needs --records"}
	}
	regions, err := pipeline.Regions(pipeline.RegionsConfig{Generators: cfg.Generators, Output: cfg.Output, Records: cfg.Records})
	if err != nil {
		return ExitError{Code: 2, Message: err.Error()}
	}
	set, _, err := record.Load(cfg.Records, record.Options{})
	if err != nil {
		return ExitError{Code: 2, Message: err.Error()}
	}

	reported := set.Records
	if cfg.Record != "" {
		reported = nil
		for _, rec := range set.Records {
			if rec.ID == cfg.Record {
				reported = append(reported, rec)
				continue
			}
			for _, alias := range rec.Aliases {
				if alias == cfg.Record {
					reported = append(reported, rec)
					break
				}
			}
		}
		if len(reported) == 0 {
			return ExitError{Code: 2, Message: "no record named " + cfg.Record}
		}
	}

	known := map[string]bool{}
	for _, rec := range set.Known {
		known[rec.ID] = true
	}
	selected := map[string]bool{}
	for _, rec := range reported {
		selected[rec.ID] = true
	}
	byRecord := map[string][]evidence.Region{}
	for _, r := range regions.Regions {
		if cfg.Record != "" && !selected[r.Record] {
			continue
		}
		if !known[r.Record] {
			return ExitError{Code: 2, Message: fmt.Sprintf("%s has regions under %s and no record file in %s", r.Record, cfg.Output, cfg.Records)}
		}
		byRecord[r.Record] = append(byRecord[r.Record], evidence.Region{Gen: r.Gen, Pkg: r.Pkg, Tier: r.Tier})
	}

	versions := map[string]int{}
	for _, rec := range reported {
		data, err := os.ReadFile(filepath.Join(cfg.Records, rec.Source))
		if err != nil {
			return ExitError{Code: 2, Message: rec.ID + ": cannot read its record file: " + err.Error()}
		}
		decoded, err := librecord.Decode(data)
		if err != nil {
			return ExitError{Code: 2, Message: rec.ID + ": cannot decode its record file: " + err.Error()}
		}
		versions[rec.ID] = decoded.Version
	}

	judge := func(rec *record.Record) (evidence.Result, error) {
		entries, err := evidence.Read(cfg.Output, rec.ID)
		if err != nil {
			return evidence.Result{}, ExitError{Code: 2, Message: err.Error()}
		}
		return evidence.Judge(rec.ID, versions[rec.ID], byRecord[rec.ID], entries), nil
	}
	results := make([]evidence.Result, len(reported))
	for i, rec := range reported {
		if results[i], err = judge(rec); err != nil {
			return err
		}
	}

	var appended []string
	if cfg.Refresh {
		for i, res := range results {
			if (res.Verdict != evidence.Stale && res.Verdict != evidence.None) || res.CodeHash == "" {
				continue
			}
			entry := evidence.Entry{Record: res.Record, Version: res.Version, CodeHash: res.CodeHash, At: evidence.Timestamp(time.Now()), Runner: evidence.Runner}
			if err := evidence.Append(cfg.Output, entry); err != nil {
				return ExitError{Code: 2, Message: err.Error()}
			}
			appended = append(appended, fmt.Sprintf("appended %s version %d %s", res.Record, res.Version, res.CodeHash))
			if results[i], err = judge(reported[i]); err != nil {
				return err
			}
		}
	}

	if !cfg.JSON {
		for _, line := range appended {
			fmt.Fprintln(out, line)
		}
	}
	if cfg.JSON {
		err = evidence.WriteJSON(out, results)
	} else {
		err = evidence.WriteText(out, results)
	}
	if err != nil {
		return err
	}
	if evidence.Failing(results) {
		return ExitError{Code: 1}
	}
	return nil
}
