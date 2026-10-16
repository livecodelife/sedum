package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// StatusAdopted marks a file that was brought under Sedum.\nconst StatusAdopted FileStatus = \"adopted\"","file":"internal/adopt/report.go","name":"StatusAdopted"}}
// StatusAdopted marks a file that was brought under Sedum.
const StatusAdopted FileStatus = "adopted"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// StatusSkipped marks a file that was left alone and not adopted.\nconst StatusSkipped FileStatus = \"skipped\"","file":"internal/adopt/report.go","name":"StatusSkipped"}}
// StatusSkipped marks a file that was left alone and not adopted.
const StatusSkipped FileStatus = "skipped"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// StatusRejected marks a file whose adoption was refused.\nconst StatusRejected FileStatus = \"rejected\"","file":"internal/adopt/report.go","name":"StatusRejected"}}
// StatusRejected marks a file whose adoption was refused.
const StatusRejected FileStatus = "rejected"

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// FileStatus is the outcome recorded for one file in the adoption report.\ntype FileStatus string","file":"internal/adopt/report.go","name":"FileStatus"}}
// FileStatus is the outcome recorded for one file in the adoption report.
type FileStatus string

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// ResidualInfo describes one non-blank residual span of a file.\ntype ResidualInfo struct {\n\tStartLine int    `json:\"start_line\"`\n\tEndLine   int    `json:\"end_line\"`\n\tFirstLine string `json:\"first_line\"`\n\tReason    Reason `json:\"reason\"`\n}","file":"internal/adopt/report.go","name":"ResidualInfo"}}
// ResidualInfo describes one non-blank residual span of a file.
type ResidualInfo struct {
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	FirstLine string `json:"first_line"`
	Reason    Reason `json:"reason"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// FileReport is the adoption outcome of one file.\ntype FileReport struct {\n\tPath      string         `json:\"path\"`\n\tPackage   string         `json:\"package,omitempty\"`\n\tRecord    string         `json:\"record,omitempty\"`\n\tStatus    FileStatus     `json:\"status\"`\n\tReason    Reason         `json:\"reason,omitempty\"`\n\tDetail    string         `json:\"detail,omitempty\"`\n\tUnits     int            `json:\"units\"`\n\tResiduals []ResidualInfo `json:\"residuals,omitempty\"`\n}","file":"internal/adopt/report.go","name":"FileReport"}}
// FileReport is the adoption outcome of one file.
type FileReport struct {
	Path      string         `json:"path"`
	Package   string         `json:"package,omitempty"`
	Record    string         `json:"record,omitempty"`
	Status    FileStatus     `json:"status"`
	Reason    Reason         `json:"reason,omitempty"`
	Detail    string         `json:"detail,omitempty"`
	Units     int            `json:"units"`
	Residuals []ResidualInfo `json:"residuals,omitempty"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Totals sums the file outcomes of a report.\ntype Totals struct {\n\tFiles           int            `json:\"files\"`\n\tAdopted         int            `json:\"adopted\"`\n\tSkipped         int            `json:\"skipped\"`\n\tRejected        int            `json:\"rejected\"`\n\tUnits           int            `json:\"units\"`\n\tSkippedByReason map[Reason]int `json:\"skipped_by_reason\"`\n}","file":"internal/adopt/report.go","name":"Totals"}}
// Totals sums the file outcomes of a report.
type Totals struct {
	Files           int            `json:"files"`
	Adopted         int            `json:"adopted"`
	Skipped         int            `json:"skipped"`
	Rejected        int            `json:"rejected"`
	Units           int            `json:"units"`
	SkippedByReason map[Reason]int `json:"skipped_by_reason"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Report is the adoption report for one run of sedum adopt.\ntype Report struct {\n\tFiles     []FileReport `json:\"files\"`\n\tTotals    Totals       `json:\"totals\"`\n\tRecording string       `json:\"recording,omitempty\"`\n}","file":"internal/adopt/report.go","name":"Report"}}
