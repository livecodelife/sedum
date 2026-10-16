package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// StatusAdopted marks a file that was adopted.\nconst StatusAdopted FileStatus = \"adopted\"","file":"internal/adopt/report.go","name":"StatusAdopted"}}
// StatusAdopted marks a file that was adopted.
const StatusAdopted FileStatus = "adopted"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// StatusSkipped marks a file that was skipped.\nconst StatusSkipped FileStatus = \"skipped\"","file":"internal/adopt/report.go","name":"StatusSkipped"}}
// StatusSkipped marks a file that was skipped.
const StatusSkipped FileStatus = "skipped"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// StatusRejected marks a file that was rejected.\nconst StatusRejected FileStatus = \"rejected\"","file":"internal/adopt/report.go","name":"StatusRejected"}}
// StatusRejected marks a file that was rejected.
const StatusRejected FileStatus = "rejected"

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// FileStatus is the outcome of adopting one file.\ntype FileStatus string","file":"internal/adopt/report.go","name":"FileStatus"}}
// FileStatus is the outcome of adopting one file.
type FileStatus string

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// ResidualInfo describes one non-blank residual span left in a file.\ntype ResidualInfo struct {\n\tStartLine int    `json:\"start_line\"`\n\tEndLine   int    `json:\"end_line\"`\n\tFirstLine string `json:\"first_line\"`\n\tReason    Reason `json:\"reason\"`\n}","file":"internal/adopt/report.go","name":"ResidualInfo"}}
// ResidualInfo describes one non-blank residual span left in a file.
type ResidualInfo struct {
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	FirstLine string `json:"first_line"`
	Reason    Reason `json:"reason"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// FileReport is the adoption outcome for one file.\ntype FileReport struct {\n\tPath      string         `json:\"path\"`\n\tPackage   string         `json:\"package,omitempty\"`\n\tRecord    string         `json:\"record,omitempty\"`\n\tStatus    FileStatus     `json:\"status\"`\n\tReason    Reason         `json:\"reason,omitempty\"`\n\tDetail    string         `json:\"detail,omitempty\"`\n\tUnits     int            `json:\"units\"`\n\tResiduals []ResidualInfo `json:\"residuals,omitempty\"`\n}","file":"internal/adopt/report.go","name":"FileReport"}}
// FileReport is the adoption outcome for one file.
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
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Totals summarises the file outcomes in a report.\ntype Totals struct {\n\tFiles           int            `json:\"files\"`\n\tAdopted         int            `json:\"adopted\"`\n\tSkipped         int            `json:\"skipped\"`\n\tRejected        int            `json:\"rejected\"`\n\tUnits           int            `json:\"units\"`\n\tSkippedByReason map[Reason]int `json:\"skipped_by_reason\"`\n}","file":"internal/adopt/report.go","name":"Totals"}}
// Totals summarises the file outcomes in a report.
type Totals struct {
	Files           int            `json:"files"`
	Adopted         int            `json:"adopted"`
	Skipped         int            `json:"skipped"`
	Rejected        int            `json:"rejected"`
	Units           int            `json:"units"`
	SkippedByReason map[Reason]int `json:"skipped_by_reason"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Report is the adoption report for a run.\ntype Report struct {\n\tFiles     []FileReport `json:\"files\"`\n\tTotals    Totals       `json:\"totals\"`\n\tRecording string       `json:\"recording,omitempty\"`\n}","file":"internal/adopt/report.go","name":"Report"}}
