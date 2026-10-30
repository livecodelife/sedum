package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/span_test.go","path":"errors"},"kwargs":"sha256:0484934529edeaf0c79ada3cc4ff57c3005a6e4cccdb0e70bf8dfa90ecb5e91a","body":"sha256:2b718237ba0a56467c8fd494cc131e942ef410e6228c03b644d70476c71862dc"}
import "errors"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/span_test.go","path":"strings"},"kwargs":"sha256:0fbfe1ef5e471e09ad7fc9f7694098a92a881e017e2db148061909ca8f58f732","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/span_test.go","path":"testing"},"kwargs":"sha256:ae443d866e814c557184998fff088e5787a3536a603cdca6db97553ec46e42a6","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/span_test.go","name":"TestCoverage"},"kwargs":"sha256:cff6728caf30f7fe1b465ba43bd878e3e7b852362de4900e90ddfaa46e7d3187","body":"sha256:f42f132095be607e3902a9d1c0f98b0edc35efdb3d611db6bc08b567fa0eb427"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/span_test.go","name":"TestValid"},"kwargs":"sha256:e0566b9fbaa6a657685f9045e399c0e8bb70e73e708cd6a12e3370bb792bc331","body":"sha256:29c3d35b283ef7afb27eab42bf872f39907cb37b1f63251fadcd746b3e038858"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/span_test.go","name":"TestCheckSpans"},"kwargs":"sha256:70c930fee8e96394b8886e0053d1dbd8136a11a6cddf455aaf752728c5d5041e","body":"sha256:9212d79899778c7c820b50237d542470e79d63c72802ea384523fbf4d6f15a56"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/span_test.go","name":"TestResidualIsBlank"},"kwargs":"sha256:7d3a7f1a1a26f0b9b4729af4f0b3147e7630e0bc2a27b9db11b9ce3a99beb096","body":"sha256:d6889a2bb63e771acb5ea805474b8a3922845d28a884a25b7f18e48c6fc67f8f"}
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
