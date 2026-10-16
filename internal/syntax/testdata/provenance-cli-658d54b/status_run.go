package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/status_run.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/status_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/status_run.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/status_run.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/status_run.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/status_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/status_run.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type statusSpec struct {\n\tPath       string `json:\"path\"`\n\tType       string `json:\"type,omitempty\"`\n\tRunCommand string `json:\"run_command,omitempty\"`\n}","file":"cmd/provenance/status_run.go","name":"statusSpec"}}
type statusSpec struct {
	Path       string `json:"path"`
	Type       string `json:"type,omitempty"`
	RunCommand string `json:"run_command,omitempty"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type statusRecord struct {\n\tID                string       `json:\"id\"`\n\tTitle             string       `json:\"title\"`\n\tType              string       `json:\"type\"`\n\tStatus            string       `json:\"status\"`\n\tVersion           int          `json:\"version\"`\n\tReason            string       `json:\"reason\"`\n\tReasonKind        string       `json:\"reason_kind,omitempty\"`\n\tCreatedAt         string       `json:\"created_at\"`\n\tAuthor            string       `json:\"author\"`\n\tAliases           []string     `json:\"aliases,omitempty\"`\n\tIntent            string       `json:\"intent\"`\n\tConstraints       []string     `json:\"constraints\"`\n\tAffectedScope     []string     `json:\"affected_scope\"`\n\tForbiddenScope    []string     `json:\"forbidden_scope\"`\n\tImplements        string       `json:\"implements,omitempty\"`\n\tImplementsVersion int          `json:\"implements_version,omitempty\"`\n\tRelated           []string     `json:\"related\"`\n\tTags              []string     `json:\"tags\"`\n\tAssociatedSpecs   []statusSpec `json:\"associated_specs\"`\n\tSealedAtSHA       string       `json:\"sealed_at_sha,omitempty\"`\n}","file":"cmd/provenance/status_run.go","name":"statusRecord"}}