// Report is the adoption report for one run of sedum adopt.
type Report struct {
	Files     []FileReport `json:"files"`
	Totals    Totals       `json:"totals"`
	Recording string       `json:"recording,omitempty"`
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// NewReport returns an empty report ready for Add.\nfunc NewReport() *Report {\n\treturn &Report{Files: []FileReport{}, Totals: Totals{SkippedByReason: map[Reason]int{}}}\n}","file":"internal/adopt/report.go","name":"NewReport"}}
// NewReport returns an empty report ready for Add.
func NewReport() *Report {
	return &Report{Files: []FileReport{}, Totals: Totals{SkippedByReason: map[Reason]int{}}}
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Add records one file's outcome and updates the totals.\nfunc (r *Report) Add(f FileReport) {\n\tr.Files = append(r.Files, f)\n\tr.Totals.Files++\n\tr.Totals.Units += f.Units\n\tswitch f.Status {\n\tcase StatusAdopted:\n\t\tr.Totals.Adopted++\n\tcase StatusSkipped:\n\t\tr.Totals.Skipped++\n\tcase StatusRejected:\n\t\tr.Totals.Rejected++\n\t}\n\tif (f.Status == StatusSkipped || f.Status == StatusRejected) && f.Reason != \"\" {\n\t\tr.Totals.SkippedByReason[f.Reason]++\n\t}\n}","file":"internal/adopt/report.go","name":"Add","receiver":"Report"}}
// Add records one file's outcome and updates the totals.
func (r *Report) Add(f FileReport) {
	r.Files = append(r.Files, f)
	r.Totals.Files++
	r.Totals.Units += f.Units
	switch f.Status {
	case StatusAdopted:
		r.Totals.Adopted++
	case StatusSkipped:
		r.Totals.Skipped++
	case StatusRejected:
		r.Totals.Rejected++
	}
	if (f.Status == StatusSkipped || f.Status == StatusRejected) && f.Reason != "" {
		r.Totals.SkippedByReason[f.Reason]++
	}
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// JSON returns the report as indented JSON with a trailing newline.\nfunc (r *Report) JSON() ([]byte, error) {\n\tb, err := json.MarshalIndent(r, \"\", \"  \")\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"marshal report: %w\", err)\n\t}\n\treturn append(b, '\\n'), nil\n}","file":"internal/adopt/report.go","name":"JSON","receiver":"Report"}}
// JSON returns the report as indented JSON with a trailing newline.
func (r *Report) JSON() ([]byte, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal report: %w", err)
	}
	return append(b, '\n'), nil
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// ResidualsOf reports the non-blank residual spans of a file, classifying each\n// as a header, a floating comment, or other residual text.\nfunc ResidualsOf(src []byte, spans []Span, comment genpkg.Comment) []ResidualInfo {\n\tvar out []ResidualInfo\n\tline := 1\n\tfirst := true\n\tunitSeen := false\n\tfor _, s := range spans {\n\t\tisFirst := first && s.Kind != KindHeader\n\t\tif s.Kind != KindHeader {\n\t\t\tfirst = false\n\t\t}\n\t\tif s.Kind == KindUnit {\n\t\t\tunitSeen = true\n\t\t}\n\t\tif s.Kind == KindResidual && !ResidualIsBlank(s) {\n\t\t\tlines := strings.Split(s.Text, \"\\n\")\n\t\t\tstart, end := -1, -1\n\t\t\tallComment := comment.Prefix != \"\"\n\t\t\tfor i, l := range lines {\n\t\t\t\tif strings.TrimSpace(l) == \"\" {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tif start < 0 {\n\t\t\t\t\tstart = i\n\t\t\t\t}\n\t\t\t\tend = i\n\t\t\t\tif !strings.HasPrefix(strings.TrimLeft(l, \" \\t\"), comment.Prefix) {\n\t\t\t\t\tallComment = false\n\t\t\t\t}\n\t\t\t}\n\t\t\treason := ReasonOtherResidual\n\t\t\tif allComment && isFirst && !unitSeen {\n\t\t\t\treason = ReasonHeader\n\t\t\t} else if allComment {\n\t\t\t\treason = ReasonFloatingComment\n\t\t\t}\n\t\t\tout = append(out, ResidualInfo{\n\t\t\t\tStartLine: line + start,\n\t\t\t\tEndLine:   line + end,\n\t\t\t\tFirstLine: strings.TrimSpace(lines[start]),\n\t\t\t\tReason:    reason,\n\t\t\t})\n\t\t}\n\t\tline += strings.Count(s.Text, \"\\n\")\n\t}\n\treturn out\n}","file":"internal/adopt/report.go","name":"ResidualsOf"}}
// ResidualsOf reports the non-blank residual spans of a file, classifying each
// as a header, a floating comment, or other residual text.
func ResidualsOf(src []byte, spans []Span, comment genpkg.Comment) []ResidualInfo {
	var out []ResidualInfo
	line := 1
	first := true
	unitSeen := false
	for _, s := range spans {
		isFirst := first && s.Kind != KindHeader
		if s.Kind != KindHeader {
			first = false
		}
		if s.Kind == KindUnit {
			unitSeen = true
		}
		if s.Kind == KindResidual && !ResidualIsBlank(s) {
			lines := strings.Split(s.Text, "\n")
			start, end := -1, -1
			allComment := comment.Prefix != ""
			for i, l := range lines {
				if strings.TrimSpace(l) == "" {
					continue
				}
				if start < 0 {
					start = i
				}
				end = i
				if !strings.HasPrefix(strings.TrimLeft(l, " \t"), comment.Prefix) {
					allComment = false
				}
			}
			reason := ReasonOtherResidual
			if allComment && isFirst && !unitSeen {
				reason = ReasonHeader
			} else if allComment {
				reason = ReasonFloatingComment
			}
			out = append(out, ResidualInfo{
				StartLine: line + start,
				EndLine:   line + end,
				FirstLine: strings.TrimSpace(lines[start]),
				Reason:    reason,
			})
		}
		line += strings.Count(s.Text, "\n")
	}
	return out
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Text returns a plain-text summary of the report.\nfunc (r *Report) Text() string {\n\tvar b strings.Builder\n\tfmt.Fprintf(&b, \"adopt: %d file(s): %d adopted, %d skipped, %d rejected, %d unit(s)\\n\", r.Totals.Files, r.Totals.Adopted, r.Totals.Skipped, r.Totals.Rejected, r.Totals.Units)\n\tif r.Totals.Skipped+r.Totals.Rejected > 0 {\n\t\tb.WriteString(\"skipped by reason:\\n\")\n\t\tvar reasons []string\n\t\tfor k := range r.Totals.SkippedByReason {\n\t\t\treasons = append(reasons, string(k))\n\t\t}\n\t\tsort.Strings(reasons)\n\t\tfor _, k := range reasons {\n\t\t\tfmt.Fprintf(&b, \"  %s: %d\\n\", k, r.Totals.SkippedByReason[Reason(k)])\n\t\t}\n\t\tfor _, k := range reasons {\n\t\t\tfmt.Fprintf(&b, \"%s:\\n\", k)\n\t\t\tfor _, f := range r.Files {\n\t\t\t\tif (f.Status != StatusSkipped && f.Status != StatusRejected) || string(f.Reason) != k {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tif len(f.Residuals) == 0 {\n\t\t\t\t\tfmt.Fprintf(&b, \"  %s\\n\", f.Path)\n\t\t\t\t}\n\t\t\t\tif len(f.Residuals) == 1 {\n\t\t\t\t\tres := f.Residuals[0]\n\t\t\t\t\tfmt.Fprintf(&b, \"  %s (lines %d-%d: %s)\\n\", f.Path, res.StartLine, res.EndLine, res.FirstLine)\n\t\t\t\t}\n\t\t\t\tif len(f.Residuals) > 1 {\n\t\t\t\t\tfor _, res := range f.Residuals {\n\t\t\t\t\t\tfmt.Fprintf(&b, \"  %s lines %d-%d: %s\\n\", f.Path, res.StartLine, res.EndLine, res.FirstLine)\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t}\n\tvar order []string\n\tbyRecord := map[string][]string{}\n\tfor _, f := range r.Files {\n\t\tif f.Status != StatusAdopted || f.Record == \"\" {\n\t\t\tcontinue\n\t\t}\n\t\tif _, ok := byRecord[f.Record]; !ok {\n\t\t\torder = append(order, f.Record)\n\t\t}\n\t\tbyRecord[f.Record] = append(byRecord[f.Record], f.Path)\n\t}\n\tif len(order) > 0 {\n\t\tb.WriteString(\"list the adopted files in their imprints so grow --check sees them:\\n\")\n\t\tfor _, id := range order {\n\t\t\tfmt.Fprintf(&b, \"  provenance edit %s\", id)\n\t\t\tfor _, p := range byRecord[id] {\n\t\t\t\tfmt.Fprintf(&b, \" --add-scope %s\", p)\n\t\t\t}\n\t\t\tb.WriteString(\"\\n\")\n\t\t}\n\t}\n\tif r.Recording != \"\" {\n\t\tfmt.Fprintf(&b, \"recording: %s\\n\", r.Recording)\n\t}\n\treturn b.String()\n}","file":"internal/adopt/report.go","name":"Text","receiver":"Report"}}
// Text returns a plain-text summary of the report.
func (r *Report) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "adopt: %d file(s): %d adopted, %d skipped, %d rejected, %d unit(s)\n", r.Totals.Files, r.Totals.Adopted, r.Totals.Skipped, r.Totals.Rejected, r.Totals.Units)
	if r.Totals.Skipped+r.Totals.Rejected > 0 {
		b.WriteString("skipped by reason:\n")
		var reasons []string
		for k := range r.Totals.SkippedByReason {
			reasons = append(reasons, string(k))
		}
		sort.Strings(reasons)
		for _, k := range reasons {
			fmt.Fprintf(&b, "  %s: %d\n", k, r.Totals.SkippedByReason[Reason(k)])
		}
		for _, k := range reasons {
			fmt.Fprintf(&b, "%s:\n", k)
			for _, f := range r.Files {
				if (f.Status != StatusSkipped && f.Status != StatusRejected) || string(f.Reason) != k {
					continue
				}
				if len(f.Residuals) == 0 {
					fmt.Fprintf(&b, "  %s\n", f.Path)
				}
				if len(f.Residuals) == 1 {
					res := f.Residuals[0]
					fmt.Fprintf(&b, "  %s (lines %d-%d: %s)\n", f.Path, res.StartLine, res.EndLine, res.FirstLine)
				}
				if len(f.Residuals) > 1 {
					for _, res := range f.Residuals {
						fmt.Fprintf(&b, "  %s lines %d-%d: %s\n", f.Path, res.StartLine, res.EndLine, res.FirstLine)
					}
				}
			}
		}
	}
	var order []string
	byRecord := map[string][]string{}
	for _, f := range r.Files {
		if f.Status != StatusAdopted || f.Record == "" {
			continue
		}
		if _, ok := byRecord[f.Record]; !ok {
			order = append(order, f.Record)
		}
		byRecord[f.Record] = append(byRecord[f.Record], f.Path)
	}
	if len(order) > 0 {
		b.WriteString("list the adopted files in their imprints so grow --check sees them:\n")
		for _, id := range order {
			fmt.Fprintf(&b, "  provenance edit %s", id)
			for _, p := range byRecord[id] {
				fmt.Fprintf(&b, " --add-scope %s", p)
			}
			b.WriteString("\n")
		}
	}
	if r.Recording != "" {
		fmt.Fprintf(&b, "recording: %s\n", r.Recording)
	}
	return b.String()
}

// /sedum:setMethod
