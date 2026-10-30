package growlog

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/render.go","path":"bytes"},"kwargs":"sha256:6b278364ab686969183af1b60befc46e273d6653ce7c99260e763fa6abbadf6b","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/render.go","path":"encoding/json"},"kwargs":"sha256:0c8be4f3d87195d187be7fa1e85b9f19326f62fe3b21ad9c38b6840da259a1d5","body":"sha256:872d7c3e9dee608e79a7f8ef57930d24b5e30f2a57d891a52549179ace45946a"}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/render.go","path":"fmt"},"kwargs":"sha256:5ec71a894310419e8f2e183cac54ef79dd53cd96814d59438f6cd530afae3439","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/render.go","path":"strings"},"kwargs":"sha256:e3e9186a84f1771b3e40aa0d6b566264fea46a7f3bb88724fdaf11b45d2fcd9e","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/render.go","path":"text/tabwriter"},"kwargs":"sha256:0e4716047b08aaafa4fc50cb6cf4c660ed5cca3c912353bd5e93b2cc91b6793e","body":"sha256:11be499c9f9e23f4957f08722e9ed4af814f9d2de02b151a195600cda7fe6b92"}
import "text/tabwriter"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/growlog/render.go","name":"JSON","receiver":"Summary"},"kwargs":"sha256:39a2c7003f30837dcb58371c221fbff042f9196b3a09f03aa61886bb86580999","body":"sha256:e6bce808e27cc6d56fff7044ea8e4e9202ff5aa70a30821e4541657c082030d2"}
// JSON renders the summary as indented JSON with a trailing newline.
func (s Summary) JSON() ([]byte, error) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

// /sedum:setMethod
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/growlog/render.go","name":"Text","receiver":"Summary"},"kwargs":"sha256:04223f9a5ff45789e1f63d8d40e9f5f335e8cf77d79c8c48da8a432cca5a4d4a","body":"sha256:078224136440c1790b120a75ad4b459e962303ae7f879ca330cf322e21bbbee9"}
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
