package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/report.go","path":"encoding/json"},"kwargs":"sha256:443b3226d4abd0e59ca177d9a2cf8ca2e50ffbb002c92650639d586c6b8b4923","body":"sha256:872d7c3e9dee608e79a7f8ef57930d24b5e30f2a57d891a52549179ace45946a"}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/report.go","path":"fmt"},"kwargs":"sha256:63d954c40b6324b4bab71201ae7f2ba84971372534c13f8a2f74659d980eaca5","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/report.go","path":"sort"},"kwargs":"sha256:4a6bf0378c0fb6520aee2b49a5d6080aa4d46bf64f7c7f3ef71e73ef03bceca9","body":"sha256:f197172617ae21ef23b1855b09aafb14571020ae81bff877d63d3ff81f4bf63c"}
import "sort"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/report.go","path":"strings"},"kwargs":"sha256:594280b9690d0729eda26ce8484c769f568242a6dc27740a725d8e2587d56b58","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/report.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:78305e1cd8833df41c4994cfeb4d9c6ebca191b03fcad13c62e340800dae3885","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/adopt/report.go","name":"StatusAdopted"},"kwargs":"sha256:6c1d8d25f0759bbcf5ff73c071d22b8bb7e70779e88ac2ea3a1b8281fabb5a34","body":"sha256:4d7b29e41a78069dde7ef03d379289b86247833541095e7d081655ce3db73c3e"}
// StatusAdopted marks a file that was brought under Sedum.
const StatusAdopted FileStatus = "adopted"

// /sedum:setConst
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/adopt/report.go","name":"StatusSkipped"},"kwargs":"sha256:3f9c677b284a3b1e5b847e70d899f918077b32a83727e014a541fca511ad78dc","body":"sha256:8630b40f2a5d012bd76e4f8fd20b3a9ab426a9f820a414fcbb70113f0a932e13"}
// StatusSkipped marks a file that was left alone and not adopted.
const StatusSkipped FileStatus = "skipped"

// /sedum:setConst
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/adopt/report.go","name":"StatusRejected"},"kwargs":"sha256:8e66b07ee0cc137321c3df03fa9809a1876863e1c4e5067c390ba9d079fa3a54","body":"sha256:08155869ec1e558d53275e5eec8366be0258edca97175b71713c08d28222aade"}
// StatusRejected marks a file whose adoption was refused.
const StatusRejected FileStatus = "rejected"

// /sedum:setConst
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/report.go","name":"FileStatus"},"kwargs":"sha256:35e7ebf6bd144af85909d1025120c606ccf95626c882b71772daef2e3049a20c","body":"sha256:801f52ab8c3e2eacc089c6300661f716d5c953934c5de4e424700c9322b7ada1"}
// FileStatus is the outcome recorded for one file in the adoption report.
type FileStatus string

// /sedum:setType
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/report.go","name":"ResidualInfo"},"kwargs":"sha256:0fdd95421b1f5d8981219d27b2fe71576e232b31cbc199b51671277aedf31b7a","body":"sha256:84f44dfc82173658057a75eab6786fbab29fba1365cc29a0c9698c3f4414e602"}
// ResidualInfo describes one non-blank residual span of a file.
type ResidualInfo struct {
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	FirstLine string `json:"first_line"`
	Reason    Reason `json:"reason"`
}

// /sedum:setType
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/report.go","name":"FileReport"},"kwargs":"sha256:a898e6a7da3a2564cf2638638a742b4ef74217105057fc143440ae3a2e3aae0d","body":"sha256:bc857bcd1db9b0bc25803ec59e03ce8c814141e7ea94d06f1da5379dafee9c22"}
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
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/report.go","name":"Totals"},"kwargs":"sha256:2536a94afc2c077c4a2474e4e7cb87767b411235a5dfdcd6127f032cd772953c","body":"sha256:de681f7905d7126fbf159aa751c0e302f2a08ec9d28a3a0238fa47fbb5790f47"}
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
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/report.go","name":"Report"},"kwargs":"sha256:6dfc01522c31fd87be5791f4532de0a7a45b0892e0a9c75d46dcebdeb3295e08","body":"sha256:77858c86b81830a34d979e1a40dc8dfc9e0d0c0dbb846439eb15ee6d6d46e4ea"}
// Report is the adoption report for one run of sedum adopt.
type Report struct {
	Files     []FileReport `json:"files"`
	Totals    Totals       `json:"totals"`
	Recording string       `json:"recording,omitempty"`
}

// /sedum:setType
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/report.go","name":"NewReport"},"kwargs":"sha256:148fcac3996fced71978cd9220b90e60b6ed208b85fb299d1b264ea1f9f43736","body":"sha256:6c4f9c12caaa871863280cd137278b2ee635b4cfb77a540bba0c4ee4ea8b54bb"}
// NewReport returns an empty report ready for Add.
func NewReport() *Report {
	return &Report{Files: []FileReport{}, Totals: Totals{SkippedByReason: map[Reason]int{}}}
}

// /sedum:setFunction
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/report.go","name":"Add","receiver":"Report"},"kwargs":"sha256:0369c42840362ef1d3c31b04c68d62dc8a961edb932817a8542993675504c9db","body":"sha256:028c537b6536d03b90263c85950a08d5dc4d66670818634a1e72470461b82211"}
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
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/report.go","name":"JSON","receiver":"Report"},"kwargs":"sha256:74e0f63367b166b67c5be0de7f49451f2578dc801085ca87c3b34a722dacd2ee","body":"sha256:33e7c833ccd65939c04e45d4151f9c242d36f76bf6dcb3a987ba549b04d27d5a"}
// JSON returns the report as indented JSON with a trailing newline.
func (r *Report) JSON() ([]byte, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal report: %w", err)
	}
	return append(b, '\n'), nil
}

// /sedum:setMethod
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/report.go","name":"ResidualsOf"},"kwargs":"sha256:e7613acd41d19060714bb40c2e88821874f58da59660d76b87523781b1fc0538","body":"sha256:e8a6cab75ab0ee986901c20f1c1d78aa929ade72bf1c66a3cb6746e91a2c90a2"}
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
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/report.go","name":"Text","receiver":"Report"},"kwargs":"sha256:b9410f9e5f809e2c942b77b8f3f99b4172dbe0b86ac077494544c165ef5ee873","body":"sha256:258d9158638510c81f6966ffa20dfb4c699f45b5730fd2e4763be69c3a9c53bd"}
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
