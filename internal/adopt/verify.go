package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/verify.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/verify.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/verify.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// StripMarkers removes every Sedum marker line from content, leaving all other lines unchanged.\nfunc StripMarkers(comment genpkg.Comment, content string) string {\n\tif comment.Prefix == \"\" {\n\t\treturn content\n\t}\n\topen := comment.Prefix + \" sedum:\"\n\tclosing := comment.Prefix + \" /sedum:\"\n\tvar b strings.Builder\n\tfor _, line := range strings.SplitAfter(content, \"\\n\") {\n\t\ttrimmed := strings.TrimLeft(line, \" \\t\")\n\t\tif strings.HasPrefix(trimmed, open) || strings.HasPrefix(trimmed, closing) {\n\t\t\tcontinue\n\t\t}\n\t\tb.WriteString(line)\n\t}\n\treturn b.String()\n}","file":"internal/adopt/verify.go","name":"StripMarkers"}}
// StripMarkers removes every Sedum marker line from content, leaving all other lines unchanged.
func StripMarkers(comment genpkg.Comment, content string) string {
	if comment.Prefix == "" {
		return content
	}
	open := comment.Prefix + " sedum:"
	closing := comment.Prefix + " /sedum:"
	var b strings.Builder
	for _, line := range strings.SplitAfter(content, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmed, open) || strings.HasPrefix(trimmed, closing) {
			continue
		}
		b.WriteString(line)
	}
	return b.String()
}

// /sedum:setFunction
// sedum:setType {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// VerifyError reports the first line where a round trip differs.\ntype VerifyError struct {\n\tLine int\n\tWant string\n\tGot  string\n}","file":"internal/adopt/verify.go","name":"VerifyError"}}
// VerifyError reports the first line where a round trip differs.
type VerifyError struct {
	Line int
	Want string
	Got  string
}

// /sedum:setType
// sedum:setMethod {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Error describes the differing line.\nfunc (e *VerifyError) Error() string {\n\treturn fmt.Sprintf(\"round trip: line %d differs: want \\\"%s\\\" got \\\"%s\\\"\", e.Line, e.Want, e.Got)\n}","file":"internal/adopt/verify.go","name":"Error","receiver":"VerifyError"}}
// Error describes the differing line.
func (e *VerifyError) Error() string {
	return fmt.Sprintf("round trip: line %d differs: want \"%s\" got \"%s\"", e.Line, e.Want, e.Got)
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Verify checks that adopted, with its Sedum markers removed, matches original line by line.\n// Blank lines are ignored and trailing whitespace is trimmed.\nfunc Verify(comment genpkg.Comment, original, adopted string) error {\n\tsignificant := func(s string) []string {\n\t\tvar out []string\n\t\tfor _, line := range strings.Split(s, \"\\n\") {\n\t\t\tline = strings.TrimRight(line, \" \\t\\r\")\n\t\t\tif line != \"\" {\n\t\t\t\tout = append(out, line)\n\t\t\t}\n\t\t}\n\t\treturn out\n\t}\n\twant := significant(original)\n\tgot := significant(StripMarkers(comment, adopted))\n\tfor i := 0; i < len(want) || i < len(got); i++ {\n\t\tw := \"<end of file>\"\n\t\tif i < len(want) {\n\t\t\tw = want[i]\n\t\t}\n\t\tg := \"<end of file>\"\n\t\tif i < len(got) {\n\t\t\tg = got[i]\n\t\t}\n\t\tif w != g {\n\t\t\treturn &VerifyError{Line: i + 1, Want: w, Got: g}\n\t\t}\n\t}\n\treturn nil\n}","file":"internal/adopt/verify.go","name":"Verify"}}
// Verify checks that adopted, with its Sedum markers removed, matches original line by line.
// Blank lines are ignored and trailing whitespace is trimmed.
func Verify(comment genpkg.Comment, original, adopted string) error {
	significant := func(s string) []string {
		var out []string
		for _, line := range strings.Split(s, "\n") {
			line = strings.TrimRight(line, " \t\r")
			if line != "" {
				out = append(out, line)
			}
		}
		return out
	}
	want := significant(original)
	got := significant(StripMarkers(comment, adopted))
	for i := 0; i < len(want) || i < len(got); i++ {
		w := "<end of file>"
		if i < len(want) {
			w = want[i]
		}
		g := "<end of file>"
		if i < len(got) {
			g = got[i]
		}
		if w != g {
			return &VerifyError{Line: i + 1, Want: w, Got: g}
		}
	}
	return nil
}

// /sedum:setFunction
