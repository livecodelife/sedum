package growlog

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/render.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/render.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/render.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/render.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/render.go","path":"text/tabwriter"}}
import "text/tabwriter"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setMethod {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// JSON renders the summary as indented JSON with a trailing newline.\nfunc (s Summary) JSON() ([]byte, error) {\n\tb, err := json.MarshalIndent(s, \"\", \"  \")\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn append(b, '\\n'), nil\n}","file":"internal/growlog/render.go","name":"JSON","receiver":"Summary"}}
// JSON renders the summary as indented JSON with a trailing newline.
func (s Summary) JSON() ([]byte, error) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Text renders the summary as a table: one row per model and a total row. A\n// run with no model shows as (none).\nfunc (s Summary) Text() string {\n\tvar b strings.Builder\n\tfmt.Fprintf(&b, \"grow log: %d line(s), %d malformed\\n\", s.Entries, s.Malformed)\n\tvar buf bytes.Buffer\n\tw := tabwriter.NewWriter(&buf, 0, 8, 2, ' ', 0)\n\tfmt.Fprintln(w, \"model\\tgrows\\tfailed\\tcalls\\tretries\\trepairs\\ttokens in/out\\tliteral\")\n\trow := func(name string, m ModelStats) {\n\t\tif name == \"\" {\n\t\t\tname = \"(none)\"\n\t\t}\n\t\tfmt.Fprintf(w, \"%s\\t%d\\t%d\\t%d\\t%d\\t%d\\t%d/%d\\t%.0f%%\\n\", name, m.Grows, m.Failed, m.Calls, m.Retries, m.Repairs, m.PromptTokens, m.CompletionTokens, m.LiteralShare()*100)\n\t}\n\tfor _, m := range s.Models {\n\t\trow(m.Model, m)\n\t}\n\trow(\"total\", s.Total)\n\tw.Flush()\n\tb.Write(buf.Bytes())\n\treturn b.String()\n}","file":"internal/growlog/render.go","name":"Text","receiver":"Summary"}}
// Text renders the summary as a table: one row per model and a total row. A
// run with no model shows as (none).
func (s Summary) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "grow log: %d line(s), %d malformed\n", s.Entries, s.Malformed)
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "model\tgrows\tfailed\tcalls\tretries\trepairs\ttokens in/out\tliteral")
	row := func(name string, m ModelStats) {
		if name == "" {
			name = "(none)"
		}
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\t%d/%d\t%.0f%%\n", name, m.Grows, m.Failed, m.Calls, m.Retries, m.Repairs, m.PromptTokens, m.CompletionTokens, m.LiteralShare()*100)
	}
	for _, m := range s.Models {
		row(m.Model, m)
	}
	row("total", s.Total)
	w.Flush()
	b.Write(buf.Bytes())
	return b.String()
}

// /sedum:setMethod
