package gitcheck

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	librecord "github.com/livecodelife/provenance/record"
)

// Options are the settings of one check. Output is the base directory, the
// directory that holds the config; Generators and Records are the directories
// class computation reads, relative to it unless absolute.
type Options struct {
	Range       string
	MessageFile string
	Config      string
	Output      string
	Generators  string
	Records     string
	Format      string
	Provenance  string
	TopLevel    string
	Staged      bool
	Paths       PathConfig
}

// Outcome counts what a check found: the findings left, the findings waived and
// the scope paths class computation skipped.
type Outcome struct {
	Unwaived int
	Waived   int
	Skipped  int
}

// Check runs provenance for each commit of the range, or for the staged change,
// waives the findings the regeneration and marker migration proofs cover, and
// writes the merged output to stdout. On an error nothing is written.
func Check(opts Options, stdout io.Writer) (Outcome, error) {
	binary, err := ResolveBinary(opts.Provenance, os.Getenv("PROVENANCE"), opts.TopLevel)
	if err != nil {
		return Outcome{}, err
	}
	env := Environment(opts.Staged, os.Environ())
	revs := []string{""}
	if !opts.Staged {
		revs, err = Commits(opts.Output, opts.Range, env)
		if err != nil {
			return Outcome{}, err
		}
	}

	var logs []Log
	var justifications [][]string
	var unwaived []Finding
	waived := 0
	skipped := map[string]bool{}
	for _, rev := range revs {
		out, err := RunProvenance(binary, CheckArgs(rev, opts.Staged, opts.MessageFile, opts.Config), opts.Output, env)
		if err != nil {
			return Outcome{}, err
		}
		log, err := ParseLog(out)
		if err != nil {
			return Outcome{}, err
		}
		results := log.Results()
		findings := make([]Finding, len(results))
		for i, r := range results {
			findings[i] = FindingOf(r)
		}
		just, err := judge(opts, env, rev, findings, skipped)
		if err != nil {
			return Outcome{}, err
		}
		for i, f := range findings {
			if just[i] != "" {
				waived++
			} else {
				unwaived = append(unwaived, f)
			}
		}
		logs = append(logs, log)
		justifications = append(justifications, just)
	}

	outcome := Outcome{Unwaived: len(unwaived), Waived: waived, Skipped: len(skipped)}
	if opts.Format == "sarif" {
		err = WriteSARIF(stdout, Merge(logs, justifications), outcome.Skipped)
	} else {
		err = WriteHuman(stdout, unwaived, outcome.Waived, outcome.Skipped)
	}
	if err != nil {
		return Outcome{}, err
	}
	return outcome, nil
}

// judge returns one justification per finding of a commit: empty when the
// finding stays, otherwise the proof and the commit that waive it. skipped
// collects the scope paths class computation skipped.
func judge(opts Options, env []string, rev string, findings []Finding, skipped map[string]bool) ([]string, error) {
	just := make([]string, len(findings))
	if len(findings) == 0 {
		return just, nil
	}
	message, err := CommitMessage(opts.Output, rev, opts.Staged, opts.MessageFile, env)
	if err != nil {
		return nil, err
	}
	tags := ParseTags(message)
	if len(tags) == 0 {
		return just, nil
	}
	base, hasBase, err := Base(opts.Output, rev, opts.Staged, env)
	if err != nil {
		return nil, err
	}
	if !hasBase {
		return just, nil
	}
	changed, err := ChangedPaths(opts.Output, rev, opts.Staged, env)
	if err != nil {
		return nil, err
	}
	classified := Classify(changed, opts.Paths)

	commitTree, cleanCommit, err := ExportTree(opts.Output, rev, opts.Staged, env)
	if err != nil {
		return nil, err
	}
	defer cleanCommit()
	commitDir, err := counterpart(opts, commitTree)
	if err != nil {
		return nil, err
	}
	commitClasses := classesOf(opts, commitDir, skipped)

	anyVersion := false
	for _, t := range tags {
		if t.HasVersion {
			anyVersion = true
		}
	}
	var baseClasses map[string]string
	if anyVersion {
		baseTree, cleanBase, err := ExportTree(opts.Output, base, false, env)
		if err != nil {
			return nil, err
		}
		defer cleanBase()
		baseDir, err := counterpart(opts, baseTree)
		if err != nil {
			return nil, err
		}
		baseClasses = classesOf(opts, baseDir, skipped)
	}

	baseReader := BaseReader(opts.Output, base, env)
	commitReader := DirReader(commitDir)
	versions := recordVersions(opts, commitReader, tags)
	regeneration := Regeneration(classified, baseReader, commitReader, commitClasses)

	word := "regeneration: "
	migrationWord := "marker migration: "
	suffix := "staged"
	if !opts.Staged {
		suffix = rev
		if len(suffix) > 7 {
			suffix = suffix[:7]
		}
	}
	for _, tag := range tags {
		migrated, _ := Migration(classified, baseReader, commitReader, baseClasses, commitClasses, tag, versions)
		for i, f := range findings {
			if just[i] != "" {
				continue
			}
			// The migration proof is the stricter one, so it names the waiver
			// when both hold.
			switch {
			case migrated && Allowed("migration", f, tag.ID):
				just[i] = migrationWord + suffix
			case regeneration && Allowed("regeneration", f, tag.ID):
				just[i] = word + suffix
			}
		}
	}
	return just, nil
}

// counterpart is the directory in an exported tree that corresponds to the base
// directory: the exported root joined with the path of Output below the git top
// level.
func counterpart(opts Options, root string) (string, error) {
	top, err := filepath.EvalSymlinks(opts.TopLevel)
	if err != nil {
		top = opts.TopLevel
	}
	out, err := filepath.EvalSymlinks(opts.Output)
	if err != nil {
		out = opts.Output
	}
	rel, err := filepath.Rel(top, out)
	if err != nil {
		return "", GitError("the output directory is not inside the git top level")
	}
	return filepath.Join(root, rel), nil
}

// classesOf computes the class map of an exported tree, rebasing the generators
// and records directories into it, and adds the skipped scope paths to skipped.
// A failure gives a nil map.
func classesOf(opts Options, dir string, skipped map[string]bool) map[string]string {
	classes, paths, err := Classes(dir, inTree(opts, dir, opts.Generators), inTree(opts, dir, opts.Records))
	if err != nil {
		return nil
	}
	for _, p := range paths {
		skipped[p] = true
	}
	return classes
}

// inTree places a directory setting in an exported tree: a relative path is
// joined to the tree's counterpart of the base directory, an absolute one that
// lies under the base directory is re-based into it, and any other is used as
// given.
func inTree(opts Options, dir, setting string) string {
	if !filepath.IsAbs(setting) {
		return filepath.Join(dir, setting)
	}
	if rel, err := filepath.Rel(opts.Output, setting); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return filepath.Join(dir, rel)
	}
	return setting
}

// recordVersions reads the current version of each tagged record from its
// record file at the commit. A record that cannot be read or decoded is absent.
func recordVersions(opts Options, commit TreeReader, tags []Tag) map[string]int {
	versions := map[string]int{}
	for _, tag := range tags {
		data, found, err := commit.Read(strings.Trim(opts.Paths.RecordsDir, "/") + "/" + tag.ID + ".yml")
		if err != nil || !found {
			continue
		}
		rec, err := librecord.Decode(data)
		if err != nil {
			continue
		}
		versions[tag.ID] = rec.Version
	}
	return versions
}
