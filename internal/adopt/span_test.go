package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:42fdf246ebb6ab5bf7b3ee907f89f14c49b671a5bc996a27f4d0f9cc2397e626","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/span_test.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:42fdf246ebb6ab5bf7b3ee907f89f14c49b671a5bc996a27f4d0f9cc2397e626","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/span_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:42fdf246ebb6ab5bf7b3ee907f89f14c49b671a5bc996a27f4d0f9cc2397e626","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/span_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:42fdf246ebb6ab5bf7b3ee907f89f14c49b671a5bc996a27f4d0f9cc2397e626","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestCoverage(t *testing.T) {\n\tsrc := []byte(\"a\\nfunc b\\n\")\n\ttests := []struct {\n\t\tname   string\n\t\tspans  []Span\n\t\toffset int\n\t\tline   int\n\t}{\n\t\t{\"exact\", []Span{{Kind: KindHeader, Text: \"a\\n\"}, {Kind: KindUnit, Text: \"func b\\n\"}}, -1, 0},\n\t\t{\"missing byte\", []Span{{Kind: KindHeader, Text: \"a\\n\"}, {Kind: KindUnit, Text: \"func \\n\"}}, 7, 2},\n\t\t{\"extra byte\", []Span{{Kind: KindHeader, Text: \"a\\n\"}, {Kind: KindUnit, Text: \"func bb\\n\"}}, 8, 2},\n\t\t{\"changed byte\", []Span{{Kind: KindHeader, Text: \"a\\n\"}, {Kind: KindUnit, Text: \"func c\\n\"}}, 7, 2},\n\t\t{\"prefix\", []Span{{Kind: KindHeader, Text: \"a\\nfunc \"}}, 7, 2},\n\t\t{\"empty spans\", nil, 0, 1},\n\t}\n\tfor _, tc := range tests {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\terr := Coverage(src, tc.spans)\n\t\t\tif tc.offset < 0 {\n\t\t\t\tif err != nil {\n\t\t\t\t\tt.Fatalf(\"unexpected error: %v\", err)\n\t\t\t\t}\n\t\t\t\treturn\n\t\t\t}\n\t\t\tvar ce *CoverageError\n\t\t\tif !errors.As(err, &ce) {\n\t\t\t\tt.Fatalf(\"want *CoverageError, got %v\", err)\n\t\t\t}\n\t\t\tif ce.Offset != tc.offset || ce.Line != tc.line {\n\t\t\t\tt.Fatalf(\"got offset %d line %d, want offset %d line %d\", ce.Offset, ce.Line, tc.offset, tc.line)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/span_test.go","name":"TestCoverage"}}
func TestCoverage(t *testing.T) {
	src := []byte("a\nfunc b\n")
	tests := []struct {
		name   string
		spans  []Span
		offset int
		line   int
	}{
		{"exact", []Span{{Kind: KindHeader, Text: "a\n"}, {Kind: KindUnit, Text: "func b\n"}}, -1, 0},
		{"missing byte", []Span{{Kind: KindHeader, Text: "a\n"}, {Kind: KindUnit, Text: "func \n"}}, 7, 2},
		{"extra byte", []Span{{Kind: KindHeader, Text: "a\n"}, {Kind: KindUnit, Text: "func bb\n"}}, 8, 2},
		{"changed byte", []Span{{Kind: KindHeader, Text: "a\n"}, {Kind: KindUnit, Text: "func c\n"}}, 7, 2},
		{"prefix", []Span{{Kind: KindHeader, Text: "a\nfunc "}}, 7, 2},
		{"empty spans", nil, 0, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Coverage(src, tc.spans)
			if tc.offset < 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			var ce *CoverageError
			if !errors.As(err, &ce) {
				t.Fatalf("want *CoverageError, got %v", err)
			}
			if ce.Offset != tc.offset || ce.Line != tc.line {
				t.Fatalf("got offset %d line %d, want offset %d line %d", ce.Offset, ce.Line, tc.offset, tc.line)
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:42fdf246ebb6ab5bf7b3ee907f89f14c49b671a5bc996a27f4d0f9cc2397e626","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestValid(t *testing.T) {\n\ttests := []struct {\n\t\tname    string\n\t\tspan    Span\n\t\twantErr bool\n\t}{\n\t\t{\"unit ok\", Span{Kind: KindUnit, Text: \"x\", Action: \"setFunction\"}, false},\n\t\t{\"unit no action\", Span{Kind: KindUnit, Text: \"x\"}, true},\n\t\t{\"unit no text\", Span{Kind: KindUnit, Action: \"setFunction\"}, true},\n\t\t{\"header ok\", Span{Kind: KindHeader, Text: \"x\"}, false},\n\t\t{\"header with action\", Span{Kind: KindHeader, Text: \"x\", Action: \"setFunction\"}, true},\n\t\t{\"residual with kwargs\", Span{Kind: KindResidual, Text: \"x\", Kwargs: map[string]any{\"a\": 1}}, true},\n\t\t{\"unknown kind\", Span{Kind: \"bogus\", Text: \"x\"}, true},\n\t}\n\tfor _, tc := range tests {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\terr := tc.span.Valid()\n\t\t\tif (err != nil) != tc.wantErr {\n\t\t\t\tt.Fatalf(\"Valid() error = %v, wantErr %v\", err, tc.wantErr)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/span_test.go","name":"TestValid"}}
func TestValid(t *testing.T) {
	tests := []struct {
		name    string
		span    Span
		wantErr bool
	}{
		{"unit ok", Span{Kind: KindUnit, Text: "x", Action: "setFunction"}, false},
		{"unit no action", Span{Kind: KindUnit, Text: "x"}, true},
		{"unit no text", Span{Kind: KindUnit, Action: "setFunction"}, true},
		{"header ok", Span{Kind: KindHeader, Text: "x"}, false},
		{"header with action", Span{Kind: KindHeader, Text: "x", Action: "setFunction"}, true},
		{"residual with kwargs", Span{Kind: KindResidual, Text: "x", Kwargs: map[string]any{"a": 1}}, true},
		{"unknown kind", Span{Kind: "bogus", Text: "x"}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.span.Valid()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Valid() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:42fdf246ebb6ab5bf7b3ee907f89f14c49b671a5bc996a27f4d0f9cc2397e626","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestCheckSpans(t *testing.T) {\n\tsrc := []byte(\"a\\nb\\n\")\n\tgood := Span{Kind: KindHeader, Text: \"a\\n\"}\n\ttests := []struct {\n\t\tname    string\n\t\tspans   []Span\n\t\twantErr string\n\t}{\n\t\t{\"ok\", []Span{good, {Kind: KindResidual, Text: \"b\\n\"}}, \"\"},\n\t\t{\"empty spans\", nil, \"no spans\"},\n\t\t{\"bad span index\", []Span{good, {Kind: KindUnit, Text: \"b\\n\"}}, \"span 1:\"},\n\t\t{\"coverage\", []Span{good}, \"coverage\"},\n\t}\n\tfor _, tc := range tests {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\terr := CheckSpans(src, tc.spans)\n\t\t\tif tc.wantErr == \"\" {\n\t\t\t\tif err != nil {\n\t\t\t\t\tt.Fatalf(\"unexpected error: %v\", err)\n\t\t\t\t}\n\t\t\t\treturn\n\t\t\t}\n\t\t\tif err == nil || !strings.Contains(err.Error(), tc.wantErr) {\n\t\t\t\tt.Fatalf(\"got error %v, want it to contain %q\", err, tc.wantErr)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/span_test.go","name":"TestCheckSpans"}}
func TestCheckSpans(t *testing.T) {
	src := []byte("a\nb\n")
	good := Span{Kind: KindHeader, Text: "a\n"}
	tests := []struct {
		name    string
		spans   []Span
		wantErr string
	}{
		{"ok", []Span{good, {Kind: KindResidual, Text: "b\n"}}, ""},
		{"empty spans", nil, "no spans"},
		{"bad span index", []Span{good, {Kind: KindUnit, Text: "b\n"}}, "span 1:"},
		{"coverage", []Span{good}, "coverage"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckSpans(src, tc.spans)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("got error %v, want it to contain %q", err, tc.wantErr)
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:42fdf246ebb6ab5bf7b3ee907f89f14c49b671a5bc996a27f4d0f9cc2397e626","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestResidualIsBlank(t *testing.T) {\n\ttests := []struct {\n\t\tname string\n\t\tspan Span\n\t\twant bool\n\t}{\n\t\t{\"blank residual\", Span{Kind: KindResidual, Text: \" \\n\\t\"}, true},\n\t\t{\"text residual\", Span{Kind: KindResidual, Text: \"x\"}, false},\n\t\t{\"blank header\", Span{Kind: KindHeader, Text: \"  \"}, false},\n\t\t{\"blank unit\", Span{Kind: KindUnit, Text: \" \", Action: \"setFunction\"}, false},\n\t}\n\tfor _, tc := range tests {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\tif got := ResidualIsBlank(tc.span); got != tc.want {\n\t\t\t\tt.Fatalf(\"ResidualIsBlank() = %v, want %v\", got, tc.want)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/span_test.go","name":"TestResidualIsBlank"}}
func TestResidualIsBlank(t *testing.T) {
	tests := []struct {
		name string
		span Span
		want bool
	}{
		{"blank residual", Span{Kind: KindResidual, Text: " \n\t"}, true},
		{"text residual", Span{Kind: KindResidual, Text: "x"}, false},
		{"blank header", Span{Kind: KindHeader, Text: "  "}, false},
		{"blank unit", Span{Kind: KindUnit, Text: " ", Action: "setFunction"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ResidualIsBlank(tc.span); got != tc.want {
				t.Fatalf("ResidualIsBlank() = %v, want %v", got, tc.want)
			}
		})
	}
}

// /sedum:setFunction
