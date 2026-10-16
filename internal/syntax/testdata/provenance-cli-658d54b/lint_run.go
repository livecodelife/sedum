package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/lint_run.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/lint_run.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/lint_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/lint_run.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/lint_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/lint_run.go","path":"github.com/livecodelife/provenance/lint"}}
import "github.com/livecodelife/provenance/lint"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type lintIssueJSON struct {\n\tRecordID string `json:\"record_id\"`\n\tField    string `json:\"field\"`\n\tMessage  string `json:\"message\"`\n\tSeverity string `json:\"severity\"`\n\tRule     string `json:\"rule\"`\n}","file":"cmd/provenance/lint_run.go","name":"lintIssueJSON"}}
type lintIssueJSON struct {
	RecordID string `json:"record_id"`
	Field    string `json:"field"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Rule     string `json:"rule"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type lintJSON struct {\n\tEnforcement string          `json:\"enforcement\"`\n\tTotal       int             `json:\"total\"`\n\tPassed      int             `json:\"passed\"`\n\tWarnings    int             `json:\"warnings\"`\n\tErrors      int             `json:\"errors\"`\n\tIssues      []lintIssueJSON `json:\"issues\"`\n}","file":"cmd/provenance/lint_run.go","name":"lintJSON"}}
type lintJSON struct {
	Enforcement string          `json:"enforcement"`
	Total       int             `json:"total"`
	Passed      int             `json:"passed"`
	Warnings    int             `json:"warnings"`
	Errors      int             `json:"errors"`
	Issues      []lintIssueJSON `json:"issues"`
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func lintRepo(r *repo, enforcement string) (*lint.Result, error) {\n\tif enforcement == \"\" {\n\t\tenforcement = r.Config.Enforcement\n\t}\n\tif enforcement == \"\" {\n\t\tenforcement = \"none\"\n\t}\n\tvar names []string\n\tfor _, s := range r.Config.SharedRepos {\n\t\tnames = append(names, s.Name)\n\t}\n\treturn lint.Lint(r.Store, lint.Options{\n\t\tRoot:         r.Root,\n\t\tEnforcement:  enforcement,\n\t\tExcludePaths: r.Config.ExcludePaths,\n\t\tSharedRepos:  names,\n\t\tManifestPath: manifestPath,\n\t})\n}","file":"cmd/provenance/lint_run.go","name":"lintRepo"}}
func lintRepo(r *repo, enforcement string) (*lint.Result, error) {
	if enforcement == "" {
		enforcement = r.Config.Enforcement
	}
	if enforcement == "" {
		enforcement = "none"
	}
	var names []string
	for _, s := range r.Config.SharedRepos {
		names = append(names, s.Name)
	}
	return lint.Lint(r.Store, lint.Options{
		Root:         r.Root,
		Enforcement:  enforcement,
		ExcludePaths: r.Config.ExcludePaths,
		SharedRepos:  names,
		ManifestPath: manifestPath,
	})
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func recordLint(res *lint.Result, id string) *lint.Result {\n\tout := &lint.Result{Records: 1, Enforcement: res.Enforcement}\n\tfor _, issue := range res.Issues {\n\t\tif issue.RecordID != id {\n\t\t\tcontinue\n\t\t}\n\t\tout.Issues = append(out.Issues, issue)\n\t\tswitch issue.Severity {\n\t\tcase lint.SeverityError:\n\t\t\tout.ErrorCount++\n\t\tcase lint.SeverityWarning:\n\t\t\tout.WarningCount++\n\t\tcase lint.SeverityHint:\n\t\t\tout.HintCount++\n\t\t}\n\t}\n\tif out.ErrorCount == 0 && out.WarningCount == 0 {\n\t\tout.PassedCount = 1\n\t}\n\treturn out\n}","file":"cmd/provenance/lint_run.go","name":"recordLint"}}
func recordLint(res *lint.Result, id string) *lint.Result {
	out := &lint.Result{Records: 1, Enforcement: res.Enforcement}
	for _, issue := range res.Issues {
		if issue.RecordID != id {
			continue
		}
		out.Issues = append(out.Issues, issue)
		switch issue.Severity {
		case lint.SeverityError:
			out.ErrorCount++
		case lint.SeverityWarning:
			out.WarningCount++
		case lint.SeverityHint:
			out.HintCount++
		}
	}
	if out.ErrorCount == 0 && out.WarningCount == 0 {
		out.PassedCount = 1
	}
	return out
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runLint(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\trecord, _ := cmd.Flags().GetString(\"record\")\n\tenforcement, _ := cmd.Flags().GetString(\"enforcement\")\n\tformat, _ := cmd.Flags().GetString(\"format\")\n\twarnOnly, _ := cmd.Flags().GetBool(\"warn\")\n\tinfoOnly, _ := cmd.Flags().GetBool(\"info\")\n\tshowAll, _ := cmd.Flags().GetBool(\"all\")\n\tif format != \"\" && format != \"human\" && format != \"json\" {\n\t\treturn fmt.Errorf(\"unknown format %q (human|json)\", format)\n\t}\n\tres, err := lintRepo(r, enforcement)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif record != \"\" {\n\t\trec, _, err := r.resolve(record)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tres = recordLint(res, rec.ID)\n\t}\n\ttotal := res.PassedCount + res.WarningCount + res.ErrorCount\n\tif format == \"json\" {\n\t\tissues := []lintIssueJSON{}\n\t\tfor _, issue := range res.Issues {\n\t\t\tissues = append(issues, lintIssueJSON{\n\t\t\t\tRecordID: issue.RecordID,\n\t\t\t\tField:    issue.Field,\n\t\t\t\tMessage:  issue.Message,\n\t\t\t\tSeverity: string(issue.Severity),\n\t\t\t\tRule:     issue.Rule,\n\t\t\t})\n\t\t}\n\t\tenc := json.NewEncoder(out)\n\t\tenc.SetIndent(\"\", \"  \")\n\t\tif err := enc.Encode(lintJSON{\n\t\t\tEnforcement: res.Enforcement,\n\t\t\tTotal:       total,\n\t\t\tPassed:      res.PassedCount,\n\t\t\tWarnings:    res.WarningCount,\n\t\t\tErrors:      res.ErrorCount,\n\t\t\tIssues:      issues,\n\t\t}); err != nil {\n\t\t\treturn err\n\t\t}\n\t} else {\n\t\tshow := func(sev lint.Severity) bool {\n\t\t\tswitch {\n\t\t\tcase showAll:\n\t\t\t\treturn true\n\t\t\tcase warnOnly:\n\t\t\t\treturn sev == lint.SeverityWarning\n\t\t\tcase infoOnly:\n\t\t\t\treturn sev == lint.SeverityHint\n\t\t\tdefault:\n\t\t\t\treturn sev == lint.SeverityError\n\t\t\t}\n\t\t}\n\t\tfmt.Fprintf(out, \"\\n✓ Linting %d provenance records (enforcement: %s)\\n\\n\", total, res.Enforcement)\n\t\tfor _, issue := range res.Issues {\n\t\t\tif !show(issue.Severity) {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tsymbol := \"→\"\n\t\t\tswitch issue.Severity {\n\t\t\tcase lint.SeverityError:\n\t\t\t\tsymbol = \"✗\"\n\t\t\tcase lint.SeverityWarning:\n\t\t\t\tsymbol = \"⚠\"\n\t\t\t}\n\t\t\tfield := \"\"\n\t\t\tif issue.Field != \"\" {\n\t\t\t\tfield = \"(\" + issue.Field + \")\"\n\t\t\t}\n\t\t\tfmt.Fprintf(out, \"  %s %-15s %s %s\\n\", symbol, issue.RecordID, issue.Message, field)\n\t\t}\n\t\tfmt.Fprintf(out, \"\\n✓ %d passed  ⚠ %d warnings  ✗ %d errors\\n\", res.PassedCount, res.WarningCount, res.ErrorCount)\n\t\tif res.ErrorCount > 0 && res.Enforcement == \"strict\" {\n\t\t\tfmt.Fprint(out, \"\\n  Hint: Associate LineSpec files with open records or set status to\\n       'implemented' if proof already exists outside of LineSpec.\\n\")\n\t\t}\n\t\tfmt.Fprint(out, \"\\n\")\n\t}\n\tif res.ErrorCount > 0 {\n\t\treturn errors.New(\"lint failed\")\n\t}\n\treturn nil\n}","file":"cmd/provenance/lint_run.go","name":"runLint"}}
func runLint(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}
	record, _ := cmd.Flags().GetString("record")
	enforcement, _ := cmd.Flags().GetString("enforcement")
	format, _ := cmd.Flags().GetString("format")
	warnOnly, _ := cmd.Flags().GetBool("warn")
	infoOnly, _ := cmd.Flags().GetBool("info")
	showAll, _ := cmd.Flags().GetBool("all")
	if format != "" && format != "human" && format != "json" {
		return fmt.Errorf("unknown format %q (human|json)", format)
	}
	res, err := lintRepo(r, enforcement)
	if err != nil {
		return err
	}
	if record != "" {
		rec, _, err := r.resolve(record)
		if err != nil {
			return err
		}
		res = recordLint(res, rec.ID)
	}
	total := res.PassedCount + res.WarningCount + res.ErrorCount
	if format == "json" {
		issues := []lintIssueJSON{}
		for _, issue := range res.Issues {
			issues = append(issues, lintIssueJSON{
				RecordID: issue.RecordID,
				Field:    issue.Field,
				Message:  issue.Message,
				Severity: string(issue.Severity),
				Rule:     issue.Rule,
			})
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(lintJSON{
			Enforcement: res.Enforcement,
			Total:       total,
			Passed:      res.PassedCount,
			Warnings:    res.WarningCount,
			Errors:      res.ErrorCount,
			Issues:      issues,
		}); err != nil {
			return err
		}
	} else {
		show := func(sev lint.Severity) bool {
			switch {
			case showAll:
				return true
			case warnOnly:
				return sev == lint.SeverityWarning
			case infoOnly:
				return sev == lint.SeverityHint
			default:
				return sev == lint.SeverityError
			}
		}
		fmt.Fprintf(out, "\n✓ Linting %d provenance records (enforcement: %s)\n\n", total, res.Enforcement)
		for _, issue := range res.Issues {
			if !show(issue.Severity) {
				continue
			}
			symbol := "→"
			switch issue.Severity {
			case lint.SeverityError:
				symbol = "✗"
			case lint.SeverityWarning:
				symbol = "⚠"
			}
			field := ""
			if issue.Field != "" {
				field = "(" + issue.Field + ")"
			}
			fmt.Fprintf(out, "  %s %-15s %s %s\n", symbol, issue.RecordID, issue.Message, field)
		}
		fmt.Fprintf(out, "\n✓ %d passed  ⚠ %d warnings  ✗ %d errors\n", res.PassedCount, res.WarningCount, res.ErrorCount)
		if res.ErrorCount > 0 && res.Enforcement == "strict" {
			fmt.Fprint(out, "\n  Hint: Associate LineSpec files with open records or set status to\n       'implemented' if proof already exists outside of LineSpec.\n")
		}
		fmt.Fprint(out, "\n")
	}
	if res.ErrorCount > 0 {
		return errors.New("lint failed")
	}
	return nil
}

// /sedum:setFunction