// Report is the adoption report for a run.
type Report struct {
	Files     []FileReport `json:"files"`
	Totals    Totals       `json:"totals"`
	Recording string       `json:"recording,omitempty"`
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// NewReport returns an empty report.\nfunc NewReport() *Report {\n\treturn &Report{Files: []FileReport{}, Totals: Totals{SkippedByReason: map[Reason]int{}}}\n}","file":"internal/adopt/report.go","name":"NewReport"}}
// NewReport returns an empty report.
func NewReport() *Report {
	return &Report{Files: []FileReport{}, Totals: Totals{SkippedByReason: map[Reason]int{}}}
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Add records one file's outcome and updates the totals.\nfunc (r *Report) Add(f FileReport) {\n\tif r.Totals.SkippedByReason == nil {\n\t\tr.Totals.SkippedByReason = map[Reason]int{}\n\t}\n\tr.Files = append(r.Files, f)\n\tr.Totals.Files++\n\tr.Totals.Units += f.Units\n\tswitch f.Status {\n\tcase StatusAdopted:\n\t\tr.Totals.Adopted++\n\tcase StatusSkipped:\n\t\tr.Totals.Skipped++\n\tcase StatusRejected:\n\t\tr.Totals.Rejected++\n\t}\n\tif (f.Status == StatusSkipped || f.Status == StatusRejected) && f.Reason != \"\" {\n\t\tr.Totals.SkippedByReason[f.Reason]++\n\t}\n}","file":"internal/adopt/report.go","name":"Add","receiver":"Report"}}
// Add records one file's outcome and updates the totals.
func (r *Report) Add(f FileReport) {
	if r.Totals.SkippedByReason == nil {
		r.Totals.SkippedByReason = map[Reason]int{}
	}
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
// sedum:setMethod {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// JSON returns the report as indented JSON with a trailing newline.\nfunc (r *Report) JSON() ([]byte, error) {\n\tb, err := json.MarshalIndent(r, \"\", \"  \")\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"encode report: %w\", err)\n\t}\n\treturn append(b, '\\n'), nil\n}","file":"internal/adopt/report.go","name":"JSON","receiver":"Report"}}
// JSON returns the report as indented JSON with a trailing newline.
func (r *Report) JSON() ([]byte, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode report: %w", err)
	}
	return append(b, '\n'), nil
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// ResidualsOf describes the non-blank residual spans of a split file. A\n// residual before any other non-header span whose lines all carry the comment\n// prefix is a header; later all-comment residuals are floating comments; all\n// other residuals are other_residual. Header spans are never reported.\nfunc ResidualsOf(src []byte, spans []Span, comment genpkg.Comment) []ResidualInfo {\n\tvar out []ResidualInfo\n\tline := 1\n\tseenAny := false\n\tfor _, sp := range spans {\n\t\tif sp.Kind == KindHeader {\n\t\t\tline += strings.Count(sp.Text, \"\\n\")\n\t\t\tcontinue\n\t\t}\n\t\tif sp.Kind == KindResidual && !ResidualIsBlank(sp) {\n\t\t\tlines := strings.Split(sp.Text, \"\\n\")\n\t\t\tfirst, last := -1, -1\n\t\t\tfor i, l := range lines {\n\t\t\t\tif strings.TrimSpace(l) != \"\" {\n\t\t\t\t\tif first < 0 {\n\t\t\t\t\t\tfirst = i\n\t\t\t\t\t}\n\t\t\t\t\tlast = i\n\t\t\t\t}\n\t\t\t}\n\t\t\tallComment := comment.Prefix != \"\"\n\t\t\tfor _, l := range lines {\n\t\t\t\tif strings.TrimSpace(l) == \"\" {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tif !strings.HasPrefix(strings.TrimLeft(l, \" \\t\\r\"), comment.Prefix) {\n\t\t\t\t\tallComment = false\n\t\t\t\t}\n\t\t\t}\n\t\t\treason := ReasonOtherResidual\n\t\t\tif allComment {\n\t\t\t\tif seenAny {\n\t\t\t\t\treason = ReasonFloatingComment\n\t\t\t\t} else {\n\t\t\t\t\treason = ReasonHeader\n\t\t\t\t}\n\t\t\t}\n\t\t\tout = append(out, ResidualInfo{\n\t\t\t\tStartLine: line + first,\n\t\t\t\tEndLine:   line + last,\n\t\t\t\tFirstLine: strings.TrimSpace(lines[first]),\n\t\t\t\tReason:    reason,\n\t\t\t})\n\t\t}\n\t\tline += strings.Count(sp.Text, \"\\n\")\n\t\tseenAny = true\n\t}\n\treturn out\n}","file":"internal/adopt/report.go","name":"ResidualsOf"}}
// ResidualsOf describes the non-blank residual spans of a split file. A
// residual before any other non-header span whose lines all carry the comment
// prefix is a header; later all-comment residuals are floating comments; all
// other residuals are other_residual. Header spans are never reported.
func ResidualsOf(src []byte, spans []Span, comment genpkg.Comment) []ResidualInfo {
	var out []ResidualInfo
	line := 1
	seenAny := false
	for _, sp := range spans {
		if sp.Kind == KindHeader {
			line += strings.Count(sp.Text, "\n")
			continue
		}
		if sp.Kind == KindResidual && !ResidualIsBlank(sp) {
			lines := strings.Split(sp.Text, "\n")
			first, last := -1, -1
			for i, l := range lines {
				if strings.TrimSpace(l) != "" {
					if first < 0 {
						first = i
					}
					last = i
				}
			}
			allComment := comment.Prefix != ""
			for _, l := range lines {
				if strings.TrimSpace(l) == "" {
					continue
				}
				if !strings.HasPrefix(strings.TrimLeft(l, " \t\r"), comment.Prefix) {
					allComment = false
				}
			}
			reason := ReasonOtherResidual
			if allComment {
				if seenAny {
					reason = ReasonFloatingComment
				} else {
					reason = ReasonHeader
				}
			}
			out = append(out, ResidualInfo{
				StartLine: line + first,
				EndLine:   line + last,
				FirstLine: strings.TrimSpace(lines[first]),
				Reason:    reason,
			})
		}
		line += strings.Count(sp.Text, "\n")
		seenAny = true
	}
	return out
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:c98f4cc867b91dee92b61349de06027c0f3fcd6d36a850f40c46fa1f500be8bc","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Text renders the report as plain text.\nfunc (r *Report) Text() string {\n\tvar b strings.Builder\n\tfmt.Fprintf(&b, \"adopt: %d file(s): %d adopted, %d skipped, %d rejected, %d unit(s)\\n\", r.Totals.Files, r.Totals.Adopted, r.Totals.Skipped, r.Totals.Rejected, r.Totals.Units)\n\tif r.Totals.Skipped+r.Totals.Rejected > 0 {\n\t\treasons := make([]string, 0, len(r.Totals.SkippedByReason))\n\t\tfor k := range r.Totals.SkippedByReason {\n\t\t\treasons = append(reasons, string(k))\n\t\t}\n\t\tsort.Strings(reasons)\n\t\tb.WriteString(\"skipped by reason:\\n\")\n\t\tfor _, k := range reasons {\n\t\t\tfmt.Fprintf(&b, \"  %s: %d\\n\", k, r.Totals.SkippedByReason[Reason(k)])\n\t\t}\n\t\tfor _, k := range reasons {\n\t\t\tfmt.Fprintf(&b, \"%s:\\n\", k)\n\t\t\tfor _, f := range r.Files {\n\t\t\t\tif (f.Status != StatusSkipped && f.Status != StatusRejected) || string(f.Reason) != k {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tswitch len(f.Residuals) {\n\t\t\t\tcase 0:\n\t\t\t\t\tfmt.Fprintf(&b, \"  %s\\n\", f.Path)\n\t\t\t\tcase 1:\n\t\t\t\t\tres := f.Residuals[0]\n\t\t\t\t\tfmt.Fprintf(&b, \"  %s (lines %d-%d: %s)\\n\", f.Path, res.StartLine, res.EndLine, res.FirstLine)\n\t\t\t\tdefault:\n\t\t\t\t\tfor _, res := range f.Residuals {\n\t\t\t\t\t\tfmt.Fprintf(&b, \"  %s lines %d-%d: %s\\n\", f.Path, res.StartLine, res.EndLine, res.FirstLine)\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t}\n\tif r.Recording != \"\" {\n\t\tfmt.Fprintf(&b, \"recording: %s\\n\", r.Recording)\n\t}\n\treturn b.String()\n}","file":"internal/adopt/report.go","name":"Text","receiver":"Report"}}
// Text renders the report as plain text.
func (r *Report) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "adopt: %d file(s): %d adopted, %d skipped, %d rejected, %d unit(s)\n", r.Totals.Files, r.Totals.Adopted, r.Totals.Skipped, r.Totals.Rejected, r.Totals.Units)
	if r.Totals.Skipped+r.Totals.Rejected > 0 {
		reasons := make([]string, 0, len(r.Totals.SkippedByReason))
		for k := range r.Totals.SkippedByReason {
			reasons = append(reasons, string(k))
		}
		sort.Strings(reasons)
		b.WriteString("skipped by reason:\n")
		for _, k := range reasons {
			fmt.Fprintf(&b, "  %s: %d\n", k, r.Totals.SkippedByReason[Reason(k)])
		}
		for _, k := range reasons {
			fmt.Fprintf(&b, "%s:\n", k)
			for _, f := range r.Files {
				if (f.Status != StatusSkipped && f.Status != StatusRejected) || string(f.Reason) != k {
					continue
				}
				switch len(f.Residuals) {
				case 0:
					fmt.Fprintf(&b, "  %s\n", f.Path)
				case 1:
					res := f.Residuals[0]
					fmt.Fprintf(&b, "  %s (lines %d-%d: %s)\n", f.Path, res.StartLine, res.EndLine, res.FirstLine)
				default:
					for _, res := range f.Residuals {
						fmt.Fprintf(&b, "  %s lines %d-%d: %s\n", f.Path, res.StartLine, res.EndLine, res.FirstLine)
					}
				}
			}
		}
	}
	if r.Recording != "" {
		fmt.Fprintf(&b, "recording: %s\n", r.Recording)
	}
	return b.String()
}

// /sedum:setMethod