type statusRecord struct {
	ID                string       `json:"id"`
	Title             string       `json:"title"`
	Type              string       `json:"type"`
	Status            string       `json:"status"`
	Version           int          `json:"version"`
	Reason            string       `json:"reason"`
	ReasonKind        string       `json:"reason_kind,omitempty"`
	CreatedAt         string       `json:"created_at"`
	Author            string       `json:"author"`
	Aliases           []string     `json:"aliases,omitempty"`
	Intent            string       `json:"intent"`
	Constraints       []string     `json:"constraints"`
	AffectedScope     []string     `json:"affected_scope"`
	ForbiddenScope    []string     `json:"forbidden_scope"`
	Implements        string       `json:"implements,omitempty"`
	ImplementsVersion int          `json:"implements_version,omitempty"`
	Related           []string     `json:"related"`
	Tags              []string     `json:"tags"`
	AssociatedSpecs   []statusSpec `json:"associated_specs"`
	SealedAtSHA       string       `json:"sealed_at_sha,omitempty"`
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func toStatusRecord(r *record.Record) statusRecord {\n\tspecs := []statusSpec{}\n\tfor _, s := range r.AssociatedSpecs {\n\t\tspecs = append(specs, statusSpec{Path: s.Path, Type: s.Type, RunCommand: s.RunCommand})\n\t}\n\treturn statusRecord{\n\t\tID:                r.ID,\n\t\tTitle:             r.Title,\n\t\tType:              string(r.EffectiveType()),\n\t\tStatus:            string(r.Status),\n\t\tVersion:           r.Version,\n\t\tReason:            r.Reason,\n\t\tReasonKind:        r.ReasonKind,\n\t\tCreatedAt:         r.CreatedAt,\n\t\tAuthor:            r.Author,\n\t\tAliases:           append([]string{}, r.Aliases...),\n\t\tIntent:            r.Intent,\n\t\tConstraints:       append([]string{}, r.Constraints...),\n\t\tAffectedScope:     append([]string{}, r.AffectedScope...),\n\t\tForbiddenScope:    append([]string{}, r.ForbiddenScope...),\n\t\tImplements:        r.Implements,\n\t\tImplementsVersion: r.ImplementsVersion,\n\t\tRelated:           append([]string{}, r.Related...),\n\t\tTags:              append([]string{}, r.Tags...),\n\t\tAssociatedSpecs:   specs,\n\t\tSealedAtSHA:       r.SealedAtSHA,\n\t}\n}","file":"cmd/provenance/status_run.go","name":"toStatusRecord"}}
func toStatusRecord(r *record.Record) statusRecord {
	specs := []statusSpec{}
	for _, s := range r.AssociatedSpecs {
		specs = append(specs, statusSpec{Path: s.Path, Type: s.Type, RunCommand: s.RunCommand})
	}
	return statusRecord{
		ID:                r.ID,
		Title:             r.Title,
		Type:              string(r.EffectiveType()),
		Status:            string(r.Status),
		Version:           r.Version,
		Reason:            r.Reason,
		ReasonKind:        r.ReasonKind,
		CreatedAt:         r.CreatedAt,
		Author:            r.Author,
		Aliases:           append([]string{}, r.Aliases...),
		Intent:            r.Intent,
		Constraints:       append([]string{}, r.Constraints...),
		AffectedScope:     append([]string{}, r.AffectedScope...),
		ForbiddenScope:    append([]string{}, r.ForbiddenScope...),
		Implements:        r.Implements,
		ImplementsVersion: r.ImplementsVersion,
		Related:           append([]string{}, r.Related...),
		Tags:              append([]string{}, r.Tags...),
		AssociatedSpecs:   specs,
		SealedAtSHA:       r.SealedAtSHA,
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func printTable(out io.Writer, enforcement string, recs []*record.Record) {\n\tfmt.Fprint(out, \"\\nPROVENANCE RECORDS\\n\\n\")\n\tif enforcement == \"\" {\n\t\tenforcement = \"none\"\n\t}\n\tfmt.Fprintf(out, \"  Enforcement: %s\\n\\n\", enforcement)\n\tfmt.Fprintf(out, \"  %-26s %-14s %-45s %s\\n\", \"ID\", \"STATUS\", \"TITLE\", \"LINESPECS\")\n\tfmt.Fprint(out, \"  \"+strings.Repeat(\"-\", 96)+\"\\n\")\n\topenNoSpecs := 0\n\tfor _, rec := range recs {\n\t\tid := fmt.Sprintf(\"%s@%d\", rec.ID, rec.Version)\n\t\tstatus := string(rec.Status)\n\t\tif rec.Status == record.StatusOpen && len(rec.AssociatedSpecs) == 0 {\n\t\t\tstatus = \"open ⚠\"\n\t\t\topenNoSpecs++\n\t\t}\n\t\ttitle := rec.Title\n\t\tif len(title) > 42 {\n\t\t\ttitle = title[:39] + \"...\"\n\t\t}\n\t\tfmt.Fprintf(out, \"  %-26s %-14s %-45s %d\\n\", id, status, title, len(rec.AssociatedSpecs))\n\t}\n\tif openNoSpecs > 0 {\n\t\tfmt.Fprintf(out, \"\\n⚠ %d open records have no associated LineSpecs\\n\", openNoSpecs)\n\t}\n\tfmt.Fprint(out, \"\\n\")\n}","file":"cmd/provenance/status_run.go","name":"printTable"}}
func printTable(out io.Writer, enforcement string, recs []*record.Record) {
	fmt.Fprint(out, "\nPROVENANCE RECORDS\n\n")
	if enforcement == "" {
		enforcement = "none"
	}
	fmt.Fprintf(out, "  Enforcement: %s\n\n", enforcement)
	fmt.Fprintf(out, "  %-26s %-14s %-45s %s\n", "ID", "STATUS", "TITLE", "LINESPECS")
	fmt.Fprint(out, "  "+strings.Repeat("-", 96)+"\n")
	openNoSpecs := 0
	for _, rec := range recs {
		id := fmt.Sprintf("%s@%d", rec.ID, rec.Version)
		status := string(rec.Status)
		if rec.Status == record.StatusOpen && len(rec.AssociatedSpecs) == 0 {
			status = "open ⚠"
			openNoSpecs++
		}
		title := rec.Title
		if len(title) > 42 {
			title = title[:39] + "..."
		}
		fmt.Fprintf(out, "  %-26s %-14s %-45s %d\n", id, status, title, len(rec.AssociatedSpecs))
	}
	if openNoSpecs > 0 {
		fmt.Fprintf(out, "\n⚠ %d open records have no associated LineSpecs\n", openNoSpecs)
	}
	fmt.Fprint(out, "\n")
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func printDetail(out io.Writer, rec *record.Record, all []*record.Record) {\n\tstatus := string(rec.Status)\n\tif rec.Status == record.StatusOpen && len(rec.AssociatedSpecs) == 0 {\n\t\tstatus = \"open ⚠\"\n\t}\n\tfmt.Fprintf(out, \"\\n%s@%d  ·  %s\\n\", rec.ID, rec.Version, status)\n\tfmt.Fprintf(out, \"%s\\n\\n\", strings.Repeat(\"─\", 60))\n\tfmt.Fprintf(out, \"Title:        %s\\n\", rec.Title)\n\tfmt.Fprintf(out, \"Author:       %s\\n\", rec.Author)\n\tfmt.Fprintf(out, \"Created:      %s\\n\", rec.CreatedAt)\n\tfmt.Fprintf(out, \"Version:      %d\", rec.Version)\n\tswitch {\n\tcase rec.ReasonKind == record.ReasonKindFix:\n\t\tfmt.Fprintf(out, \" (fix: %s)\", rec.Reason)\n\tcase rec.Reason != \"\":\n\t\tfmt.Fprintf(out, \" (%s)\", rec.Reason)\n\t}\n\tfmt.Fprintln(out)\n\tif len(rec.Aliases) > 0 {\n\t\tfmt.Fprintf(out, \"Aliases:      %s\\n\", strings.Join(rec.Aliases, \", \"))\n\t}\n\tif rec.Status == record.StatusImplemented && rec.SealedAtSHA != \"\" {\n\t\tsha := rec.SealedAtSHA\n\t\tif len(sha) > 7 {\n\t\t\tsha = sha[:7]\n\t\t}\n\t\tfmt.Fprintf(out, \"Sealed at:    %s\\n\", sha)\n\t}\n\tif rec.Implements != \"\" {\n\t\tline := rec.Implements\n\t\tif rec.ImplementsVersion != 0 {\n\t\t\tline = fmt.Sprintf(\"%s@%d\", rec.Implements, rec.ImplementsVersion)\n\t\t}\n\t\tfor _, o := range all {\n\t\t\tif o.ID == rec.Implements {\n\t\t\t\tline += \"  \" + o.Title\n\t\t\t\tbreak\n\t\t\t}\n\t\t}\n\t\tfmt.Fprintf(out, \"Implements:   %s\\n\", line)\n\t} else {\n\t\tfmt.Fprintln(out, \"Implements:   —\")\n\t}\n\tvar children []*record.Record\n\tfor _, o := range all {\n\t\tif o.Implements == rec.ID {\n\t\t\tchildren = append(children, o)\n\t\t}\n\t}\n\tif len(children) == 0 {\n\t\tfmt.Fprintln(out, \"Implemented by: —\")\n\t} else {\n\t\tfmt.Fprintln(out, \"Implemented by:\")\n\t\tfor _, o := range children {\n\t\t\tfmt.Fprintf(out, \"  · %s  [%s] %s\\n\", o.ID, o.EffectiveType(), o.Title)\n\t\t}\n\t}\n\tif len(rec.Tags) > 0 {\n\t\tfmt.Fprintf(out, \"Tags:         %s\\n\", strings.Join(rec.Tags, \", \"))\n\t}\n\tfmt.Fprintln(out)\n\tif rec.Intent != \"\" {\n\t\tfmt.Fprintln(out, \"Intent:\")\n\t\tfor _, line := range strings.Split(strings.TrimRight(rec.Intent, \"\\n\"), \"\\n\") {\n\t\t\tfmt.Fprintf(out, \"  %s\\n\", line)\n\t\t}\n\t\tfmt.Fprintln(out)\n\t}\n\tsection := func(title string, items []string) {\n\t\tif len(items) == 0 {\n\t\t\treturn\n\t\t}\n\t\tfmt.Fprintln(out, title)\n\t\tfor _, it := range items {\n\t\t\tfmt.Fprintf(out, \"  · %s\\n\", it)\n\t\t}\n\t\tfmt.Fprintln(out)\n\t}\n\tsection(\"Constraints:\", rec.Constraints)\n\tsection(\"Allowed Scope:\", rec.AffectedScope)\n\tsection(\"Forbidden Scope (explicit):\", rec.ForbiddenScope)\n\tvar paths []string\n\tfor _, s := range rec.AssociatedSpecs {\n\t\tpaths = append(paths, s.Path)\n\t}\n\tsection(\"Associated Specs:\", paths)\n}","file":"cmd/provenance/status_run.go","name":"printDetail"}}
func printDetail(out io.Writer, rec *record.Record, all []*record.Record) {
	status := string(rec.Status)
	if rec.Status == record.StatusOpen && len(rec.AssociatedSpecs) == 0 {
		status = "open ⚠"
	}
	fmt.Fprintf(out, "\n%s@%d  ·  %s\n", rec.ID, rec.Version, status)
	fmt.Fprintf(out, "%s\n\n", strings.Repeat("─", 60))
	fmt.Fprintf(out, "Title:        %s\n", rec.Title)
	fmt.Fprintf(out, "Author:       %s\n", rec.Author)
	fmt.Fprintf(out, "Created:      %s\n", rec.CreatedAt)
	fmt.Fprintf(out, "Version:      %d", rec.Version)
	switch {
	case rec.ReasonKind == record.ReasonKindFix:
		fmt.Fprintf(out, " (fix: %s)", rec.Reason)
	case rec.Reason != "":
		fmt.Fprintf(out, " (%s)", rec.Reason)
	}
	fmt.Fprintln(out)
	if len(rec.Aliases) > 0 {
		fmt.Fprintf(out, "Aliases:      %s\n", strings.Join(rec.Aliases, ", "))
	}
	if rec.Status == record.StatusImplemented && rec.SealedAtSHA != "" {
		sha := rec.SealedAtSHA
		if len(sha) > 7 {
			sha = sha[:7]
		}
		fmt.Fprintf(out, "Sealed at:    %s\n", sha)
	}
	if rec.Implements != "" {
		line := rec.Implements
		if rec.ImplementsVersion != 0 {
			line = fmt.Sprintf("%s@%d", rec.Implements, rec.ImplementsVersion)
		}
		for _, o := range all {
			if o.ID == rec.Implements {
				line += "  " + o.Title
				break
			}
		}
		fmt.Fprintf(out, "Implements:   %s\n", line)
	} else {
		fmt.Fprintln(out, "Implements:   —")
	}
	var children []*record.Record
	for _, o := range all {
		if o.Implements == rec.ID {
			children = append(children, o)
		}
	}
	if len(children) == 0 {
		fmt.Fprintln(out, "Implemented by: —")
	} else {
		fmt.Fprintln(out, "Implemented by:")
		for _, o := range children {
			fmt.Fprintf(out, "  · %s  [%s] %s\n", o.ID, o.EffectiveType(), o.Title)
		}
	}
	if len(rec.Tags) > 0 {
		fmt.Fprintf(out, "Tags:         %s\n", strings.Join(rec.Tags, ", "))
	}
	fmt.Fprintln(out)
	if rec.Intent != "" {
		fmt.Fprintln(out, "Intent:")
		for _, line := range strings.Split(strings.TrimRight(rec.Intent, "\n"), "\n") {
			fmt.Fprintf(out, "  %s\n", line)
		}
		fmt.Fprintln(out)
	}
	section := func(title string, items []string) {
		if len(items) == 0 {
			return
		}
		fmt.Fprintln(out, title)
		for _, it := range items {
			fmt.Fprintf(out, "  · %s\n", it)
		}
		fmt.Fprintln(out)
	}
	section("Constraints:", rec.Constraints)
	section("Allowed Scope:", rec.AffectedScope)
	section("Forbidden Scope (explicit):", rec.ForbiddenScope)
	var paths []string
	for _, s := range rec.AssociatedSpecs {
		paths = append(paths, s.Path)
	}
	section("Associated Specs:", paths)
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-12b1ca97","gen":"sha256:3d5eb5aff0b580425966bd18cb1e038d9fef9018ce5f4611f46a803ad6ee5bc9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runStatus(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\trecordID, err := cmd.Flags().GetString(\"record\")\n\tif err != nil {\n\t\treturn err\n\t}\n\tfilter, err := cmd.Flags().GetString(\"filter\")\n\tif err != nil {\n\t\treturn err\n\t}\n\tformat, err := cmd.Flags().GetString(\"format\")\n\tif err != nil {\n\t\treturn err\n\t}\n\tif format != \"\" && format != \"human\" && format != \"json\" {\n\t\treturn fmt.Errorf(\"unknown format %q (human|json)\", format)\n\t}\n\tids, err := r.Store.List()\n\tif err != nil {\n\t\treturn err\n\t}\n\tsort.Strings(ids)\n\tvar all []*record.Record\n\tfor _, id := range ids {\n\t\trec, _, err := r.Store.Load(id)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tall = append(all, rec)\n\t}\n\tif recordID != \"\" {\n\t\trec, _, err := r.resolve(recordID)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tif format == \"json\" {\n\t\t\tenc := json.NewEncoder(out)\n\t\t\tenc.SetIndent(\"\", \"  \")\n\t\t\treturn enc.Encode(map[string]statusRecord{\"record\": toStatusRecord(rec)})\n\t\t}\n\t\tprintDetail(out, rec, all)\n\t\treturn nil\n\t}\n\tvar filtered []*record.Record\n\tswitch {\n\tcase filter == \"\":\n\t\tfiltered = all\n\tcase strings.HasPrefix(filter, \"tag:\"):\n\t\tname := strings.TrimPrefix(filter, \"tag:\")\n\t\tfor _, rec := range all {\n\t\t\tfor _, t := range rec.Tags {\n\t\t\t\tif t == name {\n\t\t\t\t\tfiltered = append(filtered, rec)\n\t\t\t\t\tbreak\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\tdefault:\n\t\tswitch filter {\n\t\tcase \"draft\", \"open\", \"implemented\", \"deprecated\":\n\t\t\tfor _, rec := range all {\n\t\t\t\tif string(rec.Status) == filter {\n\t\t\t\t\tfiltered = append(filtered, rec)\n\t\t\t\t}\n\t\t\t}\n\t\tdefault:\n\t\t\treturn fmt.Errorf(\"unknown filter %q (draft|open|implemented|deprecated or tag:name)\", filter)\n\t\t}\n\t}\n\tif format == \"json\" {\n\t\titems := make([]statusRecord, 0, len(filtered))\n\t\tfor _, rec := range filtered {\n\t\t\titems = append(items, toStatusRecord(rec))\n\t\t}\n\t\tenc := json.NewEncoder(out)\n\t\tenc.SetIndent(\"\", \"  \")\n\t\treturn enc.Encode(map[string][]statusRecord{\"records\": items})\n\t}\n\tprintTable(out, r.Config.Enforcement, filtered)\n\treturn nil\n}","file":"cmd/provenance/status_run.go","name":"runStatus"}}
func runStatus(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}
	recordID, err := cmd.Flags().GetString("record")
	if err != nil {
		return err
	}
	filter, err := cmd.Flags().GetString("filter")
	if err != nil {
		return err
	}
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return err
	}
	if format != "" && format != "human" && format != "json" {
		return fmt.Errorf("unknown format %q (human|json)", format)
	}
	ids, err := r.Store.List()
	if err != nil {
		return err
	}
	sort.Strings(ids)
	var all []*record.Record
	for _, id := range ids {
		rec, _, err := r.Store.Load(id)
		if err != nil {
			return err
		}
		all = append(all, rec)
	}
	if recordID != "" {
		rec, _, err := r.resolve(recordID)
		if err != nil {
			return err
		}
		if format == "json" {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(map[string]statusRecord{"record": toStatusRecord(rec)})
		}
		printDetail(out, rec, all)
		return nil
	}
	var filtered []*record.Record
	switch {
	case filter == "":
		filtered = all
	case strings.HasPrefix(filter, "tag:"):
		name := strings.TrimPrefix(filter, "tag:")
		for _, rec := range all {
			for _, t := range rec.Tags {
				if t == name {
					filtered = append(filtered, rec)
					break
				}
			}
		}
	default:
		switch filter {
		case "draft", "open", "implemented", "deprecated":
			for _, rec := range all {
				if string(rec.Status) == filter {
					filtered = append(filtered, rec)
				}
			}
		default:
			return fmt.Errorf("unknown filter %q (draft|open|implemented|deprecated or tag:name)", filter)
		}
	}
	if format == "json" {
		items := make([]statusRecord, 0, len(filtered))
		for _, rec := range filtered {
			items = append(items, toStatusRecord(rec))
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(map[string][]statusRecord{"records": items})
	}
	printTable(out, r.Config.Enforcement, filtered)
	return nil
}

// /sedum:setFunction
